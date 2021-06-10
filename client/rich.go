package client

import "fmt"

/*
 *import (
 *  "github.com/charmbracelet/bubbles/progress"
 *)
 */

type Rich struct {
	percentOfAllMoney float32
	name              string
}

func (r Rich) View() string {
	return fmt.Sprintf("%s has %f", r.name, r.percentOfAllMoney)
}
