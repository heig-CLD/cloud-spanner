package client

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
)

type Rich struct {
	percentOfAllMoney float64
	name              string
	progress          *progress.Model
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
