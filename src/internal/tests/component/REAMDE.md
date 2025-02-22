# 📂 Data Directory

This directory contains all the necessary datasets and scripts to populate both **MongoDB** and **MySQL** for the Payment Registration System.

---

## 📁 Structure

```plaintext
data/
├── mongo/   # JSON files for NoSQL data (MongoDB)
│   ├── banks.json
│   ├── cards.json
│   ├── customers.json
│   ├── discounts.json
│   ├── financings.json
│   ├── purchase_monthly_payments.json
│   ├── purchase_single_payments.json
│   ├── gen.go   # Go script to generate random test data
│
├── mysql/   # SQL scripts for relational database (MySQL)
│   ├── database.sql   # Schema and initial data
```

---

## 🛠 How run the component tests

```bash
bash run_tests.sh
```

---

## 📝 Notes

- **MongoDB Data**: Stored in JSON format, designed to be used with `mongoimport`.
- **MySQL Data**: The `database.sql` file includes schema definitions and sample data.
- **Go Data Generator**: `gen.go` creates randomized test data for MongoDB. It could be easy to modify for other purposes.

---

## 🚀 Contribution & Testing

For testing, ensure that both MySQL and MongoDB services are running before executing any imports. If you encounter issues, check for:

- Connection settings in your `.env` file.
- Proper container names in the commands.
- Permissions to access the files inside the containers.

---

## 🛡️ License

GNU General Public License v3.0. Feel free to contribute or modify as needed. 🚀
