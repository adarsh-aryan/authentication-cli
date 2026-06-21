# Login CLI

A simple authentication system consisting of:

- A Go-based authentication server
- A Go CLI client for interacting with the server (register, login, logout, etc.)

## Features

- User registration
- Login / logout
- Session handling
- `whoami` command to check current user
- SQLite database (local)
- Account lock after maximum failed login attempts
- Interactive CLI with command history and auto-completion
- Optional MFA (2FA) support using Google Authenticator

## Project Structure

- `auth-server/` – backend API server
- `auth-client/` – CLI application
- `shared/` – shared types
- `docker-compose.yml` – container orchestration

## Prerequisites

Install the following:

- Go (>= 1.20 recommended)
- Docker & Docker Compose (optional)

## Setup

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd authentication-cli
```

### 2. Configure environment variables

Copy example env files:

```bash
cp .env.example auth-server/.env
cp .env.example auth-client/.env
```

Update values if needed.

## Running with Docker (Recommended)

This project uses Docker Compose but requires running the server and CLI separately.

### 1. Build the images (auth-server and auth-client)

```bash
docker compose build
```

### 1. Start the auth server

```bash
docker compose up auth-server -d --remove-orphans
```

- Runs the server in the background
- Uses SQLite for persistence inside the container

### 2. Run the CLI (interactive)

```bash
docker compose run --rm auth-client --remove-orphans
```

- Starts an interactive CLI session
- `--rm` ensures the container is removed after exit
- You can now run commands like `login`, `register`, etc.

### 4. Stop the server (when done)

```bash
docker compose down
```

## Running Manually (Without Docker)

### Start the server

```bash
cd auth-server
go run main.go
```

### Run the CLI

In a new terminal:

```bash
cd auth-client
go run main.go
```

## Usage

Once the CLI is running, you can use the following commands:

The CLI is interactive and supports:

- Command history (use arrow keys to navigate previous commands)
- Auto-completion (press Tab to complete commands and flags)

### Register

```bash
register --username=YOUR_USERNAME --password=YOUR_PASSWORD
```

### Login

```bash
login --username=YOUR_USERNAME --password=YOUR_PASSWORD
```

### Logout

```bash
logout
```

### Check current user

```bash
whoami
```

### Exit CLI

```bash
exit
```

### Enable 2FA (MFA)

```bash
enable-2fa
```

- Enables two-factor authentication for your account
- You will be prompted to scan a QR code using Google Authenticator (or a compatible app)

### Disable 2FA (MFA)

```bash
disable-2fa
```

- Disables two-factor authentication for your account

## Development Notes

- The server uses SQLite for persistence (`auth-server/data/auth.db`)
- Shared data structures are defined in `shared/types.go`
- CLI commands are implemented using a command-based structure under `auth-client/cmd/`

## Session Management

After a successful login, the CLI stores the session locally inside the client container at `data/config.json`.

- This file tracks the active session for the CLI
- It is used by `whoami` to identify the current user
- It is also used by `logout` to determine which session should be invalidated on the server
- Docker volumes are used to persist both the client session file and the server database (`auth.db`), ensuring data is retained across container restarts

## Account Locking

To improve security, accounts are automatically locked after exceeding the maximum number of failed login attempts.

- Once locked, the user will not be able to log in
- This prevents brute-force login attempts
- (If applicable) unlocking requires administrative action or resetting state (e.g., clearing the database during development)

## Troubleshooting

- Ensure the server is running before using the CLI
- Check `.env` configuration if requests fail
- Delete `auth.db` to reset local state
