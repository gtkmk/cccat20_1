package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gtkmk/cccat20_1/src"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbURL = "postgres://postgres:123456@postgres:5432/app?sslmode=disable"

type DBConnection struct {
	connection *sqlx.DB
}

func (db *DBConnection) Rows(query string, values ...any) ([]map[string]interface{}, error) {
	rows, err := db.connection.Queryx(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			continue
		}
		results = append(results, row)
	}

	return results, nil
}

func (db *DBConnection) Raw(query string, values ...any) error {
	_, err := db.connection.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func generateUUID() string {
	return uuid.New().String()
}

func isValidName(name string) bool {
	return regexp.MustCompile(`[a-zA-Z] [a-zA-Z]+`).MatchString(name)
}

func isValidEmail(email string) bool {
	return regexp.MustCompile(`^(.+)@(.+)$`).MatchString(email)
}

func isValidCarPlate(carPlate string) bool {
	return regexp.MustCompile(`[A-Z]{3}[0-9]{4}`).MatchString(carPlate)
}

func SignUpHandler(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		CPF         string `json:"cpf"`
		CarPlate    string `json:"carPlate"`
		IsPassenger bool   `json:"isPassenger"`
		IsDriver    bool   `json:"isDriver"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbConn := &DBConnection{connection: db}

	var exists bool
	rows, err := dbConn.Rows("SELECT EXISTS (SELECT 1 FROM ccca.account WHERE email = $1)", input.Email)
	if err != nil {
		log.Fatal(err)
	}

	if len(rows) > 0 {
		exists = rows[0]["exists"].(bool)
	}

	if exists {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -4})
		return
	}

	if !isValidName(input.Name) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -3})
		return
	}
	if !isValidEmail(input.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -2})
		return
	}
	if !src.ValidatePassword(input.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -5})
		return
	}
	if !src.ValidateCpf(input.CPF) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -1})
		return
	}
	if input.IsDriver && !isValidCarPlate(input.CarPlate) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": -6})
		return
	}

	id := generateUUID()
	err = dbConn.Raw("INSERT INTO ccca.account (account_id, name, email, cpf, car_plate, is_passenger, is_driver, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", id, input.Name, input.Email, input.CPF, input.CarPlate, input.IsPassenger, input.IsDriver, input.Password)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"accountId": id})
}

func GetAccountHandler(c *gin.Context) {
	accountId := c.Param("accountId")

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbConn := &DBConnection{connection: db}

	accounts, err := dbConn.Rows("SELECT * FROM ccca.account WHERE account_id = $1", accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(accounts) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, accounts[0])
}

func main() {
	r := gin.Default()
	r.POST("/signup", SignUpHandler)
	r.GET("/accounts/:accountId", GetAccountHandler)

	r.Run(":8089")
}
