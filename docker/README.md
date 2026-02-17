## BloTils Docker Image

Efficient, secure container image for BloTils, based on Go and Alpine Linux.

#### Table of Contents

- [Description](#description)
- [Quick Start](#quick-start)
- [Docker Compose](#docker-compose)
- [Environment Variables](#environment-variables)
- [Volumes & Permissions](#volumes--permissions)
- [Creating Admin User](#creating-admin-user)
- [Healthcheck](#healthcheck)
- [Building Locally](#building-locally)
- [Troubleshooting](#troubleshooting)
- [License](#license)

---

#### Description

This image runs the BloTils application using Go 1.24 and Alpine Linux. It includes:

- Built-in SQLite support
- Environment variable configuration
- SMTP email support (works with any provider)
- Robust healthcheck for production monitoring
- Non-root user execution for security

---

#### Quick Start

**1. Pull the image:**

```bash
docker pull ghcr.io/prinzpiuz/blotils:latest
```

**2. Prepare data files:**

```bash
# Create data directory
mkdir -p BloTils_Data

# Create empty database
touch BloTils_Data/BloTils.db

# Copy config file
cp [config.json](../config.json) BloTils_Data/config.json

# Set permissions for container user (UID 1001)
sudo chown -R 1001:1001 BloTils_Data/
```

**3. Create `.env` file:**

cp [env.sample](../env.sample) .env

**4. Run the container:**

```bash
docker run -d \
  --name blotils \
  -p 8000:8000 \
  --env-file .env \
  -v $(pwd)/BloTils_Data/BloTils.db:/blotils/BloTils.db \
  -v $(pwd)/BloTils_Data/config.json:/blotils/config.json \
  --security-opt no-new-privileges:true \
  --restart unless-stopped \
  ghcr.io/prinzpiuz/blotils:latest
```

**5. Create admin user:**

```bash
docker exec -it blotils /blotils/BloTils -createadmin
```

---

#### Docker Compose

**`docker-compose.yml`:**

```yaml
services:
  blotils:
    image: ghcr.io/prinzpiuz/blotils:latest
    container_name: blotils
    ports:
      - "8000:8000"
    env_file:
      - .env
    volumes:
      - ./BloTils_Data/BloTils.db:/blotils/BloTils.db
      - ./BloTils_Data/config.json:/blotils/config.json
    security_opt:
      - no-new-privileges:true
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  default:
    driver: bridge
```

**Run:**

```bash
docker compose up -d
```

---

#### Environment Variables

##### Application Settings

| Variable | Description | Default |
|:---------|:------------|:--------|
| `BT_ENV` | Environment (`dev`, `staging`, `production`) | `dev` |
| `BT_BASE_URL` | Base URL for email links | `http://localhost:8000` |
| `BT_PORT` | Server port | `8000` |
| `BT_DBLOCATION` | Database path inside container | `/blotils/BloTils.db` |

##### SMTP Settings

| Variable | Description | Default |
|:---------|:------------|:--------|
| `BT_SMTP_ENABLED` | Enable email sending | `false` |
| `BT_SMTP_HOST` | SMTP server hostname | - |
| `BT_SMTP_PORT` | SMTP server port | `587` |
| `BT_SMTP_USERNAME` | SMTP username | - |
| `BT_SMTP_PASSWORD` | SMTP password | - |
| `BT_SMTP_FROM` | From email address | - |
| `BT_SMTP_FROMNAME` | From display name | `BloTils` |

##### SMTP Provider Examples

**SendGrid:**

```bash
BT_SMTP_ENABLED=true
BT_SMTP_HOST=smtp.sendgrid.net
BT_SMTP_PORT=587
BT_SMTP_USERNAME=apikey
BT_SMTP_PASSWORD=SG.your-api-key
BT_SMTP_FROM=noreply@yourdomain.com
```

**Mailgun:**

```bash
BT_SMTP_ENABLED=true
BT_SMTP_HOST=smtp.mailgun.org
BT_SMTP_PORT=587
BT_SMTP_USERNAME=postmaster@your-domain.mailgun.org
BT_SMTP_PASSWORD=your-mailgun-password
BT_SMTP_FROM=noreply@yourdomain.com
```

**AWS SES:**

```bash
BT_SMTP_ENABLED=true
BT_SMTP_HOST=email-smtp.us-east-1.amazonaws.com
BT_SMTP_PORT=587
BT_SMTP_USERNAME=your-ses-smtp-username
BT_SMTP_PASSWORD=your-ses-smtp-password
BT_SMTP_FROM=noreply@yourdomain.com
```

**Gmail (use App Password):**

```bash
BT_SMTP_ENABLED=true
BT_SMTP_HOST=smtp.gmail.com
BT_SMTP_PORT=587
BT_SMTP_USERNAME=your@gmail.com
BT_SMTP_PASSWORD=your-app-password
BT_SMTP_FROM=your@gmail.com
```

---

#### Volumes & Permissions

##### Required Mounts

| Host Path | Container Path | Description |
|:----------|:---------------|:------------|
| `./BloTils.db` | `/blotils/BloTils.db` | SQLite database |
| `./config.json` | `/blotils/config.json` | Configuration file |

##### Important: File Permissions

The container runs as user `blotils` (UID 1001). Your mounted files must be accessible:

```bash
# Set correct ownership
sudo chown -R 1001:1001 ./BloTils_Data/

# Or set permissions
sudo chmod 666 ./BloTils_Data/BloTils.db
sudo chmod 644 ./BloTils_Data/config.json
```

##### Files Must Exist Before Starting

Docker bind mounts for **files** (not directories) must exist before starting the container:

```bash
# Create files if they don't exist
touch BloTils_Data/BloTils.db
touch BloTils_Data/config.json
```

If the files don't exist, Docker creates directories instead, causing errors.

---

#### Creating Admin User

**With running container:**

```bash
docker exec -it blotils /blotils/BloTils -createadmin
```

**With one-off container:**

```bash
docker run --rm -it \
  -v $(pwd)/BloTils_Data/BloTils.db:/blotils/BloTils.db \
  -v $(pwd)/BloTils_Data/config.json:/blotils/config.json \
  ghcr.io/prinzpiuz/blotils:latest \
  /blotils/BloTils -createadmin
```

---

#### Healthcheck

The container includes an automatic healthcheck:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1/ping || exit 1
```

**Check health status:**

```bash
docker inspect --format='{{.State.Health.Status}}' blotils
```

---

#### Building Locally

**From project root:**

```bash
# Build image
docker build -t blotils:local -f docker/Dockerfile .

# Run local image
docker run -d \
  --name blotils \
  -p 8000:8000 \
  --env-file .env \
  -v $(pwd)/BloTils_Data/BloTils.db:/blotils/BloTils.db \
  -v $(pwd)/BloTils_Data/config.json:/blotils/config.json \
  blotils:local
```

**Using Docker Compose for local development:**

```yaml
# docker-compose.local.yml
services:
  blotils:
    build:
      context: .
      dockerfile: docker/Dockerfile
    image: blotils:local
    container_name: blotils
    ports:
      - "8000:8000"
    env_file:
      - .env
    volumes:
      - ./BloTils_Data/BloTils.db:/blotils/BloTils.db
      - ./BloTils_Data/config.json:/blotils/config.json
    restart: unless-stopped
```

```bash
docker compose -f docker-compose.local.yml up --build
```

---

#### Troubleshooting

##### Permission Denied on Database

```
Error Initializing DB: unable to open database file: permission denied
```

**Fix:** Set correct ownership:

```bash
sudo chown 1001:1001 ./BloTils_Data/BloTils.db
sudo chown 1001:1001 ./BloTils_Data/config.json
```

##### Config File Becomes Directory

```
config.json is a directory, not a file
```

**Cause:** File didn't exist when container started.

**Fix:**

```bash
# Remove the wrongly created directory
rm -rf ./BloTils_Data/config.json

# Create the file
cat > ./BloTils_Data/config.json << 'EOF'
{
  "AppConfig": { "Name": "BloTils" }
}
EOF

# Restart container
docker compose down && docker compose up -d
```

##### SMTP Not Working

```
Warning: SMTP enabled but missing required config
```

**Fix:** Ensure all required SMTP variables are set:

```bash
BT_SMTP_ENABLED=true
BT_SMTP_HOST=smtp.example.com
BT_SMTP_PORT=587
BT_SMTP_USERNAME=your-username
BT_SMTP_PASSWORD=your-password
BT_SMTP_FROM=noreply@example.com
```

##### Check Logs

```bash
# View logs
docker logs blotils

# Follow logs
docker logs -f blotils

# Last 50 lines
docker logs --tail 50 blotils
```

---

#### Exposed Port

- **8000:** Main application HTTP API

---

#### License

GPL-3.0 © [prinzpiuz](https://github.com/prinzpiuz)
