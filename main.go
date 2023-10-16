package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)



func main() {
	err := godotenv.Load()
	if err!=nil{
		log.Println("Error in loading env",err)
	}
	dbConfig := DBConfig{
		Host: os.Getenv("HOST"),
		User: os.Getenv("USER"),
		Password: os.Getenv("PASS"),
		DBname: os.Getenv("DBNAME"),
		Port: os.Getenv("PORT"),
	}
	db,err:=InitializeDB(&dbConfig)
	if err!=nil{
		panic(err)
	}

    handler :=NewHandler(db)

	http.HandleFunc("/add",handler.addTask)
	http.HandleFunc("/edit",handler.editTask)
	http.HandleFunc("/delete",handler.deleteTask)
	http.HandleFunc("/list",handler.listTasks)

	log.Fatal(http.ListenAndServe(":8000",nil))
}