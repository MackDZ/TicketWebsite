package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v81"
)

type Application struct {
	port int

	//db database.Servive
}

func NewServer() *http.Server {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY") //Loads Stripe Secret Key

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewApp := &Application{
		port: port,

		//db: database.New()
	}
	log.Println("Listening on :", port)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewApp.port),
		Handler:      NewApp.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
