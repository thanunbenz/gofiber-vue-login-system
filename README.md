# Login and Registration

A sample authentication project built with Vue, Go, JWT, and PostgreSQL. It demonstrates how a frontend communicates with a backend to register users, log in, and manage refresh-token sessions.

## Features

- User registration with bcrypt password hashing
- Login with a short-lived access JWT
- Refresh token stored in an HttpOnly cookie
- Refresh-token rotation with server-side revocation
- Protected frontend routes
- PostgreSQL-backed user and session storage

## Tech Stack

- **Frontend:** Vue 3, Vue Router, TypeScript, Vite, Tailwind CSS, Axios
- **Backend:** Go, Fiber, pgx, JWT, godotenv
- **Database:** PostgreSQL
- **Tools:** Docker Compose, Playwright, Prettier, vue-tsc

## Requirements

- Go 1.26.2 or newer
- Node.js 22.18 or newer, or 24.12 or newer
- npm
- Docker Compose

## Setup

Copy the environment files:

```bash
cp Backend/.env.example Backend/.env
cp Frontend/.env.example Frontend/.env
```

Set the frontend API URL in `Frontend/.env`:

```env
API_URL=http://localhost:8080/api/v1
```

Set a strong `JWT_SECRET` and configure the PostgreSQL credentials in `Backend/.env`. The example values are intended for local development only.

## Environment Variables

### Backend

| Variable | Description | Local default |
|---|---|---|
| `PORT` | Backend HTTP port | `8080` |
| `DB_URL` | PostgreSQL connection string | See `.env.example` |
| `JWT_SECRET` | Secret used to sign access tokens | Change this value |
| `FRONTEND_URL` | Allowed CORS origin | `http://localhost:5173` |
| `COOKIE_SECURE` | Restricts refresh cookies to HTTPS | `false` |

### Frontend

| Variable | Description | Local default |
|---|---|---|
| `API_URL` | Backend API base URL | `http://localhost:8080/api/v1` |

## Run the Project

Start PostgreSQL:

```bash
docker compose -f Backend/docker-compose.yml up -d db
```

Start the backend in one terminal:

```bash
cd Backend
go run ./cmd/api
```

Start the frontend in another terminal:

```bash
cd Frontend
npm install
npm run dev
```

Open the application at <http://localhost:5173>.

Check the backend health endpoint at <http://localhost:8080/api/v1/health>.

## API Endpoints

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| `GET` | `/api/v1/health` | Check backend status | No |
| `POST` | `/api/v1/register` | Register a user | No |
| `POST` | `/api/v1/login` | Log in and create a session | No |
| `POST` | `/api/v1/refresh` | Rotate the refresh token | Refresh cookie |
| `POST` | `/api/v1/logout` | Revoke the current session | Refresh cookie |
| `GET` | `/api/v1/me` | Return the authenticated user | Access JWT |

## Authentication Notes

- Access JWTs expire after 20 minutes and are stored in frontend memory.
- Refresh tokens expire after 12 hours and are stored in HttpOnly cookies.
- PostgreSQL stores only SHA-256 hashes of refresh tokens.
- A successful refresh replaces the old refresh token with a new one.
- Logout deletes the refresh-token session from PostgreSQL.
- Production deployments must use HTTPS and set `COOKIE_SECURE=true`.

## Testing

Run backend tests and static analysis:

```bash
cd Backend
go test ./...
go vet ./...
```

Build and type-check the frontend:

```bash
cd Frontend
npm run build
```

## Project Structure

```text
Backend/
├── cmd/api/
└── internal/
    ├── config/
    ├── database/
    ├── handler/
    ├── middleware/
    ├── repository/
    └── service/

Frontend/
└── src/
    ├── router/
    ├── views/
    └── api.ts
```

## Login Flow

![Login authentication sequence diagram](docs/diagrams/sequenceDiagram.png)

## License

This project is licensed under the [MIT License](LICENSE).
