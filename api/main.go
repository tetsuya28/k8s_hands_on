package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"main/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
		log.Println(w, response)
		fmt.Println(w, response)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func main() {
	env := os.Getenv("GO_ENV")

	var err error
	if env == "docker" {
		err = godotenv.Load(".env.docker")
	} else {
		err = godotenv.Load()
	}
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
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("404")
		log.Println(w)
		log.Println(r)
	})
	r.HandleFunc("/api/ping", pingHandler).Methods("GET")

	r.Methods(http.MethodGet).Path("/api/todo").Handler(appHandler{dbConnection(dbx).getTodosHandler})
	r.Methods(http.MethodPost).Path("/api/todo").Handler(appHandler{dbConnection(dbx).postTodosHandler})
	r.Methods(http.MethodDelete).Path("/api/todo/{id}").Handler(appHandler{dbConnection(dbx).deleteTodoHandler})
	r.Methods(http.MethodOptions).Path("/api/todo").Handler(appHandler{dbConnection(dbx).optionsTodosHandler})
	r.Methods(http.MethodPost).Path("/api/todo/{id}/done").Handler(appHandler{dbConnection(dbx).postTodoStatusHandler})

	http.Handle("/", r)
	crosHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Access-Control-Allow-Headers"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
	}).Handler(r)
	log.Fatal(http.ListenAndServe(":8000", crosHandler))
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

func (a dbStore) optionsTodosHandler(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	return http.StatusOK, "ok", nil
}

func (a dbStore) deleteTodoHandler(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	_, err = a.db.Exec(`
DELETE FROM todos WHERE id = ?
	`, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusNoContent, nil, nil
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
