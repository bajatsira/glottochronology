package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==========================
// 🧩 Домен Услуги (Lang)
// ==========================
/*
// GET /api/languages — список услуг (с фильтрацией)
func (h *Handler) ApiGetLangs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/languages"})
}*/

// GET /api/languages/:id — одна услуга
func (h *Handler) ApiGetLang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET /api/languages/:id"})
}

/*
// POST /api/languages — создать новую услугу
func (h *Handler) ApiCreateLang(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "POST /api/languages"})
}
*/
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
