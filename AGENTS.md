# AGENTS.md - SnippetBox Development Guidelines

## Build & Run Commands
```bash
# Run development server
go run ./cmd/web -db=postgresql://postgres:PASSWORD@localhost:5432/snippetbox

# Run with custom port
go run ./cmd/web -port=:3000 -db=postgresql://postgres:PASSWORD@localhost:5432/snippetbox

# Database setup
docker compose up -d    # Start PostgreSQL
docker compose down     # Stop database
```

## Code Style Guidelines

### Import Organization
- Standard library imports first (grouped)
- Third-party imports second  
- Local imports third with full module paths
- Example: `"github.com/corbinlazarone/snippetbox/internal/models"`

### Naming Conventions
- Packages: lowercase, single words (`models`, `validator`)
- Structs: PascalCase (`SnippetModel`, `application`)
- Methods: PascalCase for public, camelCase for private
- Variables: camelCase (`templateCache`, `dataSource`)
- Constants: PascalCase for exported, ALL_CAPS for package-level

### Error Handling
- Use custom error types from `internal/models/errors.go`
- Return tuples `(value, error)` consistently
- Handle errors immediately, don't let them propagate silently
- Use structured logging for errors

### Project Structure
- `cmd/web/` - HTTP handlers, routing, middleware, main application
- `internal/` - Private application code (models, validation)
- `ui/` - Frontend assets (templates, CSS, JS)
- Follow dependency injection pattern through struct composition

### Testing
- No formal test suite exists - this is an educational project
- Focus on manual testing when making changes