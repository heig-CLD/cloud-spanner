package client

import (
	"cloud-spanner/shared"
	"cloud-spanner/shared/database"
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/spanner"

	tea "github.com/charmbracelet/bubbletea"
)

func StartClient() {
	// Setup spanner
	project := shared.LocalConfig()

	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.Uri())

	if err != nil {
		panic(err)
	}

	defer client.Close()

	db := db{
		store: database.NewDatabase(ctx, client),

		// The refresh rate is the same for all queries
		refreshRate: 300 * time.Millisecond,
	}

	// Tea stuff
	p := tea.NewProgram(initialModel(db), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
