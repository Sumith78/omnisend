package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Define a sample data structure for orders
type Order struct {
	ID             int    `json:"id"`
	OrderNumber    string `json:"order_number"`
	ShipmentStatus string `json:"shipment_status"`
	Email          string `json:"email"`
}

var orders []Order

func main() {
	// Initialize sample orders (replace this with your database retrieval logic)
	orders = []Order{
		{ID: 1, OrderNumber: "12345", ShipmentStatus: "In transit", Email: "sample@email.com"},
		// Add more orders as needed
	}

	// Set up API routes
	http.HandleFunc("/api/orders", getOrders)
	http.HandleFunc("/api/update-status", updateShipmentStatus)

	// Start the server
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	// Return the list of orders as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func updateShipmentStatus(w http.ResponseWriter, r *http.Request) {
	// Extract data from request (you may want to use a POST request for this)
	orderID := r.URL.Query().Get("order_id")
	newStatus := r.URL.Query().Get("new_status")

	// Find and update the shipment status
	for i, order := range orders {
		if order.OrderNumber == orderID {
			orders[i].ShipmentStatus = newStatus
			// Notify shipment status change (you need to implement this)
			NotifyShipmentStatusChange(order.Email, orderID)
			break
		}
	}

	// Respond with a success message (you may want to handle errors)
	w.Write([]byte("Shipment status updated successfully"))
}

func NotifyShipmentStatusChange(email string, orderID string) {
	from := mail.NewEmail("Your Name", "your@example.com") // Replace with your email and name
	subject := "Shipment Status Update"
	to := mail.NewEmail("", email) // Use the provided email address
	plainTextContent := fmt.Sprintf("Your order with ID %s has a new shipment status.", orderID)
	htmlContent := fmt.Sprintf("<strong>Your order with ID %s has a new shipment status.</strong>", orderID)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")) // Set your SendGrid API key as an environment variable

	response, err := client.Send(message)
	if err != nil {
		log.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully. Response:", response.StatusCode)
}