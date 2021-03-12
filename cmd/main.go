package main

import (
	"github.com/kate-network/backend/cache"
	"github.com/kate-network/backend/internal"
	"github.com/kate-network/backend/storage"
)

func main() {
	stor, err := storage.Open("root:@tcp(localhost:3306)/kate?parseTime=true")
	if err != nil {
		panic(err)
	}
	ch, err := cache.New("localhost:6379")
	if err != nil {
		panic(err)
	}

	service := internal.NewServer(stor, ch)
	service.Init()

	if err := service.Listen("127.0.0.1:7454"); err != nil {
		panic(err)
	}
}
