package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type myCustomer struct {
	Band       string
	NumTickets int
	Email      string
	Name       string
}

type PublicKey struct {
	StripeKey string `json:"key"`
}

// Loads stipe public key for use
func publicKeyHandler(w http.ResponseWriter, r *http.Request) {

	var public_keys = PublicKey{
		StripeKey: os.Getenv("STRIPE_PUBLIC_KEY"),
	}

	by, err := json.Marshal(public_keys)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(by))
	// http.Redirect(w, r, )
}

/*
Send user to stipe hosted website to checkout

	Potential TODO: make custom checkout session on this website
			- TLS and https is needed for this
*/
func createCheckoutSession(w http.ResponseWriter, r *http.Request) {

	domain := "http://" + r.Host

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			Price: stripe.String(os.Getenv("STRIPE_TICKET_PRICE_ID")),
			AdjustableQuantity: &stripe.CheckoutSessionLineItemAdjustableQuantityParams{
				Enabled: stripe.Bool(true),
				Minimum: stripe.Int64(1),
				Maximum: stripe.Int64(10),
			},
			Quantity: stripe.Int64(1),
		},
		},

		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/order/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(domain + "/order/cancel"),
	}

	s, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
	}

	http.Redirect(w, r, s.URL, http.StatusSeeOther)

}

// After a successfull checkout would retrieve email and number of tickets purchased
// to store in database for onsite ticket varification
func orderSuccess(w http.ResponseWriter, r *http.Request) {

	session_id := r.URL.Query().Get("session_id")

	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	session, err := session.Get(session_id, params)
	if err != nil {
		log.Println(err)
	}

	newCustomer := &myCustomer{
		// Band:  r.FormValue("band")
		NumTickets: int(session.LineItems.Data[0].Quantity),
		Email:      session.CustomerEmail,
		Name:       session.CustomerDetails.Name,
	}

	//TODO Db.submit(newCustomer)
	log.Printf("Order quantity:")

	w.Write([]byte("<html><body><h1>Thanks for your order, " +
		newCustomer.Name + "for purchasing " + strconv.Itoa(newCustomer.NumTickets) +
		" tickets!</h1></body></html>"))
}

// func printJson(v any) {

// 	log.Println(v)

// 	by, err := json.Marshal(v)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println(string(by), "\n")

// }
