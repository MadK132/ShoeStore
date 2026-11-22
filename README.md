
# ShoeStore — Microservices Backend

A Go-based backend project for an online shoe store, organized using a microservices architecture.

The application separates user management, product management, order processing, and email functionality into dedicated services, with an API Gateway serving as the main entry point.

## Tech Stack

- **Language:** Go (Golang)
- **Architecture:** Microservices
- **Data Serialization:** Protocol Buffers
- **Database:** MongoDB
- **Development Tools:** PowerShell and Windows batch scripts

## Architecture

The project consists of an API Gateway and five dedicated backend components.

```text
                    Client
                      |
                      v
                 API Gateway
                      |
        +-------------+-------------+
        |             |             |
        v             v             v
   User Service  Product Service  Order Service
                                    |
                                    v
                               Email Service
```

*Conceptual overview of the service responsibilities.*

## Services

### 1. API Gateway

The `api-gateway` serves as the entry point for client requests and provides a centralized layer for accessing backend services.

### 2. User Service

The `user-service` is responsible for user-related operations, keeping user management logic separate from other business domains.

### 3. Product Service

The `product-service` handles product-related functionality within the shoe store application.

### 4. Order Service

The `order-service` contains order-related business logic and separates order processing from product and user management.

### 5. Email Service

The `email-service` provides a dedicated component for email-related functionality.

## Project Structure

```text
ShoeStore/
├── api-gateway/
├── email-service/
├── order-service/
├── product-service/
├── user-service/
├── proto/
├── generate-proto.ps1
├── start-mongodb.bat
├── start.bat
├── go.mod
├── go.sum
└── README.md
```

## Key Concepts

- **Microservices Architecture** — Divides application functionality into separate services.
- **Separation of Concerns** — Each service is responsible for a specific business domain.
- **API Gateway Pattern** — Provides a centralized entry point for backend requests.
- **Protocol Buffers** — Defines structured data contracts.
- **Modularity** — Organizes backend functionality into independently maintained components.

## Development

The repository includes scripts for local development:

- `start-mongodb.bat` — MongoDB startup script.
- `generate-proto.ps1` — Protocol Buffers code generation script.
- `start.bat` — Application startup script.

### Clone the Repository

```bash
git clone https://github.com/MadK132/ShoeStore.git
cd ShoeStore
```

### Install Go Dependencies

```bash
go mod download
```

## Project Purpose

This project demonstrates backend engineering concepts through a service-oriented e-commerce application, focusing on modular architecture, service boundaries, and structured communication contracts.
