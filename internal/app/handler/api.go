package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	_ "LAB1/internal/app/auth"
	"LAB1/internal/app/ds"

	"gorm.io/gorm"

	//"LAB1/internal/app/repository"
	"LAB1/internal/app/storage"

	"github.com/gin-gonic/gin"
)

// ApiGetLangs godoc
// @Summary Получить список языков
// @Description Возвращает список языков, опционально фильтрованных по запросу.
// @Tags Languages
// @Accept  json
// @Produce  json
// @Param   query   query   string  false  "Поисковый запрос для фильтрации по имени"
// @Success 200 {array} ds.Lang "Список языков"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/languages [get]
func (h *Handler) ApiGetLangs(c *gin.Context) { // ApiGetLangs - GET /api/languages?query=...

	q := c.Query("query")
	var langs []ds.Lang
	var err error
	if q == "" {
		langs, err = h.Repository.GetLangs()
	} else {
		langs, err = h.Repository.GetLangsByName(q)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, langs)
}

// ApiGetLang godoc
// @Summary Получить язык по ID
// @Description Возвращает полную информацию о языке.
// @Tags Languages
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID языка"
// @Success 200 {object} ds.Lang "Информация о языке"
// @Failure 400 {object} object "Неверный ID"
// @Failure 404 {object} object "Язык не найден"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/languages/{id} [get]
func (h *Handler) ApiGetLang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	lang, err := h.Repository.GetLang(id)
	if err != nil {
		if gormErrNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lang)
}

// ApiCreateLang godoc
// @Summary Создать новый язык
// @Description Доступно только для Модератора/Лингвиста.
// @Tags Languages
// @Accept  json
// @Produce  json
// @Param   language  body   ds.Lang  true  "Данные нового языка"
// @Success 201 {object} ds.Lang "Успешное создание"
// @Failure 400 {object} object "Неверные данные"
// @Failure 401 {object} object "Требуется авторизация"
// @Failure 403 {object} object "Недостаточно прав (не Модератор)"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/languages [post]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiCreateLang(c *gin.Context) {
	var body ds.Lang
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repository.CreateLang(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

// ApiUpdateLang godoc
// @Summary Обновить информацию о языке
// @Description Доступно только для Модератора/Лингвиста.
// @Tags Languages
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID языка"
// @Param   updates  body   map[string]interface{}  true  "Обновляемые поля"
// @Success 200 {object} ds.Lang "Успешное обновление"
// @Failure 400 {object} object "Неверный ID или данные"
// @Failure 401 {object} object "Требуется авторизация"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/languages/{id} [put]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiUpdateLang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repository.UpdateLang(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lang, _ := h.Repository.GetLang(id)
	c.JSON(http.StatusOK, lang)
}

// ApiDeleteLang godoc
// @Summary Логически удалить (скрыть) язык
// @Description Устанавливает статус языка как 'удалён'. Доступно только для Модератора/Лингвиста.
// @Tags Languages
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID языка"
// @Success 200 {object} object "Успешное удаление"
// @Failure 400 {object} object "Неверный ID"
// @Failure 401 {object} object "Требуется авторизация"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/languages/{id} [delete]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiDeleteLang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Repository.SoftDeleteLang(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "удалён"})
}

// ApiUploadLangImage godoc
// @Summary Загрузить изображение для языка
// @Description Загрузка файла изображения. Использует multipart/form-data с полем "file". Доступно только для Модератора/Лингвиста.
// @Tags Languages
// @Accept  mpfd
// @Produce  json
// @Param   id   path   int  true  "ID языка"
// @Param   file  formData   file  true  "Файл изображения"
// @Success 200 {object} map[string]string "Ключ загруженного файла в Minio"
// @Failure 400 {object} object "Неверный ID или файл не предоставлен"
// @Failure 401 {object} object "Требуется авторизация"
// @Failure 403 {object} object "Недостаточно прав"
// @Failure 500 {object} object "Ошибка загрузки/сервера"
// @Router /api/languages/{id}/image [post]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiUploadLangImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	// генерируем имя на латинице: lang_<id>_<ts><ext>
	ext := path.Ext(header.Filename)
	objectName := fmt.Sprintf("lang_%d_%d%s", id, time.Now().Unix(), ext)

	minioClient, err := storage.NewMinio()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "minio init: " + err.Error()})
		return
	}
	if err := minioClient.Upload(c.Request.Context(), objectName, file, header.Size, header.Header.Get("Content-Type")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed: " + err.Error()})
		return
	}

	_ = h.Repository.UpdateLang(uint(id), map[string]interface{}{"photo_key": objectName})

	c.JSON(http.StatusOK, gin.H{"photo_key": objectName})
}

