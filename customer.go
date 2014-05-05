package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/golang/glog"
	"net/http"
	"strconv"
)

type Customer struct {
	Id        int
	FirstName string
	LastName  string
}

type CustomerList struct {
	Customers []*Customer
}

func (e *Customer) Load(tx Transaction) error {
	err := tx.QueryRow("SELECT first_name,last_name FROM customer WHERE id=$1", e.Id).Scan(&e.FirstName, &e.LastName)

	switch {
	case err == sql.ErrNoRows:
		log.Warningln("[CUSTOMER]: No user with that ID.")
	case err != nil:
		log.Warningf("[CUSTOMER: Unknown error reading from database, error: '%v'", err)
	}

	return err
}

func (e *Customer) Delete(tx Transaction) error {
	result, err := tx.Exec("DELETE FROM customer WHERE id = $1", e.Id)

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to delete customer '%v', error: '%v'", e.Id, err)
		return err
	}

	cnt, err := result.RowsAffected()

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to get rows affected count, error: '%v'", err)
		return err
	}

	if cnt != 1 {
		log.Warningf("Deleted an invalid number of customers '%v'", cnt)
		return errors.New(fmt.Sprintf("invalid number of customers deleted, '%v'", cnt))
	}

	log.V(INFO).Infof("Deleted '%v' customers", cnt)

	return nil
}

func (e *Customer) Save(tx Transaction) error {
	result, err := tx.Exec("UPDATE customer SET first_name=$1,last_name=$2 WHERE id = $3", e.FirstName, e.LastName, e.Id)

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to save customer '%v', error: '%v'", e.Id, err)
		return err
	}

	cnt, err := result.RowsAffected()

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to get rows affected count, error: '%v'", err)
		return err
	}

	if cnt != 1 {
		log.Warningf("Updated an invalid number of customers '%v'", cnt)
		return errors.New(fmt.Sprintf("invalid number of customers updated, '%v'", cnt))
	}

	log.V(INFO).Infof("Updated '%v' customers", cnt)

	return nil
}

func (e *Customer) Insert(tx Transaction) error {
	err := tx.QueryRow("INSERT INTO customer (id,first_name,last_name) values (nextval('customer_seq'),$1,$2) returning id", e.FirstName, e.LastName).Scan(&e.Id)

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to save customer '%v', error: '%v'", e.Id, err)
		return err
	}

	return nil
}

func (e *CustomerList) Load(tx Transaction, offset, limit int) error {
	rows, err := tx.Query("SELECT id,first_name,last_name FROM customer offset $1 limit $2", offset, limit)

	if err != nil {
		log.Warningf("[CUSTOMER]: Unknown error reading from database, error: '%v'", err)
		return err
	}

	for rows.Next() {
		customer := Customer{}
		if err := rows.Scan(&customer.Id, &customer.FirstName, &customer.LastName); err != nil {
			log.Warningf("[CUSTOMER]: error reading data from database, error: ", err)
		}
		e.Customers = append(e.Customers, &customer)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func CustomerSaveHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	e := &Customer{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(e); err != nil {
		log.Warningf("[CUSTOMER]: Unable to read json, error: '%v'", err)
		http.Error(w, "Bad json, "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Disregard any id that may be in the payload
	e.Id = id

	if err := e.Save(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Warningf("[CUSTOMER]: Unable to find customer with id '%v', error: '%v'", id, err)
			http.Error(w, "Customer not found", http.StatusNotFound)
		default:
			log.Warningf("[CUSTOMER]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	return nil
}

func CustomerCreateHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	e := &Customer{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(e); err != nil {
		log.Warningf("[CUSTOMER]: Unable to read json, error: '%v'", err)
		return err
	}

	if err := e.Insert(tx); err != nil {
		log.Warningf("[CUSTOMER]: Unable to create customer, error: '%v'", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func CustomerHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	log.V(INFO).Infof("Customer id '%v'", id)

	e := &Customer{
		Id: id,
	}

	if err := e.Load(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Warningf("[CUSTOMER]: Unable to find customer with id '%v', error: '%v'", id, err)
			http.Error(w, "Customer not found", http.StatusNotFound)
		default:
			log.Warningf("[CUSTOMER]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	b, err := json.Marshal(e)

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	return nil
}

func CustomerDeleteHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	e := &Customer{
		Id: id,
	}

	if err := e.Delete(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Warningf("[CUSTOMER]: Unable to find customer with id '%v', error: '%v'", id, err)
			http.Error(w, "Customer not found", http.StatusNotFound)
		default:
			log.Warningf("[CUSTOMER]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	return nil
}

func CustomerListHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	offset := getInt(r, "offset", 0)
	limit := getInt(r, "limit", 10)

	log.V(INFO).Infof("Offset = %v, Limit = %v", offset, limit)

	e := &CustomerList{}

	if err := e.Load(tx, offset, limit); err != nil {
		log.Warningf("[CUSTOMER]: Unable to load data from database, error: '%v'", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	b, err := json.Marshal(e)

	if err != nil {
		log.Warningf("[CUSTOMER]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	return nil
}
