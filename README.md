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
  dsn: "app-user:app-pwd@tcp(mysql:3306)/payment_registration_system?charset=utf8mb4&parseTime=True&loc=Local"
  clean: true

nosqldb:
  uri: "mongodb://app-user:app-pwd@mongodb:27017/payment_registration_system"
  database: "payment_registration_system"
  clean: true

app:
  port: 9000
  read_timeout: 15
  write_timeout: 15
  graceful_shutdown: 15
  log_path: "payment_system.log"
  is_production: false
```

3️⃣ **Run the application**

```bash
go run main.go --config=config.yml
```

📌 The API will be accessible at: **`http://localhost:<PORT>`**

> [!TIP]
> Make sure to have a MySQL and MongoDB instance running locally!

---

### 🐳 Running with Docker

1️⃣ **Build and run using Docker Compose**

```bash
docker build -t go-app .
```

Then, run the application:

```bash
docker-compose up -d --build
```

📌 The API will be available at: **`http://go-app.localhost/`**

> [!TIP]
> We create an auxiliary script if you want populate the database with some data. Run `bash src/internal/testutils/populate.sh` to do it. This is test data, so you can use it to test the API.

---

## 📡 API Endpoints

> [!NOTE]
> For each endpoint, you can choose between SQL or NoSQL storage by changing the URL path. For SQL, use `/v1/sql/` and for NoSQL, use `/v1/nosql/`.

### ✅ Bank group

- **GET** `<STORAGE>/customers/count` – Retrieves the number of customers associated with each bank.
- **POST** `<STORAGE>/promotions/add-promotion/` – Adds a new financing promotion using the request body data.
- **DELETE** `<STORAGE>/promotions/discount/{code}` – Deletes a discount promotion identified by its code.
- **PATCH** `<STORAGE>/promotions/discount/{code}` – Updates the expiration date of a discount promotion identified by its code.
- **DELETE** `<STORAGE>/promotions/financing/{code}` – Deletes a financing promotion identified by its code.
- **PATCH** `<STORAGE>/promotions/financing/{code}` – Updates the expiration date of a financing promotion identified by its code.

### ✅ Card group

- **GET** `<STORAGE>/cards/expiring-next-30-days/{month}/{year}` – Retrieves the cards that will expire in the given month and year.
- **GET** `<STORAGE>/cards/payment-summary/{cardNumber}/{month}/{year}` – Retrieves the payment summary for the given month and year.
- **GET** `<STORAGE>/cards/purchase-monthly/{cuit}/{finalAmount}/{paymentVoucher}` – Retrieves the purchase details for a given CUIT, final amount, and payment voucher.
- **GET** `<STORAGE>/cards/top` – Retrieves the top 10 cards with the highest usage.

### ✅ Promotion & Store group

- **GET** `<STORAGE>/stores/highest-revenue/{month}/{year}` – Retrieves the stores with the highest revenue for the given month and year.
- **GET** `<STORAGE>/promotions/available/{cuit}/{startDate}/{endDate}` – Retrieves the financing and discount promotions available for a store between the specified start and end dates.
- **GET** `<STORAGE>/promotions/most-used` – Retrieves the most used promotions.

---

## 📜 License

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](LICENSE) file for more details.

---

🌟 **Contributions & Feedback**

Feel free to **fork, contribute, or submit issues** to help improve this project! 🚀
