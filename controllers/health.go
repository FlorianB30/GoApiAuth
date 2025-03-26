package controllers

import (
	"auth-api-go/config"
	"auth-api-go/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get Health Code 200
func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) 
	fmt.Fprintf(w, "Requête réussie, code 200!")
}

