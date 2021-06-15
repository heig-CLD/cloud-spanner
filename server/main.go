package server

import (
	"cloud-spanner/shared"
	spanner "cloud.google.com/go/spanner"
	"context"
	_ "fmt"
)

func StartServer() {
	project := shared.LocalConfig()
	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.Uri())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	deleteDBContent(ctx, client)

	_, err = createUsers(ctx, client, 20, 10000)

	if err != nil {
		println(err.Error())
	}

	showUsers(ctx, client)
	showItems(ctx, client)
	showOffers(ctx, client)

}
