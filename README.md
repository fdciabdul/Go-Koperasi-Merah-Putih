# Go Koperasi Management System

Sistem manajemen koperasi berbasis Go dengan arsitektur clean dan fitur lengkap untuk operasional koperasi di Indonesia.

## 🚀 Features

| Module | Description | Status |
|--------|-------------|--------|
| **Authentication** | User registration, login, JWT auth | ✅ Complete |
| **Koperasi Management** | CRUD koperasi, member management | ✅ Complete |
| **Simpan Pinjam** | Savings & loans products, transactions | ✅ Complete |
| **Product Management** | Inventory, suppliers, sales, purchases | ✅ Complete |
| **Klinik** | Healthcare services, patients, medicines | ✅ Complete |
| **Financial** | Chart of accounts, journals, reports | ✅ Complete |
| **PPOB** | Payment Point Online Bank services | ✅ Complete |
| **Payment Gateway** | Midtrans & Xendit integration | ✅ Complete |
| **Analytics** | Cassandra-based analytics | 🔄 In Progress |
| **Audit Logging** | Complete system audit trail | ✅ Complete |

## 🏗️ Architecture

### Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend** | Go 1.19+, Gin Web Framework | REST API server |
| **Database** | PostgreSQL 13+ | Primary data storage |
| **Analytics** | Apache Cassandra | Big data analytics |
| **Cache** | Redis (optional) | Session & cache management |
| **ORM** | GORM v2 | Database operations |
| **Authentication** | JWT + bcrypt | Security layer |
| **Payment** | Midtrans, Xendit | Payment processing |

## 🛠️ Installation & Setup

### Prerequisites

