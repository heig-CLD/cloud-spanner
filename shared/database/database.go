package database

import (
	"cloud-spanner/shared"
	"cloud.google.com/go/spanner"
	"context"
)

type database struct {
	ctx    context.Context
	client *spanner.Client
}

type Database interface {
	// Clear removes the contents of the database, by deleting all the rows in
	// the Users, Items and Offers tables.
	Clear() error

	// GetUsersRichest returns the richest users. Only `limit` users are
	// returned.
	GetUsersRichest(limit int) ([]shared.User, error)

	GetMoneyTotal() (int64, error)
	GetMoneyMax() (int64, error)
	GetMoneyMin() (int64, error)

	AddUsers(users []shared.User) error
}

func NewDatabase(ctx context.Context, client *spanner.Client) Database {
	return &database{
		ctx:    ctx,
		client: client,
	}
}

func (db *database) Clear() error {
	t := func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		return txn.BufferWrite([]*spanner.Mutation{
			spanner.Delete("Users", spanner.AllKeys()),
			spanner.Delete("Items", spanner.AllKeys()),
			spanner.Delete("Offers", spanner.AllKeys()),
		})
	}
	_, err := db.client.ReadWriteTransaction(db.ctx, t)
	return err
}

func (db *database) GetUsersRichest(limit int) ([]shared.User, error) {
	statement := spanner.Statement{
		SQL: "SELECT * FROM Users ORDER BY Money DESC LIMIT @limit",
		Params: map[string]interface{}{
			"limit": limit,
		}}
	users := make([]shared.User, 0)
	err := db.client.
		Single().
		Query(db.ctx, statement).
		Do(func(row *spanner.Row) error {
			var user shared.User
			err := row.ToStruct(&user)
			users = append(users, user)
			return err
		})
	return users, err
}

func (db *database) getMoneyAggregate(operator string) (int64, error) {
	query := "SELECT " + operator + "(Money) FROM Users"
	transaction := db.client.Single()
	iterator := transaction.Query(db.ctx, spanner.Statement{SQL: query})
	res, err := iterator.Next()
	if err != nil {
		return 0, err
	}
	var sum int64
	err = res.Column(0, &sum)
	return sum, err
}

func (db *database) GetMoneyTotal() (int64, error) {
	return db.getMoneyAggregate("SUM")
}

func (db *database) GetMoneyMax() (int64, error) {
	return db.getMoneyAggregate("MAX")
}

func (db *database) GetMoneyMin() (int64, error) {
	return db.getMoneyAggregate("MIN")
}

func (db *database) AddUsers(users []shared.User) error {
	t := func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		for _, user := range users {
			// Create user
			mut, err := spanner.InsertOrUpdateStruct("Users", user)
			if err != nil {
				return err
			}

			// Give the user a car
			// mut, err = createItem(cars, user.Id)
			// if err != nil {
			//	return err
			//}

			err = txn.BufferWrite([]*spanner.Mutation{mut})
			if err != nil {
				return err
			}
		}
		return nil
	}
	_, err := db.client.ReadWriteTransaction(db.ctx, t)
	return err
}
