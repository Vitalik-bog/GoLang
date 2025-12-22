package entity

// Русские названия полей для Swagger/документации
const (
    FieldServiceNameRU = "Название сервиса"
    FieldPriceRU       = "Стоимость подписки (руб/мес)"
    FieldUserIDRU      = "ID пользователя (UUID)"
    FieldStartDateRU   = "Дата начала подписки (ММ-ГГГГ)"
    FieldEndDateRU     = "Дата окончания подписки (ММ-ГГГГ, опционально)"
    FieldCreatedAtRU   = "Дата создания"
    FieldUpdatedAtRU   = "Дата обновления"
    FieldTotalCostRU   = "Общая стоимость"
)

// Русские описания для API
var SubscriptionDescriptions = map[string]string{
    "Create": "Создание новой подписки на сервис",
    "Get":    "Получение информации о подписке",
    "Update": "Обновление данных подписки",
    "Delete": "Удаление подписки",
    "List":   "Получение списка подписок с фильтрацией",
    "TotalCost": "Расчет общей стоимости подписок за период",
}
