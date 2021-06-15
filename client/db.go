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

	refreshRate time.Duration
}

type msgUser []Rich
type msgTotalMoney int64
type msgRichest int64
type msgPoorest int64

func (db db) retrieveUsers() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		users := db.getUsers()
		richPeople := usersToRiches(users)
		return msgUser(richPeople)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrieveTotalMoney() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := shared.AggregateMoney(shared.Sum, db.ctx, db.client)
		return msgTotalMoney(money)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrieveRichest() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := shared.AggregateMoney(shared.Max, db.ctx, db.client)
		return msgRichest(money)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrievePoorest() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := shared.AggregateMoney(shared.Min, db.ctx, db.client)
		return msgPoorest(money)
	}

	return tea.Tick(db.refreshRate, retrieve)
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
