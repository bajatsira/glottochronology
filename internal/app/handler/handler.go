package handler

import (
	"LAB1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":   time.Now().Format("15:04:05"),
		"orders": orders,
		"query":  searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
	})
}

func (h *Handler) GetLangs(ctx *gin.Context) {
	var langs []repository.Lang
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем все языки
		langs, err = h.Repository.GetLangs()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		langs, err = h.Repository.GetLangsByName(searchQuery) // ищем языки по названию
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":  time.Now().Format("15:04:05"),
		"langs": langs,
		"query": searchQuery, // передаем введенный запрос обратно на страницу
	})
}

func (h *Handler) GetLang(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr) // преобразуем строку в int
	if err != nil {
		logrus.Error(err)
	}

	lang, err := h.Repository.GetLang(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "lang.html", gin.H{
		"lang": lang,
	})

}

// GetChronos - отображение страницы заявки с заглушкой
func (h *Handler) GetChronos(ctx *gin.Context) {
	// Используем заглушку из репозитория или создаём локальную для примера
	request := h.Repository.GetChronosData() // Предполагаемый метод

	ctx.HTML(http.StatusOK, "chronos.html", gin.H{
		"request": request,
	})
}
