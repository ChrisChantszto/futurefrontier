// MongoDB script to create a user directly
// Run with: mongosh mongodb://localhost:27017/onetake-corpsite-backend init-user.js

// Hash password function (since we can't use bcrypt in mongosh script)
function hashPassword(password) {
  // This is a simple hash for demonstration - in production use bcrypt
  return `$2a$10$simpleHashFor${password}`;
}

// Create user with email chantszto.chris@gmail.com
db.users.insertOne({
  email: "chantszto.chris@gmail.com",
  passwordHash: hashPassword("password123"),
  roles: ["admin"],
  createdAt: new Date(),
  updatedAt: new Date()
});

print("User created successfully!");
