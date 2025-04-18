
<a href="https://en.wiktionary.org/wiki/smol"><img align="left" src="assets/smol-kitten.jpg" alt="a smol cat by Ron whisky" width="150" height="100" /></a>

**smolgit** offers a minimalist [git](https://git-scm.com/) server, making it perfect for small teams or individual developers. Its minimal simple and just works. It's perfect for those who value simplicity and efficiency in their workflow. Small memory footprint, one binary to go.

<!-- toc -->
- [Features](#features)
- [Preview](#preview)
- [Getting Started](#getting-started)
  - [Install](#install)
  - [Run](#run)
  - [Config](#config)
  - [Docker](#docker)
- [Prerequisites](#prerequisites)
- [Built with](#built-with)
- [Contribution](#contribution)
<!-- /toc -->

### Features

1. **git operations** - easily perform `pull`, `push`, `clone` and `fetch` operations.
1. **repository visualization** - browse files, view logs, explore the commit, branch and tag lists.
1. **user management** - simple user management, add users with `ssh-keys` to `config.yaml`.
1. **permissions** - assign persmissions to user.
1. **ligh-dark** - web theme based on your system settings.
1. **basic-auth** - web basic auth middleware.

### Preview

<p align="center">
   <img src="assets/web_1.png" alt="screenshot" width="700" />
</p>
<p align="center">
   <img src="assets/web_2.png" alt="screenshot" width="700" />
</p>


### Getting Started

#### Install

1. Download binary from [ release page ](https://github.com/roman-kiselenko/smolgit/releases).
1. Generate default `config.yaml` file with command `./smolgit config > config.yaml`.
    - Use [`yq`](https://github.com/mikefarah/yq) for inline changes `./smolgit config | yq '.server.disabled = true' > config.yaml`
1. Run `./smolgit`

```shell
$> ./smolgit
10:08AM INF set loglevel level=DEBUG
10:08AM INF version version=main-a4f6438
10:08AM INF initialize web server addr=:3080
10:08AM INF initialize ssh server addr=:3081
10:08AM INF start server brand=smolgit address=:3080
10:08AM INF starting SSH server addr=:3081
```

#### Config

Generate default `config.yaml` file with command `./bin/smolgit config > config.yaml`.

```yaml
log:
  # Color log output
  color: true
  # Log as json
  json: false
  # Log level (INFO, DEBUG, TRACE, WARN)
  level: DEBUG
server:
  # Disable web server
  disabled: false
  # Enable basic http auth
  auth:
    enabled: false
    # Credentials for basic auth
    accounts:
      - login: user2
        password: bar
      - login: user1
        password: foo
  # Web server address
  addr: ":3080"
  # Navbar brand string
  brand: "smolgit"
ssh:
  # SSH server address
  addr: ":3081"
git:
  # Folder to save git repositories
  path: /tmp/smolgit
  # Base for clone string formating
  # (e.g. ssh://git@my-git-server.lan/myuser/project.git)
  base: "git@my-git-server.lan"
  users:
    # User name used for folder in git.path
  - name: "bob"
    # Permissions, wildcard or regex
    # User to check access for other repositories
    # '*' - access for all repositories
    # 'admin' - access for admin's repositories
    # '(admin|billy)' - access for admin's and billy's repositories
    permissions: "*"
    keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCq9rD9b8tYyuSLsTECHCn... developer@mail.com
```

cli options:

```shell
$> ./smolgit --help
Usage of ./smolgit:
  -config string
        path to config (default "./config.yaml")
```

#### Docker

In order to run `smolgit` in docker there is the [`Dockerfile`](/Dockerfile).

1. Build image `make build-docker`
1. Generate `config.yaml` file `make config-docker`, it'll create `config.yaml` in the current directory and mount it for docker.
1. Run `smolgit` in docker:

```shell
$> make run-docker
docker run -it -p 3080:3080 -p 3081:3081 -v /path-to-smolgit-project/smolgit/:/etc/smolgit smolgit
3:53PM INF set loglevel level=DEBUG
3:53PM INF version version=dev
3:53PM INF initialize web server addr=:3080
3:53PM INF initialize ssh server addr=:3081
3:53PM INF start server brand=smolgit address=:3080
3:53PM INF starting SSH server addr=:3081
```

### Prerequisites

- git

### Built with

:heart:

- [golang](https://go.dev/)
- [gin](https://github.com/gin-gonic/gin)
- [go-git](https://github.com/go-git/go-git)
- [pico](https://picocss.com/docs)
- [gossh](https://github.com/gliderlabs/ssh)

### Local development

- [golang](https://go.dev/)
- [yq](https://mikefarah.gitbook.io/yq)
- [bats](https://bats-core.readthedocs.io/en/stable/)

### Contribution

Contributions are more than welcome! Thank you!
