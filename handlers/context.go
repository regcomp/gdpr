package handlers

var STX *ServiceContext

type ServiceContext struct {
	Testing string
	// db, certs, and keys
}

func CreateServiceContext(getenv func(string) string) *ServiceContext {
	// other context setup goes here, like getting certs/keys

	return &ServiceContext{
		Testing: "CONTEXT",
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
