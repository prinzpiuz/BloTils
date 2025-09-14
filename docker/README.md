## BloTils

Efficient, secure container image for the BloTils, based on Go and Alpine Linux.

#### Table of Contents

- [Description](#description)
- [How To Use](#how-to-use)
- [Environment Variables](#environment-variables)
- [Exposed Port](#exposed-port)
- [Healthcheck](#healthcheck)
- [Entrypoint](#entrypoint-and-cmd)
- [Volumes](#volumes)
- [Building Locally](#building-locally)
- [License](#license)

***

#### Description

This image runs the BloTils application using Go 1.24 and Alpine Linux. It includes built-in support for SQLite, environment variable configuration, and a robust healthcheck for production monitoring.

***

#### How To Use

Pull the image from GitHub Container Registry:

```sh
docker pull ghcr.io/prinzpiuz/blotils:latest
```

Run with configuration and database files as mounts:

```sh
docker run -d \
  --name blotils \
  -p 8000:8000 \
  -e BT_MAILGUNPASSWORD=your_mailgun_password \
  -e BT_ENV=production \
  -v $(pwd)/BloTils.db:/blotils/BloTils.db \
  -v $(pwd)/config.json:/blotils/config.json \
  ghcr.io/prinzpiuz/blotils:latest
```

Or use with Docker Compose:

```yaml
services:
  blotils:
    image: ghcr.io/prinzpiuz/blotils:latest
    ports:
      - 8000:8000
    environment:
      BT_MAILGUNPASSWORD: <your_mailgun_password>
      BT_ENV: production
    volumes:
      - ./BloTils.db:/blotils/BloTils.db
      - ./config.json:/blotils/config.json
```

***

#### Building Locally

```sh
docker build -t blotils -f docker/Dockerfile .
```

See the Dockerfile for build stages, `go build` options, and included files.

***

#### Environment Variables

| Name | Description | Example |
| :-- | :-- | :-- |
| BT_MAILGUNPASSWORD | MailGun Password | xxxxx |
| BT_ENV | App environment (`dev`, `production`, etc.) | production |

***

#### Exposed Port

- **8000:** Main application HTTP API.

***

#### Creating an Admin User

```sh
docker run --rm \
  -v $(pwd)/BloTils.db:/blotils/BloTils.db \
  -v $(pwd)/config.json:/blotils/config.json \
  ghcr.io/prinzpiuz/blotils:latest createAdmin
```

or with a running container

```sh
docker exec -it blotils /blotils/BloTils createAdmin
```

***

#### Healthcheck

Container includes an automatic healthcheck for the `/api/v1/ping` endpoint:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1/ping || exit 1
```

***

#### Entrypoint and CMD

- Entrypoint runs `/blotils/entrypoint.sh`, which manages permissions.
- Default command starts the BloTils app binary.

***

#### Volumes

- Mount your SQLite database file, config file, and other resources.
- Required:
  - `BloTils.db` at `/blotils/BloTils.db`
  - `config.json` at `/blotils/config.json`

***

#### License

GPL-3.0 © [prinzpiuz](https://github.com/prinzpiuz)
