package client

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	Game "github.com/marktai/Tic-Tac-Toe-Squared-Server/src/game"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// "/games/63714/ws"
func Ws(host string, id uint) {
	flag.Parse()
	log.SetFlags(0)

	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/games/%d/ws", host, id), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				addToOutput(fmt.Sprint("read:", err))
				break
			}
			if string(message) == "Changed" {
				refreshBoardGlobals()
				update()
				addToOutput("Updated game")
			}

		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func GetStringArray(host string, id uint) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/games/%d/string", host, id))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	type boardResponse struct {
		Board []string
		Error string
	}

	b := &boardResponse{}

	err = json.Unmarshal(body, b)

	if err != nil {
		return nil, errors.New(string(body))
	}

	if b.Error != "" {
		return nil, errors.New(b.Error)
	}

	if b.Board != nil {
		return b.Board, nil
	}

	return nil, errors.New("Not recognized response")
}

func GetGame(host string, id uint) (*Game.GameInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/games/%d", host, id))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	dummyMap := make(map[string]*Game.GameInfo)

	err = json.Unmarshal(body, &dummyMap)

	if err != nil {
		return nil, errors.New(string(body))
	}

	g := dummyMap["Game"]

	return g, nil
}

func GetAllGames(host string) ([]uint, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/games", host))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	dummyMap := make(map[string][]uint)

	err = json.Unmarshal(body, &dummyMap)

	if err != nil {
		return nil, errors.New(string(body))
	}

	return dummyMap["Games"], nil
}

func MakeMove(host string, id, player, box, square uint) error {
	resp, err := http.Post(fmt.Sprintf("https://%s/games/%d?Player=%d&Box=%d&Square=%d", host, id, player, box, square), "empty", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)

	respMap := make(map[string]string)

	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return errors.New(string(body))
	}

	if errText, ok := respMap["Error"]; ok {
		return errors.New(errText)
	}

	return nil
}

func MakeGame(host string, player1, player2 uint) (uint, error) {
	resp, err := http.Post(fmt.Sprintf("https://%s/games?Player1=%d&Player2=%d", host, player1, player2), "empty", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	type createResponse struct {
		ID    uint
		Error string
	}

	respStruct := createResponse{}

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return 0, errors.New(string(body))
	}

	if respStruct.Error != "" {
		return 0, errors.New(respStruct.Error)
	}

	return respStruct.ID, nil
}
