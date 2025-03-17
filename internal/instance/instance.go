package instance

import "mailer/internal/storage"

type Instance struct {
	Storage *storage.Storage
}

func CreateInstance() *Instance {
	st := storage.CreateStorage()
	st.MigrateDatabase()

	return &Instance{
		Storage: st,
	}
}
