# Hybrid Store Architecture Guide

## Overview

The AI Interview Backend now features a **Hybrid Store Architecture** that automatically detects and switches between memory and database backends based on environment configuration. This provides maximum flexibility for development, testing, and production deployments.

## Backend Auto-Detection

The system automatically chooses the backend based on the `DATABASE_URL` environment variable:

- **Memory Backend**: Used when `DATABASE_URL` is not set or empty
- **Database Backend**: Used when `DATABASE_URL` is provided

## Configuration

### Required Environment Variables
- **None** - The system works out of the box with sensible defaults

### Optional Environment Variables
- `DATABASE_URL` - PostgreSQL connection string (enables database backend)
- `PORT` - Server port (defaults to 8080)

### Example Configurations

#### Development Mode (Memory Backend)
```bash
# No configuration needed - uses defaults
PORT=3001
```

#### Production Mode (Database Backend)
```bash
DATABASE_URL=postgresql://user:password@localhost:5432/ai_interview
PORT=8080
```

## Backend Features

### Memory Backend
- ✅ Zero setup - works immediately
- ✅ Fast performance for development/testing
- ✅ Full feature parity with database backend
- ✅ Pagination, filtering, and sorting
- ✅ Complete CRUD operations
- ⚠️ Data lost on restart
- ⚠️ No persistence across deployments

### Database Backend
- ✅ Persistent data storage
- ✅ Production-ready scalability
- ✅ ACID transactions
- ✅ Advanced querying capabilities
- ✅ Automatic schema migrations
- ⚠️ Requires PostgreSQL setup
- ⚠️ Additional configuration needed

## Unified Interface

Both backends implement the same interface, providing:

1. **Interview Management**
   - Create, read, update interviews
   - Advanced listing with pagination, filtering, sorting
   - Status tracking

2. **Evaluation System**
   - Store and retrieve evaluations
   - Link evaluations to interviews
   - Support for both traditional Q&A and chat evaluations

3. **Chat Sessions**
   - Real-time chat session management
   - Message history storage
   - Session lifecycle tracking

4. **Data Integrity**
   - Consistent data models
   - Proper error handling
   - Transaction support (database mode)

## Usage Examples

### Development Workflow
```bash
# Start development server (memory backend)
go run main.go

# Output: "Using in-memory store backend (set DATABASE_URL for database mode)"
```

### Production Deployment
```bash
# Set database URL and start
export DATABASE_URL="postgresql://user:pass@db:5432/ai_interview"
go run main.go

# Output: "Using PostgreSQL database backend"
```

### Testing
```bash
# Unit tests (memory backend)
go test ./...

# Integration tests with database
export DATABASE_URL="postgresql://test:test@localhost:5432/ai_interview_test"
go test ./...
```

## Migration Path

### From Memory to Database
1. Set `DATABASE_URL` environment variable
2. Restart application
3. Data will be persisted going forward
4. Previous memory data will be lost (expected)

### From Database to Memory (Debugging)
1. Unset or remove `DATABASE_URL`
2. Restart application
3. Fresh memory store will be created
4. Database data remains intact

## Implementation Details

### Store Initialization
```go
// Auto-detects backend and initializes appropriate store
err = data.InitGlobalStore()
if err != nil {
    log.Fatalf("failed to initialize store: %v", err)
}

// Check which backend is active
if data.GlobalStore.GetBackend() == data.BackendDatabase {
    fmt.Println("Using PostgreSQL database backend")
} else {
    fmt.Println("Using in-memory store backend")
}
```

### Method Compatibility
All methods maintain identical signatures across backends:
```go
// Same interface for both backends
func GetInterviewsWithOptions(options ListInterviewsOptions) (*ListInterviewsResult, error)
func CreateInterview(interview *Interview) error
func AddChatMessage(sessionID string, message *ChatMessage) error
```

## Benefits

1. **Developer Experience**: No database setup required for local development
2. **Testing Efficiency**: Fast unit tests with memory backend
3. **Production Ready**: Seamless transition to database backend
4. **Flexibility**: Easy switching between backends without code changes
5. **Gradual Adoption**: Teams can adopt database backend when ready

## Best Practices

### Development
- Use memory backend for quick prototyping
- Run tests with memory backend for speed
- Use database backend for integration testing

### Production
- Always use database backend with `DATABASE_URL`
- Set up proper database monitoring
- Configure backup and recovery procedures
- Use connection pooling for high load

### Testing
- Unit tests: Memory backend (fast)
- Integration tests: Database backend (realistic)
- E2E tests: Database backend (production-like)

## Troubleshooting

### Common Issues

**Problem**: "failed to initialize store"
**Solution**: Check DATABASE_URL format and database connectivity

**Problem**: Data not persisting
**Solution**: Verify DATABASE_URL is set for database backend

**Problem**: Slow performance
**Solution**: Use memory backend for development, database for production

**Problem**: Connection errors
**Solution**: Ensure PostgreSQL is running and accessible

### Health Checks
```go
// Check store health
if err := data.GlobalStore.Health(); err != nil {
    log.Printf("Store health check failed: %v", err)
}
```

## Future Enhancements

- [ ] Redis backend for distributed caching
- [ ] MongoDB backend for document storage
- [ ] Automatic failover between backends
- [ ] Hot-swapping backends without restart
- [ ] Backend-specific optimizations
- [ ] Metrics and monitoring per backend
