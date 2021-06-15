package server

import (
	"bufio"
	"cloud-spanner/shared"
	"cloud.google.com/go/spanner"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func deleteDBContent(ctx context.Context, client *spanner.Client) {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {

		deleteUsers := spanner.Delete("Users", spanner.AllKeys())
		deleteItems := spanner.Delete("Items", spanner.AllKeys())
		deleteOffers := spanner.Delete("Offers", spanner.AllKeys())

		mutations := []*spanner.Mutation{deleteUsers, deleteItems, deleteOffers}

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

func createItem(cars []BrandCars, userId []byte) (*spanner.Mutation, error) {
	id, _ := uuid.New().MarshalBinary()

	randIndex := rand.Intn(len(cars))
	brandCar := cars[randIndex]
	randIndex = rand.Intn(len(brandCar.Models))

	brand := brandCar.Brand
	model := brandCar.Models[randIndex]

	return spanner.InsertOrUpdateStruct("Items", shared.Item{Id: id, Description: brand + " " + model, UserId: userId})
}

type BrandCars struct {
	Brand  string   `json:"brand"`
	Models []string `json:"models"`
}

func getCars() ([]BrandCars, error) {
	file, _ := ioutil.ReadFile("shared/brand-cars.json")

	var data []BrandCars
	err := json.Unmarshal(file, &data)

	return data, err
}

func createUsers(ctx context.Context, client *spanner.Client, n int, maxMoney int64) (commitTimestamp time.Time, err error) {
	cars, err := getCars()
	if err != nil {
		panic(err)
	}

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	return client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		users := randomUsers(n, maxMoney)
		var mutations []*spanner.Mutation
		for _, user := range users {
			// Create user
			mut, err := spanner.InsertOrUpdateStruct("Users", user)
			if err != nil {
				return err
			}
			mutations = append(mutations, mut)

			// Give the user a car
			mut, err = createItem(cars, user.Id)
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

/// TaxesAmount indicates what percentage of a user's net worth may be transferred as part of a single transaction. A
/// random amount, inferior to this percentage, will be taken.
const TaxesAmount = 0.4

func Transfer(from []byte, to []byte, ctx context.Context, client spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		return transfer(from, to, ctx, txn)
	})
	return err
}

func TransferRandomly(ctx context.Context, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		it := txn.Query(ctx, spanner.Statement{SQL: "SELECT * FROM Users TABLESAMPLE RESERVOIR (2 ROWS)"})
		users := make([]shared.User, 0)
		err := it.Do(func(r *spanner.Row) error {
			var user shared.User
			err := r.ToStruct(&user)
			if err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
		if err != nil {
			return err
		}
		if len(users) != 2 {
			return errors.New("not enough users")
		}
		return transfer(users[0].Id, users[1].Id, ctx, txn)
	})
	return err
}

/// transfer takes a random amount of money from the person with the from identifier, and transfers it to the person
/// with the to identifier.
func transfer(from []byte, to []byte, ctx context.Context, txn *spanner.ReadWriteTransaction) error {
	fetchUser := func(id []byte) (shared.User, error) {
		row, err := txn.ReadRow(ctx, "Users", spanner.Key{id}, []string{"Id", "Money", "Name"})
		if err != nil {
			return shared.User{}, err
		}
		var user shared.User
		err = row.ToStruct(&user)
		return user, err
	}

	// Fetch both users individually.
	user1, err := fetchUser(from)
	if err != nil {
		return err
	}

	user2, err := fetchUser(to)
	if err != nil {
		return err
	}

	mostMoney, leastMoney := user1, user2
	if user1.Money < user2.Money {
		mostMoney = user2
		leastMoney = user1
	}

	// Calculate how the money should be transferred.
	transferred := int64(float64(mostMoney.Money) * TaxesAmount * rand.Float64())
	mostMoney.Money -= transferred
	leastMoney.Money += transferred

	m1, err := spanner.UpdateStruct("Users", mostMoney)
	if err != nil {
		return err
	}
	m2, err := spanner.UpdateStruct("Users", leastMoney)
	if err != nil {
		return err
	}

	// Create a mutation with the updates.
	return txn.BufferWrite([]*spanner.Mutation{m1, m2})
}

func randomUsers(n int, maxMoney int64) []shared.User {
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

func idAsString(bytes []byte) string {
	id1 := binary.BigEndian.Uint64(bytes[0:8])
	id2 := binary.BigEndian.Uint64(bytes[8:16])

	s1 := strconv.FormatUint(id1, 10)
	s2 := strconv.FormatUint(id2, 10)

	return s1 + s2
}

func showUsers(ctx context.Context, client *spanner.Client) {
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users ORDER BY Money DESC"))
	iterator.Do(func(row *spanner.Row) error {
		var user shared.User
		row.ToStruct(&user)
		println("User - Name: " + user.Name + " Money: " + strconv.FormatInt(user.Money, 10) + " Id: " + idAsString(user.Id))
		return nil
	})
}

func showItems(ctx context.Context, client *spanner.Client) {
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Items"))
	iterator.Do(func(row *spanner.Row) error {
		var item shared.Item
		row.ToStruct(&item)

		println("Item - Description: " + item.Description + " UserId: " + idAsString(item.UserId) + " Id: " + idAsString(item.Id))
		return nil
	})
}

func showOffers(ctx context.Context, client *spanner.Client) {
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Offers"))
	iterator.Do(func(row *spanner.Row) error {
		var offer shared.Offer
		row.ToStruct(&offer)
		println("Offer - Price: " + strconv.FormatInt(offer.Price, 10) + " Id: " + idAsString(offer.Id))
		return nil
	})
}
