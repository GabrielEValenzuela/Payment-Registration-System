import random
import string
from datetime import datetime, timedelta
import pymysql
from pymongo import MongoClient

AMOUNT_BANKS = 10
AMOUNT_CUSTOMERS = 50
AMOUNT_DISCOUNTS = 5

# Configuration
DB_CONFIG = {
    "mysql": {
        "host": "localhost",
        "user": "testuser",  # Replace with test user
        "password": "testpassword",  # Replace with test password
        "database": "payment_registration_system",
    },
    "mongodb": {
        "host": "localhost",
        "port": 27017,
        "database": "payment_registration_system",
        "host": "localhost",
        "username": "testuser",
        "password": "testpassword",
        "authSource": "payment_registration_system",
    },
}

# Sample Data
BANK_NAMES = ["Galicia", "BBVA", "Galicia Mas", "Santander"]
CUSTOMER_NAMES = [
    "German Castiglione",
    "Jeremias Herrera",
    "Ignazio March",
    "Ramiro Gershkovich",
    "Graciana Endrizzi",
    "Yara Fanucci",
    "Crisanna Toledo",
    "Carilla Benítez",
    "Vicenta Gimenez",
    "Eliana Pereyra"]
PROMOTION_TITLES = ["10% Off", "Zero Interest", "Cashback 2%"]
STORE_NAMES = ["Falabela", "COTO", "Kiosko 365", "Tecno Compro"]


def connect_mysql():
    return pymysql.connect(**DB_CONFIG["mysql"])


def connect_mongo():
    client = MongoClient(
        host=DB_CONFIG["mongodb"]["host"],
        port=DB_CONFIG["mongodb"]["port"],
        username=DB_CONFIG["mongodb"]["username"],
        password=DB_CONFIG["mongodb"]["password"],
        authSource=DB_CONFIG["mongodb"]["authSource"],
    )
    return client[DB_CONFIG["mongodb"]["database"]]


def generate_random_string(length=10):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))


def populate_banks(mysql_conn, mongo_db, num_banks=10):
    print("Populating Banks...")
    cursor = mysql_conn.cursor()
    for i in range(num_banks):
        name = random.choice(BANK_NAMES)
        cuit = f"{30}{str(i + 1).zfill(11)}{random.randint(0, 9)}"
        address = f"{random.randint(1, 999)} Bank Street"
        telephone = f"{random.randint(1000000000, 9999999999)}"
        created_at = datetime.now()

        # Insert into MySQL
        cursor.execute(
            "INSERT INTO BANKS (name, cuit, address, telephone, created_at) VALUES (%s, %s, %s, %s, %s)",
            (name, cuit, address, telephone, created_at),
        )

        # Insert into MongoDB
        mongo_db.banks.insert_one({
            "name": name,
            "cuit": cuit,
            "address": address,
            "telephone": telephone,
            "created_at": created_at,
        })

    mysql_conn.commit()
    print(f"Inserted {num_banks} banks.")


def populate_customers(mysql_conn, mongo_db, num_customers=50):
    print("Populating Customers...")
    cursor = mysql_conn.cursor()
    for i in range(num_customers):
        name = random.choice(CUSTOMER_NAMES)
        dni = str(random.randint(4000000, 50000000))
        cuit = f"{random.randint(20, 27)}{dni}{random.randint(0, 9)}"
        address = f"{random.randint(1, 999)} Customer Avenue"
        telephone = f"{random.randint(1000000000, 9999999999)}"
        entry_date = datetime.now()
        created_at = datetime.now()
        updated_at = created_at

        # Insert into MySQL
        cursor.execute(
            "INSERT INTO CUSTOMERS (complete_name, dni, cuit, address, telephone, entry_date, created_at, updated_at) VALUES (%s, %s, %s, %s, %s, %s,%s, %s)",
            (name, dni, cuit, address, telephone,
             entry_date, created_at, updated_at),
        )

        # Insert into MongoDB
        mongo_db.customers.insert_one({
            "complete_name": name,
            "dni": dni,
            "cuit": cuit,
            "address": address,
            "telephone": telephone,
            "entry_date": entry_date,
            "created_at": created_at,
            "updated_at": updated_at,
        })

    mysql_conn.commit()
    print(f"Inserted {num_customers} customers.")


