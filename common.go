package main

import (
	log "github.com/golang/glog"
	"net/http"
	"strconv"
)

func getInt(r *http.Request, key string, defaultVal int) int {
	v := r.FormValue(key)
	if v == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(v)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		return defaultVal
	}

	return val
}
