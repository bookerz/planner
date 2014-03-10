package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Employee struct {
	Id        int
	FirstName string
	LastName  string
}

func (e *Employee) Load(tx *sql.Tx) error {
	err := tx.QueryRow("SELECT first_name,last_name FROM employee WHERE id=$1", e.Id).Scan(&e.FirstName, &e.LastName)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("[EMPLOYEE]: No user with that ID.")
	case err != nil:
		log.Printf("[EMPLOYEE]: Unknown error reading from database: '%v'", err)
	}

	return err
}

func EmployeeHandler(w http.ResponseWriter, r *http.Request, tx *sql.Tx, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Printf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	e := &Employee{
		Id: id,
	}

	if err = e.Load(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Printf("[EMPLOYEE]: Unable to to find employee with id '%v', error: '%v'", id, err)
			http.Error(w, "Employee not found", http.StatusNotFound)
		default:
			log.Printf("[EMPLOYEE]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	b, err := json.Marshal(e)

	if err != nil {
		log.Printf("[EMPLOYEE]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	return nil
}
