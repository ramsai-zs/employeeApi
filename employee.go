package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

func connect() *sql.DB {
	var err error
	db, err = sql.Open("mysql", "sai:password@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Println(err)
	}
	return db
}

type Employee struct {
	ID      int
	Name    string
	Address string
}

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		post(w, r)
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		deleteId(w, r)
		w.WriteHeader(http.StatusOK)
	case http.MethodPut:
		put(w, r)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func deleteId(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("ID")
	_, err := db.Query("delete from employee where ID = ?", query)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func get(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select * from employee;")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer rows.Close()
	var emp []Employee
	for rows.Next() {
		var e Employee
		err = rows.Scan(&e.ID, &e.Name, &e.Address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		emp = append(emp, e)
	}

	query := r.URL.Query().Get("ID")
	if query != "" {
		for k, v := range emp {
			if strconv.Itoa(v.ID) == query {
				res, err := json.Marshal(emp[k])
				if err != nil {
					log.Fatalln(err)
				}
				_, err = w.Write(res)
				if err != nil {
					log.Println(err)
				}

				return
			}
		}
	}

	res, err := json.Marshal(emp)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(res)
}

func post(w http.ResponseWriter, r *http.Request) {
	var e Employee
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &e)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec("INSERT INTO employee (id,name,city)values (?,?,?)", e.ID, e.Name, e.Address)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Wrong"))
	}
	w.Write([]byte("Success"))
}

func put(w http.ResponseWriter, r *http.Request) {
	var e Employee
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &e)
	_, err := db.Exec("UPDATE employee SET name = ? where Id=?", e.Name, e.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Success"))
}

func main() {
	defer db.Close()
	http.HandleFunc("/employee", router)
	log.Fatalln(http.ListenAndServe(":3036", nil))
}
