package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var products []Product
var lastID int
var productByCode = make(map[string]Product)

func loadJSONFile(path string) error {
	var (
		err       error
		file      *os.File
		byteValue []byte
	)
	file, err = os.Open(path)
	if err != nil {
		return errors.New("File not found")
	}
	defer file.Close()
	byteValue, err = io.ReadAll(file)
	if err != nil {
		return errors.New("Cannot read file")
	}
	err = json.Unmarshal(byteValue, &products)
	if err == nil {
		lastID = products[len(products)-1].ID
		for _, product := range products {
			productByCode[product.CodeValue] = product
		}
	}
	if err != nil {
		return errors.New("Cannot load json")
	}
	return nil
}

func GetAllProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	for _, product := range products {
		idProduct, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if product.ID == idProduct {
			c.JSON(http.StatusOK, gin.H{"product": product})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func GetProduct(c *gin.Context) {
	price := c.Query("price")
	if price == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Price is required"})
		return
	}
	priceProduct, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var filteredProducts []Product
	for _, v := range products {
		if v.Price > priceProduct {
			filteredProducts = append(filteredProducts, v)
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"products": filteredProducts})
}

func validateProduct(product Product) bool {
	isCorrect := true
	if product.Name == "" {
		isCorrect = false
	}
	if product.Quantity == 0 {
		isCorrect = false
	}
	if product.CodeValue == "" {
		isCorrect = false
	}
	if product.Expiration == "" {
		isCorrect = false
	}
	if product.Price == float64(0) {
		isCorrect = false
	}
	return isCorrect
}

func CreateProduct(c *gin.Context) {
	var newProduct Product
	newProduct.IsPublished = false
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate product fields
	if !validateProduct(newProduct) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is missing required values"})
		return
	}

	// Validate code value
	if _, codeValueExists := productByCode[newProduct.CodeValue]; codeValueExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code value already exists"})
		return
	}

	// Validate expiration date
	_, err := time.Parse("02/01/2006", newProduct.Expiration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration date format"})
		return
	}

	lastID++
	newProduct.ID = lastID
	products = append(products, newProduct)
	productByCode[newProduct.CodeValue] = newProduct
	c.JSON(http.StatusCreated, newProduct)
}

func main() {
	var err error
	err = loadJSONFile("./exercise1/products.json")
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.GET("/products", GetAllProducts)
	router.GET("/products/:id", GetProductByID)
	router.GET("/products/search", GetProduct)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/products", CreateProduct)
	router.Run(":8080")
}
