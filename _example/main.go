package main

import (
	"github.com/dzen-it/fstorage"
	"github.com/dzen-it/fstorage/storage"
)

const _100Kb = 10 << 10 * 100

func main() {
	s, err := storage.NewFileStorage("./data", _100Kb)
	if err != nil {
		panic(err)
	}
	server := fstorage.NewServer(s)
	server.Start(":8080")
}
