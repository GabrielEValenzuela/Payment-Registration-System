# Payment Registration System

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white) ![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white) ![MongoDB](https://img.shields.io/badge/MongoDB-4ea94b?style=for-the-badge&logo=mongodb&logoColor=white) ![Docker](https://img.shields.io/badge/Docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Fiber](https://img.shields.io/badge/Fiber-000000?style=for-the-badge&logo=fiber&logoColor=white)

---

## What is Payment Registration System? 🤔

**Payment Registration System** is a robust API designed to handle payment processing, transaction management, and promotional campaigns using **Go**, **MySQL**, and **MongoDB**. It is built with **Fiber** for high-performance API management and deployed using **Docker** for scalability and reliability.

### 🚀 Why Payment Registration System?

- **Hybrid SQL & NoSQL storage:** MySQL for structured financial data and MongoDB for flexible transactions.
- **High-performance API:** Built with Fiber for fast request processing.
- **Scalable architecture:** Dockerized environment for deployment flexibility.
- **Rate-limiting & Security:** Prevents excessive API requests and secures transaction data.

---

## ✨ Features

✔️ **Manage Promotions** – Add, extend, and delete financing promotions.
✔️ **Transaction Handling** – Store purchase information including single and installment payments.
✔️ **Payment Summarization** – Generate total payment summaries per month.
✔️ **Card Management** – Retrieve expiring cards and top usage statistics.
✔️ **Store Analytics** – Identify the highest revenue-generating stores.
✔️ **Comprehensive API** – RESTful endpoints for financial data access and management.

---

## 🛠 Technologies Used

- 🐹 **Go (Golang)** – Core API implementation.
- ⚡ **Fiber** – High-performance web framework.
- 🐘 **MySQL** – Structured data storage for transactions.
- 🍃 **MongoDB** – NoSQL storage for flexible payment data.
- 🐳 **Docker** – Containerized deployment.

---

## 🚀 Installation and Setup

### 🔧 Local Setup

1️⃣ **Clone the repository**

```bash
git clone https://github.com/GabrielEValenzuela/Payment-Registration-System.git
cd Payment-Registration-System
```

2️⃣ **Set up environment variables**

Create a `config.yml` file with:

```yml
sqldb:
  dsn: "testuser:testpassword@tcp(127.0.0.1:3306)/payment-registration-db?charset=utf8mb4&parseTime=True&loc=Local"
  clean: true

nosqldb:
  uri: "mongodb://testuser:testpassword@localhost:27017/payment_registration_system?authSource=admin"
  database: "payment_registration_system"
  clean: true

app:
  port: 8080
  read_timeout: 15
  write_timeout: 15
  graceful_shutdown: 15
  log_path: "debug_payment.log"
  is_production: false
```

3️⃣ **Run the application**

```bash
go run main.go --config=config.yml
```

📌 The API will be accessible at: **`http://localhost:8080`**

> [!TIP]
> Make sure to have a MySQL and MongoDB instance running locally!

---

### 🐳 Running with Docker

1️⃣ **Build and run using Docker Compose**

```bash
docker-compose up -d --build
```

📌 The API will be available at: **`http://localhost:8080`**

---

## 📡 API Endpoints

### ✅ Promotion Management
- **POST** `/v1/sql/promotions/financing` – Add a financing promotion to a bank.
- **PATCH** `/v1/sql/promotions/financing/{code}` – Extend the validity of a financing promotion.
- **PATCH** `/v1/sql/promotions/discount/{code}` – Extend the validity of a discount promotion.
- **DELETE** `/v1/sql/promotions/financing/{code}` – Remove a financing promotion.
- **DELETE** `/v1/sql/promotions/discount/{code}` – Remove a discount promotion.

### ✅ Payment & Card Handling
- **POST** `/v1/sql/cards/summary/` – Get the total payment summary for a card.
- **POST** `/v1/sql/cards/expiring/` – Retrieve expiring cards within 30 days.
- **GET** `/v1/sql/cards/top` – Get the 10 most-used cards.
- **POST** `/v1/sql/cards/purchase/monthly/` – Retrieve purchase information including installments.

### ✅ Store & Analytics
- **GET** `/v1/sql/stores/highest-revenue/{month}/{year}` – Get the store with the highest revenue.
- **GET** `/v1/sql/promotions/most-used` – Retrieve the most used promotion.
- **GET** `/v1/sql/banks/customers/count` – Get the number of customers per bank.

---

## 📜 License

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](LICENSE) file for more details.

---

🌟 **Contributions & Feedback**

Feel free to **fork, contribute, or submit issues** to help improve this project! 🚀

