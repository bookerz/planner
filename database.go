package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Must adapt to HttpRouter library?
func RunInTransaction(f func(w http.ResponseWriter, r *http.Request, tx *sql.Tx, vars map[string]string) error) func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) {

		tx, err := db.Begin()

		if err != nil {
			log.Printf("[TRANSACTION]: Unable to create database transaction. error: '%v'", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = f(w, r, tx, vars)

		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("[TRANSACTION]: Unable to rollback transaction, error: '%v'", err)
			}
			log.Printf("[TRANSACTION]: Rolling back transaction, error: '%v'", err)
			return
		}

		err = tx.Commit()

		if err != nil {
			log.Printf("[TRANSACTION]: Unable to commit transaction, error: '%v'", err)
		}
	}
}
