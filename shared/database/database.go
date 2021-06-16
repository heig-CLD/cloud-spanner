package database

import (
	"cloud-spanner/shared"
	"cloud.google.com/go/spanner"
	"context"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type database struct {
	ctx    context.Context
	client *spanner.Client
}

type Transfer struct {
	Amount    int64
	Timestamp time.Time
	FromId    []byte
	From      string
	ToId      []byte
	To        string
}

type Database interface {
	// Clear removes the contents of the database, by deleting all the rows in
	// the Users, Items and Offers tables.
	Clear() error

	GetUsersCount() (int64, error)

	// GetUsersRichest returns the richest users. Only `limit` users are
	// returned.
	GetUsersRichest(limit int) ([]shared.User, error)

	GetMoneyTotal() (int64, error)
	GetMoneyMax() (int64, error)
	GetMoneyMin() (int64, error)

	AddUsers(users []shared.User) error

	GetTransfersCount(bound spanner.TimestampBound) (int64, error)
	GetTransfersLatest(limit int, bound spanner.TimestampBound) ([]Transfer, error)

	TransferRandomly() error
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
			spanner.Delete("Transfers", spanner.AllKeys()),
			spanner.Delete("Users", spanner.AllKeys()),
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

func (db *database) GetUsersCount() (int64, error) {
	query := "SELECT COUNT(*) FROM Users"
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

func (db *database) GetTransfersCount(bound spanner.TimestampBound) (int64, error) {
	query := "SELECT COUNT(*) FROM Transfers"
	transaction := db.client.Single().WithTimestampBound(bound)
	iterator := transaction.Query(db.ctx, spanner.Statement{SQL: query})
	res, err := iterator.Next()
	if err != nil {
		return 0, err
	}
	var sum int64
	err = res.Column(0, &sum)
	return sum, err
}

func (db *database) GetTransfersLatest(limit int, bound spanner.TimestampBound) ([]Transfer, error) {
	t1 := Transfer{
		Amount:    123,
		Timestamp: time.Now(),
		From:      "Salut",
		FromId:    make([]byte, 0),
		To:        "Le monde",
		ToId:      make([]byte, 0),
	}
	t2 := Transfer{
		Amount:    125,
		Timestamp: time.Now(),
		From:      "Marcel",
		FromId:    make([]byte, 0),
		To:        "Rémi",
		ToId:      make([]byte, 0),
	}
	return []Transfer{t1, t2}, nil
}

// TaxesAmount indicates what percentage of a user's net worth may be
// transferred as part of a single transaction. A  random amount, inferior to
// this percentage, will be taken.
const TaxesAmount = 0.4

// transfer takes a random amount of money from the person with the from
// identifier, and transfers it to the person with the to identifier.
func transfer(from shared.User, to shared.User, txn *spanner.ReadWriteTransaction) error {

	// Figure out the transfer direction.
	mostMoney, leastMoney := from, to
	if from.Money < to.Money {
		mostMoney = to
		leastMoney = from
	}

	// Calculate how the money should be transferred.
	transferred := int64(float64(mostMoney.Money) * TaxesAmount * rand.Float64())
	mostMoney.Money -= transferred
	leastMoney.Money += transferred

	m1, err := spanner.UpdateStruct("Users", mostMoney)
	if err != nil {
		return err
	}
	m2, err := spanner.UpdateStruct("Users", leastMoney)
	if err != nil {
		return err
	}
	cols := []string{"Id", "Amount", "FromUserId", "ToUserId", "AtTimestamp"}
	id, _ := uuid.New().MarshalBinary()

	m3 := spanner.Insert("Transfers", cols, []interface{}{
		id, transferred, mostMoney.Id, leastMoney.Id, spanner.CommitTimestamp,
	})

	// Create a mutation with the updates.
	return txn.BufferWrite([]*spanner.Mutation{m1, m2, m3})
}

func (db *database) TransferRandomly() error {
	t := func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		it := txn.Query(ctx, spanner.Statement{
			SQL: "SELECT * FROM Users TABLESAMPLE RESERVOIR (2 ROWS)",
		})
		users := make([]shared.User, 0)
		err := it.Do(func(r *spanner.Row) error {
			var user shared.User
			err := r.ToStruct(&user)
			if err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
		if err != nil {
			return err
		}
		if len(users) != 2 {
			return errors.New("not enough users")
		}
		return transfer(users[0], users[1], txn)
	}
	_, err := db.client.ReadWriteTransaction(db.ctx, t)
	return err
}
