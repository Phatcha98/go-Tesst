package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	////////////////////////////////// เชื่อมต่อกับฐานข้อมูล PostgreSQL /////////////////////////////////////
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// เช็คการเชื่อมต่อ
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	/////////////////////////////////// กำหนดตัวแปล และรับข้อมูลinput จาก user /////////////////////////////
	var name string
	var age int
	fmt.Println("Enter name:")
	fmt.Scanf("%s ", &name)
	fmt.Println("Enter age:")
	fmt.Scanf("%d", &age)
	////////////////////////////////  เพิ่มข้อมูลเข้า sql ///////////////////////////////////////////////////
	_, err = db.Exec("INSERT INTO people(name, age) VALUES($1, $2)", name, age)
	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted data into the database successfully!")
}
