package client

import "time"

type transaction struct {
	from      string
	to        string
	amount    int64
	timestamp time.Time
}

type Transactions struct {
	strong       []transaction
	strongAmount int64
	staleAmount  int64
}

func (transactions Transactions) View() string {
	return ""
}
