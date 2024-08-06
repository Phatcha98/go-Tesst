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

type Machine struct {
	ID           int    `json:"id"`
	Line         string `json:"line"`
	Module       int    `json:"module"`
	Machine_name string `json:"machine_name"`
	Machine      string `json:"machine"`
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/machine", handleMachine)

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

func handleMachine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMachine(w, r)
	case http.MethodPost:
		createMachine(w, r)
	case http.MethodPut, http.MethodPatch:
		updateMachine(w, r)
	case http.MethodDelete:
		deleteMachine(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getMachine(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, line, module, machine_name, machine FROM public.machine_list")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var Machines []Machine
	for rows.Next() {
		var m Machine
		err := rows.Scan(&m.ID, &m.Line, &m.Module, &m.Machine_name, &m.Machine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Machines = append(Machines, m)
	}
	json.NewEncoder(w).Encode(Machines)
}

func createMachine(w http.ResponseWriter, r *http.Request) {
	var m Machine
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("INSERT INTO public.machine_list (line, module, machine_name, machine) VALUES ($1, $2, $3, $4)", m.Line, m.Module, m.Machine_name, m.Machine)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateMachine(w http.ResponseWriter, r *http.Request) {
	var m Machine
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("UPDATE public.machine_list SET line=$1, module=$2, machine_name=$3, machine=$4 WHERE id=$5", m.Line, m.Module, m.Machine_name, m.Machine, m.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteMachine(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("line")
	strmodule := r.URL.Query().Get("module")

	module, err := strconv.Atoi(strmodule)
	if err != nil {
		http.Error(w, "Invalid module", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM public.machine_list WHERE line=$1 AND module=$2", line, module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
