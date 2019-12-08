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

func getTemp() (float32, int) {
	var temp float32
	var heating int
	db, err := sql.Open("mysql", "user:passwd@tcp(mysql:3306)/temp")
	if err != nil {
		log.Fatal(err)
	}

	q := `SELECT temp, heat FROM temp ORDER BY temp.when DESC LIMIT 1`
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&temp, &heating); err != nil {
			log.Fatal(err)
		}
	}
	return (temp * 9 / 5) + 32, heating
}

// HomeHandler renders the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var heatStr string
	tmpl, err := template.ParseFiles("/home.html")
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	temp, heating := getTemp()

	if heating == 1 {
		heatStr = "&#128293;"
	}
	err = tmpl.Execute(w, fmt.Sprintf(" %s %6.1f", heatStr, temp))

}

func main() {
	var dir string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
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
