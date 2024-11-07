package main

import (
	"fmt"
	"github.com/Nelwhix/numeris/handlers"
	"github.com/Nelwhix/numeris/pkg"
	"github.com/Nelwhix/numeris/pkg/middlewares"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/go-playground/validator/v10"
	gHandlers "github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ServerPort = ":8080"
)

var validate *validator.Validate

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	logFile := filepath.Join("logs", "numeris.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := pkg.CreateNewLogger(f)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	conn, err := pkg.CreateDbConn()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	model := models.Model{
		Conn: conn,
	}

	handler := handlers.Handler{
		Model:     model,
		Logger:    logger,
		Validator: validate,
	}

	m := middlewares.AuthMiddleware{
		Model: model,
	}

	r := http.NewServeMux()

	// Guest Routes
	r.HandleFunc("POST /api/v1/auth/signup", handler.SignUp)
	r.HandleFunc("POST /api/v1/auth/login", handler.Login)

	// Auth routes
	r.Handle("GET /api/v1/invoices/widgets", m.Register(handler.GetInvoiceWidgetsData))
	//r.Handle("POST /api/games", m.Register(handler.CreateNewGame))

	fmt.Printf("Numeris started at http://localhost:%s\n", ServerPort)

	err = http.ListenAndServe(ServerPort, gHandlers.CombinedLoggingHandler(os.Stdout, r))
	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
