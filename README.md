<p align="center">
  <img src="https://github.com/prinzpiuz/BloTils/blob/staging/static/favicon/android-chrome-512x512.png?raw=true" alt="BloTils Logo" width="200"/><br/>
  <b>Blog-Utils, Common utilities for your blogging platform.</b>
</p>

## BloTils

Efficient and extensible blog utilities built with Go, powered by SQLite and designed for static and dynamic content management.

#### Table of Contents

- [Features](#features)
- [TryOut](#quick-docker-try-out)
- [Docs](https://github.com/prinzpiuz/BloTils/wiki/BloTils-Docs)
- [Docker Docs](docker/README.md)
- [Want To Run Locally?](#want-to-run-locally)
- [Want To Deploy?](https://github.com/prinzpiuz/BloTils/wiki/How-To-Deploy-Yourself)

#### Features

- [x] Likes Counter
- [ ] Comments
- [ ] Simple Analytics
- [ ] Email Subscription
- [ ] Forms
- [ ] Polls
- [ ] Dynamic Content
- [ ] App For Blog

***

#### Quick Docker Try-Out

- Create `BloTils.db` sqlite DB and copy `config.json` and Run

```sh
docker run -d --name blotils -p 8000:8000 -e BT_MAILGUNPASSWORD=your_sendgrid_api_key -e BT_ENV=production -v ./BloTils.db:/blotils/BloTils.db -v ./config.json:/blotils/config.json ghcr.io/prinzpiuz/blotils:latest
```

```sh
docker exec -it blotils /blotils/BloTils createAdmin
```

More about docker based setup can be found in

***

#### Want To Run Locally?

- Clone the repository

  ```sh
   git clone https://github.com/prinzpiuz/BloTils.git
   cd BloTils
   ```

- Source the env variables

  ```sh
  BT_MAILGUNPASSWORD="XXX"
  BT_ENV="production"
  BT_DBLOCATION="./BloTils.db"
  ```

- Suggested GO version go1.24.1

    ```sh
    go run blotils.go
    ```

***

#### License

GPL-3.0 © [prinzpiuz](https://github.com/prinzpiuz)
