# 🧾 Invoice Generator for Freelancers

A self-hosted, developer-friendly backend API for freelancers to manage clients and generate invoices (with PDF export). Built with **Go**, **PostgreSQL**, **Echo**, **Swagger**, and **MinIO**.

---

## 🚀 Features

- 🧑‍💼 **User Authentication** (JWT)
- 👥 **Client Management** (CRUD)
- 💸 **Invoice Management** (CRUD, partial updates)
- 📄 **PDF Invoice Generation** using HTML templates
- 🧾 **Swagger/OpenAPI Docs**
- 🛡️ Secure & modular architecture (repository + service layers)

---

## 📦 Tech Stack

| Component     | Tool                    |
| ------------- | ----------------------- |
| Language      | Go (Golang)             |
| Web Framework | Echo                    |
| DB            | PostgreSQL (via GORM)   |
| PDF Engine    | chromedp + Go templates |
| Auth          | JWT                     |
| Docs          | Swag + Swagger UI       |
| Dev Tools     | Docker, docker-compose  |

---

## ⚙️ Setup

### 1. Clone the repo

```bash
git clone https://github.com/hutamy/invoice-generator.git
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
.
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
