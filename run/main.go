package main

import (
	"fmt"
	"net/http"

	"github.com/sistemakreditasi/backend-akreditasi/routes"
)

func main() {
	// Define the routes and start the server
	http.HandleFunc("/", routes.URL)
	port := ":8080"
	fmt.Println("Server started at: http://localhost" + port)
	http.ListenAndServe(port, nil)
}
