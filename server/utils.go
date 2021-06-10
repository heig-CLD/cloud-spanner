package server

import (
	_ "bufio"
	"cloud.google.com/go/spanner"
	"context"
	_ "os"
)

func DeleteDBContent(ctx context.Context, client *spanner.Client) {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		mut := spanner.Delete("Users", spanner.AllKeys())

		mutations := []*spanner.Mutation{mut}

		err := transaction.BufferWrite(mutations)
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		panic(err)
	}
}