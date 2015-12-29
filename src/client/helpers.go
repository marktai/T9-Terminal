package client

import (
	"fmt"
	Game "github.com/marktai/Tic-Tac-Toe-Squared-Server/src/game"
	"strconv"
	"strings"
)

func stringtoUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	return uint(i), err
}

func parseHeader(game *Game.GameInfo, player uint) string {
	turnString := ""
	if game.Players[game.Turn/10] == player {
		turnString = "Your Turn"
	} else {
		turnString = "Other Player's Turn"
	}
	box, ok := boxToString[game.Turn%10]
	if !ok {
		box = fmt.Sprintf("%d", game.Turn%10)
		parMap["output"].Text = fmt.Sprint(boxToString)
	}
	out := fmt.Sprintf("ID: %d | Player: %d | %s | Box: %s", game.GameID, player, turnString, box)
	return out
}

func parseMoveHistory(game *Game.GameInfo) [18]string {

	var moves [18]string
	for i, move := range game.MoveHistory {
		if move == 127 {
			break
		}
		moves[i] = fmt.Sprintf("%d Ago) B:%s, S:%s", i+1, boxToString[move/9], boxToString[move%9])
	}

	return moves
}

func getGameAndString(host string, gameid, playerid uint) (*Game.GameInfo, string, error) {
	lines, err := GetStringArray(host, gameid)
	if err != nil {
		return nil, "", err
	}
	boardText := ""
	start := true
	for _, line := range lines {
		if !start {
			boardText += "\n"
		} else {
			start = false
		}
		line = strings.Replace(line, "x", "[x](fg-red)", -1)
		line = strings.Replace(line, "o", "[o](fg-green)", -1)
		boardText += line
	}

	game, err := GetGame(host, gameid)
	if err != nil {
		return nil, "", err
	}

	return game, boardText, nil
}
