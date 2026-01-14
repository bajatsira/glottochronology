package handler

import (
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

/*
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
	router.POST("/languages/:id/lexicon", h.UpdateLangLexiconForm)   // добавляб слова в заявку
	router.POST("/lang-calculation/set-base/:id", h.SetBaseLanguage) // выбор базового языка

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
	api.GET("/lang-calculation", h.ApiGetGlottos)           // GET список заявок (фильтр по статусу и дате)
	api.GET("/lang-calculation/:id", h.ApiGetGlotto)        // GET одна заявка с услугами
	api.GET("/lang-calculation/cart", h.ApiGetCartIcon)     //не работает из-за минио           // GET иконка корзины (черновик и кол-во услуг)
	api.PUT("/lang-calculation/:id", h.ApiUpdateGlotto)     // PUT изменить поля заявки
	api.PUT("/lang-calculation/:id/form", h.ApiFormGlotto)  // PUT сформировать заявку (создатель)
	api.POST("/lang-calculation/:id/form", h.ApiFormGlotto) // POST сформировать заявку (создатель)

	api.PUT("/lang-calculation/:id/complete", h.ApiCompleteGlotto) // PUT завершить/отклонить заявку (модератор)
	router.POST("/api/lang-calculation/:id/complete", h.ApiCompleteGlotto)

	api.DELETE("/lang-calculation/:id", h.ApiDeleteGlotto) // DELETE удалить заявку (создатель)

	// --- Домен m-m (GlottoLanguage) ---
	api.POST("/lang-calculation/:id/langs", h.ApiAddServiceToGlotto)        // POST добавить услугу в заявку (m-m)
	api.PUT("/lang-calculation/:id/langs", h.ApiUpdateServiceInGlotto)      // PUT изменить параметры связи (m-m)
	api.DELETE("/lang-calculation/:id/langs", h.ApiDeleteServiceFromGlotto) // DELETE удалить услугу из заявки (m-m)

	// --- Домен Пользователь ---
	api.POST("/users/register", h.ApiRegisterUser) // POST регистрация нового пользователя
	api.GET("/users/me", h.ApiGetCurrentUser)      // GET данные текущего пользователя
	api.PUT("/users/me", h.ApiUpdateCurrentUser)   // PUT обновить данные пользователя

	// --- Домен Аутентификация ---
	api.POST("/auth/login", h.ApiLogin)   // POST аутентификация
	api.POST("/auth/logout", h.ApiLogout) // POST деавторизация

}*/

/*

func (h *Handler) RegisterHandler(router *gin.Engine) {

	// 1. ПРИМЕНЕНИЕ ГЛОБАЛЬНОГО MIDDLEWARE АУТЕНТИФИКАЦИИ
	// Он обрабатывает JWT и Cookie и устанавливает CurrentUser в контекст.
	router.Use(AuthMiddleware(h.Repository))

	// 2. SWAGGER ROUTE
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//router.GET("/", h.GetOrders)
	//router.GET("/order/:id", h.GetOrder)

	router.GET("/languages", h.GetLangs)
	//r.GET("/order/:id", handler.GetOrder) // вот наш новый обработчик
	router.GET("/lang/:id", h.GetLang)
	//router.GET("/chronos", h.GetChronos)

	// Эти фронтенд-маршруты требуют авторизации, так как работают с личными черновиками
	router.GET("/lang-calculation/draft", RequireAuth(), h.GetDraft) // просмотр содержимого текущей заявки
	router.POST("/lang-calculation/add/:id", RequireAuth(), h.AddLanguageToDraft)
	router.POST("/lang-calculation/delete/:id", RequireAuth(), h.DeleteGlotto) // логическое удаление заявки
	router.GET("/lang-calculation/draft/:id", RequireAuth(), h.GetDraftByID)
	router.POST("/languages/:id/lexicon", RequireAuth(), h.UpdateLangLexiconForm)   // добавляб слова в заявку
	router.POST("/lang-calculation/set-base/:id", RequireAuth(), h.SetBaseLanguage) // выбор базового языка

	// --- API маршруты ---
	api := router.Group("/api")

	// --- Домен Услуги (Lang) ---
	// GETs оставляем публичными, как "чтение-получение данных"
	api.GET("/languages", h.ApiGetLangs)    // GET список услуг (с фильтрацией)
	api.GET("/languages/:id", h.ApiGetLang) // GET одна услуга

	// CRUD методы (Только для Модератора/Лингвиста)
	api.POST("/languages", RequireLinguist(), h.ApiCreateLang)                // POST добавить новую услугу
	api.PUT("/languages/:id", RequireLinguist(), h.ApiUpdateLang)             // PUT изменить услугу
	api.DELETE("/languages/:id", RequireLinguist(), h.ApiDeleteLang)          // DELETE удалить услугу
	api.POST("/languages/:id/image", RequireLinguist(), h.ApiUploadLangImage) // POST изображение

	// --- Домен Заявки (LangCalculation) ---
	// GET список заявок: закрыт для гостей (RequireAuth)
	api.GET("/lang-calculation", RequireAuth(), h.ApiGetGlottos) // GET список заявок (фильтр по статусу и дате)

	// Остальные GET/CRUD требуют авторизации
	api.GET("/lang-calculation/:id", RequireAuth(), h.ApiGetGlotto)        // GET одна заявка с услугами
	api.GET("/lang-calculation/cart", h.ApiGetCartIcon)                    // GET иконка корзины (оставляем без Auth, если логика корзины простая)
	api.PUT("/lang-calculation/:id", RequireAuth(), h.ApiUpdateGlotto)     // PUT изменить поля заявки
	api.PUT("/lang-calculation/:id/form", RequireAuth(), h.ApiFormGlotto)  // PUT сформировать заявку (создатель)
	api.POST("/lang-calculation/:id/form", RequireAuth(), h.ApiFormGlotto) // POST сформировать заявку (создатель)

	// Завершение заявки: Только для Модератора/Лингвиста
	api.PUT("/lang-calculation/:id/complete", RequireLinguist(), h.ApiCompleteGlotto)

	// Удаление дубликата маршрута, который был вне api группы:
	// router.POST("/api/lang-calculation/:id/complete", h.ApiCompleteGlotto)

	api.DELETE("/lang-calculation/:id", RequireAuth(), h.ApiDeleteGlotto) // DELETE удалить заявку (создатель)

	// --- Домен m-m (GlottoLanguage) ---
	// Все операции с содержимым заявки требуют авторизации
	api.POST("/lang-calculation/:id/langs", RequireAuth(), h.ApiAddServiceToGlotto)        // POST добавить услугу в заявку (m-m)
	api.PUT("/lang-calculation/:id/langs", RequireAuth(), h.ApiUpdateServiceInGlotto)      // PUT изменить параметры связи (m-m)
	api.DELETE("/lang-calculation/:id/langs", RequireAuth(), h.ApiDeleteServiceFromGlotto) // DELETE удалить услугу из заявки (m-m)

	// --- Домен Пользователь ---
	api.POST("/users/register", h.ApiRegisterUser)              // POST регистрация (Публичный)
	api.GET("/users/me", RequireAuth(), h.ApiGetCurrentUser)    // GET данные текущего пользователя
	api.PUT("/users/me", RequireAuth(), h.ApiUpdateCurrentUser) // PUT обновить данные пользователя

	// --- Домен Аутентификация ---
	api.POST("/auth/login", h.ApiLogin)                  // POST аутентификация (Публичный)
	api.POST("/auth/logout", RequireAuth(), h.ApiLogout) // POST деавторизация (Требует авторизации для очистки сессии)
}
*/

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// --- Swagger ---
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Домен Языки ---
	langs := router.Group("/api/langs")
	{
		langs.GET("/", h.ApiGetLangs)
		langs.GET("/:id", h.ApiGetLang)
		langs.POST("/", RequireLinguist(), h.ApiCreateLang)
		langs.PUT("/:id", RequireLinguist(), h.ApiUpdateLang)
		langs.DELETE("/:id", RequireLinguist(), h.ApiDeleteLang)
	}

	// --- Домен Заявки + Корзина ---
	api := router.Group("/api")
	api.GET("/lang-calculation", RequireAuth(), h.ApiGetGlottos)
	api.GET("/lang-calculation/draft/count", RequireAuth(), h.ApiGetGlottoDraftCount) // <--- ВОТ НОВЫЙ МАРШРУТ
	api.POST("/lang-calculation/:id/langs", RequireAuth(), h.ApiAddServiceToGlotto)
	api.PUT("/lang-calculation/:id/complete", RequireLinguist(), h.ApiCompleteLangCalculation)
	api.DELETE("/lang-calculation/:id", RequireAuth(), h.ApiDeleteGlotto)

	// --- Домен Аутентификация ---
	api.POST("/auth/login", h.ApiLogin)
	api.POST("/auth/logout", RequireAuth(), h.ApiLogout)

	// --- Домен Пользователи (для регистрации) ---
	users := router.Group("/users")
	{
		users.POST("/register", h.ApiRegisterUser)
	}
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	//router.Static("/styles", "./styles")
	//r.Static("/static", "./resources")

	router.Static("/styles", "resources/styles")
	router.Static("/img", "resources/img")

}

/*
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
