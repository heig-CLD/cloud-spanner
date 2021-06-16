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

func (db db) tick(f func() tea.Msg) tea.Cmd {
	return tea.Tick(db.refreshRate, func(time.Time) tea.Msg {
		return f()
	})
}

func (db db) retrieveUsers() tea.Msg {
	users, _ := db.store.GetUsersRichest(20)
	richPeople := usersToRiches(users)
	return msgUser(richPeople)
}

func (db db) retrieveTotalUsers() tea.Msg {
	users, _ := db.store.GetUsersCount()
	return msgTotalUsers(users)
}

func (db db) retrieveTotalMoney() tea.Msg {
	money, _ := db.store.GetMoneyTotal()
	return msgTotalMoney(money)
}

func (db db) retrieveRichest() tea.Msg {
	money, _ := db.store.GetMoneyMax()
	return msgRichest(money)
}

func (db db) retrievePoorest() tea.Msg {
	money, _ := db.store.GetMoneyMin()
	return msgPoorest(money)
}

func (db db) retrieveStrongTransactionsCount() tea.Msg {
	amount, _ := db.store.GetTransfersCount(spanner.StrongRead())
	return msgStrongTransactionTotal(amount)
}

func (db db) retrieveStaleTransactionsCount() tea.Msg {
	amount, _ := db.store.GetTransfersCount(spanner.ExactStaleness(15 * time.Second))
	return msgStaleTransactionTotal(amount)
}

func (db db) retrieveTransactions() tea.Msg {
	transfers, _ := db.store.GetTransfersLatest(17, spanner.StrongRead())

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
