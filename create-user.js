// MongoDB script to create a user
// Run with: mongosh mongodb://localhost:27017/onetake-corpsite-backend create-user.js

const bcrypt = require('bcryptjs');

db.users.insertOne({
  email: "email",
  passwordHash: bcrypt.hashSync("pw", 10),
  roles: ["admin"],
  createdAt: new Date(),
  updatedAt: new Date()
});

print("User created successfully!");
