// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	// "flag"
	"fmt"
	"log"
	// "net/url"
	"reflect"
	"strings"
	// "time"
	"strconv"
	"github.com/gizak/termui"
	Game "github.com/marktai/Tic-Tac-Toe-Squared-Server/src/game"
)

var parMap = make(map[string]*termui.Par)

var boxToStringTranslator = make(map[uint]string)
var stringToBoxTranslator = make(map[string]uint)

var host string
var gameid string
var playerid string

var state = 0

func initTranslators() {

	for i := uint(0); i < 9; i++ {
		out := ""
		height := i / 3
		if height == 0 {
			out += "top-"
		} else if height == 1 {
			out += "middle-"
		} else if height == 2 {
			out += "bottom-"
		}

		width := i % 3
		if width == 0 {
			out += "left"
		} else if width == 1 {
			out += "middle"
		} else if width == 2 {
			out += "right"
		}

		boxToStringTranslator[i] = out
		stringToBoxTranslator[out] = i

	}

	boxToStringTranslator[9] = "anywhere"
}

func parseHeader(game *Game.Game, player uint) string {
	turnString := ""
	if game.Players[game.Turn/10] == player {
		turnString = "Your Turn"
	} else {
		turnString = "Other Player's Turn"
	}
	box, ok := boxToStringTranslator[game.Turn%10]
	if !ok {
			box = fmt.Sprintf("%d", game.Turn%10)
			parMap["output"].Text = fmt.Sprint(boxToStringTranslator)
	}
	out := fmt.Sprintf("ID: %d | Player: %d | %s | Box: %s", game.GameID, player, turnString, box)
	return out

}

func setupBody() {

	prompt := termui.NewPar("")
	prompt.Height = 1
	prompt.Border = false
	parMap["prompt"] = prompt

	input := termui.NewPar("")
	input.Height = 5
	input.BorderLabel = "Input"
	input.BorderFg = termui.ColorYellow
	parMap["input"] = input

	output := termui.NewPar("")
	output.Height = 6
	output.BorderLabel = "Output"
	output.BorderFg = termui.ColorGreen
	output.Display = false
	parMap["output"] = output

	board := termui.NewPar("")
	board.Height = 23
	board.Width = 37
	board.BorderLabel = "Board"
	board.BorderFg = termui.ColorRed
	board.Display = false
	parMap["board"] = board

	// build layout
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, parMap["prompt"], parMap["input"]),
			termui.NewCol(6, 0, parMap["output"]),
		),
		termui.NewRow(
			termui.NewCol(12, 0, board),
		),
	)
	changeState(0)
}

func refreshBoard(host, gameid, playerid string) {

	lines, err := GetStringArray(host, gameid)
	if err != nil {
		log.Println(err)
	} else {
		boardText := ""
		for _, line := range lines {
			boardText += line + "\n"
		}

		parMap["board"].Text = boardText
	}

	game, err := GetGame(host, gameid)
	if err != nil {
		parMap["output"].Text = fmt.Sprint(err)
	} else {
		//parMap["output"].Text = fmt.Sprint(game)

		playerInt, err := strconv.Atoi(playerid)
		if err != nil {
				// this not good
		}
		parMap["board"].BorderLabel = "Board | " + parseHeader(game, uint(playerInt))
}
}

func changeState(inp int) {
	switch inp {
	case 0:
		parMap["prompt"].Text = "Host?"
	case 1:
		parMap["prompt"].Text = "Game ID?"
	case 2:
		parMap["prompt"].Text = "Player ID?"
	case 3:
		parMap["prompt"].Text = "Command (r, h, i, m, p, s, q)?"
	}
	state = inp
}

func parseInput(inp string) {
	inp = strings.ToLower(inp)

	switch state {
	case 0:
		host = inp
		if host == "" {
			host = "localhost:8080"
		}
		changeState(1)
	case 1:
		gameid = inp
		if gameid == "" {
			gameid = "63714"
		}
		changeState(2)
	case 2:
		playerid = inp
		if playerid == "" {
			playerid = "0"
		}
		changeState(3)
		refreshBoard(host, gameid, playerid)
	case 3:
		switch inp {
		case "r", "refresh":
			refreshBoard(host, gameid, playerid)
	case "m", "move":
//			makeMove(host, gameid, playerid)
	case "q", "quit", ":q":
			termui.StopLoop()
		}
	}

	//parMap["output"].Text = inp
}

func nothing(inp string) {

}
func setupHandlers() {

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd", func(ev termui.Event) {
		if kbdEvent, ok := ev.Data.(termui.EvtKbd); ok {

			keyStr := kbdEvent.KeyStr
			switch {
			case keyStr == "<enter>":
				inp := parMap["input"].Text
				parseInput(inp)
				parMap["input"].Text = ""
			case keyStr == "<space>":
				parMap["input"].Text += " "
			case keyStr == "C-8":
				// for some reason, backspace is C-8
				if len(parMap["input"].Text) > 0 {
					parMap["input"].Text = parMap["input"].Text[:len(parMap["input"].Text)-1]
				}
			case strings.Contains(keyStr, "<"):

			default:
				parMap["input"].Text += keyStr
			}

			termui.Body.Align()
			termui.Render(termui.Body)
		} else {
			dataType := reflect.TypeOf(ev.Data)
			parMap["input"].Text += fmt.Sprint(ev.Data)
			parMap["output"].Text += fmt.Sprint(dataType)
			termui.Body.Align()
			termui.Render(termui.Body)
		}

	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Render(termui.Body)
	})
}

func UI() {
//		test()
	initTranslators()
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	setupBody()

	// calculate layout
	termui.Body.Align()
	termui.Render(termui.Body)

	setupHandlers()
	termui.Loop()
}

func test() {

		game, err := GetGame("localhost:8080", "63714")
	if err != nil {
		log.Println(err)
	} else {
			log.Println(game)
	}
}
