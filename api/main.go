package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"main/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var (
	dbx *sqlx.DB
)

type dbStore struct {
	db *sqlx.DB
}

func dbConnection(db *sqlx.DB) *dbStore {
	return &dbStore{db: db}
}

type appHandler struct {
	h func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (a appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	status, res, err := a.h(w, r)
	if err != nil {
		respondJSON(w, status, err)
	}
	respondJSON(w, status, res)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file. %s", err)
	}

	dbx, err := dbClient()
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
		return
	}

	defer dbx.Close()

	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler).Methods("GET")

	r.Methods(http.MethodPost).Path("/todo").Handler(appHandler{dbConnection(dbx).postTodosHandler})
	r.Methods(http.MethodGet).Path("/todo").Handler(appHandler{dbConnection(dbx).getTodosHandler})
	r.Methods(http.MethodPost).Path("/todo/{id}/done").Handler(appHandler{dbConnection(dbx).postTodoStatusHandler})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func dbClient() (*sqlx.DB, error) {
	datasource := os.Getenv("DATASOURCE")
	if datasource == "" {
		log.Fatal("Cannot get datasource for database.")
	}

	return sqlx.Open("mysql", datasource)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}

func (a dbStore) postTodosHandler(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	todo := &model.Todo{}
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return http.StatusBadRequest, nil, err
	}

	_, err := a.db.Exec(`
INSERT INTO todos (name, is_done) VALUES (?, ?)
	`, todo.Name, todo.IsDone)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, todo, nil
}

func (a dbStore) getTodosHandler(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	todos := []model.Todo{}
	err := a.db.Select(&todos, `
SELECT * FROM todos
	`)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &todos, nil
}

func (a dbStore) postTodoStatusHandler(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	todo := model.Todo{}
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	err = a.db.Get(&todo, `
SELECT * FROM todos WHERE id = ?
	`, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, err = a.db.Exec(`
UPDATE todos SET is_done = ? WHERE id = ?
	`, !todo.IsDone, todo.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &todo, nil
}
