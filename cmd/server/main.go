package main

import "inventory/internal/network"

func main() {
	app := network.App{}
	app.Initialize(network.DbUser, network.DbPass, network.DbHost, network.DbName)
	app.Run("localhost:8080")
}
