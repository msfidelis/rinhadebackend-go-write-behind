package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
)

var onceDB sync.Once
var dbInstance *sql.DB

func GetDB() *sql.DB {
	onceDB.Do(func() {
		var err error

		user := os.Getenv("DATABASE_USER")
		pass := os.Getenv("DATABASE_PASSWORD")
		host := os.Getenv("DATABASE_HOST")
		port := os.Getenv("DATABASE_PORT")
		schema := os.Getenv("DATABASE_DB")

		connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, schema)
		dbInstance, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
		}

		// Verifica a conexão
		err = dbInstance.Ping()
		if err != nil {
			log.Fatalf("Erro ao estabelecer uma conexão com o banco de dados: %v", err)
		}
	})
	return dbInstance
}
