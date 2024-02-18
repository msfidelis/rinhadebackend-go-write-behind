package routines

import (
	"fmt"
	"log"
	"main/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver do banco de dados
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Driver de arquivos
)

func DatabaseMigration() {
	db := utils.GetDB()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Certifique-se de que este caminho está correto e acessível
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations aplicadas com sucesso")
}
