# 🏎️ F1 Race Leaderboard

A cost-effective, serverless API for fetching live Formula 1 race data — built on AWS with Go, Terraform/OpenTofu, and Nix.

## Architecture

```
Client → API Gateway → Lambda (Go) → DynamoDB (cache) → Jolpica F1 API
                           │
                    VPC (private subnets)
                           │
                    NAT Gateway (outbound)
                           │
                    DynamoDB VPC Endpoint (free)
```

| Component         | Purpose                              |
| ----------------- | ------------------------------------ |
| API Gateway (v2)  | HTTP entry point, CORS, routing      |
| Lambda (Go/arm64) | Business logic, caching strategy     |
| DynamoDB          | TTL-based cache (PAY_PER_REQUEST)    |
| VPC + NAT         | Network isolation, outbound internet |
| VPC Endpoint      | Free DynamoDB access from VPC        |

## API Endpoints

| Route                         | Description               |
| ----------------------------- | ------------------------- |
| `GET /`                       | Health check + route list |
| `GET /standings`              | Driver championship       |
| `GET /standings/drivers`      | Driver championship       |
| `GET /standings/constructors` | Constructor championship  |
| `GET /results`                | Latest race results       |
| `GET /schedule`               | Current season schedule   |

All responses include:

- `X-Data-Source: cache | api` header indicating data origin
- CORS headers for browser access
- 5-minute TTL cache to minimize external API calls

## Prerequisites

- [Nix](https://nixos.org/download) (recommended) **or** Go 1.21+, OpenTofu/Terraform 1.5+, AWS CLI
- AWS account with credentials configured

## Quick Start

```bash
# Enter dev environment
nix develop

# Run tests
make test

# Build Lambda binary
make build

# Deploy to AWS
make init
make deploy
```

## Development

```bash
# Format code
make fmt

# Validate infrastructure (OpenTofu)
make validate

# Preview changes
make plan

# Tear down
make destroy
```

## Configuration

| Variable             | Default     | Description              |
| -------------------- | ----------- | ------------------------ |
| `aws_region`         | `eu-west-1` | AWS deployment region    |
| `environment`        | `dev`       | Environment tag          |
| `cache_ttl_seconds`  | `300`       | Cache TTL (seconds)      |
| `lambda_memory_size` | `128`       | Lambda memory (MB)       |
| `lambda_timeout`     | `10`        | Lambda timeout (seconds) |

Override via `terraform/terraform.tfvars` (works with both OpenTofu and Terraform):

```hcl
aws_region        = "us-east-1"
cache_ttl_seconds = 600
```

## Data Source

Uses the [Jolpica Ergast-compatible API](https://api.jolpi.ca/ergast/f1/) — free, no authentication required.

## License

MIT
