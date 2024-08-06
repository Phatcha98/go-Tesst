package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	DBHost     = "localhost"
	DBPort     = 5432
	DBUser     = "postgres"
	DBPassword = "postgres"
	DBName     = "postgres"
)

var db *sql.DB

type Employees struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Position string `json:"position"`
	Wage     int    `json:"wage"`
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/employees", handleEmployees)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func handleEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEmployees(w, r)
	case http.MethodPost:
		createEmployees(w, r)
	case http.MethodPut, http.MethodPatch:
		updateEmployees(w, r)
	case http.MethodDelete:
		deleteEmployees(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getEmployees(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, name, age, country, position, wage FROM public.employees")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var Employeess []Employees
	for rows.Next() {
		var e Employees
		err := rows.Scan(&e.ID, &e.Name, &e.Age, &e.Country, &e.Position, &e.Wage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Employeess = append(Employeess, e)
	}
	json.NewEncoder(w).Encode(Employeess)
}

func createEmployees(w http.ResponseWriter, r *http.Request) {
	var e Employees
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("INSERT INTO public.employees(name, age, country, position, wage) VALUES ($1, $2, $3, $4, $5)", e.Name, e.Age, e.Country, e.Position, e.Wage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateEmployees(w http.ResponseWriter, r *http.Request) {
	var e Employees
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("UPDATE public.employees SET name=$1, age=$2, country=$3, position=$4, wage=$5 WHERE id=$6", e.Name, e.Age, e.Country, e.Position, e.Wage, e.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteEmployees(w http.ResponseWriter, r *http.Request) {
	Name := r.URL.Query().Get("name")
	strage := r.URL.Query().Get("age")

	age, err := strconv.Atoi(strage)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM public.employees WHERE name=$1 AND age=$2", Name, age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
