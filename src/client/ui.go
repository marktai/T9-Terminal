// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	// "flag"
	"fmt"
	"github.com/gizak/termui"
	Game "github.com/marktai/Tic-Tac-Toe-Squared-Server/src/game"
	"reflect"
	"strings"
)

var parMap = make(map[string]*termui.Par)
var linesMap = make(map[string]*lines)

var boxToString = make(map[uint]string)
var stringToBox = make(map[string]uint)

var host string
var gameid uint
var playerid uint

var players [2]uint

var globalGame *Game.GameInfo
var box uint
var square uint

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

		boxToString[i] = out

		// stringToBox has both with and without dashes in the string
		stringToBox[out] = i
		stringToBox[strings.Replace(out, "-", "", -1)] = i

	}

	boxToString[9] = "anywhere"
}

func setupBody() {

	height := termui.TermHeight() - 23

	prompt := termui.NewPar("")
	prompt.Height = 1
	prompt.Border = false
	parMap["prompt"] = prompt

	input := termui.NewPar("")
	input.Height = 3
	input.BorderLabel = "Input"
	input.BorderFg = termui.ColorYellow
	parMap["input"] = input

	moveHistory := termui.NewPar("")
	moveHistory.Height = height - 4
	moveHistory.BorderLabel = "Move History"
	moveHistory.BorderFg = termui.ColorBlue
	parMap["moveHistory"] = moveHistory
	linesMap["moveHistory"] = &lines{[]string{}, 0}

	output := termui.NewPar("")
	output.Height = height
	output.BorderLabel = "Output"
	output.BorderFg = termui.ColorGreen
	parMap["output"] = output
	linesMap["output"] = &lines{[]string{}, 0}

	board := termui.NewPar("")
	board.Height = 23
	board.Width = 37
	board.BorderLabel = "Board"
	board.BorderFg = termui.ColorRed
	parMap["board"] = board

	// build layout
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, parMap["prompt"], parMap["input"], parMap["moveHistory"]),
			termui.NewCol(6, 0, parMap["output"]),
		),
		termui.NewRow(
			termui.NewCol(12, 0, parMap["board"]),
		),
	)
	changeState(0)
}

func addToOutput(s string) {
	linesMap["output"].Add(s)
	for linesMap["output"].Length() > parMap["output"].Height-2 {
		linesMap["output"].Down()
	}
}

func setOutput(s string) {
	linesMap["output"].Set(s)
}

func clearOutput() {
	linesMap["output"].Clear()
}

func adjustDimensions() {
	termui.Body.Width = termui.TermWidth()
	height := termui.TermHeight() - 23
	parMap["moveHistory"].Height = height - 4
	parMap["output"].Height = height
}

func update() {
	parMap["moveHistory"].Text = linesMap["moveHistory"].String()
	parMap["output"].Text = linesMap["output"].String()
	termui.Body.Align()
	termui.Render(termui.Body)
}

func refreshMoveHistory(game *Game.GameInfo) {
	moves := parseMoveHistory(game)
	linesMap["moveHistory"].Clear()
	for _, move := range moves {
		linesMap["moveHistory"].Add(move)
	}
}

func refreshBoard(host string, gameid, playerid uint) error {
	game, board, err := getGameAndString(host, gameid, playerid)
	if err != nil {
		addToOutput(err.Error())
		if err.Error() == "Game not found" {
			changeState(1)
		} else if strings.Contains(err.Error(), "no such host") {
			changeState(0)
		}
		return err
	}

	parMap["board"].BorderLabel = parseHeader(game, playerid)
	parMap["board"].Text = board

	refreshMoveHistory(game)

	globalGame = game
	return nil
}

func displayInfo(game *Game.GameInfo) {

	if game == nil {
		addToOutput("Game is nil")
		return
	}

	addToOutput(fmt.Sprintf("Game ID: %d", game.GameID))
	addToOutput(fmt.Sprintf("Players: %d, %d", game.Players[0], game.Players[1]))
	addToOutput(fmt.Sprintf("Started: %s", game.Started.String()))
	addToOutput(fmt.Sprintf("Modified: %s", game.Modified.String()))
}

