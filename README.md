# NFTScout Bot

A Go-based NFT discovery and monitoring bot that aggregates data from multiple sources including OpenSea and Alchemy to provide comprehensive NFT insights and collection tracking.

## Features

- **Multi-Source Data Aggregation**: Fetches NFT data from OpenSea and Alchemy APIs for comprehensive coverage
- **Real-time Monitoring**: Continuously monitors NFT collections and transactions
- **MongoDB Integration**: Stores collection data, transactions, and error logs
- **Robust Error Handling**: Comprehensive error logging and recovery mechanisms
- **Configurable**: Easy configuration through environment variables

## Architecture

NFTScout uses a multi-source architecture to provide reliable and comprehensive NFT data:

- **OpenSea API**: Primary source for collection metadata and marketplace data
- **Alchemy API**: Secondary source for blockchain transaction data and NFT metadata
- **MongoDB**: Persistent storage for collections, transactions, and error logs
- **Go Runtime**: High-performance concurrent processing

## Prerequisites

- Go 1.22.2 or higher
- MongoDB instance (local or cloud)
- OpenSea API key
- Alchemy API key

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd NFTScoutBot
```

2. Install dependencies:
```bash
go mod tidy
```

3. Copy the environment configuration:
```bash
cp .env.example .env
```

4. Configure your environment variables in `.env`:
```env
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGODB_DATABASE=nftscout
MONGODB_COLLECTION_TRANSACTION=transactions
MONGODB_COLLECTION_ERROR=errors

# API Keys and Endpoints
OPENSEA_API_KEY=your_opensea_api_key_here
OPENSEA_BASE_URL=https://api.opensea.io/api/v1
ALCHEMY_API_KEY=your_alchemy_api_key_here
ALCHEMY_BASE_URL=https://eth-mainnet.g.alchemy.com/v2
```

## Usage

Run the bot:
```bash
go run main.go
```

The bot will:
1. Connect to MongoDB and verify collections
2. Initialize API connections to OpenSea and Alchemy
3. Start monitoring and fetching NFT data
4. Store results in the configured MongoDB database

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MONGO_URI` | MongoDB connection string | Yes |
| `MONGODB_DATABASE` | Database name (default: nftscout) | Yes |
| `MONGODB_COLLECTION_TRANSACTION` | Collection for transactions | Yes |
| `MONGODB_COLLECTION_ERROR` | Collection for error logs | Yes |
| `OPENSEA_API_KEY` | OpenSea API key | Yes |
| `OPENSEA_BASE_URL` | OpenSea API base URL | Yes |
| `ALCHEMY_API_KEY` | Alchemy API key | Yes |
| `ALCHEMY_BASE_URL` | Alchemy API base URL | Yes |

## Data Sources

### OpenSea API
- Collection metadata
- Floor prices and statistics
- Recent sales data
- Marketplace listings

### Alchemy API
- Blockchain transaction data
- NFT metadata and ownership
- Contract information
- Token transfers

## Database Schema

### Collections
- Collection metadata and statistics
- Floor prices and volume data
- Last updated timestamps

### Transactions
- Transaction hashes and details
- NFT transfers and sales
- Timestamp and block information

### Errors
- Error logs and stack traces
- API response errors
- System failures and recovery

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.