// otto (заявки)
/*
// ApiGetGlottos - GET /api/glottos?status=&date_from=&date_to=
func (h *Handler) ApiGetGlottos(c *gin.Context) {
	status := c.Query("status")
	df := c.Query("date_from") // yyyy-mm-dd
	dt := c.Query("date_to")
	var dateFrom, dateTo *time.Time
	if df != "" {
		t, err := time.Parse("2006-01-02", df)
		if err == nil {
			dateFrom = &t
		}
	}
	if dt != "" {
		t, err := time.Parse("2006-01-02", dt)
		if err == nil {
			dateTo = &t
		}
	}
	glottos, err := h.Repository.GetGlottosFiltered(status, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, glottos)
}
*/

// ApiGetGlottos godoc
// @Summary Получить список заявок
// @Description Получает список заявок. Гость: 401. Создатель: только свои заявки. Модератор: все заявки.
// @Tags Заявки (LangCalculation)
// @Produce  json
// @Param   status  query   string  false  "Фильтр по статусу заявки (для модератора)"
// @Success 200 {array} ds.LangCalculation "Успешное получение списка заявок"
// @Failure 401 {object} object "Требуется аутентификация"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/lang-calculation [get]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiGetGlottos(c *gin.Context) {
	v, _ := c.Get(CtxUserKey)
	currentUser := v.(*CurrentUser)

	status := c.Query("status")

	if currentUser.IsLinguist {
		glottos, err := h.Repository.GetGlottosFiltered(status, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, glottos)
		return
	}

	glottos, err := h.Repository.GetLangCalculationsByResearcher(currentUser.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, glottos)
}

// ApiGetGlottoDraftCount godoc
// @Summary Получить количество языков в черновике (корзине)
// @Description Возвращает количество языков в текущем черновике пользователя. Требуется авторизация.
// @Tags Заявки (LangCalculation)
// @Produce  json
// @Success 200 {object} object{count=int} "Количество языков в черновике"
// @Failure 401 "Требуется аутентификация"
// @Router /api/lang-calculation/draft/count [get]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiGetGlottoDraftCount(c *gin.Context) {
	v, _ := c.Get(CtxUserKey)
	currentUser := v.(*CurrentUser)

	// Используем GetLangCountForUser с ID текущего пользователя
	count := h.Repository.GetLangCountForUser(currentUser.ID)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// ApiGetGlotto godoc
// @Summary Получить заявку по ID
// @Description Возвращает детальную информацию о заявке. Требуется авторизация.
// @Tags Заявки (LangCalculation)
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID заявки"
// @Success 200 {object} ds.LangCalculation
// @Failure 400 {object} object "Неверный ID"
// @Failure 401 {object} object "Требуется авторизация"
// @Failure 404 {object} object "Заявка не найдена"
// @Router /api/lang-calculation/{id} [get]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiGetGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	g, err := h.Repository.GetGlottoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

// ApiGetCartIcon godoc
// @Summary Получить статус черновика/корзины
// @Description Возвращает ID черновика и количество языков в нем для текущего пользователя.
// @Tags Заявки (LangCalculation)
// @Accept  json
// @Produce  json
// @Success 200 {object} object "ID черновика и количество"
// @Router /api/lang-calculation/cart [get]
func (h *Handler) ApiGetCartIcon(c *gin.Context) {
	// creator singleton: пока используем 1 (как в проекте)
	researcherID := uint(1)
	g, err := h.Repository.GetDraftByResearcher(researcherID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"draft_id": nil, "count": 0})
		return
	}
	// count:
	count := h.Repository.GetLangCount()
	c.JSON(http.StatusOK, gin.H{"draft_id": g.ID, "count": count})
}

