package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Person struct{
	gorm.Model

	Name string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book

}

type Book struct{
	gorm.Model

	Title string
	Author string
	CallNumber int32 `gorm:"unique_index"`
	PersonID int
}

var dbX *gorm.DB
var dbErr error

var(
	person = &Person{Name:"Yasin", Email: "yasinkhan4008@gmail.com"}
	books = []Book{
		{Title: "Book1", Author: "Yasin", CallNumber: 1, PersonID: 1},
		{Title: "Book2", Author: "Shejan", CallNumber: 2, PersonID: 1},
	}
)

func main(){
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// dialect:= os.Getenv("DIALECT")
	host:= os.Getenv("HOST")
	dbport:= os.Getenv("DBPORT")
	user:= os.Getenv("DBUSER")
	dbName:= os.Getenv("NAME")
	dbPassword:= os.Getenv("PASSWORD")

	connectionStr:= fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, dbPassword, dbport)

	db, dbErr := gorm.Open(postgres.Open(connectionStr),&gorm.Config{})
	if(dbErr != nil){
		fmt.Println("error connecting to DB")
		return
	}

	dbObj, bdBbjErr:= db.DB()
	if(bdBbjErr != nil){
		fmt.Println("Experienced some internal problem")
		return
	}

	dbX = db

	fmt.Println("connected to DB")

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	defer dbObj.Close()

	router := mux.NewRouter()

	router.HandleFunc("/people", getPeople).Methods("GET")

	http.ListenAndServe(os.Getenv("PORT"), router)

}

func getPeople(head http.ResponseWriter, req *http.Request){
	var people []Person

	dbX.Find(&people)

	json.NewEncoder(head).Encode(&people)

}