package main

import (
	// --- THIS PART IS CHANGED ---
	"backend/blockchain" // Was "agrochain/backend/blockchain"
	"backend/handlers"   // Was "agrochain/backend/handlers"
	// --- END OF CHANGE ---
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("FATAL: Error loading .env file. Make sure it exists in the backend directory.")
	}
	log.Println("Successfully loaded .env file.")

	fabricService, err := blockchain.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize Fabric service: %v", err)
	}
	defer fabricService.Close()

	router := gin.Default()
	h := handlers.NewProduceHandler(fabricService.Contract)

	router.POST("/api/produce", h.CreateProduce)
	router.GET("/api/produce/:id", h.GetProduce)
	router.GET("/api/produce/:id/history", h.GetProduceHistory)
	router.PUT("/api/produce/:id/transfer", h.TransferProduce)

	fmt.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}