func changeState(inp int) {
	switch inp {
	case 0:
		parMap["prompt"].Text = "Host?"
	case 6:
		parMap["prompt"].Text = "Create game or join one (c, j)?"
	case 1:
		parMap["prompt"].Text = "Game ID?"
	case 2:
		parMap["prompt"].Text = "Player ID?"
	case 3:
		parMap["prompt"].Text = "Command (r, m, i, c, p, s, q)?"
	case 4:
		parMap["prompt"].Text = "Box?"
	case 5:
		parMap["prompt"].Text = "Square?"
	case 7:
		parMap["prompt"].Text = "Player 1 ID?"
	case 8:
		parMap["prompt"].Text = "Player 2 ID?"

	}
	state = inp

}

func parseInput(inp string) {
	inp = strings.ToLower(inp)

	switch state {
	case 0: // getting host
		if inp == "" {
			inp = "localhost:8080"
		}
		host = inp
		changeState(6)

	case 6:
		if inp == "" {
			inp = "j"
		}

		switch inp {
		case "j", "join":
			changeState(1)

			games, err := GetAllGames(host)
			if err != nil {
				addToOutput(err.Error())
				if strings.Contains(err.Error(), "no such host") {
					changeState(0)
				}
			} else {
				addToOutput("Avaliable Games:")
				for _, game := range games {
					addToOutput(fmt.Sprint(game))
				}
			}

		case "c", "create":
			changeState(7)
		}

	case 1: // getting game id
		if inp == "" {
			inp = "63714"
		}

		var err error
		if gameid, err = stringtoUint(inp); err != nil {
			addToOutput("Bad Game ID")
		} else {
			changeState(2)
		}

	case 2: // getting player id
		if inp == "" {
			inp = "0"
		}

		var err error
		if playerid, err = stringtoUint(inp); err != nil {
			addToOutput("Bad Player ID")
		} else {
			changeState(3)
			refreshBoard(host, gameid, playerid)
			displayInfo(globalGame)
		}

	case 3: // getting generic command
		switch inp {
		case "r", "refresh":
			refreshBoard(host, gameid, playerid)
			addToOutput("Refreshed")

		case "m", "move":
			if globalGame.Players[globalGame.Turn/10] == playerid {
				if tempBox := globalGame.Turn % 10; tempBox == 9 {
					changeState(4)
				} else {
					box = tempBox
					changeState(5)
				}
			} else {
				addToOutput("Not your turn")
			}

		case "i", "info": // display info
			displayInfo(globalGame)

		case "p", "player": // switch players
			changeState(2)
		case "s", "switch": // switch game
			changeState(1)

		case "c", "clear": // switch game
			clearOutput()

		case "q", "quit", ":q":
			termui.StopLoop()
		}

	case 4, 5: // getting box or square

		goodInput := false
		var err error
		var tempNum uint

		if inp == "b" || inp == "back" {
			changeState(3)
		} else {
			// if number input like 1
			if tempNum, err = stringtoUint(inp); err == nil {
				goodInput = true
			} else {
				var ok bool
				// if word input like top middle
				if tempNum, ok = stringToBox[strings.Replace(inp, " ", "", -1)]; ok {
					goodInput = true
				} else {
					addToOutput("Bad Position")
				}
			}
		}

		if goodInput {
			if state == 4 {
				box = tempNum
				changeState(5)
			} else if state == 5 {
				square = tempNum

				err := MakeMove(host, gameid, playerid, box, square)
				if err != nil {
					addToOutput(err.Error())
				}
				changeState(3)
				refreshBoard(host, gameid, playerid)
			}
		}
	case 7:
		if inp == "" {
			inp = "0"
		}

		var err error
		if players[0], err = stringtoUint(inp); err != nil {
			addToOutput("Bad Player ID")
		} else {
			changeState(8)
		}
	case 8:
		if inp == "" {
			inp = "1"
		}

		var err error
		if players[1], err = stringtoUint(inp); err != nil {
			addToOutput("Bad Player ID")
		} else {
			id, err := MakeGame(host, players[0], players[1])
			if err != nil {
				changeState(6)
				addToOutput(err.Error())
			} else {
				changeState(3)
				refreshBoard(host, id, players[0])
				displayInfo(globalGame)
			}
		}

	}
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

			case keyStr == "<up>":
				for _, linesItem := range linesMap {
					linesItem.Up()
				}
			case keyStr == "<down>":
				for _, linesItem := range linesMap {
					linesItem.Down()
				}

			case strings.Contains(keyStr, "<"):

			default:
				parMap["input"].Text += keyStr
			}

			update()
		} else {
			dataType := reflect.TypeOf(ev.Data)
			addToOutput("event type of " + fmt.Sprint(dataType))
			update()
		}

	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		adjustDimensions()
		update()
	})
}

func UI() {
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
