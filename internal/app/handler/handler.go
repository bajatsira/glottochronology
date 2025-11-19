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
	router.GET("/lang-calculation/draft", h.GetDraft) // просмотр содержимого текущей заявки
	router.POST("/lang-calculation/add/:id", h.AddLanguageToDraft)
	router.POST("/lang-calculation/delete/:id", h.DeleteGlotto) // логическое удаление заявки
	router.GET("/lang-calculation/draft/:id", h.GetDraftByID)
	// --- API маршруты ---
	api := router.Group("/api")

	// --- Домен Услуги (Lang) ---
	api.GET("/languages", h.ApiGetLangs)                   // GET список услуг (с фильтрацией)
	api.GET("/languages/:id", h.ApiGetLang)                // GET одна услуга
	api.POST("/languages", h.ApiCreateLang)                // POST добавить новую услугу (без изображения)
	api.PUT("/languages/:id", h.ApiUpdateLang)             // PUT изменить услугу
	api.DELETE("/languages/:id", h.ApiDeleteLang)          // DELETE удалить услугу (+удаление изображения)
	api.POST("/languages/:id/image", h.ApiUploadLangImage) // POST добавить/заменить изображение услуги (Minio)
	//api.POST("/glottos/add/:language_id", h.ApiAddServiceToDraft) // POST добавить услугу в заявку-черновик

	// --- Домен Заявки (LangCalculation) ---
	api.GET("/LangCalculation", h.ApiGetGlottos)                  // GET список заявок (фильтр по статусу и дате)
	api.GET("/LangCalculation/:id", h.ApiGetGlotto)               // GET одна заявка с услугами
	api.GET("/LangCalculation/cart", h.ApiGetCartIcon)            //не работает из-за минио           // GET иконка корзины (черновик и кол-во услуг)
	api.PUT("/LangCalculation/:id", h.ApiUpdateGlotto)            // PUT изменить поля заявки
	api.PUT("/LangCalculation/:id/form", h.ApiFormGlotto)         // PUT сформировать заявку (создатель)
	api.PUT("/LangCalculation/:id/complete", h.ApiCompleteGlotto) // PUT завершить/отклонить заявку (модератор)
	api.DELETE("/LangCalculation/:id", h.ApiDeleteGlotto)         // DELETE удалить заявку (создатель)

	// --- Домен m-m (GlottoLanguage) ---
	api.POST("/LangCalculation/:id/langs", h.ApiAddServiceToGlotto)        // POST добавить услугу в заявку (m-m)
	api.PUT("/LangCalculation/:id/langs", h.ApiUpdateServiceInGlotto)      // PUT изменить параметры связи (m-m)
	api.DELETE("/LangCalculation/:id/langs", h.ApiDeleteServiceFromGlotto) // DELETE удалить услугу из заявки (m-m)
	/*
		// --- Домен Пользователь ---
		api.POST("/users/register", h.ApiRegisterUser) // POST регистрация нового пользователя
		api.GET("/users/me", h.ApiGetCurrentUser)      // GET данные текущего пользователя
		api.PUT("/users/me", h.ApiUpdateCurrentUser)   // PUT обновить данные пользователя

		// --- Домен Аутентификация ---
		api.POST("/auth/login", h.ApiLogin)   // POST аутентификация
		api.POST("/auth/logout", h.ApiLogout) // POST деавторизация*/

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

/*
// RegisterAPI регистрирует REST API с префиксом /api
func (h *Handler) RegisterAPI(router *gin.Engine) {
    api := router.Group("/api")
    {
        // --- Домен Услуги (Lang) ---
        api.GET("/languages", h.ApiGetLangs)                 // GET список услуг (с фильтрацией)
        api.GET("/languages/:id", h.ApiGetLang)              // GET одна услуга
       // api.POST("/languages", h.ApiCreateLang)              // POST добавить новую услугу (без изображения)
        api.PUT("/languages/:id", h.ApiUpdateLang)           // PUT изменить услугу
        api.DELETE("/languages/:id", h.ApiDeleteLang)        // DELETE удалить услугу (+удаление изображения)
        api.POST("/languages/:id/image", h.ApiUploadLangImage) // POST добавить/заменить изображение услуги (Minio)
        api.POST("/glottos/add/:language_id", h.ApiAddServiceToDraft) // POST добавить услугу в заявку-черновик

        // --- Домен Заявки (Glotto) ---
        api.GET("/glottos", h.ApiGetGlottos)                 // GET список заявок (фильтр по статусу и дате)
        api.GET("/glottos/:id", h.ApiGetGlotto)              // GET одна заявка с услугами
        api.GET("/glottos/cart", h.ApiGetCartIcon)           // GET иконка корзины (черновик и кол-во услуг)
        api.PUT("/glottos/:id", h.ApiUpdateGlotto)           // PUT изменить поля заявки
        api.PUT("/glottos/:id/form", h.ApiFormGlotto)        // PUT сформировать заявку (создатель)
        api.PUT("/glottos/:id/complete", h.ApiCompleteGlotto)// PUT завершить/отклонить заявку (модератор)
        api.DELETE("/glottos/:id", h.ApiDeleteGlotto)        // DELETE удалить заявку (создатель)

        // --- Домен m-m (GlottoLanguage) ---
        api.POST("/glottos/:id/services", h.ApiAddServiceToGlotto)   // POST добавить услугу в заявку (m-m)
        api.PUT("/glottos/:id/services", h.ApiUpdateServiceInGlotto) // PUT изменить параметры связи (m-m)
        api.DELETE("/glottos/:id/services", h.ApiDeleteServiceFromGlotto) // DELETE удалить услугу из заявки (m-m)

        // --- Домен Пользователь ---
        api.POST("/users/register", h.ApiRegisterUser)       // POST регистрация нового пользователя
        api.GET("/users/me", h.ApiGetCurrentUser)            // GET данные текущего пользователя
        api.PUT("/users/me", h.ApiUpdateCurrentUser)         // PUT обновить данные пользователя

        // --- Домен Аутентификация ---
        api.POST("/auth/login", h.ApiLogin)                  // POST аутентификация
        api.POST("/auth/logout", h.ApiLogout)                // POST деавторизация
    }
}

*/