def populate_discounts(mysql_conn, mongo_db, num_discounts=5):
    print("Populating Discounts...")
    cursor = mysql_conn.cursor()

    for i in range(num_discounts):
        code = str('DIS-') + str(i).zfill(3)
        # Generate random values for discount data
        title = random.choice(PROMOTION_TITLES)
        # Random percentage (e.g., 5.0% to 50.0%)
        discount_percentage = round(random.uniform(5.0, 50.0), 2)
        store = random.choice(STORE_NAMES)
        cuit_store = f"{30}{str(i + 5).zfill(11)}{random.randint(0, 9)}"
        start_date = datetime.now()
        end_date = start_date + timedelta(days=random.randint(1, 30))
        price_cap = round(random.uniform(100.0, 1000.0), 2)  # Random price cap
        only_cash = random.choice([True, False])
        comments = f"Discount for {store}"
        created_at = datetime.now()
        updated_at = created_at
        is_deleted = random.choice([True, False])
        bank_id = random.randint(1, AMOUNT_BANKS)

        # Insert into MySQL
        cursor.execute(
            """
            INSERT INTO DISCOUNTS (
                code, promotion_title, discount_percentage, name_store, cuit_store, 
                validity_start_date, validity_end_date, price_cap, only_cash, 
                comments, bank_id, created_at, updated_at, is_deleted
            ) 
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """,
            (
                code, title, discount_percentage, store, cuit_store,
                start_date, end_date, price_cap, only_cash,
                comments, bank_id, created_at, updated_at, is_deleted
            ),
        )

        # Insert into MongoDB
        mongo_db.discounts.insert_one({
            "code": code,
            "promotion_title": title,
            "discount_percentage": discount_percentage,
            "name_store": store,
            "cuit_store": cuit_store,
            "validity_start_date": start_date,
            "validity_end_date": end_date,
            "price_cap": price_cap,
            "only_cash": only_cash,
            "comments": comments,
            "created_at": created_at,
            "updated_at": updated_at,
            "is_deleted": is_deleted,
            "bank_id": bank_id,
        })

    # Commit the MySQL transaction
    mysql_conn.commit()
    print(f"Inserted {num_discounts} discounts.")


def populate_banks_customers(mysql_conn, mongo_db):
    print("Populating Banks, Customers, and Relationships...")
    cursor = mysql_conn.cursor()

    # Predefined bank and customer IDs
    bank_ids = list(range(1, AMOUNT_BANKS + 1))
    customer_ids = list(range(1, AMOUNT_CUSTOMERS + 1))

    # Populate relationships in customers_banks
    # Use exactly AMOUNT_CUSTOMERS
    for customer_id in random.sample(customer_ids, AMOUNT_CUSTOMERS):
        # Each customer relates to 1-3 banks
        assigned_banks = random.sample(bank_ids, random.randint(1, 3))
        for bank_id in assigned_banks:
            # Insert into MySQL
            cursor.execute(
                """
                INSERT INTO customers_banks (customer_entity_sql_id, bank_entity_sql_id) 
                VALUES (%s, %s)
                """,
                (customer_id, bank_id)
            )

    customers = list(mongo_db.customers.find({}, {"_id": 1}))
    banks = list(mongo_db.banks.find({}, {"_id": 1}))

    for _ in range(AMOUNT_CUSTOMERS):
        customer = random.choice(customers)["_id"]
        bank = random.choice(banks)["_id"]

        mongo_db.customers_banks.insert_one({
            "customer_id": customer,
            "bank_id": bank
        })

    # Commit the MySQL transaction
    mysql_conn.commit()
    print(f"Inserted relationships for {
          len(customer_ids)} customers and {len(bank_ids)} banks.")


def main():
    # Connect to databases
    mysql_conn = connect_mysql()
    mongo_db = connect_mongo()

    # Populate data
    populate_banks(mysql_conn, mongo_db, num_banks=10)
    populate_customers(mysql_conn, mongo_db, num_customers=50)
    populate_discounts(mysql_conn, mongo_db, num_discounts=5)
    populate_banks_customers(mysql_conn, mongo_db)

    # Close connections
    mysql_conn.close()
    print("Data population complete.")


if __name__ == "__main__":
    main()
