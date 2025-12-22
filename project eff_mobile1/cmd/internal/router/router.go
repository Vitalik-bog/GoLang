package router

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "effective_mobile_service/internal/handler"
    "effective_mobile_service/internal/repository/postgres"
    "effective_mobile_service/internal/usecase"
)

func SetupRouter(db *sql.DB) *gin.Engine {
    router := gin.Default()
    
    // Инициализация зависимостей
    subscriptionRepo := postgres.NewSubscriptionRepository(db)
    subscriptionUseCase := usecase.NewSubscriptionUseCase(subscriptionRepo)
    subscriptionHandler := handler.NewSubscriptionHandler(subscriptionUseCase)
    
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Документация
        v1.GET("/", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "service": "Сервис управления подписками",
                "version": "1.0",
                "language": "ru",
                "endpoints": []string{
                    "GET   /health - Проверка здоровья",
                    "POST  /subscriptions - Создать подписку",
                    "GET   /subscriptions - Список подписок",
                    "GET   /subscriptions/:id - Получить подписку",
                    "PUT   /subscriptions/:id - Обновить подписку", 
                    "DELETE /subscriptions/:id - Удалить подписку",
                    "GET   /subscriptions/total-cost - Общая стоимость",
                },
            })
        })
        
        // Health check
        v1.GET("/health", func(c *gin.Context) {
            if err := db.Ping(); err != nil {
                c.JSON(500, gin.H{
                    "status":  "unhealthy",
                    "message": "Ошибка подключения к базе данных",
                })
                return
            }
            c.JSON(200, gin.H{
                "status":  "healthy",
                "message": "Сервис работает нормально",
                "service": "Управление подписками",
            })
        })
        
        // Подписки
        subscriptions := v1.Group("/subscriptions")
        {
            subscriptions.POST("", subscriptionHandler.CreateSubscription)
            subscriptions.GET("", subscriptionHandler.ListSubscriptions)
            subscriptions.GET("/total-cost", subscriptionHandler.GetTotalCost)
            subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
            subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
            subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
        }
    }
    
    return router
}
