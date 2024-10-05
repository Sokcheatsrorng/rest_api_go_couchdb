package main

import (
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/delivery"
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/repository/couchdb"
	routes "github.com/Sokcheatsrorng/go-clean-architecture/internal/router"
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/service"
)
	

func main(){
	//init and crate database
	dbName := "users"
	db:= couchdb.InitDB(dbName);

	userRepo := couchdb.NewUserRepository(db)
	userService := service.NewUserService(*userRepo)
	userHandler := delivery.NewUserHandler(userService)

	r := routes.InitRoutes(userHandler)
	r.Run(":9090")
	
}