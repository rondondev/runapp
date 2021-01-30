package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/rondondev/runapp/api"
	db "github.com/rondondev/runapp/db/sqlc"
	"github.com/rondondev/runapp/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./", "app")
	if err != nil {
		log.Fatal("cannot load configs: ", err)
	}

	fmt.Println(config)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.New(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
