package main

import (
	"encoding/json"
	"fmt"
	"strconv" // Import the string conversion package
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct{ contractapi.Contract }

type Produce struct {
	ID        string    `json:"id"`
	Crop      string    `json:"crop"`
	Quantity  int       `json:"quantity"`
	Owner     string    `json:"owner"`
	Timestamp time.Time `json:"timestamp"`
}

// --- THIS FUNCTION IS CORRECTED ---
func (s *SmartContract) CreateProduce(ctx contractapi.TransactionContextInterface, id string, crop string, quantityStr string, owner string) error {
	exists, err := s.ProduceExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the produce %s already exists", id)
	}

	// Convert the quantity from a string to an integer
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return fmt.Errorf("quantity must be a numeric string: %w", err)
	}

	produce := Produce{
		ID:        id,
		Crop:      crop,
		Quantity:  quantity,
		Owner:     owner,
		Timestamp: time.Now(),
	}
	produceJSON, err := json.Marshal(produce)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, produceJSON)
}
// --- END OF CORRECTION ---

// (The rest of your smartcontract.go file remains the same)
func (s *SmartContract) ReadProduce(ctx contractapi.TransactionContextInterface, id string) (*Produce, error) {
	produceJSON, err := ctx.GetStub().GetState(id)
	if err != nil { return nil, fmt.Errorf("failed to read from world state: %v", err) }
	if produceJSON == nil { return nil, fmt.Errorf("the produce %s does not exist", id) }
	var produce Produce
	if err = json.Unmarshal(produceJSON, &produce); err != nil { return nil, err }
	return &produce, nil
}
func (s *SmartContract) TransferProduce(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	produce, err := s.ReadProduce(ctx, id)
	if err != nil { return err }
	produce.Owner = newOwner
	produce.Timestamp = time.Now()
	produceJSON, err := json.Marshal(produce)
	if err != nil { return err }
	return ctx.GetStub().PutState(id, produceJSON)
}
func (s *SmartContract) GetProduceHistory(ctx contractapi.TransactionContextInterface, id string) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil { return nil, fmt.Errorf("failed to get history for key: %v", err) }
	defer resultsIterator.Close()
	var history []string
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil { return nil, err }
		history = append(history, string(response.Value))
	}
	return history, nil
}
func (s *SmartContract) ProduceExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	produceJSON, err := ctx.GetStub().GetState(id)
	if err != nil { return false, fmt.Errorf("failed to read from world state: %v", err) }
	return produceJSON != nil, nil
}
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil { fmt.Printf("Error creating chaincode: %v", err); return }
	if err := chaincode.Start(); err != nil { fmt.Printf("Error starting chaincode: %v", err) }
}