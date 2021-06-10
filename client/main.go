package client

import (
	"cloud-spanner/shared"
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/spanner"

	tea "github.com/charmbracelet/bubbletea"
)

func StartClient() {
	// Setup spanner stuff
	project := shared.LocalConfig()

	ctx := context.TODO()
	client, err := spanner.NewClient(ctx, project.Uri())

	if err != nil {
		panic(err)
	}

	defer client.Close()

	db := db{
		ctx:    ctx,
		client: client,
	}

	// Tea stuff
	p := tea.NewProgram(initialModel(db), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
