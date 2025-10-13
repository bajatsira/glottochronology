package handler

import (
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	//router.GET("/", h.GetOrders)
	//router.GET("/order/:id", h.GetOrder)

	router.GET("/languages", h.GetLangs)
	//r.GET("/order/:id", handler.GetOrder) // вот наш новый обработчик
	router.GET("/lang/:id", h.GetLang)
	//router.GET("/chronos", h.GetChronos)

}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	//router.Static("/styles", "./styles")
	//r.Static("/static", "./resources")

	router.Static("/styles", "resources/styles")
	router.Static("/img", "resources/img")

}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
