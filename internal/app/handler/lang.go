package handler

import (
	"net/http"
	"strconv"
	"time"

	"encoding/json"
	"strings"

	//"errors"
	//"gorm.io/gorm"

	"LAB1/internal/app/auth"
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

	var words []string
	if len(lang.Lexicon) > 0 {
		json.Unmarshal(lang.Lexicon, &words)
	}

	ctx.HTML(http.StatusOK, "lang.html", gin.H{
		"lang":       lang,
		"lexiconCSV": strings.Join(words, ", "),
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

	//researcherID := uint(1) // пока хардкодим
	researcherID := auth.GetCreatorID()

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

	LangCalculation, err := h.Repository.GetDraftByResearcher(researcherID)
	if err != nil {
		logrus.Error("Ошибка при получении черновика:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить черновик заявки"})
		return
	}

	ctx.HTML(http.StatusOK, "chronos.html", gin.H{
		"LangCalculation": LangCalculation,
		"cart_count":      h.Repository.GetLangCount(),
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

func (h *Handler) GetDraftByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/languages")
		return
	}

	draft, err := h.Repository.GetGlottoByID(uint(id))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/languages")
		return
	}

	ctx.HTML(http.StatusOK, "chronos.html", gin.H{
		"LangCalculation": draft,
	})
}

/*
// ApiGetLangs - возвращает JSON список языков, поддерживает query=name и optional filters
func (h *Handler) ApiGetLangs(c *gin.Context) {
	q := c.Query("query")
	var langs []ds.Lang
	var err error
	if q == "" {
		langs, err = h.Repository.GetLangs()
	} else {
		langs, err = h.Repository.GetLangsByName(q)
	}
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, langs)
}*/

// === POST /api/languages ===
func (h *Handler) CreateLanguage(c *gin.Context) {
	var lang ds.Lang
	if err := c.ShouldBindJSON(&lang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := h.Repository.CreateLanguage(&lang); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lang)
}

/*
// === PUT /api/languages/:id ===
func (h *Handler) UpdateLanguage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var lang ds.Lang
	if err := c.ShouldBindJSON(&lang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	lang.ID = uint(id)

	if err := h.Repository.UpdateLanguage(&lang); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	c.JSON(http.StatusOK, lang)
}*/

// === DELETE /api/languages/:id ===
func (h *Handler) DeleteLanguage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.Repository.DeleteLanguage(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted"})
}

/*type CreateLangRequest struct {
	Name        string `json:"name" binding:"required"`
	Family      string `json:"family"`
	Subgroup    string `json:"subgroup"`
	Description string `json:"description"`
}*/

// ApiCreateLang
/*func (h *Handler) ApiCreateLang(c *gin.Context) {
	var req CreateLangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lang := ds.Lang{
		Name:        req.Name,
		Family:      req.Family,
		Subgroup:    req.Subgroup,
		Description: req.Description,
	}
	if err := h.Repository.CreateLang(&lang); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, lang)
}*/
