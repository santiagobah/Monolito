package models

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func Conectar_db() {
	host := "db"
	port := 5432
	user := "admin"
	password := "admin123"
	dbname := "minipackettracerdb"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatal("no se pudo establecer la conexión: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("no se pudo conectar a la bd: ", err)
	}

	fmt.Println("Conexión exitosa a la aase de datos")
}
