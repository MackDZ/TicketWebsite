package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var Mytemplates = template.Must(template.ParseGlob("templates/*"))

func (app *Application) RegisterRoutes() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/", func(r chi.Router) {
		r.Get("/", app.getWebpage("Home Page"))
		r.Get("/ticket", app.getWebpage("Ticket Page"))
		r.Get("/ticket/order#confirmed", nil)
		r.Get("/public_keys", publicKeyHandler)

		r.Post("/create-checkout-session", createCheckoutSession)
	})

	r.Route("/order", func(r chi.Router) {
		r.Get("/success", orderSuccess)
	})

	return r

}

func (app *Application) getWebpage(name string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host)

		Mytemplates.ExecuteTemplate(w, name, nil)
	}
}
