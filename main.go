package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/aledevpro/interactive-presentation/api"
	db "github.com/aledevpro/interactive-presentation/db/sqlc"
	"github.com/aledevpro/interactive-presentation/util"
)


func main() {
	config:=util.Config{
		Environment:  "production",
		DBDriver:     "postgres",
		DBSource:     "postgresql://root:XM5IPpSNjRggVpS22Tp2@interactive-presentations.chihgxkzxezm.eu-north-1.rds.amazonaws.com:5432/interactive_presentations",
		MigrationURL: "file://db/migration",
	}
	conn, err := sql.Open(config.DBDriver,  config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()


	m, err := migrate.New(
		config.MigrationURL,
		config.DBSource, 
	)
	if err != nil {
		log.Fatal("migration error --> ", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply.")
		} else {
			log.Fatal(err)
		}
	}

	log.Println("Migrations applied successfully.")
	store := db.NewStore(conn)

	server, err := api.NewServer(store)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = server.Start(":8080")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}