// ApiAddServiceToGlotto godoc
// @Summary Добавить язык в черновик (корзину)
// @Description Добавляет язык в текущий черновик пользователя. Если черновика нет, он создается. Требуется авторизация.
// @Tags Заявки (LangCalculation)
// @Param   id   path   int  true  "ID языка, который нужно добавить"
// @Success 204 "Успешное добавление"
// @Failure 400 "Неверный ID"
// @Failure 401 "Требуется авторизация"
// @Failure 500 "Ошибка сервера/репозитория"
// @Router /api/lang-calculation/{id}/langs [post]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiAddServiceToGlotto(c *gin.Context) {
	langIDStr := c.Param("id")
	langID, err := strconv.Atoi(langIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Получаем текущего пользователя из контекста
	v, _ := c.Get(CtxUserKey)
	currentUser := v.(*CurrentUser)

	// Используем ID текущего пользователя
	if err := h.Repository.AddServiceToDraft(currentUser.ID, uint(langID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ApiUpdateServiceInGlotto godoc
// @Summary Обновить параметры языка в заявке
// @Description Изменение параметров добавленного языка в заявке (заглушка в реализации). Требуется авторизация.
// @Tags Заявки (LangCalculation)
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID заявки (LangCalculation)"
// @Param   updates  body   map[string]interface{}  true  "Обновляемые параметры"
// @Success 200 {object} object "Статус обновления"
// @Failure 400 "Неверный ID или данные"
// @Failure 401 "Требуется авторизация"
// @Router /api/lang-calculation/{id}/langs [put]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiUpdateServiceInGlotto(c *gin.Context) {
	// пример: body { "language_id": 5, "value": "...", "position": 2 }
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// реализуй логику по твоим данным — заглушка:
	c.JSON(http.StatusOK, gin.H{"status": "ok", "updated": body})
}

// ApiDeleteServiceFromGlotto godoc
// @Summary Удалить язык из заявки
// @Description Удаляет язык из конкретной заявки. Требуется авторизация.
// @Tags Заявки (LangCalculation)
// @Param   id   path   int  true  "ID заявки (LangCalculation)"
// @Param   language_id  query   int  true  "ID языка, который нужно удалить"
// @Success 204 "Успешное удаление"
// @Failure 400 "Неверный ID или language_id"
// @Failure 401 "Требуется авторизация"
// @Failure 500 "Ошибка сервера"
// @Router /api/lang-calculation/{id}/langs [delete]
// @Security ApiKeyAuth
// @Security CookieAuth
// ApiDeleteServiceFromGlotto — DELETE /api/glottos/:id/services?language_id=5
func (h *Handler) ApiDeleteServiceFromGlotto(c *gin.Context) {
	glottoIDStr := c.Param("id")
	glottoID, err := strconv.Atoi(glottoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	langIDStr := c.Query("language_id")
	langID, err := strconv.Atoi(langIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "language_id required"})
		return
	}
	// raw SQL delete from glotto_languages where glotto_id = ? and language_id = ?
	if err := h.Repository.DeleteGlottoLanguage(uint(glottoID), uint(langID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ApiUpdateGlotto godoc
// @Summary Обновить поля заявки
// @Description Обновление несистемных полей заявки. Требуется авторизация (Создатель или Модератор).
// @Tags Заявки (LangCalculation)
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID заявки"
// @Param   updates  body   map[string]interface{}  true  "Обновляемые поля"
// @Success 200 {object} ds.LangCalculation "Обновленная заявка"
// @Failure 400 "Неверный ID или данные"
// @Failure 401 "Требуется авторизация"
// @Failure 500 "Ошибка сервера"
// @Router /api/lang-calculation/{id} [put]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiUpdateGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// protection: не разрешать менять системные поля
	delete(updates, "id")
	delete(updates, "researcher_id")
	delete(updates, "linguist_id")
	delete(updates, "date_create")
	delete(updates, "date_finish")
	delete(updates, "status")
	if err := h.Repository.UpdateGlotto(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g, _ := h.Repository.GetGlottoByID(uint(id))
	c.JSON(http.StatusOK, g)
}

// ApiFormGlotto godoc
// @Summary      Сформировать (отправить) заявку
// @Description  Изменяет статус черновика на 'сформирован', устанавливает дату формирования. Проверяет, что заявка не пуста.
// @Tags         Заявки (LangCalculation)
// @Param        id   path      int  true  "ID заявки"
// @Success      204  "Успешное формирование"
// @Failure      400  {object}  object{error=string} "Ошибка валидации (например, пустая заявка)"
// @Failure      403  {object}  object{error=string} "Доступ запрещен (вы не владелец)"
// @Failure      401  "Требуется аутентификация"
// @Router       /api/lang-calculation/{id}/form [put]
// @Security     ApiKeyAuth
// @Security     CookieAuth
func (h *Handler) ApiFormGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	v, _ := c.Get(CtxUserKey)
	currentUser := v.(*CurrentUser)

	if err := h.Repository.FormGlotto(uint(id), currentUser.ID); err != nil {
		if strings.Contains(err.Error(), "empty request") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ApiCompleteGlotto — PUT /api/glottos/:id/complete with body { "action": "завершить"|"отклонить" } -- это для авторизации уже
/*func (h *Handler) ApiCompleteGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Action string `json:"action"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	moderatorID := uint(1) // заглушка
	if err := h.Repository.CompleteLangCalculation(uint(id), moderatorID, body.Action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}*/

func (h *Handler) ApiCompleteGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	moderatorID := uint(1)

	action := "завершить"

	if err := h.Repository.CompleteLangCalculation(uint(id), moderatorID, action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ApiSetBaseLanguage godoc
// @Summary      Назначить базовый язык для заявки
// @Description  Устанавливает один из языков в черновике как базовый. Доступно только для создателя.
// @Tags         Заявки (LangCalculation)
// @Param        calculation_id path int true "ID заявки (черновика)"
// @Param        language_id    path int true "ID языка, который должен стать базовым"
// @Success      204 "Базовый язык успешно назначен"
// @Failure      403 "Доступ запрещен (вы не владелец или заявка не черновик)"
// @Failure      401 "Требуется аутентификация"
// @Router       /api/lang-calculation/{calculation_id}/base/{language_id} [put]
// @Security     ApiKeyAuth
// @Security     CookieAuth
func (h *Handler) ApiSetBaseLanguage(c *gin.Context) {
	calcIDStr := c.Param("id")
	langIDStr := c.Param("language_id")

	calcID, err1 := strconv.Atoi(calcIDStr)
	langID, err2 := strconv.Atoi(langIDStr)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	// Получаем текущего пользователя из контекста.
	v, _ := c.Get(CtxUserKey)
	currentUser := v.(*CurrentUser)

	// Вызываем новый метод репозитория.
	if err := h.Repository.SetBaseLanguage(uint(calcID), uint(langID), currentUser.ID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ApiCompleteLangCalculation godoc
// @Summary      Завершить или отклонить заявку (для Модератора)
// @Description  Обновляет статус заявки на 'завершён' или 'отклонён'. Доступно только для Модератора/Лингвиста. При завершении, производит расчет.
// @Tags         Заявки (LangCalculation)
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "ID заявки"
// @Param        action  body      object{action=string}  true  "Действие: 'завершить' или 'отклонить'"
// @Success      204     "Успешное обновление статуса"
// @Failure      400     {object}  object "Неверный ID или данные"
// @Failure      401     {object}  object "Требуется аутентификация"
// @Failure      403     {object}  object "Недостаточно прав (не Модератор)"
// @Failure      500     {object}  object "Ошибка сервера"
// @Router       /api/lang-calculation/{id}/complete [put]
// @Security     ApiKeyAuth
// @Security     CookieAuth
func (h *Handler) ApiCompleteLangCalculation(c *gin.Context) {

	if c.IsAborted() {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Action string `json:"action"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, _ := c.Get(CtxUserKey)
	current := v.(*CurrentUser)

	if err := h.Repository.CompleteLangCalculation(uint(id), current.ID, body.Action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ApiDeleteGlotto godoc
// @Summary Удалить заявку (логическое удаление)
// @Description Логически удаляет заявку. Требуется авторизация.
// @Tags Заявки (LangCalculation)
// @Param   id   path   int  true  "ID заявки"
// @Success 204 "Успешное удаление"
// @Failure 400 "Неверный ID"
// @Failure 401 "Требуется авторизация"
// @Failure 500 "Ошибка сервера"
// @Router /api/lang-calculation/{id} [delete]
// @Security ApiKeyAuth
// @Security CookieAuth
func (h *Handler) ApiDeleteGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// логическое удаление через raw sql (есть метод DeleteDraftSQL)
	if err := h.Repository.DeleteDraftSQL(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func gormErrNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound || (err != nil && err.Error() == "record not found")
}

// UpdateLangLexiconForm godoc
// @Summary Обновить лексикон языка (через форму)
// @Description Принимает данные лексикона в формате CSV из HTML-формы и сохраняет.
// @Tags Languages
// @Accept application/x-www-form-urlencoded <-- ИСПРАВЛЕНО
// @Produce  html
// @Param   id   path   int  true  "ID языка"
// @Param   lexicon_csv  formData   string  true  "Лексикон в формате CSV"
// @Success 303 "Перенаправление на страницу языка после успешного обновления"
// @Failure 400 "Неверный ID или пустой лексикон"
// @Failure 500 "Ошибка сервера"
// @Router /languages/{id}/lexicon [post]
func (h *Handler) UpdateLangLexiconForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id")
		return
	}

	csv := c.PostForm("lexicon_csv")
	if csv == "" {
		c.String(http.StatusBadRequest, "lexicon is empty")
		return
	}

	parts := strings.Split(csv, ",")
	var words []string
	for _, p := range parts {
		w := strings.TrimSpace(p)
		if w != "" {
			words = append(words, w)
		}
	}

	b, err := json.Marshal(words)
	if err != nil {
		c.String(http.StatusInternalServerError, "marshal error")
		return
	}

	if err := h.Repository.UpdateLang(uint(id), map[string]interface{}{
		"lexicon": b,
	}); err != nil {
		c.String(http.StatusInternalServerError, "db error: "+err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/lang/"+idStr)
}

/*
// GET /api/languages — список услуг (с фильтрацией)
func (h *Handler) ApiGetLangs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/languages"})
}*/
/*
// GET /api/languages/:id — одна услуга
func (h *Handler) ApiGetLang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/languages/:id"})
}

/*
// POST /api/languages — создать новую услугу
func (h *Handler) ApiCreateLang(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "POST /api/languages"})
}

// PUT /api/languages/:id — изменить услугу
func (h *Handler) ApiUpdateLang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT /api/languages/:id"})
}

// DELETE /api/languages/:id — удалить услугу (+удаление изображения)
func (h *Handler) ApiDeleteLang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DELETE /api/languages/:id"})
}

// POST /api/languages/:id/image — добавить/заменить изображение услуги
func (h *Handler) ApiUploadLangImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/languages/:id/image"})
}

// POST /api/glottos/add/:language_id — добавить услугу в заявку-черновик
func (h *Handler) ApiAddServiceToDraft(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/glottos/add/:language_id"})
}



// GET /api/glottos — список заявок (фильтр по статусу и дате)
func (h *Handler) ApiGetGlottos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/glottos"})
}

// GET /api/glottos/:id — одна заявка с услугами
func (h *Handler) ApiGetGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/glottos/:id"})
}

// GET /api/glottos/cart — иконка корзины (черновик + количество услуг)
func (h *Handler) ApiGetCartIcon(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/glottos/cart"})
}

// PUT /api/glottos/:id — изменить поля заявки
func (h *Handler) ApiUpdateGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT /api/glottos/:id"})
}

// PUT /api/glottos/:id/form — сформировать заявку (создатель)
func (h *Handler) ApiFormGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT /api/glottos/:id/form"})
}

// PUT /api/glottos/:id/complete — завершить/отклонить заявку (модератор)
func (h *Handler) ApiCompleteGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT /api/glottos/:id/complete"})
}

// DELETE /api/glottos/:id — удалить заявку
func (h *Handler) ApiDeleteGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DELETE /api/glottos/:id"})
}


// POST /api/glottos/:id/services — добавить услугу в заявку
func (h *Handler) ApiAddServiceToGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/glottos/:id/services"})
}

// PUT /api/glottos/:id/services — изменить параметры связи (m-m)
func (h *Handler) ApiUpdateServiceInGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT /api/glottos/:id/services"})
}

// DELETE /api/glottos/:id/services — удалить услугу из заявки
func (h *Handler) ApiDeleteServiceFromGlotto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DELETE /api/glottos/:id/services"})
}
*/
/*
// POST /api/users/register — регистрация нового пользователя
func (h *Handler) ApiRegisterUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "POST /api/users/register"})
}
*/
// GET /api/users/me — получить данные текущего пользователя
func (h *Handler) ApiGetCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/users/me"})
}

// PUT /api/users/me — обновить данные пользователя
func (h *Handler) ApiUpdateCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT /api/users/me"})
}

// ==========================
// 🔐 Домен Аутентификация
// ==========================
/*
// POST /api/auth/login — аутентификация
func (h *Handler) ApiLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/auth/login"})
}

// POST /api/auth/logout — деавторизация
func (h *Handler) ApiLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/auth/logout"})
}
*/
