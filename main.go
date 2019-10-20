package main

import "shitty.moe/satelit-project/satelit-scraper/server"

func main() {
	// TODO: move to config
	err := server.Serve(10700)
	if err != nil {
		panic(err)
	}
}
