package main

import "github.com/NaturalSelectionLabs/RSS3-PreGod/task/service"

func main() {
	go service.SubscribeEns()
}
