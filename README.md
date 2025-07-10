# URL Crawler

A comprehensive URL Crawler application with real-time crawling analysis, visualization, and Firecrawl integration. Built with Go backend, React frontend, and MySQL database.

> **Initial Project Scaffolding was done with** [go-blueprint](https://github.com/Melkeydev/go-blueprint) - A CLI tool that allows users to spin up a quick Go project.

## Demo Video

### 🎬 See the URL Crawler in Action

_See the URL Crawler in action: real-time crawling, interactive dashboard, mobile-responsive design, and queue-based processing with Firecrawl integration._

<div align="center">
  <img src="./video/demo-video.gif" alt="URL Crawler Demo" style="width:100%;max-width:800px;border-radius:8px;box-shadow:0 4px 8px rgba(0,0,0,0.1);">
</div>

## 🚀 Features

- **Real-time URL crawling** with Firecrawl integration
- **Interactive dashboard** with charts and analytics
- **Queue-based processing** with configurable workers
- **Authentication & rate limiting** for API security
- **Mobile-responsive UI** with modern design
- **Docker containerization** for easy development and deployment

## 🛠️ Tech Stack

- **Backend**: Go + Echo framework
- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Database**: MySQL with JSON support
- **Queue System**: In-memory with configurable workers
- **Crawler**: Firecrawl API integration
- **Development**: Docker Compose

## 📋 Prerequisites

- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience commands)
- [Firecrawl API Key](https://firecrawl.dev) (for crawler functionality)

## 🚀 Quick Start with Docker

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd url-crawler
```

### 2. Environment Configuration

The `.env` file is already configured with default values. Update these key variables if needed:

```bash
# Required: Get your API key from https://firecrawl.dev
FIRECRAWL_API_KEY=your-firecrawl-api-key

# Optional: Customize database credentials
URL_CRAWLER_DB_ROOT_PASSWORD=rootpassword123
URL_CRAWLER_DB_USERNAME=crawler_user
URL_CRAWLER_DB_PASSWORD=crawler_password123
```

### 3. Start All Services

```bash
# Docker development (recommended) - auto-applies migrations
make docker-up

# Or using Docker Compose directly
docker-compose up -d

# Local development (requires local MySQL)
make db-setup     # Run once to create database
make db-migrate   # Apply migrations
make run         # Start backend and frontend
```

### 4. Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/api/health (health check)

## 🐳 Docker Services

| Service  | Port | Description                    |
| -------- | ---- | ------------------------------ |
| Frontend | 5173 | React app with Vite dev server |
| Backend  | 8080 | Go API with hot reload         |
| MySQL    | 3306 | Database with auto-migrations  |

## 📁 Project Structure

```
url-crawler/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── database/      # Database connection & migrations
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # Authentication & rate limiting
│   ├── models/        # Data models
│   ├── server/        # Server setup & routes
│   └── services/      # Business logic & crawler
├── frontend/          # React frontend application
├── docker-compose.yml # Docker services configuration
├── Dockerfile     # Development Docker image
├── .env              # Environment variables
└── Makefile          # Development commands
```

## 🔧 Development Commands

### Docker Management

```bash
# Start all services (auto-applies migrations)
make docker-up

# Stop all services
make docker-down

# View logs
docker-compose logs
docker-compose logs -f backend    # Follow specific service

# Restart specific service
docker-compose restart backend
docker-compose restart frontend
docker-compose restart mysql
```

### Local Development (Alternative)

```bash
# Database setup (run once)
make db-setup               # Create database and user
make db-migrate             # Apply migrations

# Build the application
make build

# Run backend locally (requires local MySQL)
make run

# Run with live reload
make watch

# Run tests
make test

# Integration tests
make itest
```

### Database Management

```bash
# Docker development
make docker-up              # Start all containers (auto-applies migrations)
make docker-down            # Stop all containers

# Local development
make db-setup               # Setup database and user (run once)
make db-migrate             # Apply migrations

# Database access
docker-compose exec mysql mysql -u crawler_user -p url_crawler  # Docker
mysql -u crawler_user -p url_crawler                           # Local

# View database logs
docker-compose logs mysql
```

## 🗄️ Database Migrations

### Overview

The project uses SQL migrations located in `internal/database/migrations_mysql.sql`. Migrations are automatically applied when the MySQL container starts.

### Migration File Structure

```sql
-- Create tables
CREATE TABLE IF NOT EXISTS crawl_results (
    id VARCHAR(36) PRIMARY KEY,
    original_url TEXT NOT NULL,
    final_url TEXT,
    status_code INT,
    content_type VARCHAR(255),
    -- ... other fields
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_crawl_results_created_at ON crawl_results(created_at);
-- ... other indexes
```

### How Migrations Work

1. **Automatic Application**: Migrations run automatically when MySQL container starts
2. **Docker Volume Mount**: Migration file is mounted as `/docker-entrypoint-initdb.d/init.sql`
3. **Idempotent**: Uses `IF NOT EXISTS` to safely run multiple times
4. **Single File**: All migrations are in one file for simplicity

### Managing Migrations

#### For Docker Development (Recommended)

```bash
# Migrations auto-apply when containers start
make docker-up

# View current database schema (Docker)
docker-compose exec mysql mysql -u crawler_user -p url_crawler -e "SHOW TABLES;"

# Describe a specific table (Docker)
docker-compose exec mysql mysql -u crawler_user -p url_crawler -e "DESCRIBE crawl_results;"
```

#### For Local Development

```bash
# Setup database and user (run once)
make db-setup

# Run migrations
make db-migrate

# Check migration file
cat internal/database/migrations_mysql.sql
```

## 🔑 Environment Variables

Key configuration variables in `.env`:

```bash
# Server Configuration
PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=30s

# Database Configuration
URL_CRAWLER_DB_HOST=mysql
URL_CRAWLER_DB_PORT=3306
URL_CRAWLER_DB_USERNAME=crawler_user
URL_CRAWLER_DB_PASSWORD=crawler_password123
URL_CRAWLER_DB_DATABASE=url_crawler

# Crawler Configuration
CRAWLER_TIMEOUT=30s
CRAWLER_USER_AGENT=URL-Crawler/1.0
CRAWLER_MAX_REDIRECTS=5

# Authentication
AUTH_REQUIRED=true
API_KEY_DEV=dev-api-key-12345

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=60

# Frontend Configuration
VITE_API_BASE_URL=http://localhost:8080
VITE_API_KEY=dev-api-key-12345

# External Services
FIRECRAWL_API_KEY=your-firecrawl-api-key
```

## 🐛 Troubleshooting

### Common Issues

1. **Frontend crypto.hash error**

   ```bash
   # Fixed by using Node 20+ in docker-compose.yml
   image: node:20-alpine
   ```

2. **Backend build errors**

   ```bash
   # Check Air configuration in .air.toml
   cmd = "go build -o ./main ./cmd/api/main.go"
   ```

3. **Database connection failed**

   ```bash
   # Ensure MySQL is healthy before starting backend
   docker-compose ps
   # Restart backend after MySQL is ready
   docker-compose restart backend
   ```

4. **Port conflicts**

   ```bash
   # Check if ports are already in use
   lsof -i :8080
   lsof -i :5173
   lsof -i :3306
   ```

5. **Migration syntax errors**
   ```bash
   # Check for duplicate column errors in MySQL logs
   docker-compose logs mysql | grep ERROR
   # Fix: Use IF NOT EXISTS in migration statements
   # Example: ALTER TABLE table ADD COLUMN IF NOT EXISTS new_col VARCHAR(255);
   ```

### Logs and Debugging

```bash
# View all logs
docker-compose logs

# Follow specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mysql

# Check container status
docker-compose ps

# Access container shell
docker-compose exec backend sh
docker-compose exec mysql sh
```

## 🔄 Development Workflow

1. **Start Development Environment**

   ```bash
   make docker-up
   ```

2. **Make Changes**

   - Backend: Edit files in `internal/` or `cmd/` (hot reload enabled)
   - Frontend: Edit files in `frontend/src/` (Vite hot reload)

3. **View Changes**

   - Frontend: http://localhost:5173 (auto-refresh)
   - Backend: Changes reflected immediately via Air

4. **Database Changes**
   - Update migrations in `internal/database/migrations_mysql.sql`
   - Use idempotent SQL (`IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`)
   - Test migration:
     - Docker: `docker-compose restart mysql`
     - Local: `make db-migrate`
   - Verify changes:
     - Docker: `docker-compose exec mysql mysql -u crawler_user -p url_crawler -e "SHOW TABLES;"`
     - Local: `mysql -u crawler_user -p url_crawler -e "SHOW TABLES;"`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with `make docker-up`
5. Submit a pull request
