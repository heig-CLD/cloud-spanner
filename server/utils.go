package server

import (
	"bufio"
	_ "bufio"
	"cloud-spanner/shared"
	"cloud.google.com/go/spanner"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	_ "os"
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

	rand.Seed(time.Now().UTC().UnixNano())
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

/// transfer takes a random amount of money from the person with the from identifier, and transfers it to the person
/// with the to identifier.
func transfer(from []byte, to []byte, ctx context.Context, client spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		fetchUser := func(id []byte) (shared.User, error) {
			// TODO : Fetch the right columns.
			row, err := txn.ReadRow(ctx, "Users", spanner.Key{}, []string{})
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

		// TODO : Update the database.
		println(user1.Money + user2.Money)
		return nil
	})
	return err
}

func randomUsers(n int, maxMoney int64) []shared.User {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
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
	iterator := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users ORDER BY Money Desc"))
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
