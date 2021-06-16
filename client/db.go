package client

import (
	"cloud-spanner/shared"
	"cloud-spanner/shared/database"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type db struct {
	store database.Database

	refreshRate time.Duration
}

type msgUser []Rich
type msgTotalMoney int64
type msgTotalUsers int64
type msgRichest int64
type msgPoorest int64

func (db db) retrieveUsers() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		users, _ := db.store.GetUsersRichest(20)
		richPeople := usersToRiches(users)
		return msgUser(richPeople)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrieveTotalUsers() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		// TODO: Change this to the correct method
		users, _ := db.store.GetMoneyTotal()
		return msgTotalUsers(users)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrieveTotalMoney() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := db.store.GetMoneyTotal()
		return msgTotalMoney(money)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrieveRichest() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := db.store.GetMoneyMax()
		return msgRichest(money)
	}

	return tea.Tick(db.refreshRate, retrieve)
}

func (db db) retrievePoorest() tea.Cmd {
	retrieve := func(t time.Time) tea.Msg {
		money, _ := db.store.GetMoneyMin()
		return msgPoorest(money)
	}

	return tea.Tick(db.refreshRate, retrieve)
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