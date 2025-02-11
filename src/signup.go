package src

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SignupInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Cpf         string `json:"cpf"`
	Password    string `json:"password"`
	IsDriver    bool   `json:"isDriver"`
	IsPassenger bool   `json:"isPassenger"`
	CarPlate    string `json:"carPlate"`
}

func SignupHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		ctx := context.Background()
		var exists bool

		err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM ccca.account WHERE email = $1)", input.Email).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			return
		}
		if exists {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -4}) // Já existe
			return
		}

		if len(input.Name) < 3 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -3}) // Nome inválido
			return
		}
		if len(input.Email) < 5 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -2}) // Email inválido
			return
		}
		if len(input.Password) < 8 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -5}) // Senha inválida
			return
		}
		if input.IsDriver && len(input.CarPlate) != 7 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -6}) // Placa inválida
			return
		}

		id := uuid.New().String()
		_, err = db.ExecContext(ctx, "INSERT INTO ccca.account (account_id, name, email, cpf, car_plate, is_passenger, is_driver, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			id, input.Name, input.Email, input.Cpf, input.CarPlate, input.IsPassenger, input.IsDriver, input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"accountId": id})
	}
}

func GetAccountHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountId := c.Param("accountId")
		var account SignupInput
		ctx := context.Background()

		err := db.QueryRowContext(ctx, "SELECT name, email, cpf, car_plate, is_passenger, is_driver, password FROM ccca.account WHERE account_id = $1", accountId).
			Scan(&account.Name, &account.Email, &account.Cpf, &account.CarPlate, &account.IsPassenger, &account.IsDriver, &account.Password)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
			return
		}
		c.JSON(http.StatusOK, account)
	}
}
