<div style="text-align:center">
   <img src="assets/smol-kitten.jpg" alt="a smol cat by Ron whisky" width="200" />
   <br/>
   <small><a href="https://en.wiktionary.org/wiki/smol#/media/File:Kitten_(06)_by_Ron.jpg">A smol cat.</a></small>
</div>

**smolgit** offers a minimalist [git](https://git-scm.com/) server, making it perfect for small teams or individual developers. Its minimal simple and just works. Small memory footprint, one binary to go. No fancy pancy javascript frameworks.

<!-- toc -->
- [Motivation](#motivation)
- [Usage](#usage)
- [Getting Started](#getting-started)
  - [Install](#install)
  - [Run](#run)
  - [Config](#config)
  - [Docker](#docker)
- [Prerequisites](#prerequisites)
- [Built with](#built-with)
- [Contribution](#contribution)
- [Donations](#donations)
<!-- /toc -->

### Motivation
**smolgit** is perfect for those who value simplicity and efficiency in their workflow. The goal is to eliminate the bloat and offer a tool that is easy to set up, use, and maintain.

### Usage

1. **git operations**: Easily perform `pull`, `push`, `clone` and `fetch` operations.
1. **repository visualization**: Browse files, view logs, explore the commit and branch and tag lists.
1. **user management**: Simple user management, add users with `ssh-keys` to `config.yaml`.


### Getting Started

#### Install

Download binary from release page.

#### Build

Build from source by cloning repository and run `make build` to build the binary

#### Run

```shell
$> ./bin/smolgit
4:37PM INF set loglevel level=DEBUG
4:37PM INF version version=main-c77299f
4:37PM INF initialize ssh server addr=:3080
4:37PM INF start server brand=smolgit address=:3080
4:37PM INF starting SSH server addr=:3081
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
$> ./bin/smolgit --help
Usage of ./bin/smolgit:
  -config string
        path to config (default "./config.yaml")
```

#### Docker

TBD

### Prerequisites

- git

## Built with

:heart:

- [golang](https://go.dev/)
- [gin](https://github.com/gin-gonic/gin)
- [pico](https://picocss.com/docs)
- [gossh](https://github.com/gliderlabs/ssh)

## Contribution

Contributions are more than welcome! Thank you!
