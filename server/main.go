package server

import (
	"bufio"
	"cloud-spanner/shared"
	spanner "cloud.google.com/go/spanner"
	"context"
	"fmt"
	_ "fmt"
	"os"
	"time"
)

func StartServer() {

	n := 20
	maxMoney := int64(10000)

	project := shared.LocalConfig()
	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.Uri())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	scanner := bufio.NewScanner(os.Stdin)
	helpText := "Available commands are:\n" +
		"- \"init\": Cleans DB and populates it\n" +
		"- \"go\": Launches simulation\n" +
		"- \"start\": Cleans DB, populates it and starts simulation. Equivalent to <init, launch>\n" +
		"- \"show\": Shows DB content\n" +
		"- \"clear\": Clears DB content\n" +
		"- \"stop\": Stops simulation"

	fmt.Println(helpText)
	for {
		fmt.Print("$ ")
		scanner.Scan()
		switch scanner.Text() {
		case "init":
			initDB(ctx, client, n, maxMoney)
		case "launch":
			launch(ctx, client)
		case "test":
			test(ctx, client)
		case "start":
			start(ctx, client, n, maxMoney)
		case "show":
			showDB(ctx, client)
		case "clear":
			clearDB(ctx, client)
		case "stop":
			stop(ctx, client)
		default:
			fmt.Println("Unrecognized command... " + helpText)
		}
	}
}

func initDB(ctx context.Context, client *spanner.Client, n int, maxMoney int64) {
	deleteDBContent(ctx, client)
	_, err := createUsers(ctx, client, n, maxMoney)
	if err != nil {
		panic(err)
	}
}

func launch(ctx context.Context, client *spanner.Client) {
	go func() {
		for {
			<-time.After(100 * time.Millisecond)
			if err := TransferRandomly(ctx, client); err != nil {
				fmt.Println(err)
			}
			fmt.Println(".")
		}
	}()
}

func test(ctx context.Context, client *spanner.Client) {
	println(shared.AggregateMoney(ctx, client))
}

func showDB(ctx context.Context, client *spanner.Client) {
	fmt.Println("DB Content:")
	showUsers(ctx, client)
	showItems(ctx, client)
	showOffers(ctx, client)
}

func clearDB(ctx context.Context, client *spanner.Client) {
	fmt.Println("Clearing DB...")
	deleteDBContent(ctx, client)
	fmt.Println("Cleared DB...")
}

func start(ctx context.Context, client *spanner.Client, n int, maxMoney int64) {
	initDB(ctx, client, n, maxMoney)
	launch(ctx, client)
}

func stop(ctx context.Context, client *spanner.Client) {

}
