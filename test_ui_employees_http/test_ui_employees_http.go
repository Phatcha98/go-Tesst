package main

import (
	"fmt"
	"net/http"
)

func submitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	age := r.FormValue("age")
	country := r.FormValue("country")
	position := r.FormValue("position")
	wage := r.FormValue("wage")

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %s\n", age)
	fmt.Printf("Country: %s\n", country)
	fmt.Printf("Position: %s\n", position)
	fmt.Printf("Wage: %s\n", wage)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/submit", submitForm)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
