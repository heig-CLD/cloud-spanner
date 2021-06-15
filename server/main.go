package server

import (
	"bufio"
	"cloud-spanner/shared"
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"os"
	"time"
)

const (
	refresh  = 100 * time.Millisecond
	n        = 20
	maxMoney = int64(10000)
)

func StartServer() {

	project := shared.LocalConfig()
	background := context.Background()

	launchContext, launchCancel := context.WithCancel(background)

	client, err := spanner.NewClient(background, project.Uri())
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
			initDB(background, client, n, maxMoney)
		case "launch":
			launchCancel()
			launchContext, launchCancel = context.WithCancel(background)
			launch(launchContext, client)
		case "start":
			launchCancel()
			launchContext, launchCancel = context.WithCancel(background)
			start(background, launchContext, client, n, maxMoney)
		case "show":
			showDB(background, client)
		case "clear":
			clearDB(background, client)
		case "stop":
			stop(launchCancel)
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
	fmt.Println("DB cleared & repopulated")
}

func launch(ctx context.Context, client *spanner.Client) {
	fmt.Println("Launching simulation")
	go func() {
		for {
			select {
			case <-time.After(refresh):
				if err := TransferRandomly(ctx, client); err != nil {
					fmt.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func start(initCtx context.Context, launchCtx context.Context, client *spanner.Client, n int, maxMoney int64) {
	initDB(initCtx, client, n, maxMoney)
	launch(launchCtx, client)
}

func stop(cancel func()) {
	fmt.Println("Stopping simulation")
	cancel()
}

func showDB(ctx context.Context, client *spanner.Client) {
	fmt.Println("DB Content:")
	showUsers(ctx, client)
	showItems(ctx, client)
	showOffers(ctx, client)
}

func clearDB(ctx context.Context, client *spanner.Client) {
	deleteDBContent(ctx, client)
	fmt.Println("Cleared DB...")
}
