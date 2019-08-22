package main

import (
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/micro/go-micro"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
)

func main() {
	srv := micro.NewService(
		micro.Name("shippy.user.service"),
	)

	srv.Init()

	postgresCfg := DefaultPostgresConfig()
	dbhost := os.Getenv("DB_HOST")
	if dbhost != "" {
		postgresCfg.Host = dbhost
	}

	dbport := os.Getenv("DB_PORT")
	if dbport != "" {
		var err error
		postgresCfg.Port, err = strconv.Atoi(dbport)
		if err != nil {
			log.Panic(err)
		}
	}

	db, err := gorm.Open(postgresCfg.Dialect(), postgresCfg.ConnectionInfo())
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&proto.User{})
	repo := userRepository{db}
	service := userService{
		repo:         &repo,
		tokenService: &TokenService{&repo},
	}
	proto.RegisterUserServiceHandler(srv.Server(), &service)

	if err = srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
