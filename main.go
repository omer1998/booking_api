package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	apiserver "github.com/omer1998/booking_api/api_server"
	"github.com/omer1998/booking_api/config"
	"github.com/omer1998/booking_api/database"
)

func main() {
	cxt := context.Background()
	if err := run(cxt); err != nil {
		panic(err)
	}
}

func run(cxt context.Context) error {

	myConfig, err := config.New()
	if err != nil {
		return err
	}
	config.PrintConfig(myConfig)
	fmt.Println(myConfig.GetConnectionDbUrl())
	context, cancel := signal.NotifyContext(cxt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	connPool := database.ConnectPool(context, myConfig)
	db := database.NewPostgresDbPool(connPool, context)
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	fmt.Println(db)
	apiServer := apiserver.NewApiServer(":8000", db, context)
	if err := apiServer.Run(); err != nil {
		return err
	}
	// fmt.Println("Server is running on port 8000")
	return nil
}
