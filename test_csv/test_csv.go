package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	////////////////////////////////////////// เชื่อมต่อกับฐานข้อมูล PostgreSQL ////////////////////////////////////////////////////
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	/////////////////////////////////////////  กำหนด path ของไฟล์ CSV //////////////////////////////////////////////////////////
	filePath := "D:\\test-go\\test_csv\\1.csv"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//////////////////////////////////////  สร้าง reader โดยไม่คำนึงถึง header ///////////////////////////////////////////////////////
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	////////////////////////////////////  อ่านและข้าม header row แรกของไฟล์ CSV ////////////////////////////////////////////////////
	if _, err := reader.Read(); err != nil {
		panic(err)
	}
	///////////////////////////////////  อ่านและเพิ่มข้อมูลลงในฐานข้อมูล SQL //////////////////////////////////////////////////////////
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		line := record[0]
		module := record[1]
		machineName := record[2]
		machine := record[3]

		_, err = db.Exec("INSERT INTO public.machine_list(line, \"module\", machine_name, machine) VALUES($1, $2, $3, $4)", line, module, machineName, machine)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Inserted data into the database successfully!")

}
