db = db.getSiblingDB('payment_registration_system');

db.createUser({
  user: "app-user",
  pwd: "app-pwd",
  roles: [{ role: "readWrite", db: "payment_registration_system" }]
});

print("MongoDB user and database created.");