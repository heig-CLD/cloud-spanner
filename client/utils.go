package client

import (
	"cloud-spanner/shared"
	"context"

	"cloud.google.com/go/spanner"
)

func getUsers(ctx context.Context, client *spanner.Client) []shared.User {
	var users []shared.User

	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users ORDER BY Money DESC"))
	iterator.Do(func(row *spanner.Row) error {
		var user shared.User
		row.ToStruct(&user)
		users = append(users, user)
		return nil
	})

	return users
}

func usersToRiches(users []shared.User) []Rich {
	var mostMoney int64 = 0

	for _, u := range users {
		if u.Money > mostMoney {
			mostMoney = u.Money
		}
	}

	var riches []Rich
	for _, u := range users {
		riches = append(riches, InitializeRich(u.Name, float64(u.Money)/float64(mostMoney)))
	}

	return riches
}
