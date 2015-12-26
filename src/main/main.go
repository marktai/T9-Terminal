package main

import (
	"client"
	"flag"
)

func main() {

	var refresh uint

	flag.UintVar(&refresh, "refresh", 100, "Refresh rate in milliseconds")

	flag.Parse()

	client.UI(refresh)
}
