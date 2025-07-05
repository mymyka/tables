# 📊 Tables

> **Type-safe PostgreSQL schema to Go code generator**

Transform your PostgreSQL database schemas into elegant, type-safe Go code with zero boilerplate.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/mymyka/tables)

---

## ✨ Features

- 🔒 **Type Safety** - Generate strongly-typed Go aliases from PostgreSQL columns
- 🚀 **Zero Boilerplate** - Automatic code generation with minimal configuration
- 📦 **Package Organization** - Clean, structured output with proper Go conventions
- 🎯 **Column Mapping** - Easy-to-use column name constants for query building
- 🔄 **Incremental Updates** - Regenerate only when schema changes
- 🛡️ **Null Safety** - Proper handling of nullable columns with pointer types

---

## 🚀 Quick Start

### Installation

```bash
go install github.com/mymyka/tables/cmd/tables@latest
```

### Basic Usage

```bash
tables --db "host=localhost port=5432 user=postgres password=postgres dbname=mydb sslmode=disable" --output gen/tables
```

---

## 📖 How It Works

### Input: PostgreSQL Schema
```sql
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    google_user_id  VARCHAR(100) UNIQUE,
    hashed_password VARCHAR(255)
);
```

### Output: Type-Safe Go Code
```go
package users

import "time"

// Type aliases for compile-time safety
type Id = int32
type Username = string
type FirstName = string
type LastName = string
type CreatedAt = time.Time
type Email = string
type GoogleUserId = *string      // Nullable
type HashedPassword = *string    // Nullable

// Column name constants
type usersColumnNames struct {
    Id             string
    Username       string
    FirstName      string
    LastName       string
    CreatedAt      string
    Email          string
    GoogleUserId   string
    HashedPassword string
}

var C = usersColumnNames{
    Id:             "id",
    Username:       "username",
    FirstName:      "first_name",
    LastName:       "last_name",
    CreatedAt:      "created_at",
    Email:          "email",
    GoogleUserId:   "google_user_id",
    HashedPassword: "hashed_password",
}

var Table = "users"
```

---

## 💡 Usage Examples

### Building Models
```go
import "your-project/gen/tables/users"

type UserModel struct {
    Id        users.Id
    Username  users.Username
    FirstName users.FirstName
    LastName  users.LastName
    Email     users.Email
    CreatedAt users.CreatedAt
}
```

### Type-Safe Queries
```go
// Using column constants prevents typos
query := fmt.Sprintf(`
    SELECT %s, %s, %s 
    FROM %s 
    WHERE %s = $1
`, 
    users.C.Id, 
    users.C.Username, 
    users.C.Email,
    users.Table,
    users.C.Username,
)

// Type-safe scanning
var user UserModel
err := db.QueryRow(query, "john_doe").Scan(
    &user.Id,
    &user.Username, 
    &user.Email,
)
```

### Working with Nullable Fields
```go
// Nullable fields are properly typed as pointers
var googleId users.GoogleUserId
if user.GoogleUserId != nil {
    googleId = *user.GoogleUserId
}
```

---

## ⚙️ Configuration

### Command Line Options

| Flag | Description | Required | Default |
|------|-------------|----------|---------|
| `--db` | PostgreSQL connection string | ✅ | - |
| `--output` | Output directory for generated code | ✅ | - |
| `--exclude` | Comma-separated list of tables to exclude | ❌ | - |
| `--include` | Comma-separated list of tables to include | ❌ | All tables |
| `--package-prefix` | Prefix for generated package names | ❌ | - |

### Connection String Format
```
host=localhost port=5432 user=username password=password dbname=database sslmode=disable
```

---

## 🏗️ Project Structure

After running Tables, your project structure will look like:

```
your-project/
├── gen/
│   └── tables/
│       ├── users/
│       │   └── users.go
│       ├── orders/
│       │   └── orders.go
│       └── products/
│           └── products.go
├── main.go
└── go.mod
```

---

## 🔧 Type Mapping

| PostgreSQL Type | Go Type | Nullable Go Type |
|----------------|---------|------------------|
| `SERIAL`, `INTEGER` | `int32` | `*int32` |
| `BIGSERIAL`, `BIGINT` | `int64` | `*int64` |
| `VARCHAR`, `TEXT` | `string` | `*string` |
| `BOOLEAN` | `bool` | `*bool` |
| `TIMESTAMP` | `time.Time` | `*time.Time` |
| `DATE` | `time.Time` | `*time.Time` |
| `DECIMAL`, `NUMERIC` | `float64` | `*float64` |
| `UUID` | `string` | `*string` |
| `JSONB` | `[]byte` | `*[]byte` |

---

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by the need for type-safe database interactions in Go
- Built with ❤️ for the Go community
- Special thanks to all contributors

---

<div align="center">
  <strong>Made with 📊 by developers, for developers</strong>
</div>