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
	"strconv"
	"time"
)

func deleteDBContent(ctx context.Context, client *spanner.Client) {
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

func createUsers(ctx context.Context, client *spanner.Client, n int, maxMoney int64) (commitTimestamp time.Time, err error) {
	return client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		users := randomUsers(n, maxMoney)
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

func randomUsers(n int, maxMoney int64) []shared.User {
	rand.Seed(20)
	names, err := getAllNames()
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	var people []shared.User
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(len(names))
		name := names[randIndex]
		randMoney := rand.Int63n(maxMoney)

		id, _ := uuid.New().MarshalBinary()

		people = append(people, shared.User{Id: id, Name: name, Money: randMoney})
	}

	return people
}

func showUsers(ctx context.Context, client *spanner.Client) {
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users ORDER BY Money Desc"))
	iterator.Do(func(row *spanner.Row) error {
		var user shared.User
		row.ToStruct(&user)
		println("User - Name: " + user.Name + " Money: " + strconv.FormatInt(user.Money, 10))
		return nil
	})
}

func showItems(ctx context.Context, client *spanner.Client) {
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Items"))
	iterator.Do(func(row *spanner.Row) error {
		var item shared.Item
		row.ToStruct(&item)
		println("Item - Description: " + item.Description)
		return nil
	})
}

func showOffers(ctx context.Context, client *spanner.Client) {
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Offers"))
	iterator.Do(func(row *spanner.Row) error {
		var offer shared.Offer
		row.ToStruct(&offer)
		println("Offer - Price: " + strconv.FormatInt(offer.Price, 10))
		return nil
	})
}
