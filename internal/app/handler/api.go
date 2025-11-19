package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	_ "LAB1/internal/app/auth"
	"LAB1/internal/app/ds"

	"gorm.io/gorm"

	//"LAB1/internal/app/repository"
	"LAB1/internal/app/storage"

	"github.com/gin-gonic/gin"
)

// NOTE: предполагается, что Handler уже содержит поле Repository *repository.Repository

// ----------------- Lang (услуги) -----------------

// ApiGetLangs - GET /api/languages?query=...
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, langs)
}

// ApiGetLang - GET /api/languages/:id
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

// ApiCreateLang - POST /api/languages
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

// ApiUpdateLang - PUT /api/languages/:id
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

// ApiDeleteLang - PUT/DELETE logical delete -> status = 'удалён'
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

// ApiUploadLangImage - POST /api/languages/:id/image (form file "file")
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

	// обновим запись в БД (photo_key поле — добавь в ds.Lang если нужно)
	_ = h.Repository.UpdateLang(uint(id), map[string]interface{}{"photo_key": objectName})

	c.JSON(http.StatusOK, gin.H{"photo_key": objectName})
}

// ----------------- Glotto (заявки) -----------------
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

func (h *Handler) ApiGetGlottos(c *gin.Context) {

	// 1. Извлекаем пользователя из контекста
	v, exists := c.Get(CtxUserKey)

	// currentUser — указатель, может быть nil
	var currentUser *CurrentUser
	if exists {
		currentUser = v.(*CurrentUser) // ← ВАЖНО: тип *CurrentUser
	}

	status := c.Query("status")

	// -------------------------
	// 2. Модератор → видит все
	// -------------------------
	if currentUser != nil && currentUser.IsModerator {
		glottos, err := h.Repository.GetGlottosFiltered(status, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, glottos)
		return
	}

	// ---------------------------------------------------
	// 3. Гость (нет авторизации) → 401 Unauthorized
	// ---------------------------------------------------
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// ---------------------------------------------------------
	// 4. Создатель → получает только свои собственные заявки
	// ---------------------------------------------------------
	glottos, err := h.Repository.GetLangCalculationsByResearcher(currentUser.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, glottos)
}

// ApiGetGlotto - GET /api/glottos/:id
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

// ApiGetCartIcon - GET /api/glottos/cart (черновик текущего пользователя)
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

// ApiAddServiceToGlotto (POST) — алиас на AddServiceToDraft
func (h *Handler) ApiAddServiceToGlotto(c *gin.Context) {
	langIDStr := c.Param("id")
	langID, err := strconv.Atoi(langIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	researcherID := uint(1)
	if err := h.Repository.AddServiceToDraft(researcherID, uint(langID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ApiUpdateServiceInGlotto — PUT /api/glottos/:id/services
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

// ApiUpdateGlotto — PUT /api/glottos/:id (изменение полей заявки)
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

// ApiFormGlotto — PUT /api/glottos/:id/form
func (h *Handler) ApiFormGlotto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	researcherID := uint(1)
	if err := h.Repository.FormGlotto(uint(id), researcherID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ApiCompleteGlotto — PUT /api/glottos/:id/complete with body { "action": "завершить"|"отклонить" }
func (h *Handler) ApiCompleteGlotto(c *gin.Context) {
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
	moderatorID := uint(1) // заглушка — в реале из JWT
	if err := h.Repository.CompleteGlotto(uint(id), moderatorID, body.Action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Завершить или отклонить расчёт языкового родства (только модератор)
// @Tags lang-calculation
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID расчёта"
// @Param body body CompleteLangCalculationRequest true "Действие"
// @Success 204
// @Failure 400,401,403,404 {object} map[string]string
// @Router /api/lang-calculations/{id}/complete [put]
func (h *Handler) ApiCompleteLangCalculation(c *gin.Context) {
	// 1. Получаем ID расчёта из пути
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid calculation id"})
		return
	}

	// 2. Парсим тело: action должен быть "завершить" или "отклонить"
	var body struct {
		Action string `json:"action" binding:"required,oneof=завершить отклонить"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action is required and must be 'завершить' or 'отклонить'"})
		return
	}

	// 3. Получаем текущего пользователя из middleware
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	isModeratorRaw, exists := c.Get("is_moderator")
	if !exists || !isModeratorRaw.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only moderator can complete calculation"})
		return
	}

	moderatorID := userIDRaw.(uint)

	// 4. Вызываем метод репозитория (переименованный под новую сущность)
	err = h.Repository.CompleteLangCalculation(uint(id), moderatorID, body.Action)
	if err != nil {
		// Можно более детально обрабатывать ошибки, если хочешь
		switch err.Error() {
		case "calculation not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "расчёт не найден"})
		case "invalid status transition":
			c.JSON(http.StatusBadRequest, gin.H{"error": "расчёт можно завершить только в статусе 'на модерации'"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 5. Успешно — ничего не возвращаем
	c.Status(http.StatusNoContent)
}

// ApiDeleteGlotto — DELETE /api/glottos/:id
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

// ----------------- Helpers -----------------
func gormErrNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound || (err != nil && err.Error() == "record not found")
}

// ==========================
// 🧩 Домен Услуги (Lang)
// ==========================
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

// ==========================
// 🧾 Домен Заявки (Glotto)
// ==========================

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

// ==========================
// 🔗 Домен m-m (GlottoLanguage)
// ==========================

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

// ==========================
// 👤 Домен Пользователь
// ==========================

// POST /api/users/register — регистрация нового пользователя
func (h *Handler) ApiRegisterUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "POST /api/users/register"})
}

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

// POST /api/auth/login — аутентификация
func (h *Handler) ApiLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/auth/login"})
}

// POST /api/auth/logout — деавторизация
func (h *Handler) ApiLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST /api/auth/logout"})
}
*/
