package shared

import (
	"cloud.google.com/go/spanner"
	"context"
)

type Aggregation int

const (
	Sum Aggregation = iota
	Max
	Min
)

func (agg Aggregation) toSQL() string {
	switch agg {
	case Sum:
		return "SUM"
	case Max:
		return "MAX"
	case Min:
		return "MIN"
	default:
		panic("Unknown aggregation function !!!")
	}
}

func AggregateMoney(agg Aggregation, ctx context.Context, client *spanner.Client) (int64, error) {
	query := "SELECT " + agg.toSQL() + "(Money) FROM Users"
	transaction := client.ReadOnlyTransaction()
	iterator := transaction.Query(ctx, spanner.Statement{SQL: query})
	res, err := iterator.Next()
	if err != nil {
		return 0, err
	}
	var sum int64
	err = res.Column(0, &sum)
	return sum, err
}
