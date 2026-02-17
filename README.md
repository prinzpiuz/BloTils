<p align="center">
  <img src="https://github.com/prinzpiuz/BloTils/blob/staging/static/favicon/android-chrome-512x512.png?raw=true" alt="BloTils Logo" width="200"/><br/>
  <b>Blo</b>g-u<b>Tils</b>, <b>Common utilities for your blogging platform.</b>
</p>

## BloTils

Efficient and extensible blog utilities built with Go, powered by SQLite and designed for static and dynamic content management.

#### Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Documentation](#documentation)
- [Running Locally](#running-locally)
- [Docker Deployment](#docker-deployment)
- [License](#license)

---

#### Features

- [x] Likes/Clap Counter
- [ ] Comments
- [ ] Simple Analytics
- [ ] Email Subscription
- [ ] Forms
- [ ] Polls
- [ ] Dynamic Content
- [ ] App For Blog

---

#### Quick Start

**Using Docker (Recommended):**

```bash
# 1. Create data directory and files
mkdir -p BloTils_Data
touch BloTils_Data/BloTils.db
cp config.json BloTils_Data/config.json

# 2. Set permissions for container user (UID 1001)
sudo chown -R 1001:1001 BloTils_Data/

# 3. copy .env file
cp env.sample .env

# 4. Run container
docker run -d \
  --name blotils \
  -p 8000:8000 \
  --env-file .env \
  -v $(pwd)/BloTils_Data/BloTils.db:/blotils/BloTils.db \
  -v $(pwd)/BloTils_Data/config.json:/blotils/config.json \
  ghcr.io/prinzpiuz/blotils:latest

# 5. Create admin user
docker exec -it blotils /blotils/BloTils -createadmin
```

Visit `http://localhost:8000` to access BloTils.

---

#### Configuration

BloTils uses a hybrid configuration system:

1. **`config.json`** — Default settings (non-sensitive)
2. **Environment variables** — Overrides and secrets (prefixed with `BT_`)

**Configuration Precedence:** Environment Variables > config.json > Defaults

##### Environment Variables

| Variable | Description | Default |
|:---------|:------------|:--------|
| `BT_ENV` | Environment (`dev`, `staging`, `production`) | `dev` |
| `BT_BASE_URL` | Application base URL (for email links) | `http://localhost:8000` |
| `BT_DBLOCATION` | SQLite database path | `./BloTils.db` |
| `BT_PORT` | Server port | `8000` |
| `BT_SMTP_ENABLED` | Enable email sending | `false` |
| `BT_SMTP_HOST` | SMTP server host | - |
| `BT_SMTP_PORT` | SMTP server port | `587` |
| `BT_SMTP_USERNAME` | SMTP username | - |
| `BT_SMTP_PASSWORD` | SMTP password | - |
| `BT_SMTP_FROM` | From email address | - |
| `BT_SMTP_FROMNAME` | From display name | `BloTils` |

---

#### Documentation

- [Full Documentation](https://github.com/prinzpiuz/BloTils/wiki/BloTils-Docs)
- [Docker Deployment Guide](docker/README.md)
- [Self-Hosting Guide](https://github.com/prinzpiuz/BloTils/wiki/How-To-Deploy-Yourself)
- [API Reference](https://github.com/prinzpiuz/BloTils/wiki/API-Reference)

---

#### Running Locally

**Prerequisites:**

- Go 1.24.1 or later
- SQLite 3.35+

**Steps:**

```bash
# 1. Clone the repository
git clone https://github.com/prinzpiuz/BloTils.git
cd BloTils

# 2. Update values in config.json

# 3. Set environment variables (or create .env file)
cp env.sample .env

# 4. Run the application
go run bloTils.go

# 5. Create admin user (in another terminal)
go run bloTils.go -createadmin
```

Visit `http://localhost:8000`

---

#### Docker Deployment

For production deployments, see the [Docker README](docker/README.md).

**Quick Docker Compose setup:**

```yaml
services:
  blotils:
    image: ghcr.io/prinzpiuz/blotils:latest
    ports:
      - "8000:8000"
    env_file:
      - .env
    volumes:
      - ./BloTils_Data/BloTils.db:/blotils/BloTils.db
      - ./BloTils_Data/config.json:/blotils/config.json
    restart: unless-stopped
```

---

#### License

GPL-3.0 © [prinzpiuz](https://github.com/prinzpiuz)
