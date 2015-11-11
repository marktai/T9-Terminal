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
	"net/url"
	"time"
)

// "/games/63714/ws"
func Ws(host, path string) {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: host, Path: path}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
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
	resp, err := http.Get(fmt.Sprintf("http://%s/games/%d/string", host, id))
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
		return nil, err
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
	resp, err := http.Get(fmt.Sprintf("http://%s/games/%d", host, id))
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

	g := dummyMap["Game"]

	if err != nil {
		return nil, err
	}

	return g, nil
}

func GetAllGames(host string) ([]uint, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/games", host))
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
		return nil, err
	}

	return dummyMap["Games"], nil
}

func MakeMove(host string, id, player, box, square uint) error {
	resp, err := http.Post(fmt.Sprintf("http://%s/games/%d?Player=%d&Box=%d&Square=%d", host, id, player, box, square), "empty", nil)
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
		return err
	}

	if errText, ok := respMap["Error"]; ok {
		return errors.New(errText)
	}

	return nil
}

func MakeGame(host string, player1, player2 uint) (uint, error) {
	resp, err := http.Post(fmt.Sprintf("http://%s/games?Player1=%d&Player2=%d", host, player1, player2), "empty", nil)
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
		return 0, err
	}

	if respStruct.Error != "" {
		return 0, errors.New(respStruct.Error)
	}

	return respStruct.ID, nil
}
