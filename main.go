package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/cccat20_1/src"
)

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:123456@localhost:5432/app")
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	r.POST("/signup", src.SignupHandler(db))
	r.GET("/accounts/:accountId", src.GetAccountHandler(db))

	r.Run(":3000")
}
