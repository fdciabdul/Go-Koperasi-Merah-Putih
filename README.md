# Koperasi Merah Putih - Digital Cooperative Management System

A comprehensive digital platform for cooperative management in Indonesia, built with Go, PostgreSQL, and Apache Cassandra.

## Features

### Core Functionality
- **Multi-tenant Architecture**: Support for multiple cooperatives
- **User Registration with Payment**: Members must pay simpanan pokok before activation
- **PPOB (Payment Point Online Bank)**: Bill payment services
- **Savings & Loans Management**: Member savings and loan products
- **Clinic Management**: Healthcare services for members
- **Financial Management**: Chart of accounts and journaling
- **Payment Gateway Integration**: Midtrans and Xendit support

### Architecture
- **PostgreSQL**: Transactional data (GORM)
- **Cassandra**: Analytics and logging data (GoCQL)
- **Go Gin**: REST API framework
- **Multi-database**: Hybrid architecture for optimal performance

## Project Structure

```
koperasi-merah-putih/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── database/              # Database connections
│   ├── models/                # Data models
│   │   ├── postgres/          # PostgreSQL models (GORM)
│   │   └── cassandra/         # Cassandra models (GoCQL)
│   ├── repository/            # Data access layer
│   │   ├── postgres/          # PostgreSQL repositories
│   │   └── cassandra/         # Cassandra repositories
│   ├── services/              # Business logic layer
│   ├── handlers/              # HTTP handlers
│   ├── routes/                # Route definitions
│   ├── middleware/            # HTTP middleware
│   └── gateway/               # Payment gateway integrations
├── scripts/                   # Database scripts
├── .env.example              # Environment variables template
├── docker-compose.yml        # Docker services
├── Dockerfile               # Application container
└── Makefile                # Build and development commands
```

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Apache Cassandra 4.1+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd koperasi-merah-putih
   ```

2. **Install dependencies**
   ```bash
   make install-deps
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database configurations
   ```

4. **Start databases with Docker**
   ```bash
   make docker-compose-up
   ```

5. **Initialize databases**
   ```bash
   make postgres-init
   make cassandra-init
   ```

6. **Run the application**
   ```bash
   make run
   ```

### Using Docker

1. **Start all services**
   ```bash
   make docker-compose-up
   ```

2. **Build and run application**
   ```bash
   make docker-build
   make docker-run
   ```

## API Endpoints

### User Registration Flow
```
POST /api/v1/users/register
→ Creates user registration with payment requirement
→ Returns payment link for simpanan pokok

PUT /api/v1/users/verify-payment/{payment_id}
→ Verifies payment completion
→ Updates registration status to payment_verified

PUT /api/v1/users/registrations/{id}/approve
→ Admin approves registration
→ Creates anggota_koperasi and user records
```

### PPOB Services
```
GET /api/v1/ppob/kategoris
→ List PPOB categories

GET /api/v1/ppob/kategoris/{kategori_id}/produks
→ List products by category

POST /api/v1/ppob/transactions
→ Create PPOB transaction with payment
```

### Payment Callbacks
```
POST /api/v1/payments/midtrans/callback
POST /api/v1/payments/xendit/callback
→ Handle payment gateway callbacks
```

## Business Process Flows

### User Registration Flow
1. User fills registration form
2. System creates `user_registrations` record with status `pending_payment`
3. System creates `payment_transactions` record for simpanan pokok
4. User completes payment via gateway
5. Gateway sends callback, payment status updated to `paid`
6. System updates registration status to `payment_verified`
7. Admin approves registration
8. System creates `anggota_koperasi` and `users` records
9. System creates `simpanan_pokok_transaksi` record
10. System creates journal entries

### PPOB Transaction Flow
1. Customer selects PPOB product and fills form
2. System validates product and calculates total amount (price + admin fee)
3. System creates `payment_transactions` record
4. Customer completes payment via gateway
5. Gateway sends callback, payment status updated to `paid`
6. System processes PPOB transaction to provider
7. System updates `ppob_transaksi` status based on provider response
8. System creates journal entries for revenue recognition
9. System handles settlement based on configuration

## Database Schema

### PostgreSQL (Transactional Data)
- Multi-tenant tables with proper relationships
- Master data (wilayah, koperasi, anggota)
- Financial data (COA, jurnal, transaksi)
- PPOB and payment data
- User management and registration

### Cassandra (Analytics Data)
- Transaction logs for real-time analytics
- Performance metrics
- User activity logs
- Error logs
- Monthly aggregated facts

## Configuration

### Environment Variables
```bash
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=koperasi_user
POSTGRES_PASSWORD=koperasi_password
POSTGRES_DB=koperasi_merah_putih

# Cassandra
CASSANDRA_HOSTS=127.0.0.1
CASSANDRA_KEYSPACE=koperasi_analytics

# Payment Gateways
MIDTRANS_SERVER_KEY=your-midtrans-server-key
XENDIT_SECRET_KEY=your-xendit-secret-key

# PPOB Provider
PPOB_PROVIDER_URL=https://api.ppob-provider.com
PPOB_API_KEY=your-ppob-api-key
```

## Development Commands

```bash
# Development
make run                # Run application
make dev                # Run with hot reload
make test               # Run tests
make lint               # Run linter

# Database
make migrate-up         # Run migrations
make seed               # Seed database
make cassandra-init     # Initialize Cassandra schema

# Docker
make docker-compose-up  # Start all services
make docker-build       # Build application image
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please contact the development team or create an issue in the repository.