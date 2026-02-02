# 🌍 Country Search API

**Country-Search-API** is a RESTful service written in Go that provides country information by leveraging the public **REST Countries API** (https://restcountries.com/). It demonstrates clean architecture practices, robust caching, concurrency safety, and comprehensive testing.

---

## ✨ Features

- 🔎 Search for countries by name
- ⚡ Custom in-memory caching to reduce external API calls
- 🧵 Safe concurrent access with race-condition protection
- ⏱ Context-based timeout handling
- 🛑 Graceful shutdown support
- 🧪 Extensive unit, integration, and race-condition testing

---

## 🏗 Architecture Highlights

This project demonstrates:

- Dependency management using handlers and a custom HTTP client
- Context propagation for request timeouts and cancellations
- Graceful server shutdown for production readiness
- Thread-safe cache design with concurrent access support

---

## 🧪 Testing Strategy

### Unit Tests

- Cache implementation  
- Custom HTTP client  
- Service layer  
- API handlers  

### Race Condition Tests

- Concurrent cache access  
- Multiple simultaneous API requests  

### Integration Tests

- End-to-end API behavior validation  

---

## 🧰 Testing Libraries Used

- https://github.com/stretchr/testify  
- https://github.com/vektra/mockery  

---

## 🚀 Getting Started

### Run the Application

```bash
go run .
```

## 🔌 API Endpoint

Search for a country by name:
```bash
curl "http://host:port/api/countries/search?name=India"
```

## 🏗 Build the Project

```bash
go build
```

## 🧪 Run Tests

Run all tests:

```bash
go test ./...
```

Run tests with race detection:

```bash
go test ./... -race
```

Run benchmarks:

```bash
go test ./... -bench=.
```

Run tests with verbose output and coverage:

```bash
go test ./... -v -cover
```

Generate a coverage profile:

```bash
go test ./... -coverprofile=coverage.out
```
