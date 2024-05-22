package main

import "inventory/internal/network"

func main() {
	app := network.App{}
	app.Initialize()
	app.Run("localhost:8080")
}
