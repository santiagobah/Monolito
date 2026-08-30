package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// aqui guardamos la conexion global a la bd, ya se que no es lo mas correcto
// pero para el proyecto asi nos dijeron que estaba bien (monolito jaja)
var DB *sql.DB

func ConectarDB() {
	// TODO: sacar esto a variables de entorno bien, por ahora hardcodeado
	host := "db"
	port := 5432
	user := "admin"
	password := "admin123"
	dbname := "redsim"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("no se pudo abrir la conexion: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("no se pudo conectar a la bd: ", err)
	}

	fmt.Println("conectado a la base de datos ez")
}
