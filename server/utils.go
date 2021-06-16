package server

import (
	"bufio"
	"cloud-spanner/shared"
	"cloud-spanner/shared/database"
	"cloud.google.com/go/spanner"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

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
		return transfer(users[0], users[1], txn)
	})
	return err
}

/// transfer takes a random amount of money from the person with the from identifier, and transfers it to the person
/// with the to identifier.
func transfer(from shared.User, to shared.User, txn *spanner.ReadWriteTransaction) error {

	// Figure out the transfer direction.
	mostMoney, leastMoney := from, to
	if from.Money < to.Money {
		mostMoney = to
		leastMoney = from
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

func showUsers(store database.Database) {
	users, err := store.GetUsersRichest(20)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s when reading users.\n", err.Error())
		return
	}
	for _, user := range users {
		println("User - Name: " + user.Name + " Money: " + strconv.FormatInt(user.Money, 10) + " Id: " + idAsString(user.Id))
	}
}
