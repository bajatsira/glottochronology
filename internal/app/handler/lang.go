package handler

import (
	"net/http"
	"strconv"
	"time"

	"LAB1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetLangs(ctx *gin.Context) {
	var langs []ds.Lang
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
