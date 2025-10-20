# Snippet Box

A server-side rendered web application built in Go, following [*Let's Go* by Alex Edwards](https://lets-go.alexedwards.net/). Snippet Box allows users to create, store, and share snippets of text with a clean, user-friendly interface.

## Prerequisites

* Go (1.21 or higher)
* PostgreSQL (15 or higher)
* Docker and Docker Compose
* A terminal or command-line interface

## Run Locally

1. **Install PostgreSQL:** Download from [postgresql.org/download](https://postgresql.org/download) (version 15+ recommended).

2. **Set Up Environment:** Copy `.env.example` to `.env` in the project root and add your PostgreSQL password:
```
   DB_PASSWORD=your_secure_password
```

3. **Clone the Repository:**
```bash
   git clone <repository-url>
   cd snippetbox
```

4. **Start the Database:** Run the following to create the `snippetbox` database with tables and sample data:
```bash
   docker compose up -d
```

5. **Run the Application:** Start the server, replacing `YOUR_PASSWORD` with your PostgreSQL password:
```bash
   go run ./cmd/web -db=postgresql://postgres:YOUR_PASSWORD@localhost:5432/snippetbox
```

   The app runs on `http://localhost:4000`.

6. **Optional: Change Port:** Use the `-port` flag to change the default port:
```bash
   go run ./cmd/web -port=:3000 -db=postgresql://postgres:YOUR_PASSWORD@localhost:5432/snippetbox
```

7. **Optional: Secure Database User:** Create a dedicated user for `snippetbox`:
```sql
   CREATE USER snippetbox_user WITH PASSWORD 'secure_password';
   GRANT CONNECT ON DATABASE snippetbox TO snippetbox_user;
   GRANT USAGE ON SCHEMA public TO snippetbox_user;
   GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO snippetbox_user;
```

   Update the connection string:
```bash
   go run ./cmd/web -db=postgresql://snippetbox_user:secure_password@localhost:5432/snippetbox
```

## Accessing the Application

Open `http://localhost:4000` (or your custom port) in your browser.

## Stopping the Application

* **Stop the server:** Press `Ctrl+C` in the terminal.
* **Stop the database:** Run `docker compose down`.

## Troubleshooting

* **Database connection error:** Verify PostgreSQL is running, the password matches `.env`, and the `snippetbox` database exists.
* **Port conflict:** Ensure port `4000` (or your custom port) is free or use a different port.
* **Docker issues:** Confirm Docker is running and check logs with `docker compose logs`.

## Security Note

Avoid using the `postgres` user in production. Use a dedicated user with limited permissions, as shown in step 7.
