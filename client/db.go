package client

import (
	"cloud-spanner/shared"
	"cloud-spanner/shared/database"
	"time"

	"cloud.google.com/go/spanner"
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
type msgTransactions []transaction
type msgStrongTransactionTotal int64
type msgStaleTransactionTotal int64
type msgReschedule func() tea.Cmd

func (db db) tick(f func() tea.Msg) tea.Cmd {
	return tea.Tick(db.refreshRate, func(time.Time) tea.Msg {
		return f
	})
}

func (db db) retrieveUsers() tea.Msg {
	users, err := db.store.GetUsersRichest(20)
	if err != nil {
		return db.tick(db.retrieveUsers)
	}

	richPeople := usersToRiches(users)
	return msgUser(richPeople)
}

func (db db) retrieveTotalUsers() tea.Msg {
	users, err := db.store.GetUsersCount()
	if err != nil {
		return db.tick(db.retrieveTotalUsers)
	}

	return msgTotalUsers(users)
}

func (db db) retrieveTotalMoney() tea.Msg {
	money, err := db.store.GetMoneyTotal()
	if err != nil {
		return db.tick(db.retrieveTotalMoney)
	}

	return msgTotalMoney(money)
}

func (db db) retrieveRichest() tea.Msg {
	money, err := db.store.GetMoneyMax()
	if err != nil {
		return db.tick(db.retrieveRichest)
	}
	return msgRichest(money)
}

func (db db) retrievePoorest() tea.Msg {
	money, err := db.store.GetMoneyMin()
	if err != nil {
		return db.tick(db.retrievePoorest)
	}

	return msgPoorest(money)
}

func (db db) retrieveTransactions() tea.Msg {
	transfers, err := db.store.GetTransfersLatest(20, spanner.StrongRead())
	if err != nil {
		return db.tick(db.retrieveTransactions)
	}

	transactions := []transaction{}

	for _, t := range transfers {
		transactions = append(transactions, transfersToTransaction(t))
	}

	return msgTransactions(transactions)
}

func transfersToTransaction(transfer database.Transfer) transaction {
	return transaction{
		from:      transfer.FromName,
		to:        transfer.ToName,
		amount:    transfer.Amount,
		timestamp: transfer.Timestamp,
	}
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
