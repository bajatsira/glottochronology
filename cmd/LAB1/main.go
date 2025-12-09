package main

import (
	"fmt"

	"LAB1/internal/app/auth"
	"LAB1/internal/app/config"
	"LAB1/internal/app/dsn"
	"LAB1/internal/app/handler"
	"LAB1/internal/app/repository"
	"LAB1/internal/pkg"

	"github.com/joho/godotenv"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		logrus.Warnf("error loading .env: %v", err)
	}

	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	/*if err := auth.InitRedis(os.Getenv("REDIS_ADDR")); err != nil {
		logrus.Warnf("redis not initialized: %v", err)
	}*/

	if err := auth.InitRedis(os.Getenv("localhost:6379")); err != nil {
		logrus.Fatalf("redis not initialized: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()

}
