# Micro GeoIP

A minimalistic REST API that returns the country from an IP address using MaxMind's GeoLite2 database.

## Features

- 🌍 **GeoIP Lookup**: Get country information from IP addresses
- 🔄 **Auto-updates**: Downloads GeoIP database monthly from MaxMind or DB-IP
- 🆓 **Free Database**: Uses DB-IP.com free database when no MaxMind API key is provided
- 🏠 **Caller IP Detection**: Uses caller's IP when no IP parameter is provided
- 🔒 **Security Option**: Can be configured to block IP parameter and always use caller IP
- 🐳 **Containerized**: Docker support with multi-architecture builds
- 🚀 **CI/CD**: GitHub Actions for automated testing and deployment to GHCR
- ⚡ **Lightweight**: Minimal resource usage and fast response times

## API Endpoints

### Health Check
```
GET /health
```

### GeoIP Lookup
```
GET /                    # Uses caller IP
GET /geoip               # Uses caller IP or ?ip parameter
GET /geoip?ip=8.8.8.8   # Looks up specific IP
GET /geoip/8.8.8.8      # Looks up specific IP
```

### Response Format
```json
{
  "ip": "8.8.8.8",
  "country": "United States",
  "country_code": "US"
}
```

Error response:
```json
{
  "ip": "invalid-ip",
  "country": "",
  "country_code": "",
  "error": "Invalid IP address"
}
```

## Configuration

### Environment Variables
- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: 0.0.0.0)
- `MAXMIND_API_KEY`: MaxMind API key for database downloads (optional)
- `GEOIP_DB_PATH`: Path to GeoIP database file (default: ./data/GeoLite2-Country.mmdb)
- `GEOIP_UPDATE_INTERVAL`: Update interval (default: 720h = 30 days)
- `MAXMIND_DOWNLOAD_URL`: MaxMind download URL (default: https://download.maxmind.com/app/geoip_download)
- `DBIP_DOWNLOAD_URL`: DB-IP download URL template (default: https://download.db-ip.com/free/dbip-country-lite-%s.mmdb.gz)
- `PREFER_DBIP`: Prefer DB-IP over MaxMind even if API key is available (default: false)
- `BLOCK_IP_PARAM`: Block IP parameter and always use caller IP (default: false)

### Configuration File
Create a `config.yaml` file (see `config.yaml.example`):

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

geoip:
  maxmind_api_key: "your-maxmind-api-key-here"  # Optional
  database_path: "./data/GeoLite2-Country.mmdb"
  update_interval: "720h"
  maxmind_url: "https://download.maxmind.com/app/geoip_download"
  dbip_url: "https://download.db-ip.com/free/dbip-country-lite-%s.mmdb.gz"
  prefer_dbip: false

security:
  block_ip_param: false
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- MaxMind API key (optional - will use free DB-IP database if not provided)

### Local Development
1. Clone the repository
2. Copy `config.yaml.example` to `config.yaml` and configure
3. **Optional**: Set your MaxMind API key for better accuracy:
   ```bash
   export MAXMIND_API_KEY=your-api-key-here
   ```
4. Run the application:
   ```bash
   go run main.go
   ```
   
   **Note**: If no MaxMind API key is provided, the service will automatically use the free DB-IP database.

### Using Docker

#### Build locally:
```bash
docker build -t micro-geoip .
```

#### Run with Docker:
```bash
# With MaxMind API key (recommended)
docker run -d \
  -p 8080:8080 \
  -e MAXMIND_API_KEY=your-api-key-here \
  --name micro-geoip \
  micro-geoip

# Without API key (uses free DB-IP database)
docker run -d \
  -p 8080:8080 \
  --name micro-geoip \
  micro-geoip
```

#### Use pre-built image:
```bash
# With MaxMind API key
docker run -d \
  -p 8080:8080 \
  -e MAXMIND_API_KEY=your-api-key-here \
  --name micro-geoip \
  ghcr.io/yourusername/micro_geoip:latest

# Without API key (uses free DB-IP database)
docker run -d \
  -p 8080:8080 \
  --name micro-geoip \
  ghcr.io/yourusername/micro_geoip:latest
```

### Kubernetes Deployment
Use the generated Kubernetes manifests from the GitHub Actions workflow or create your own:

```yaml
# Optional: Create secret for MaxMind API key
apiVersion: v1
kind: Secret
metadata:
  name: geoip-secret
stringData:
  api-key: "your-maxmind-api-key-here"
---
# Apply the deployment manifest from the release artifacts
# Note: Service works without the secret (will use DB-IP free database)
```

## Development

### Running Tests
```bash
go test -v ./...
```

### Running with Coverage
```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Building
```bash
go build -o micro_geoip .
```

## Security Features

### IP Parameter Blocking
When `BLOCK_IP_PARAM=true` or `security.block_ip_param: true`, the service will:
- Ignore any IP parameters in requests
- Always use the caller's IP address
- Useful for preventing IP enumeration attacks

### Client IP Detection
The service detects client IPs using:
1. `X-Forwarded-For` header (first IP)
2. `X-Real-IP` header
3. Request RemoteAddr as fallback

## Database Sources

### DB-IP (Default/Free)
- **Free database** from DB-IP.com
- No API key required
- Monthly updates available 
- Good accuracy for most use cases
- Used automatically when no MaxMind API key is provided

### MaxMind GeoLite2 (Optional/Enhanced)
- **Enhanced accuracy** with MaxMind's GeoLite2 Country database
- Requires a free MaxMind account and API key
- Monthly automatic updates
- Database cached locally for fast lookups
- Get your API key: https://www.maxmind.com/en/accounts/current/license-key

### Database Selection Priority
1. If `PREFER_DBIP=true`: Always use DB-IP
2. If MaxMind API key is provided and `PREFER_DBIP=false`: Use MaxMind
3. If MaxMind API key is not provided: Use DB-IP (free)
4. If MaxMind download fails: Fallback to DB-IP

## License

This project may use data from:
- **DB-IP.com**: See [DB-IP License](https://db-ip.com/db/lite.php) for their free database terms
- **MaxMind GeoLite2**: See [MaxMind License](https://www.maxmind.com/en/geolite2/eula) terms when using their database

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

The CI pipeline will run tests, linting, and security checks automatically.