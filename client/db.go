package client

import (
	"cloud-spanner/shared"
	"context"
	"time"

	"cloud.google.com/go/spanner"
	tea "github.com/charmbracelet/bubbletea"
)

type db struct {
	ctx    context.Context
	client *spanner.Client
}

type userMsg []Rich
type moneyMsg int64

func (db db) retrieveUsers() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		users := db.getUsers()
		richPeople := usersToRiches(users)
		return userMsg(richPeople)
	}

	return tea.Tick(300*time.Millisecond, retrieve)
}

func (db db) retrieveMoney() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := shared.AggregateMoney(db.ctx, db.client)
		return moneyMsg(money)
	}

	return tea.Tick(300*time.Millisecond, retrieve)
}

func (db db) getUsers() []shared.User {
	var users []shared.User

	iterator := db.client.Single().Query(db.ctx, spanner.NewStatement("SELECT * FROM Users ORDER BY Money DESC"))
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
