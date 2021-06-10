package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/charmbracelet/bubbles/progress"
)

type Rich struct {
	percentOfAllMoney float64
	name              string
	progress          *progress.Model
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

func RandomRichPeople(n int) []Rich {
	rand.Seed(20)
	names, err := getAllNames()
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	people := []Rich{}
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(len(names))
		name := names[randIndex]
		randPerc := rand.Float64()

		people = append(people, InitializeRich(name, randPerc))
	}

	return people
}

func InitializeRich(name string, percent float64) Rich {
	rich := Rich{
		percentOfAllMoney: percent,
		name:              name,
	}

	prog, err := progress.NewModel(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	if err != nil {
		fmt.Println("Could not initialize progress model:", err)
		os.Exit(1)
	}

	rich.progress = prog

	return rich
}

func (r Rich) View() string {
	str := fmt.Sprintf("%s has %s", r.name, r.progress.View(r.percentOfAllMoney))
	return str
}
