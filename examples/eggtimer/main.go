// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gioui.org/app"
	"github.com/mheremans/goui"
	"github.com/mheremans/goui/examples/eggtimer/views"
	"github.com/mheremans/goui/types"
)

var wg sync.WaitGroup
var window *goui.Window

func main() {
	window = goui.NewWindow("Egg Timer")
	window.SetSize(types.NewWindowSize(400, 600))
	window.OnClose = func() {
		fmt.Println("Window Closed")
	}
	closeChan, err := window.Show(views.NewTimerView())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-closeChan
		os.Exit(0)
	}()

	app.Main()

}
