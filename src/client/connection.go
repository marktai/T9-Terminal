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

// "/game/63714/ws"
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

func GetStringArray(host, id string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/game/%s/string", host, id))
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

func GetGame(host, id string) (*Game.Game, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/game/%s", host, id))
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


	
	dummyMap :=make(map[string]*Game.Game)

	err = json.Unmarshal(body, &dummyMap)

	g := dummyMap["Game"]

	if err != nil {
		return nil, err
	}

	return g, nil
}
