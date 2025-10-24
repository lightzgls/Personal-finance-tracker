# Personal Finance Tracker

A web-based personal finance management application built with Go and PostgreSQL. Track your income, expenses, and account balances with an intuitive interface.
<img width="1918" height="926" alt="image" src="https://github.com/user-attachments/assets/5735fa3e-4c4d-4e46-bd1e-667917d5bbd1" />


## Features

- **Account Management**: Add and manage multiple financial sources (bank accounts, credit cards, cash, etc.)
- **Transaction Tracking**: Record income and expenses with detailed categorization
- **Financial Summary Dashboard**: 
  - View current total balance across all active accounts
  - Track monthly income and expenses
  - Recent transaction history
- **Transaction Categories**: Organize transactions with custom categories (Food, Transportation, Salary, etc.)
- **Balance Validation**: Ensures sufficient funds before recording expense transactions
- **Active/Inactive Accounts**: Toggle account status without losing transaction history

## Tech Stack

### Backend
- **Go (Golang)**: Main programming language
- **pgx/v5**: PostgreSQL driver and toolkit
- **gorilla/schema**: Form data decoder

### Frontend
- **HTML5**: Template structure
- **CSS3**: Styling and layout
- **Go Templates**: Server-side rendering

### Database
- **PostgreSQL**: Relational database for data persistence

### Additional Tools
- **godotenv**: Environment variable management
- **google/uuid**: UUID generation for database records

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/personal-finance-tracker.git
   cd personal-finance-tracker
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   
   Create a PostgreSQL database:
   ```sql
   CREATE DATABASE finance;
   ```

   Create the required tables:
   ```sql
   CREATE TABLE ACCOUNT (
       SOURCE_ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       BALANCE NUMERIC(19,1) NOT NULL,
       SOURCE_NAME VARCHAR(100) UNIQUE NOT NULL,
       SOURCE_TYPE VARCHAR(50) NOT NULL,
       CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       DESCRIPTION TEXT,
       IS_ACTIVE BOOLEAN DEFAULT TRUE
   );

   CREATE TABLE CATEGORY (
       CATEGORY_ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       CATEGORY_NAME VARCHAR(100) NOT NULL,
       CATEGORY_TYPE VARCHAR(20) NOT NULL CHECK (CATEGORY_TYPE IN ('income', 'expense'))
   );

   CREATE TABLE TRANSACTION (
       TRANSACTION_ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       CATEGORY_ID UUID REFERENCES CATEGORY(CATEGORY_ID),
       AMOUNT NUMERIC(19,1) NOT NULL,
       TRANSACTION_DATE TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       SOURCE_NAME VARCHAR(100) REFERENCES ACCOUNT(SOURCE_NAME),
       DESCRIPTION TEXT
   );
   ```

4. **Configure environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   DATABASE_URL=postgres://username:password@localhost:5432/finance
   ```

5. **Run the application**
   ```bash
   go run cmd/main/main.go
   ```

6. **Access the application**
   
   Open your browser and navigate to:
   ```
   http://localhost:8080/home
   ```

## Project Structure

```
Personal-finance-tracker/
├── cmd/
│   └── main/
│       └── main.go              # Application entry point
├── database/
│   └── database.go              # Database connection setup
├── handler/
│   └── handler.go               # HTTP request handlers
├── model/
│   └── model.go                 # Data structures and models
├── repository/
│   └── transaction_repo.go      # Database operations
├── templates/
│   └── home.html                # HTML template for UI
├── static/
│   └── css/
│       └── style.css            # Stylesheet
├── .env                         # Environment variables (not committed)
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
└── README.md                    # Project documentation
```

## Overview

### Architecture

The application follows a layered architecture pattern:

- **Handler Layer**: Processes HTTP requests and responses
- **Repository Layer**: Manages database operations and queries
- **Model Layer**: Defines data structures
- **Database Layer**: Handles database connections

### Key Components

#### 1. Account Management
- Add new financial sources with initial balance
- Track multiple account types (savings, checking, credit card, cash)
- Maintain account status (active/inactive)

#### 2. Transaction Management
- Record income and expense transactions
- Automatic balance updates
- Category-based organization
- Transaction history with timestamps

#### 3. Dashboard
- Real-time balance calculation across all active accounts
- Monthly income/expense summary
- Recent transaction list with details

### Database Design

- **ACCOUNT**: Stores financial sources and their balances
- **CATEGORY**: Defines transaction categories (income/expense)
- **TRANSACTION**: Records all financial transactions with references to accounts and categories

### Error Handling

The application includes comprehensive error handling for:
- Negative amounts validation
- Insufficient balance checks
- Duplicate account prevention
- Database connection issues

## API Endpoints

- `GET /home` - Main dashboard with summary and transactions
- `GET /Balances` - View all account balances
- `POST /AddTransaction` - Add a new transaction
- `POST /AddSource` - Add a new financial source

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
