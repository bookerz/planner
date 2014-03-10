package main

import (
	"database/sql"
	"log"
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
