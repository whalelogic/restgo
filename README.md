
# Employee API 👨‍💼👩‍⚕️👨🏽‍🏭

A lightweight RESTful API for managing employee records built with **Go** (standard library `net/http`) and **Redis**.

Repository: `github.com/whalelogic/restgo`

## Environment Variables

| Variable | Description | Default Value |
| :--- | :--- | :--- |
| `REDIS_ADDR` | Network address (host:port) of the Redis instance | `localhost:6379` |

### Setting the Environment Variable (macOS & Linux)
```bash
# Inline execution
REDIS_ADDR="127.0.0.1:6379" go run main.go

# Session-wide export
export REDIS_ADDR="127.0.0.1:6379"

```

## API Endpoints Reference

The service listens on port `:8080`.

**Method**

**Endpoint**

**Description**

`GET`

`/employees`

List all employees

`GET`

`/employees/{id}`

Get employee by ID

`POST`

`/employees`

Create a new employee

`PUT`

`/employees/{id}`

Fully update/replace an employee

`PATCH`

`/employees/{id}`

Partially update an employee

`DELETE`

`/employees/{id}`

Delete an employee

## cURL Examples

### 1. Get All Employees

Bash

```
curl -X GET http://localhost:8080/employees

```

### 2. Get Employee by ID

Bash

```
curl -X GET http://localhost:8080/employees/emp_01J8

```

### 3. Create Employee

Bash

```
curl -X POST http://localhost:8080/employees \
  -H "Content-Type: application/json" \
  -d '{"name": "Alex Smith", "role": "DevOps", "email": "alex@example.com"}'

```

### 4. Update Employee (Full Replacement)

Bash

```
curl -X PUT http://localhost:8080/employees/emp_01J8 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alex Smith", "role": "Lead DevOps", "email": "alex@example.com"}'

```

### 5. Patch Employee (Partial Update)

Bash

```
curl -X PATCH http://localhost:8080/employees/emp_01J8 \
  -H "Content-Type: application/json" \
  -d '{"role": "Principal DevOps"}'

```

### 6. Delete Employee

Bash

```
curl -X DELETE http://localhost:8080/employees/emp_01J8

```

## Internal Routing Architecture

The routing maps directly to the `employee.EmployeeHandler` struct via Go's native standard library routing:

Go

```
mux := http.NewServeMux()
mux.HandleFunc("GET /employees", h.List)
mux.HandleFunc("GET /employees/{id}", h.Get)
mux.HandleFunc("POST /employees", h.Create)
mux.HandleFunc("PUT /employees/{id}", h.Update)
mux.HandleFunc("PATCH /employees/{id}", h.Patch)
mux.HandleFunc("DELETE /employees/{id}", h.Delete)

```

