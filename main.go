package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func getTemp() float32 {
	var temp float32
	db, err := sql.Open("mysql", "user:password@tcp(mysql:3306)/temp")
	if err != nil {
		log.Fatal(err)
	}

	q := `SELECT temp FROM temp ORDER BY temp.when DESC LIMIT 1`
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&temp); err != nil {
			log.Fatal(err)
		}
	}
	return (temp * 9 / 5) + 32
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("/home.html")
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, fmt.Sprintf("%6.1f", getTemp()))

}

func main() {

	var dir string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.HandleFunc("/", HomeHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
