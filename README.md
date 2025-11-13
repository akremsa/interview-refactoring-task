 # Technical Interview Task: REST API Refactoring

## Overview
This repository contains a simple REST API server with intentional design flaws. The candidate's task is to identify these issues and refactor the code following proper design patterns and principles.

## Current Implementation
The `main.go` file contains a basic REST API with one endpoint `/save-data` that accepts a byte array and saves it to different storage types (file, database, or memory).

## Task Instructions for Candidate
1. **Review the code** in `main.go` and identify design pattern violations and code quality issues
2. **Refactor the code** to follow best practices and design patterns
3. **Explain your changes** and the reasoning behind them

## What to Test the Candidate For

### 1. Factory Pattern Violation
**Problem**: The `DataStorage.SaveData()` method uses hard-coded if-else statements to determine storage type.
**Expected Solution**: 
- Create a `Storage` interface
- Implement concrete storage classes (FileStorage, DatabaseStorage, MemoryStorage)
- Create a StorageFactory to instantiate the appropriate storage type

### 2. Single Responsibility Principle Violations
**Problems**:
- `DataStorage` class handles multiple storage types
- `DataHandler.HandleSaveData()` does too many things (parsing, validation, storage, response)
- `Request` struct mixes data and metadata

**Expected Solutions**:
- Separate concerns into different classes/functions
- Create dedicated validator
- Create service layer for business logic

### 3. Open/Closed Principle Violation
**Problem**: Adding new storage types requires modifying existing code
**Expected Solution**: Use interfaces and factory pattern to make the system extensible

### 4. Dependency Injection Issues
**Problems**:
- Hard-coded dependencies
- No constructor injection
- Tight coupling
- **Database connection created inside handler method** (major violation)

**Expected Solutions**:
- Use constructor injection
- Pass dependencies through constructors
- Create proper service layer
- **Move database initialization to main() and inject into handler**

### 5. Performance and Resource Management Issues
**Problems**:
- **Database connection created on every HTTP request** (Critical Issue!)
- **Expensive operations in handler method**
- **No connection pooling or reuse**

**Expected Solutions**:
- Initialize database connection once at startup
- Reuse connections across requests
- Implement proper resource management

**Key Learning Point**: This is a major performance and design issue that candidates should immediately identify. Creating database connections inside request handlers is extremely inefficient and violates dependency injection principles.

### 6. Error Handling Issues
**Problems**:
- Inconsistent error handling
- No proper HTTP status codes
- Missing validation

### 7. Other Issues to Identify
- No HTTP method validation
- Hard-coded responses
- No configuration management
- Missing graceful shutdown
- No middleware usage

## Running the Application

### Original Code
```bash
go run main.go
```

Test with curl:
```bash
# Save to file
curl -X POST http://localhost:8080/save-data \
  -H "Content-Type: application/json" \
  -d '{"data":"SGVsbG8gV29ybGQ=","storage_type":"file"}'

# Save to database (mock)
curl -X POST http://localhost:8080/save-data \
  -H "Content-Type: application/json" \
  -d '{"data":"SGVsbG8gV29ybGQ=","storage_type":"database"}'
```

### Expected Refactored Solution
The `solution.go` file contains a properly refactored version showing:
- Factory pattern implementation
- Proper separation of concerns
- Dependency injection
- Interface-based design
- Better error handling

## Evaluation Criteria

### Excellent (Senior Level)
- Identifies all major design pattern violations
- Implements proper factory pattern
- Uses interfaces for polymorphism
- Applies dependency injection
- Separates concerns properly
- Adds proper validation and error handling

### Good (Mid Level)
- Identifies most design issues
- Implements basic factory pattern
- Uses some interfaces
- Improves separation of concerns
- Adds basic validation

### Needs Improvement (Junior Level)
- Identifies some obvious issues
- Makes basic improvements
- May miss factory pattern opportunity
- Limited understanding of design principles

## Follow-up Questions
1. How would you add a new storage type (e.g., Redis)?
2. How would you add authentication to this API?
3. How would you implement logging and monitoring?
4. How would you handle concurrent requests?
5. How would you add configuration management?

## Time Allocation
- **Code Review**: 10-15 minutes
- **Refactoring**: 30-45 minutes  
- **Discussion**: 15-20 minutes
- **Total**: 60-80 minutes
