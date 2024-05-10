package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DBUsername = "reak"
	DBPassword = "reak"
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "demo"
)

type FormData struct {
	Name  string
	Email string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	name := r.FormValue("name")
	email := r.FormValue("email")

	// Save form data to database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS forms (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), email VARCHAR(255))")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO forms (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to a thank you page
	http.Redirect(w, r, "/thankyou", http.StatusSeeOther)
}

func thankyouHandler(w http.ResponseWriter, r *http.Request) {
	// Render a thank you message
	w.Write([]byte("Thank you for submitting the form!"))
}

func main() {
	http.HandleFunc("/submit", formHandler)
	http.HandleFunc("/thankyou", thankyouHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
