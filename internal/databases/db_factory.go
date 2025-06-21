package databases

import (
	"os"

	"github.com/alejandro-albiol/athenai/internal/databases/interfaces"
	"github.com/alejandro-albiol/athenai/internal/databases/services"
)

func NewDBService() interfaces.DBService {
    dbType := os.Getenv("DB_TYPE")
    dsn := os.Getenv("DB_DSN")
    switch dbType {
    case "mysql":
        return services.NewMySQLService(dsn)
    case "postgres":
        return services.NewPostgresService(dsn)
    default:
        panic("Unsupported DB_TYPE")
    }
}