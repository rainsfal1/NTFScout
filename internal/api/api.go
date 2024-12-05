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
	Name        string `json:"name"`
	Contract    string `json:"contract"`
	Description string `json:"description"`
	FloorPrice  string `json:"floor_price"`
}

type Transaction struct {
	Hash     string `json:"hash"`
	To       string `json:"to"`
	Value    string `json:"value"`
	NftCount string `json:"nft_count"`
	CallData string `json:"call_data"`
	EthValue string `json:"eth_value"`
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
type AlchemyContract struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	TotalSupply string `json:"totalSupply"`
	Deployer    string `json:"deployerAddress"`
}

type AlchemyResponse struct {
	Contracts []AlchemyContract `json:"contracts"`
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
	
	// For now, return demo data if no API keys are configured
	if config.AlchemyAPIKey == "" && config.OpenSeaAPIKey == "" {
		return getDemoCollections(), nil
	}
	
	// Implement actual API calls here
	return []Collection{}, nil
}

func getDemoCollections() []Collection {
	return []Collection{
		{
			Name:        "Demo Collection 1",
			Contract:    "0x1234567890123456789012345678901234567890",
			Description: "Demo NFT collection for testing",
			FloorPrice:  "0.1",
		},
		{
			Name:        "Demo Collection 2", 
			Contract:    "0x0987654321098765432109876543210987654321",
			Description: "Another demo collection",
			FloorPrice:  "0.05",
		},
	}
}

func fetchFromOpenSea(ctx context.Context, config *ApiConfig) ([]Collection, error) {
	if config.OpenSeaAPIKey == "" {
		return nil, errors.New("OpenSea API key not configured")
	}
	
	// Get trending collections on Base chain
	url := fmt.Sprintf("%s/collections?chain=base&order_by=seven_day_volume&limit=20", config.OpenSeaBaseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("X-API-KEY", config.OpenSeaAPIKey)
	req.Header.Set("Accept", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenSea API error: %s", response.Status)
	}
	
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	
	var openSeaResp OpenSeaResponse
	if err := json.Unmarshal(body, &openSeaResp); err != nil {
		return nil, err
	}
	
	// Convert OpenSea format to our Collection format
	var collections []Collection
	for _, osCol := range openSeaResp.Collections {
		if len(osCol.Contracts) > 0 {
			collection := Collection{
				Contract:   osCol.Contracts[0].Address,
				Deployer:   "", // OpenSea doesn't provide deployer info
				Name:       osCol.Name,
				TotalMints: strconv.Itoa(osCol.TotalSupply),
				IsReported: []string{}, // Assume not reported
			}
			collections = append(collections, collection)
		}
	}
	
	return collections, nil
}

func fetchFromAlchemy(ctx context.Context, config *ApiConfig) ([]Collection, error) {
	if config.AlchemyAPIKey == "" {
		return nil, errors.New("Alchemy API key not configured")
	}
	
	// Get recently deployed NFT contracts on Base
	url := fmt.Sprintf("https://base-mainnet.g.alchemy.com/nft/v3/%s/getContractsForOwner?owner=0x0000000000000000000000000000000000000000&withMetadata=true&pageSize=20", config.AlchemyAPIKey)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Accept", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		log.Printf("Alchemy API response: %s", string(body))
		return nil, fmt.Errorf("Alchemy API error: %s", response.Status)
	}
	
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	
	log.Printf("Alchemy API response: %s", string(body))
	
	// For now, return empty since we need to understand the actual API structure
	// We'll use OpenSea as primary source and add Alchemy later
	return []Collection{}, nil
}

func filterCollections(collections []Collection) []Collection {
	// Remove duplicates and filter out reported collections
	seen := make(map[string]bool)
	var filtered []Collection
	
	for _, col := range collections {
		// Skip if already seen
		if seen[col.Contract] {
			continue
		}
		seen[col.Contract] = true
		
		// Skip if reported
		if len(col.IsReported) > 0 {
			continue
		}
		
		// Skip if contract address is empty
		if col.Contract == "" {
			continue
		}
		
		filtered = append(filtered, col)
	}
	
	return filtered
}

func GetTransaction(ctx context.Context, collection Collection) ([]Transaction, error) {
	// For now, create mock transaction data for testing
	mockTransaction := map[string]interface{}{
		"To":       collection.Contract,
		"CallData": "0x", // Would need to analyze contract ABI for actual mint function
		"NftCount": "1",  // Default to minting 1 NFT
		"EthValue": "0",  // Assume free mint initially
	}
	
	// Convert mock transaction to Transaction format
	transaction := Transaction{
		To:       mockTransaction["To"].(string),
		CallData: mockTransaction["CallData"].(string),
		NftCount: mockTransaction["NftCount"].(string),
		EthValue: mockTransaction["EthValue"].(string),
	}
	
	// If you have specific contracts you know how to interact with,
	// you could add more sophisticated transaction building here
	
	return []Transaction{transaction}, nil
}

// Helper function to create a mock collection for testing
func CreateMockCollection() Collection {
	return Collection{
		Contract:   "0x1234567890123456789012345678901234567890",
		Deployer:   "0x0987654321098765432109876543210987654321",
		Name:       "Test Collection",
		TotalMints: "100",
		IsReported: []string{},
	}
}

// Create demo collections for testing
func createDemoCollections() []Collection {
	return []Collection{
		{
			Contract:   "0x1234567890123456789012345678901234567890",
			Deployer:   "0x0987654321098765432109876543210987654321",
			Name:       "Demo NFT Collection #1",
			TotalMints: "150",
			IsReported: []string{},
		},
		{
			Contract:   "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			Deployer:   "0x1111111111111111111111111111111111111111",
			Name:       "Base Art Collection",
			TotalMints: "75",
			IsReported: []string{},
		},
		{
			Contract:   "0x9999999999999999999999999999999999999999",
			Deployer:   "0x2222222222222222222222222222222222222222",
			Name:       "Free Mint Collection",
			TotalMints: "500",
			IsReported: []string{},
		},
	}
}
