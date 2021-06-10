package server

import (
	spanner "cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

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

func StartServer() {
	fmt.Println("This is the server ppl")
	project := localConfig()

	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.uri())

	if err != nil {
		panic(err)
	}

	defer client.Close()

	client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		uuidAlice, _ := uuid.New().MarshalBinary()
		columns := []string{"Id", "Name", "Money"}
		values := []interface{}{uuidAlice, "Alice", 400}

		mutations := []*spanner.Mutation{
			spanner.InsertOrUpdate("Users", columns, values),
		}

		err := transaction.BufferWrite(mutations)
		if err != nil {
			return err
		}
		return nil
	})

	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users"))
	iterator.Do(func(row *spanner.Row) error {
		println("Row size: " + strconv.Itoa(row.Size()))
		var name string
		var money int

		row.ColumnByName("Name", &name)
		row.ColumnByName("Money", &money)

		println("Name: " + name + " Money: " + strconv.Itoa(money))
		return nil
	})
}
