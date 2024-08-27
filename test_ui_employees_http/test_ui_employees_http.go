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
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "postgres"
)

type Employee struct {
    Country string `json:"country"`
    Wage    int    `json:"wage"`
}

func submitForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    name := r.FormValue("name")
    ageStr := r.FormValue("age")
    country := r.FormValue("country")
    position := r.FormValue("position")
    wageStr := r.FormValue("wage")

    log.Printf("Received form data: name=%s, age=%s, country=%s, position=%s, wage=%s", name, ageStr, country, position, wageStr)

    age, err := strconv.Atoi(ageStr)
    if err != nil {
        log.Printf("Error converting age to int: %v", err)
        http.Error(w, "Invalid age value", http.StatusBadRequest)
        return
    }

    wage, err := strconv.Atoi(wageStr)
    if err != nil {
        log.Printf("Error converting wage to int: %v", err)
        http.Error(w, "Invalid wage value", http.StatusBadRequest)
        return
    }

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Printf("Error connecting to database: %v", err)
        http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Printf("Error pinging database: %v", err)
        http.Error(w, "Unable to reach database", http.StatusInternalServerError)
        return
    }

    // Check if the data already exists
    var exists bool
    checkQuery := `
        SELECT EXISTS(
            SELECT 1 FROM public.employees WHERE name=$1 AND age=$2 AND country=$3 AND position=$4 AND wage=$5
        )`
    err = db.QueryRow(checkQuery, name, age, country, position, wage).Scan(&exists)
    if err != nil {
        log.Printf("Error checking existing data: %v", err)
        http.Error(w, "Unable to check existing data", http.StatusInternalServerError)
        return
    }

    if exists {
        log.Println("Data already exists")
        fmt.Fprintln(w, "Error: have data already")
        return
    }

    // Insert the new employee data
    sqlStatement := `
        INSERT INTO public.employees (name, age, country, position, wage)
        VALUES ($1, $2, $3, $4, $5)`
    _, err = db.Exec(sqlStatement, name, age, country, position, wage)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Unable to execute query", http.StatusInternalServerError)
        return
    }

    log.Println("Data inserted successfully")
    fmt.Fprintf(w, "Employee added successfully")
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Printf("Error connecting to database: %v", err)
        http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT country, wage FROM public.employees")
    if err != nil {
        log.Printf("Error querying database: %v", err)
        http.Error(w, "Unable to query database", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var employees []Employee
    for rows.Next() {
        var employee Employee
        if err := rows.Scan(&employee.Country, &employee.Wage); err != nil {
            log.Printf("Error scanning row: %v", err)
            http.Error(w, "Unable to scan row", http.StatusInternalServerError)
            return
        }
        employees = append(employees, employee)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(employees)
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })

    http.HandleFunc("/submit", submitForm)
    http.HandleFunc("/api/employees", getEmployees)

    fmt.Println("Server is running on http://10.17.77.190:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}