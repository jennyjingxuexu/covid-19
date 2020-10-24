package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"covid-19/internal/config"
	"covid-19/internal/database"
	"covid-19/internal/handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// preferrably setup env when spinning up the http server instead of loading on app start
	if err := godotenv.Load(); err != nil {
		panic("Cannot Start App - .env Loading failed")
	}
	env := config.LoadEnv()

	orm := database.GetOrm(env.Db)
	user := handler.NewUserProvider(orm)

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/users", user.CreateUser()).Methods("POST")
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
