package controllers

import (
	"fmt"
	"net/http"
)

// Get Health Code 200
func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) 
	fmt.Fprintf(w, "Requête réussie, code 200!")
}

