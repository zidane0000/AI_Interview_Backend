# AI Interview Backend

A robust Go backend service for an AI-driven interview platform. Features flexible data storage (memory/PostgreSQL), comprehensive REST API, and AI-powered interview evaluation.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL (optional - uses memory backend if not configured)

### Run with Memory Backend (No Database)
```bash
git clone https://github.com/zidane0000/AI_Interview_Backend
cd AI_Interview_Backend
go mod tidy
go run main.go
```
Server starts at `http://localhost:8080`

### Run with PostgreSQL Backend
```bash
# Set database connection
export DATABASE_URL="postgres://username:password@localhost:5432/ai_interview"
go run main.go
```

## 📋 Features

- **Flexible Storage**: Auto-detects memory or PostgreSQL backend
- **RESTful API**: Complete interview and evaluation endpoints
- **AI Integration**: OpenAI/Gemini support with mock fallback
- **Comprehensive Testing**: Unit and E2E test coverage
- **Production Ready**: Docker support, environment configuration

## 📖 Documentation

- **[API Reference](docs/openapi.yaml)** - Complete API specification
- **[Documentation Hub](docs/README.md)** - All documentation navigation
- **[Architecture Details](architecture/README.md)** - System design and patterns

## 🧪 Testing

### Quick Test Commands
```bash
# Run all tests
go test ./...

# Run unit tests only (fast)
go test $(go list ./... | Where-Object { $_ -notlike "*e2e*" })

# Run E2E tests (requires running server)  
go test ./e2e/...
```

### Test Coverage
- Unit tests for all components
- E2E tests for complete workflows
- Comprehensive error handling tests

## 🏗️ Architecture

```
├── api/          # HTTP handlers, routing, middleware
├── business/     # Core business logic
├── data/         # Database models and repositories  
├── ai/           # AI provider integrations
├── config/       # Environment configuration
├── docs/         # Documentation and API specs
└── e2e/          # End-to-end tests
```

## 🔧 Configuration

The application auto-detects the backend based on environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | No | *none* | PostgreSQL connection (enables DB backend) |
| `PORT` | No | `8080` | HTTP server port |
| `AI_API_KEY` | No | *none* | OpenAI/Gemini API key (uses mock if not set) |

**Examples:**
```bash
# Development (memory backend)
PORT=8080

# Production (PostgreSQL backend)  
DATABASE_URL=postgres://user:pass@host:5432/ai_interview
PORT=8080
AI_API_KEY=sk-your-openai-key
```
