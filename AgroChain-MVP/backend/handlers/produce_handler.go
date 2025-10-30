package handlers

import (
	"fmt"
	"log"
	"net/http"

	// --- THIS LINE IS CORRECTED ---
	"github.com/gin-gonic/gin"
	// --- END OF CORRECTION ---
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type ProduceHandler struct {
	Contract *client.Contract
}

func NewProduceHandler(contract *client.Contract) *ProduceHandler {
	return &ProduceHandler{Contract: contract}
}

type CreateProduceRequest struct {
	Crop     string `json:"crop" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
}

type TransferProduceRequest struct {
	NewOwner string `json:"newOwner" binding:"required"`
}

func (h *ProduceHandler) CreateProduce(c *gin.Context) {
	var req CreateProduceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	produceID := uuid.New().String()

	_, err := h.Contract.SubmitTransaction("CreateProduce", produceID, req.Crop, fmt.Sprintf("%d", req.Quantity), req.Owner)
	if err != nil {
		log.Printf("ERROR: Failed to submit transaction: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to submit transaction: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Produce created successfully", "produceId": produceID})
}

func (h *ProduceHandler) GetProduce(c *gin.Context) {
	produceID := c.Param("id")
	evaluateResult, err := h.Contract.EvaluateTransaction("ReadProduce", produceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %v", err)})
		return
	}
	c.Data(http.StatusOK, "application/json", evaluateResult)
}

func (h *ProduceHandler) TransferProduce(c *gin.Context) {
	produceID := c.Param("id")
	var req TransferProduceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Contract.SubmitTransaction("TransferProduce", produceID, req.NewOwner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to submit transaction: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Produce %s transferred to %s", produceID, req.NewOwner)})
}

func (h *ProduceHandler) GetProduceHistory(c *gin.Context) {
	produceID := c.Param("id")
	evaluateResult, err := h.Contract.EvaluateTransaction("GetProduceHistory", produceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %v", err)})
		return
	}
	c.Data(http.StatusOK, "application/json", evaluateResult)
}