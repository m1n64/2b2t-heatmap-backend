# Golang Service Template

A batteries‑included starter template for building Go microservices.

**Stack**

* **Go 1.25**
* **Gin** (HTTP framework)
* **Air** (live‑reload during development)
* **Delve** (remote debugger, runs in dev mode)
* **PostgreSQL**, **Redis**, **RabbitMQ** (optional deps wired via Docker Compose)
* **GORM** as the default ORM
* Dockerfile with **dev** and **prod** stages
* `make` commands to orchestrate common tasks

---

## What’s inside

```
.
├── Dockerfile
├── docker-compose.yml
├── docker-compose.override.yml     # dev overrides (Air, Delve ports)
├── cmd
│   ├── seed/                       # seeds (go run cmd/seed)
│   └── start/                      # service entrypoint (main)
├── internal/                       # your internal application code
├── pkg/                            # shared packages (di, utils, …)
├── makefile                        # all project automation
├── README.md
└── tmp/                            # Air build artifacts
```

---

## Prerequisites

* **Docker** and **Docker Compose**
* **GNU make**
* (Optional) Go toolchain if you plan to run things outside Docker

---

## Getting started

### 1) Use this template

```bash
git clone https://github.com/m1n64/go-service-template your-service
cd your-service
```

### 2) Rename the Go module

Update the module path from the template name to your own module path.

```bash
sed -i'' -e 's|module golang-service-template|module github.com/your-org/your-service|' go.mod
```

Update any imports accordingly.

### 3) Set up environment variables

Copy the example environment file and adjust values:

```bash
cp .env.example .env
```

Provide credentials for Postgres, Redis, RabbitMQ, etc.

### 4) Run in development (Air + Delve)

In dev mode, the app is started by **Air** and Delve runs in headless debug mode.

```bash
make up   # starts app + db + redis + rabbitmq with dev overrides
```

* Service listens on **:8000**
* Delve listens on **:5864**

---

## Makefile commands

| Command                           | What it does                                                                   |
|-----------------------------------|--------------------------------------------------------------------------------|
| `make help`                       | Show all available commands with descriptions.                                 |
| `make up`                         | Start the full stack in **dev** mode (Air + Delve, Postgres, Redis, RabbitMQ). |
| `make prod`                       | Build and run the **production** image/container (no Air, no Delve).           |
| `make stop`                       | Stop running containers without removing them.                                 |
| `make down`                       | Stop and remove containers, networks, and volumes.                             |
| `make restart`                    | Restart the entire stack.                                                      |
| `make restart-container c=<name>` | Restart a single container by name.                                            |
| `make stop-container c=<name>`    | Stop a single container by name.                                               |
| `make logs`                       | Tail logs from the app container.                                              |
| `make bash`                       | Open an interactive shell in the app container.                                |
| `make seed`                       | Run seeders inside the app container.                                          |
| `make psql`                       | Open a `psql` shell in the Postgres container.                                 |
| `make redis`                      | Open a Redis CLI in the Redis container.                                       |
| `make rabbitmq`                   | Access RabbitMQ CLI or management UI.                                          |

---

## Development details

### Air configuration

Air rebuilds to `./tmp/main` and runs Delve in dev mode. Failures prevent running stale binaries.

### Delve configuration

Runs headless in the container in dev mode:

```
dlv exec ./tmp/main \
  --continue \
  --listen=:5864 \
  --api-version=2 \
  --log=true \
  --headless=true \
  --accept-multiclient
```

Requires:

* `cap_add: [SYS_PTRACE]`
* `security_opt: ["seccomp=unconfined"]`

### Database and messaging

* **PostgreSQL** via **GORM** (`pkg/utils/db.go`).
* **Redis** helpers (`pkg/utils/redis.go`).
* **RabbitMQ** helpers (`pkg/utils/rabbitmq.go`).

---

## Troubleshooting

* **`exit code 127`**: likely incorrect `full_bin` in Air config.
* **`Failed to sync logger`**: benign in Docker.
* **`operation not permitted`** with Delve: ensure correct container capabilities.
* **glibc/musl mismatch**: build with `CGO_ENABLED=0` or use a glibc image.

---

## License

This project is licensed under the terms of the **LICENSE** file.

---

Happy shipping! 🚀
