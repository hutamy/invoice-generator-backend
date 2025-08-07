# 🧾 Invoice Generator for Freelancers

A self-hosted, developer-friendly backend API for freelancers to manage clients and generate invoices (with PDF export). Built with **Go**, **PostgreSQL**, **Echo**, and **Swagger**.

---

## 🚀 Features

- 🧑‍💼 **User Authentication** (JWT)
- 👥 **Client Management** (CRUD)
- 💸 **Invoice Management** (CRUD)
- 📄 **PDF Invoice Generation** using HTML templates
- 🧾 **Swagger/OpenAPI Docs**
- 🛡️ Secure & modular architecture (repository + service layers)
- 🆓 **Public Invoice Generator** (no login, instant PDF generation without data storage)
- 🔐 **Authenticated Mode** for saving invoices and clients

---

## ⚙️ Setup

### 1. Clone the repo

```bash
git clone https://github.com/hutamy/invoice-generator-backend.git
cd invoice-generator
```

### 2. Set up .env

```
cp .env.example .env
# fill in DB, JWT_SECRET
```

### 3. Run with docker compose

```
docker-compose up --build
```

## 📚 API Documentation

Visit: `http://localhost:8080/swagger/index.html`

## 💡 Project Structure

```
├── cmd/                  # Main application entrypoint
├── config/               # Configuration files and helpers
├── controllers/          # HTTP handlers
├── docs/                 # Swagger/OpenAPI docs
├── middleware/           # Middleware for JWT
├── models/               # GORM models
├── repositories/         # DB access layer
├── routes/               # HTTP routes
├── services/             # Business logic
├── templates/            # HTML templates for PDF invoices
├── utils/                # Shared utilities and packages
├── scripts/              # Helper scripts (e.g., DB migrations)
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🧭 How It Works

- **Public Mode**: Anyone can POST invoice data and receive a PDF (no auth required)
- **Authenticated Mode**: Logged-in users can save clients, manage invoices, and view history

## 🧪 Example Usage

### Sign Up (Register)

```bash
curl --location 'http://localhost:8080/v1/public/auth/sign-up' \
--header 'Content-Type: application/json' \
--header 'X-API-Key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDk3MzI5MzcsImlhdCI6MTc0OTQ3MzczNywidXNlcl9pZCI6Mn0.YOw9hjXh_Llj11bWqOMnAkwSWb2TdH23ppsSX7g3aPo' \
--data-raw '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "yourpassword",
    "adderss": "some street",
    "phone_number": "1234567890",
    "bank_name": "bank name",
    "bank_account_name": "Jane Doe",
    "bank_account_number": "1234567890"
}'
```

### Sign In

```bash
curl --location 'http://localhost:8080/v1/public/auth/sign-in' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "jane@example.com",
    "password": "yourpassword",
}'
```

### Me

```bash
curl --location 'http://localhost:8080/v1/protected/me' \
--header 'Authorization: Bearer <token>'
```

### Refresh Token

```bash
curl --location 'http://localhost:8080/v1/protected/auth/refresh-token' \
--header 'Authorization: Bearer <token>'
--data-raw '{
    "refresh_token": <token>
}'
```

### Create Client

```bash
curl --location 'http://localhost:8080/v1/protected/clients' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <token>' \
--data-raw '{
  "name": "Client Name",
  "email": "client@example.com",
  "phone": "1234567890",
  "address": "some street"
}'
```

### Get All Clients

```bash
curl --location 'http://localhost:8080/v1/protected/clients' \
--header 'Authorization: Bearer <token>' \
```

### Get Client By ID

```bash
curl --location 'http://localhost:8080/v1/protected/clients/1' \
--header 'Authorization: Bearer <token>' \
```

### Update Client

```bash
curl --location --request PUT 'http://localhost:8080/v1/protected/clients/1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <token>' \
--data-raw '{
  "name": "Client Name",
  "email": "client@example.com",
  "phone": "1234567890",
  "address": "some street"
}'
```

### Delete Client

```bash
curl --location --request DELETE 'http://localhost:8080/v1/protected/clients/1' \
--header 'Authorization: Bearer <token>'
```

### Create Invoice

```bash
curl --location 'http://localhost:8080/v1/protected/invoices' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <token>' \
--data '{
    "client_id": 1,
    "invoice_number": "INV 30/VI/2025",
    "due_date": "2025-06-30",
    "notes": "Make payment befor 30 days",
    "tax_rate": 10,
    "items": [
        {
            "description": "Item Description",
            "quantity": 1,
            "unit_price": 100
        }
    ]
}'
```

### Get All Invoices

```bash
curl --location 'http://localhost:8080/v1/protected/invoices' \
--header 'Authorization: Bearer <token>' \
```

### Get Invoice By ID

```bash
curl --location 'http://localhost:8080/v1/protected/invoices/1' \
--header 'Authorization: Bearer <token>' \
```

### Update Invoice

```bash
curl --location --request PUT 'http://localhost:8080/v1/protected/invoices/1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <token>' \
--data '{
    "client_id": 1,
    "invoice_number": "INV 30/VI/2025",
    "due_date": "2025-06-30",
    "notes": "Make payment befor 30 days",
    "tax_rate": 10,
    "status": "sent",
    "items": [
        {
            "id": 1,
            "description": "Item Description",
            "quantity": 1,
            "unit_price": 100
        }
    ]
}'
```

### Update Invoice Status

```bash
curl --location --request PATCH 'http://localhost:8080/v1/protected/invoices/1/status' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <token>' \
--data '{
    "status": "sent"
}'
```

### Delete Invoice

```bash
curl --location --request DELETE 'http://localhost:8080/v1/protected/invoices/1' \
--header 'Authorization: Bearer <token>'
```

### Generate PDF

```bash
curl --location --request POST 'http://localhost:8080/v1/protected/invoices/1/pdf' \
--header 'Authorization: Bearer <token>'
```

### Generate Public PDF

```bash
curl --location 'http://localhost:8080/v1/public/invoices/generate-pdf' \
--header 'Content-Type: application/json' \
--data-raw '{
    "invoice_number": "INV 30/VI/2025",
    "due_date": "2025-06-30",
    "notes": "Make payment befor 30 days",
    "issue_date": "2025-06-01",
    "sender": {
        "name": "Jane Doe",
        "email": "jane@example.com",
        "adderss": "some street",
        "phone_number": "1234567890",
        "bank_name": "bank name",
        "bank_account_name": "Jane Doe",
        "bank_account_number": "1234567890"
    },
    "recipient": {
        "name": "Client Name",
        "email": "client@example.com",
        "phone": "1234567890",
        "address": "some street"
    },
    "items": [
        {
          "description": "Item Description",
          "quantity": 1,
          "unit_price": 100
        }
    ],
    "tax_rate": 10
}'
```
