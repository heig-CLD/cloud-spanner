package server

import (
	spanner "cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

type user struct {
	Id    []byte `spanner:"Id"`
	Name  string `spanner:"Name"`
	Money int64  `spanner:"Money"`
}

type gcloudConfig struct {
	project  string
	instance string
	database string
}

func (config gcloudConfig) uri() string {
	return fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		config.project,
		config.instance,
		config.database,
	)
}

func localConfig() gcloudConfig {
	return gcloudConfig{
		project:  "noted-episode-316407",
		instance: "test-instance",
		database: "test-database",
	}
}

func DeleteDBContent(ctx context.Context, client *spanner.Client) {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		mut := spanner.Delete("Users", spanner.AllKeys())

		mutations := []*spanner.Mutation{mut}

		err := transaction.BufferWrite(mutations)
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		panic(err)
	}
}

func StartServer() {
	fmt.Println("This is the server ppl")
	project := localConfig()

	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.uri())

	if err != nil {
		panic(err)
	}

	defer client.Close()

	DeleteDBContent(ctx, client)

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		uuidAlice, _ := uuid.New().MarshalBinary()

		mut, err := spanner.InsertOrUpdateStruct("Users", user{
			Id:    uuidAlice,
			Name:  "Alice",
			Money: 400,
		})

		if err != nil {
			return err
		}

		mutations := []*spanner.Mutation{mut}

		err = transaction.BufferWrite(mutations)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		println(err.Error())
	}

	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users"))
	iterator.Do(func(row *spanner.Row) error {
		var user user
		row.ToStruct(&user)
		println("Name: " + user.Name + " Money: " + strconv.FormatInt(user.Money, 10) + " Id: " + string(user.Id))
		return nil
	})
}
