package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func Saludar(c *gin.Context) {
	var p Person
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	message := fmt.Sprintf("Hola %v %v", p.Name, p.Lastname)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func main() {
	router := gin.Default()
	// Ejercicio 1
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})
	// Ejercicio 2
	router.POST("/saludo", Saludar)
	router.Run(":8080")
}
