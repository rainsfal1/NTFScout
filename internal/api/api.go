package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApiConfig struct {
	OpenSeaAPIKey   string
	AlchemyAPIKey   string
	OpenSeaBaseURL  string
	AlchemyBaseURL  string
}

type Collection struct {
	Contract   string   `json:"contract"`
	Deployer   string   `json:"deployer"`
	Name       string   `json:"name"`
	TotalMints string   `json:"totalMints"`
	IsReported []string `json:"isReported"`
}

type Transaction struct {
	To       string `json:"to"`
	CallData string `json:"callData"`
	NftCount string `json:"nftCount"`
	EthValue string `json:"ethValue"`
}

type CollectionResponse struct {
	Collections []Collection `json:"collections"`
}

type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// OpenSea API response structures
type OpenSeaCollection struct {
	Collection string `json:"collection"`
	Name       string `json:"name"`
	Contracts  []struct {
		Address string `json:"address"`
	} `json:"contracts"`
	TotalSupply int `json:"total_supply"`
}

type OpenSeaResponse struct {
	Collections []OpenSeaCollection `json:"collections"`
}

// Alchemy API response structures
type AlchemyResponse struct {
	Contracts []struct {
		Address     string `json:"address"`
		Name        string `json:"name"`
		TotalSupply string `json:"totalSupply"`
		TokenType   string `json:"tokenType"`
		OpenSeaMetadata struct {
			FloorPrice     interface{} `json:"floorPrice"`
			CollectionName string      `json:"collectionName"`
		} `json:"openSeaMetadata"`
		IsSpam bool `json:"isSpam"`
	} `json:"contracts"`
	TotalCount int    `json:"totalCount"`
	PageKey    string `json:"pageKey"`
}

func GetApiConfig() ApiConfig {
	return ApiConfig{
		OpenSeaAPIKey:  os.Getenv("OPENSEA_API_KEY"),
		AlchemyAPIKey:  os.Getenv("ALCHEMY_API_KEY"),
		OpenSeaBaseURL: os.Getenv("OPENSEA_BASE_URL"),
		AlchemyBaseURL: os.Getenv("ALCHEMY_BASE_URL"),
	}
}

func FetchCollection(ctx context.Context) ([]Collection, error) {
	config := GetApiConfig()
	
	// Demo mode: if no API keys are configured, use mock data
	if config.OpenSeaAPIKey == "" && config.AlchemyAPIKey == "" {
		log.Println("No API keys configured, running in DEMO mode with mock data")
		return createDemoCollections(), nil
	}
	
	// Combine data from both OpenSea and Alchemy
	var allCollections []Collection
	
	// Fetch from OpenSea (trending collections on Base)
	openSeaCollections, err := fetchFromOpenSea(ctx, config)
	if err != nil {
		log.Printf("OpenSea fetch error: %v", err)
	} else {
		allCollections = append(allCollections, openSeaCollections...)
	}
	
	// Fetch from Alchemy (newly deployed contracts)
	alchemyCollections, err := fetchFromAlchemy(ctx, config)
	if err != nil {
		log.Printf("Alchemy fetch error: %v", err)
	} else {
		allCollections = append(allCollections, alchemyCollections...)
	}
	
	if len(allCollections) == 0 {
		log.Println("No collections from APIs, falling back to demo data")
		return createDemoCollections(), nil
	}
	
	// Filter out reported collections and duplicates
	filteredCollections := filterCollections(allCollections)
	
	log.Printf("Fetched %d collections from multiple sources", len(filteredCollections))
	return filteredCollections, nil
}

func fetchFromOpenSea(ctx context.Context, config ApiConfig) ([]Collection, error) {
	// OpenSea implementation would go here
	// For now, return empty slice as OpenSea API requires specific endpoints
	return []Collection{}, nil
}

func fetchFromAlchemy(ctx context.Context, config ApiConfig) ([]Collection, error) {
	if config.AlchemyAPIKey == "" {
		return []Collection{}, fmt.Errorf("Alchemy API key not configured")
	}

	// Use a dummy address for testing - in production this would be dynamic
	url := fmt.Sprintf("%s/%s/getContractsForOwner?owner=0x0000000000000000000000000000000000000000&withMetadata=true&pageSize=20", 
		config.AlchemyBaseURL, config.AlchemyAPIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var alchemyResp AlchemyResponse
	if err := json.NewDecoder(resp.Body).Decode(&alchemyResp); err != nil {
		return nil, err
	}

	// Log the response for debugging
	respBytes, _ := json.Marshal(alchemyResp)
	log.Printf("Alchemy API response: %s", string(respBytes))

	var collections []Collection
	for _, contract := range alchemyResp.Contracts {
		if !contract.IsSpam && contract.Name != "" {
			collections = append(collections, Collection{
				Contract:   contract.Address,
				Name:       contract.Name,
				TotalMints: contract.TotalSupply,
				Deployer:   contract.Address, // Using contract address as deployer for now
				IsReported: []string{},
			})
		}
	}

	return collections, nil
}

func filterCollections(collections []Collection) []Collection {
	seen := make(map[string]bool)
	var filtered []Collection
	
	for _, collection := range collections {
		if !seen[collection.Contract] && len(collection.IsReported) == 0 {
			seen[collection.Contract] = true
			filtered = append(filtered, collection)
		}
	}
	
	return filtered
}

func createDemoCollections() []Collection {
	return []Collection{
		{
			Contract:   "0x1234567890123456789012345678901234567890",
			Name:       "Demo NFT Collection 1",
			TotalMints: "1000",
			Deployer:   "0xdemo1234567890123456789012345678901234567890",
			IsReported: []string{},
		},
		{
			Contract:   "0x0987654321098765432109876543210987654321",
			Name:       "Demo NFT Collection 2",
			TotalMints: "500",
			Deployer:   "0xdemo0987654321098765432109876543210987654321",
			IsReported: []string{},
		},
	}
}

func GetTransaction(ctx context.Context, contract string) ([]Transaction, error) {
	// Return demo transactions for now
	return []Transaction{
		{
			To:       contract,
			CallData: "0x",
			NftCount: "1",
			EthValue: "0.01",
		},
	}, nil
}
