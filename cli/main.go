// CLOB/cli/main.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"CLOB/actions"
	"CLOB/genesis"
	"CLOB/storage"
	"CLOB/vm"
)

func main() {
	// Define a flag for the genesis configuration file
	genesisFile := flag.String("genesis", "", "Path to genesis JSON configuration file")
	flag.Parse()

	// Read genesis configuration
	var genesisConfig []byte
	var err error
	if *genesisFile != "" {
		genesisConfig, err = ioutil.ReadFile(*genesisFile)
		if err != nil {
			log.Fatalf("Failed to read genesis file: %v", err)
		}
	} else {
		// Use default genesis configuration
		genesisInstance := genesis.Default()
		genesisConfig, err = json.Marshal(genesisInstance)
		if err != nil {
			log.Fatalf("Failed to marshal default genesis configuration: %v", err)
		}
	}

	// Initialize the VM with genesis configuration
	engine, err := vm.NewMatchingEngineVM(genesisConfig)
	if err != nil {
		log.Fatalf("Failed to initialize VM: %v", err)
	}

	fmt.Println("Order Book Matching Engine Initialized with Genesis Configuration.")

	// Example: Adding a new order via CLI (can be extended to accept user inputs)
	buyOrder := &storage.Order{
		ID:        "cli_buy_1",
		Side:      storage.Buy,
		Price:     102.0,
		Quantity:  20,
		Timestamp: time.Now().UTC(),
		OrderType: storage.Limit,
	}

	addBuyOrderAction := &actions.AddOrderAction{Order: buyOrder}
	if err := engine.ExecuteAction(addBuyOrderAction); err != nil {
		fmt.Printf("Error adding buy order: %v\n", err)
	} else {
		fmt.Println("Buy order added successfully!")
	}

	// Add more CLI interactions as needed
	// For example, parsing user commands to add/cancel orders, view order book, etc.

	// Example: Display remaining orders
	fmt.Println("Current Orders in the Order Book:")
	for id, order := range engine.OrderBook.OrderMap {
		fmt.Printf("Order ID: %s, Side: %s, Quantity: %.2f, Price: %.2f, Type: %s\n",
			id, order.Side, order.Quantity, order.Price, order.OrderType)
	}

	// Prevent the CLI from exiting immediately (for demonstration)
	// In a real application, you'd implement a loop to accept user commands
	os.Exit(0)
}
