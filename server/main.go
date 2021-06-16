package server

import (
	"bufio"
	"cloud-spanner/shared"
	"cloud-spanner/shared/database"
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

	server := create(client)

	for {
		fmt.Print("$ ")
		scanner.Scan()
		switch scanner.Text() {
		case "init":
			server.init()
		case "launch":
			server.launch()
		case "start":
			server.start()
		case "show":
			server.show()
		case "clear":
			server.clear()
		case "stop":
			server.stop()
		default:
			fmt.Println("Unrecognized command... " + helpText)
		}
	}
}

type server struct {
	background context.Context

	context    context.Context
	cancelFunc context.CancelFunc

	client *spanner.Client
}

// create returns a new server
func create(spanner *spanner.Client) server {
	bg := context.Background()
	ctx, can := context.WithCancel(bg)
	return server{
		background: bg,
		context:    ctx,
		cancelFunc: can,
		client:     spanner,
	}
}

// withDatabase executes an operation with the current database
func (s *server) withDatabase(op func(db database.Database)) {
	op(database.NewDatabase(s.context, s.client))
}

func (s *server) populate() {
	s.withDatabase(func(db database.Database) {
		users := randomUsers(n, maxMoney)
		_ = db.AddUsers(users)
		fmt.Printf("Added %d new users to the database.\n", n)
	})
}

func (s *server) init() {
	s.withDatabase(func(db database.Database) {
		_ = db.Clear()
		users := randomUsers(n, maxMoney)
		_ = db.AddUsers(users)
	})
}

func (s *server) launch() {
	s.cancel()
	fmt.Println("Launching simulation")
	ctx := s.context
	go func() {
		for {
			select {
			case <-time.After(refresh):
				_ = TransferRandomly(ctx, s.client)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *server) cancel() {
	s.cancelFunc()
	ctx, can := context.WithCancel(s.background)
	s.context = ctx
	s.cancelFunc = can
}

func (s *server) start() {
	s.init()
	s.launch()
}

func (s *server) stop() {
	fmt.Println("Stopping simulation")
	s.cancel()
}

func (s *server) clear() {
	s.withDatabase(func(db database.Database) {
		_ = db.Clear()
		fmt.Println("Cleared DB...")
	})
}

func (s *server) show() {
	fmt.Println("DB Content:")
	showUsers(database.NewDatabase(s.context, s.client))
	showItems(s.context, s.client)
	showOffers(s.context, s.client)
}
