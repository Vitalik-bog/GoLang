package ru

const (
    // Общие сообщения
    HealthCheckSuccess = "Сервис работает нормально"
    HealthCheckFailed  = "Проблемы с подключением к базе данных"
    
    // Подписки
    SubscriptionCreated     = "Подписка успешно создана"
    SubscriptionUpdated     = "Подписка обновлена"
    SubscriptionDeleted     = "Подписка удалена"
    SubscriptionNotFound    = "Подписка не найдена"
    SubscriptionListEmpty   = "Список подписок пуст"
    
    // Ошибки валидации
    InvalidJSON            = "Неверный формат JSON"
    MissingRequiredFields  = "Отсутствуют обязательные поля"
    InvalidServiceName     = "Название сервиса должно быть от 1 до 100 символов"
    InvalidPrice           = "Цена должна быть положительным числом"
    InvalidUserID          = "Неверный формат ID пользователя"
    InvalidDate            = "Неверный формат даты (используйте ММ-ГГГГ)"
    InvalidSubscriptionID  = "Неверный формат ID подписки"
    
    // База данных
    DatabaseError          = "Ошибка базы данных"
    ConnectionError        = "Ошибка подключения к базе данных"
    
    // Успешные операции
    OperationSuccess       = "Операция выполнена успешно"
    TotalCostCalculated    = "Общая стоимость рассчитана"
)
