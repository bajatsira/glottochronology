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

func (h *Handler) GetLanguageById(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	language, err := h.Repository.GetLanguageByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "language.page.tmpl", language)
}

// AddLanguageToDraft — добавляет язык в текущую заявку исследователя (черновик)
func (h *Handler) AddLanguageToDraft(ctx *gin.Context) {
	langIDStr := ctx.Param("id")
	langID, err := strconv.Atoi(langIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID языка"})
		logrus.Error("Некорректный ID:", err)
		return
	}

	researcherID := uint(1) // пока хардкодим

	err = h.Repository.AddServiceToDraft(researcherID, uint(langID))
	if err != nil {
		logrus.Error("Ошибка при добавлении языка в заявку:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить язык в заявку"})
		return
	}

	ctx.Redirect(http.StatusFound, "/languages")
}

// GetDraft — показывает содержимое текущей заявки (черновика)
func (h *Handler) GetDraft(ctx *gin.Context) {
	researcherID := uint(1) // временно захардкодим

	glotto, err := h.Repository.GetDraftByResearcher(researcherID)
	if err != nil {
		logrus.Error("Ошибка при получении черновика:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить черновик заявки"})
		return
	}

	ctx.HTML(http.StatusOK, "chronos.html", gin.H{
		"glotto":     glotto,
		"cart_count": h.Repository.GetLangCount(),
	})
}

/*
// DeleteGlotto — логически удаляет заявку (меняет статус на "удалён")
func (h *Handler) DeleteGlotto(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		logrus.Error("Некорректный ID:", err)
		return
	}

	err = h.Repository.DeleteDraftSQL(uint(id))
	if err != nil {
		logrus.Error("Ошибка при удалении заявки:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить заявку"})
		return
	}

	ctx.Redirect(http.StatusFound, "/langs")
}
*/

// DeleteGlotto — логически удаляет заявку (меняет статус на "удалён")
func (h *Handler) DeleteGlotto(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Некорректный ID заявки:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заявки"})
		return
	}

	// Попытка удалить заявку
	if err := h.Repository.DeleteDraftSQL(uint(id)); err != nil {
		logrus.Error("Ошибка при удалении заявки:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить заявку"})
		return
	}

	// Редирект на список языков после удаления
	ctx.Redirect(http.StatusSeeOther, "/languages")
}
