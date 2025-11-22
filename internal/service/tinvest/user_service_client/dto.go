package user_service_client

type AccountType int32

const (
	AccountType_ACCOUNT_TYPE_UNSPECIFIED AccountType = 0 //Тип аккаунта не определён.
	AccountType_ACCOUNT_TYPE_TINKOFF     AccountType = 1 //Брокерский счёт Т-Инвестиций.
	AccountType_ACCOUNT_TYPE_TINKOFF_IIS AccountType = 2 //ИИС.
	AccountType_ACCOUNT_TYPE_INVEST_BOX  AccountType = 3 //Инвесткопилка.
	AccountType_ACCOUNT_TYPE_INVEST_FUND AccountType = 4 //Фонд денежного рынка.
)

type AccessLevel int32

const (
	AccessLevel_ACCOUNT_ACCESS_LEVEL_UNSPECIFIED AccessLevel = 0 //Уровень доступа не определён.
	AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS AccessLevel = 1 //Полный доступ к счёту.
	AccessLevel_ACCOUNT_ACCESS_LEVEL_READ_ONLY   AccessLevel = 2 //Доступ с уровнем прав «только чтение».
	AccessLevel_ACCOUNT_ACCESS_LEVEL_NO_ACCESS   AccessLevel = 3 //Доступа нет.
)

type AccountStatus int32

const (
	AccountStatus_ACCOUNT_STATUS_UNSPECIFIED AccountStatus = 0 //Статус счёта не определён.
	AccountStatus_ACCOUNT_STATUS_NEW         AccountStatus = 1 //Новый, в процессе открытия.
	AccountStatus_ACCOUNT_STATUS_OPEN        AccountStatus = 2 //Открытый и активный счёт.
	AccountStatus_ACCOUNT_STATUS_CLOSED      AccountStatus = 3 //Закрытый счёт.
	AccountStatus_ACCOUNT_STATUS_ALL         AccountStatus = 4 //Все счета.
)

type GetAccountsResponse struct {
	Accounts []Account
}

type Account struct {
	Id          string
	Name        string
	Type        AccountType
	AccessLevel AccessLevel
	Status      AccountStatus
}
