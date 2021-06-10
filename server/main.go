package server

import (
	"cloud-spanner/shared"
	spanner "cloud.google.com/go/spanner"
	"context"
	_ "fmt"
	"strconv"
)

func StartServer() {
	project := shared.LocalConfig()

	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.Uri())

	if err != nil {
		panic(err)
	}

	defer client.Close()

	DeleteDBContent(ctx, client)

	_, err = createUsers(ctx, client, 100)

	if err != nil {
		println(err.Error())
	}

	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users"))
	iterator.Do(func(row *spanner.Row) error {
		var user shared.User
		row.ToStruct(&user)
		println("Name: " + user.Name + " Money: " + strconv.FormatInt(user.Money, 10))
		return nil
	})
}
