package server

import (
	spanner "cloud.google.com/go/spanner"
	"context"
	"fmt"
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

	/*client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Users (Key, Email) VALUES
				('alice', 'alice@heig-vd.ch'),
				('bob', 'bob@heig-vd.ch'),
				('charlie', 'charlie@heig-vd.ch')`,
		}
		updated, err := transaction.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows inserted", updated)
		return nil
	})*/

	row, err := client.Single().ReadRow(ctx, "Users", spanner.Key{"alice"}, []string{"Email"})
	if err != nil {
		panic(err)
	}

	println("3")

	println(row.Size())
	println(row.ColumnName(0))

	var email string
	row.ColumnByName("Email", &email)
	println(email)
}
