package api

import (
	"LAB1/internal/app/handler"
	"LAB1/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/languages", handler.GetLangs)
	//r.GET("/order/:id", handler.GetOrder) // вот наш новый обработчик
	r.GET("/lang/:id", handler.GetLang)
	r.GET("/chronos", handler.GetChronos)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
