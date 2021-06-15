package shared

import (
	"cloud.google.com/go/spanner"
	"context"
)

const (
	Sum = "SUM"
	Max = "MAX"
	Min = "MIN"
)

// TODO : Provide an way to choose the aggregate (max, min, sum)
func AggregateMoney(ctx context.Context, client *spanner.Client) (int64, error) {
	transaction := client.ReadOnlyTransaction()
	iterator := transaction.Query(ctx, spanner.Statement{SQL: "SELECT SUM(Money) FROM Users"})
	res, err := iterator.Next()
	if err != nil {
		return 0, err
	}
	var sum int64
	err = res.Column(0, &sum)
	return sum, err
}
