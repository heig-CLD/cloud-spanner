package server

import (
	"bufio"
	"cloud-spanner/shared"
	"cloud-spanner/shared/database"
	"encoding/binary"
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
