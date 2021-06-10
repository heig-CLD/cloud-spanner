package server

import (
	"bufio"
	_ "bufio"
	"cloud-spanner/shared"
	"cloud.google.com/go/spanner"
	"context"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"os"
	_ "os"
	"time"
)

func DeleteDBContent(ctx context.Context, client *spanner.Client) {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		mut := spanner.Delete("Users", spanner.AllKeys())

		mutations := []*spanner.Mutation{mut}

		err := transaction.BufferWrite(mutations)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func createUsers(ctx context.Context, client *spanner.Client, n int) (commitTimestamp time.Time, err error) {
	return client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		users := RandomUsers(n)
		var mutations []*spanner.Mutation
		for _, u := range users {
			mut, err := spanner.InsertOrUpdateStruct("Users", u)

			if err != nil {
				return err
			}

			mutations = append(mutations, mut)
		}

		err = transaction.BufferWrite(mutations)
		if err != nil {
			return err
		}
		return nil
	})
}

func getAllNames() ([]string, error) {
	file, err := os.Open("shared/names.txt")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var allNames []string
	for scanner.Scan() {
		allNames = append(allNames, scanner.Text())
	}

	return allNames, nil
}

func RandomUsers(n int) []shared.User {
	rand.Seed(20)
	names, err := getAllNames()
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	var people []shared.User
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(len(names))
		name := names[randIndex]
		randMoney := rand.Int63n(10000)

		id, _ := uuid.New().MarshalBinary()

		people = append(people, shared.User{Id: id, Name: name, Money: randMoney})
	}

	return people
}
