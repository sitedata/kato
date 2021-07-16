package bcode

// tenant application 11000~11099
var (
	//ErrApplicationNotFound -
	ErrApplicationNotFound = newByMessage(404, 11001, "application not found")
	//ErrApplicationExist -
	ErrApplicationExist = newByMessage(400, 11002, "application already exists")
	//ErrCreateNeedCorrectAppID
	ErrCreateNeedCorrectAppID = newByMessage(404, 11003, "create service needs correct application ID")
	//ErrUpdateNeedCorrectAppID
	ErrUpdateNeedCorrectAppID = newByMessage(404, 11004, "update service needs correct application ID")
	//ErrDeleteDueToBindService
	ErrDeleteDueToBindService = newByMessage(400, 11005, "the application cannot be deleted because there are bound services")

	ErrK8sServiceNameExists = newByMessage(400, 11006, "kubernetes service name already exists")
)

// app config group 11100~11199
var (
	//ErrApplicationConfigGroupExist -
	ErrApplicationConfigGroupExist = newByMessage(409, 11101, "application config group already exist")
	//ErrConfigGroupServiceExist -
	ErrConfigGroupServiceExist = newByMessage(409, 11102, "config group under this service already exists")
	//ErrConfigItemExist -
	ErrConfigItemExist = newByMessage(409, 11103, "config item under this config group already exist")
	//ErrServiceNotFound -
	ErrServiceNotFound = newByMessage(404, 11104, "this service ID cannot be found under this application")
)
