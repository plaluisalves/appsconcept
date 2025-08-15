# FizzBuzz API with Metrics

A go based REST API that implements the classic FizzBuzz algorithm with comprehensive metrics tracking and OpenAPI 3.0 documentation. Built with the Gin web framework and featuring both custom and prometheus metrics support.

---

## Features

### Core Functionality
- **FizzBuzz Algorithm**: Generates sequences where numbers divisible by specific integers are replaced with custom strings
- **Flexible Parameters**: Supports custom divisors (int1, int2) and replacement strings (str1, str2)
- **Configurable Limits**: Generate sequences up to a limit

### Metrics & Monitoring
- **Dual Metrics Support**: Choose between custom in-memory metrics or Prometheus integration

### API Documentation
- **OpenAPI 3.0 (Swagger)**: Auto-generated API documentation

---

## How it Works

The FizzBuzz algorithm generates a sequence from 1 to the specified limit, replacing numbers that are:
- Divisible by `int1` with `str1`
- Divisible by `int2` with `str2`
- Divisible by both with the concatenation of `str1` + `str2`

For example, with `int1=3`, `int2=5`, `str1="Fizz"`, `str2="Buzz"`, `limit=15`:
```
1, 2, Fizz, 4, Buzz, Fizz, 7, 8, Fizz, Buzz, 11, Fizz, 13, 14, FizzBuzz
```

The API tracks metrics for each unique combination of parameters.

---

## Project Structure

```
APPSCONCEPT/
├── cmd/
│   └── web/
│       └── main.go              # Application entry point
├── internal/
│   ├── app.go                   # Main application setup and routing
│   ├── constants.go             # API endpoint constants
│   ├── http_handler_fizzbuzz.go # FizzBuzz endpoint handler
│   ├── htpp_handler_metrics.go  # Metrics endpoint handler
│   └── metrics/
│       ├── metrics.go           # Metrics service implementation
│       └── template.html        # Metrics dashboard template
├── tests/
│   ├── fake/                    # Fake implementations for testing
│   └── mocks/                   # Mock implementations for testing
├── utils/
│   ├── loggin.go               # Logging utilities
│   └── logwriter.go            # Custom log writer
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
└── makefile                    # Build and development commands
```

---

## API Endpoints

### GET /fizzbuzz
Generates a FizzBuzz sequence based on query parameters.

**Query Parameters:**
- `int1` (required): First divisor
- `int2` (required): Second divisor  
- `limit` (required): Sequence length limit
- `str1` (required): String to replace numbers divisible by int1
- `str2` (required): String to replace numbers divisible by int2

**Example Request:**
```
GET /fizzbuzz?int1=3&int2=5&limit=15&str1=Fizz&str2=Buzz
```

**Example Response:**
```json
["1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"]
```

### GET /metrics
- Returns a real-time dashboard showing metrics for all FizzBuzz requests.
- Shows total counts for each unique parameter combination

### GET /swagger
- Interactive OpenAPI documentation with Swagger UI.

---

## How to start

### Requirements
- Go 1.24.1 or higher

### Core Dependencies
- **Gin**: HTTP web framework
- **Tonic**: OpenAPI documentation generator
- **Prometheus**: Metrics collection (optional)
- **Testify**: Testing framework

### Run 

```bash
make run
```

The API will be available at `http://localhost:8080`

---

## Metrics

The application tracks request metrics in two ways:

* **Default (Custom):** In-memory storage with HTML dashboard
```go
a.metrics = metrics.New(metrics.WithCustomSolution())
```
* **Prometheus:** Standard Prometheus metrics for monitoring systems
```go
a.metrics = metrics.New(metrics.WithPrometheus())
```

Both track the same data: request counts for different parameter combinations and must not be used simultaneously.












