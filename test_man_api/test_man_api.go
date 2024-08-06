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

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/people", handlePeople)

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

func handlePeople(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPeople(w, r)
	case http.MethodPost:
		createPerson(w, r)
	case http.MethodPut, http.MethodPatch:
		updatePerson(w, r)
	case http.MethodDelete:
		deletePerson(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPeople(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM public.people")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		people = append(people, p)
	}
	json.NewEncoder(w).Encode(people)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("INSERT INTO public.people (name, age) VALUES ($1, $2)", p.Name, p.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("UPDATE public.people SET name=$1, age=$2 WHERE id=$3", p.Name, p.Age, p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	ageStr := r.URL.Query().Get("age")

	///////////////////////// แปลงตัวเลข อายุจาก สตริง ที่มาจาก http เป็นตัวเลข //////////////
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}
	////////////////////////////////////////////////////////////////////////////////

	_, err = db.Exec("DELETE FROM public.people WHERE name=$1 AND age=$2", name, age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