| Requirement | Version | Installation |
|-------------|---------|--------------|
| **Go** | 1.19+ | [Download Go](https://golang.org/dl/) |
| **PostgreSQL** | 13+ | [Download PostgreSQL](https://www.postgresql.org/download/) |
| **Git** | Latest | [Download Git](https://git-scm.com/downloads) |

### Quick Start

#### For Unix/Linux/macOS:
```bash
# 1. Clone Repository
git clone <repository-url>
cd go_koperasi

# 2. Install Dependencies
go mod download

# 3. Configure Environment
cp .env.example .env
# Edit .env with your database credentials

# 4. Setup Database
createdb koperasi_db
make migrate-fresh

# 5. Run Application
make run
```

#### For Windows:
```cmd
REM 1. Clone Repository
git clone <repository-url>
cd go_koperasi

REM 2. Install Dependencies
go mod download

REM 3. Configure Environment
copy .env.example .env
REM Edit .env with your database credentials

REM 4. Setup Database
createdb koperasi_db
make.bat migrate-fresh

REM 5. Run Application
make.bat run
```

#### One-Command Setup:
```bash
# Unix/Linux/macOS
make quick-start

# Windows
make.bat quick-start
```

## 📋 Available Commands

This project supports both **Unix/Linux/macOS** (Makefile) and **Windows** (make.bat) environments.

### Command Usage

| Platform | Usage | Example |
|----------|-------|---------|
| **Unix/Linux/macOS** | `make <command>` | `make run` |
| **Windows** | `make.bat <command>` | `make.bat run` |

### Development Commands

| Command | Description | Unix | Windows |
|---------|-------------|------|---------|
| `help` | Show all available commands | `make help` | `make.bat help` |
| `build` | Build the application | `make build` | `make.bat build` |
| `run` | Run the application | `make run` | `make.bat run` |
| `test` | Run all tests | `make test` | `make.bat test` |
| `fmt` | Format code | `make fmt` | `make.bat fmt` |
| `lint` | Lint code | `make lint` | `make.bat lint` |
| `dev` | Hot reload development | `make dev` | `make.bat dev` |

### Database Commands

| Command | Description | Unix | Windows |
|---------|-------------|------|---------|
| `migrate` | Run GORM auto-migrations | `make migrate` | `make.bat migrate` |
| `seed` | Run database seeders | `make seed` | `make.bat seed` |
| `migrate-fresh` | Drop, migrate, and seed | `make migrate-fresh` | `make.bat migrate-fresh` |
| `migrate-drop` | Drop all tables and migrate | `make migrate-drop` | `make.bat migrate-drop` |
| `dev-setup` | Complete development setup | `make dev-setup` | `make.bat dev-setup` |

### Tool Commands

| Command | Description | Unix | Windows |
|---------|-------------|------|---------|
| `install-tools` | Install development tools | `make install-tools` | `make.bat install-tools` |
| `quick-start` | Complete setup for new developers | `make quick-start` | `make.bat quick-start` |
| `clean` | Clean build artifacts | `make clean` | `make.bat clean` |
| `env-info` | Show environment information | `make env-info` | `make.bat env-info` |

## 🗄️ Database Migration

### GORM Auto-Migration

This project uses **GORM's auto-migration** instead of SQL files:

```go
// Run migrations
go run cmd/migrate/main.go

// Fresh migration (drop + migrate + seed)
go run cmd/migrate/main.go -fresh

// Drop tables and migrate
go run cmd/migrate/main.go -drop
```

### Migration Features

| Feature | Description |
|---------|-------------|
| **Auto-Migration** | GORM automatically creates/updates tables |
| **Model-Based** | Migrations based on Go struct models |
| **Index Creation** | Automatic index creation for performance |
| **Constraint Addition** | Custom business rule constraints |
| **Seeder Integration** | Automatic seeding after migration |

### Migration Command Options

| Flag | Description | Example |
|------|-------------|---------|
| `-drop` | Drop all tables before migration | `go run cmd/migrate/main.go -drop` |
| `-seed` | Run seeders after migration | `go run cmd/migrate/main.go -seed` |
| `-fresh` | Drop, migrate, and seed | `go run cmd/migrate/main.go -fresh` |

## 🔌 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/users/register` | User registration | Public |
| `POST` | `/api/v1/auth/login` | User login | Public |
| `PUT` | `/api/v1/users/verify-payment/:id` | Verify payment | Public |

### Koperasi Management

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/koperasi` | Create koperasi | SuperAdmin |
| `GET` | `/api/v1/koperasi` | List koperasi | Authenticated |
| `GET` | `/api/v1/koperasi/:id` | Get koperasi details | Authenticated |
| `PUT` | `/api/v1/koperasi/:id` | Update koperasi | Admin |

### Product Management

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/produk` | Create product | Admin |
| `GET` | `/api/v1/produk/:koperasi_id` | List products | Authenticated |
| `POST` | `/api/v1/produk/purchase-order` | Create purchase order | Admin |
| `POST` | `/api/v1/produk/penjualan` | Create sales transaction | Authenticated |

### Financial Management

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/financial/jurnal` | Create journal entry | Financial |
| `GET` | `/api/v1/financial/:id/neraca-saldo` | Trial balance | Financial |
| `GET` | `/api/v1/financial/:id/laba-rugi` | Profit & loss | Financial |

## 🏗️ Modular Routes Architecture

Routes are organized in a modular structure for better maintainability:

```
internal/routes/
├── routes.go              # Main orchestrator
└── modules/               # Domain-specific modules
    ├── auth_routes.go     # Authentication & payments
    ├── koperasi_routes.go # Koperasi management
    ├── produk_routes.go   # Product management
    ├── financial_routes.go # Financial operations
    └── ...                # Other domain routes
```

### Benefits

| Benefit | Description |
|---------|-------------|
| **Separation of Concerns** | Each domain has its own route file |
| **Maintainability** | Easy to locate and modify endpoints |
| **Scalability** | Easy to add new domain modules |
| **Team Collaboration** | Reduces conflicts when working on different features |

## 🔐 Authentication & Authorization

### Role-Based Access Control (RBAC)

| Role | Permissions | Access Level |
|------|-------------|--------------|
| **SuperAdmin** | Full system access | All operations |
| **Admin** | Koperasi management | Koperasi-specific |
| **Financial** | Financial operations | Financial modules |
| **User** | Basic operations | Limited access |
| **Operator** | Data entry | Specific modules |

## 📊 Business Features

### Indonesian Compliance

| Feature | Description | Implementation |
|---------|-------------|----------------|
| **NIAK Generation** | Automatic cooperative ID | Algorithm-based |
| **NIK Validation** | Indonesian ID validation | 16-digit validation |
| **Regional Data** | Complete Indonesian regions | Provinsi → Kelurahan |
| **KBLI Integration** | Business classification | Standard compliance |

### Product Management

| Feature | Description | Benefits |
|---------|-------------|----------|
| **12 Product Categories** | Food, beverages, livestock, etc. | Organized inventory |
| **Barcode Support** | EAN-13 generation | Efficient tracking |
| **Supplier Management** | Multi-supplier support | Cost optimization |
| **Perishable Tracking** | Expiry date management | Waste reduction |
| **Purchase Orders** | Complete procurement workflow | Organized purchasing |
| **Sales Transactions** | POS-style sales processing | Easy transactions |
| **Stock Movement** | Real-time inventory tracking | Accurate stock levels |

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/services/... -v

# Run with coverage
go test ./... -cover
```

## 🚀 Deployment

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `koperasi_db` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `JWT_SECRET` | JWT signing key | `your-secret-key` |

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.

## 📞 Support

- 📧 Email: support@example.com
- 🐛 Issues: GitHub Issues
- 💬 Discussions: GitHub Discussions