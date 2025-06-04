# 🧾 Invoice Generator for Freelancers

A backend service that helps freelancers and solo entrepreneurs manage their invoices, send them to clients, and track payments — all in one place.

Built with Go, PostgreSQL, and clean architecture principles.

---

## 🔧 Features

- **Authentication**
  - Sign-up / login with JWT
- **Client Management**
  - Add, edit, and delete clients
- **Invoice Management**
  - Create, send, and manage invoices
  - Status tracking: Draft, Sent, Paid, Overdue
- **PDF Generation**
  - Download invoice as branded PDF
- **Email Integration**
  - Send invoice via email (SendGrid / Mailgun)
- **Payment Integration**
  - Accept payments via Stripe or PayPal
- **Dashboard (optional)**
  - View total income, unpaid invoices, overdue status

---

## 🧱 Tech Stack

| Layer            | Technology                     |
| ---------------- | ------------------------------ |
| Language         | Go (Golang)                    |
| Framework        | Echo                           |
| Database         | PostgreSQL                     |
| ORM / SQL Mapper | GORM                           |
| Auth             | JWT                            |
| PDF Generation   | gofpdf / pdfcpu                |
| Email Service    | SendGrid / Mailgun             |
| Payments         | Stripe / PayPal                |
| Observability    | OpenTelemetry + Datadog        |
| Deployment       | Docker, GitHub Actions (CI/CD) |

---

## 📐 Architecture

The project follows **Clean Architecture** with a **Modular Monolith** approach:
