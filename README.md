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

1. **git operations**: Easily perform `pull`, `push`, and `fetch` operations.
1. **repository visualization**: Browse files, view logs, explore the commit and branch and tag lists.
1. **file content**: See the content of non-binary files.
1. **user management**: Simple user management, create users with `ssh-keys` to allow them `push\pull`


### Getting Started

#### Install

Download binary from release page.

#### Build

Build from source by cloning repository and run `make build` to build the binary

#### Run

```shell
$> ./bin/smolgit server
12:30PM INF set loglevel level=DEBUG
12:30PM INF found config file config=./config.yaml
12:30PM INF version version=main-fbebd48
12:30PM INF initialize sqlite database db_path=./database.sqlite
12:30PM INF initialize ssh server addr=:3081
12:30PM INF start server brand=smolgit address=:3080
12:30PM INF starting SSH server addr=:3081
```

#### Config

If directory contains `config.yaml` file, it will be used (can be overrided with cli options).

```yaml
log:
  color: true
  json: false
  level: DEBUG
server:
  addr: ":3080"
  brand: "smolgit"
ssh:
  addr: ":3081"
git:
  path: /tmp/smolgit
  base: "git@my-git-server.lan"
db:
  path: ./database.sqlite
root:
  login: "admin"
  password: "root"
  keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCq9rD9b8tYyuSLsTECHCnx0/KV0tHhqgj2xHjuKqcyhKqyntPqHG+kd5EBFCZEW5zrzmY4ErKlA1aLcrGkwPsEjc5nKM9EvbzDdIqxQpPDl8oL6cxu0rQxwpYDDXSQxRsP6113rg6KLxfQZbvT3GkUv9deRMD1/CiPxZ0vMMv09U226kzaaf0q7F/Vo0iXGElwPl5Q73sXSUIqPr/Fcz4PEUwYRCXeMw07+ziksa9HwFApdP/JEbrzGa7BHmrx4qANOIbrxom8Csr7Eb0YnCnx08boxykFO5l7lLWYzU6ji6wTIGtL/52kJiHz764RnAMKQHFDEhFejj5qQi/NLlnhw0x2ZZh0Oc2O6oNQ+f3PTyygnBNoffRhZTEU0jgj3gatoy4my77pRbRhf7hcRT8E0/qZVLWSEZEe8R46Fd5KxuFi5s92JRV92VhG7L+91J7tczi5OeKTDUZFTTYhGcBv3rzAlDaiNvLUoBvTEJXysUZsaBJyHqZ6Wt9CyN2Llgk= developer@mail.com
```

cli options:

```shell
$> ./bin/smolgit --help
# ...
GLOBAL OPTIONS:
   --ssh_addr value                         Address for SSH server (default: ":3081")
   --log_color                              Color output of logs, works only with text logs (default: true)
   --log_json                               Output logs in json format (default: false)
   --log_level value                        Set log level (INFO, DEBUG, WARN, TRACE) (default: "DEBUG")
   --server_addr value                      Address for Web server (default: ":3080")
   --server_brand value                     Header used in web (default: "smolgit")
   --git_base value                         Format 'git clone ...' link (default: "git@my-git-server.lan")
   --git_path value                         Set directory for git repositories (default: "./tmp/smolgit")
   --db_path value                          Path to save database file (default: "./database.sqlite")
   --root_login value                       Admin user will be create at first start
   --root_password value                    Admin password will be create at first start
   --root_keys value [ --root_keys value ]  Admin ssh-keys will be create at first start
   --help, -h                               show help (default: false)
   --version, -v                            print the version (default: false)
```

#### Docker

TBD

### Prerequisites

- go (>= 1.22+)
- git
- sqlite


## Built with

- [golang](https://go.dev/)
- [gin](https://github.com/gin-gonic/gin)
- [pico](https://picocss.com/docs)
- [sqlite](https://www.sqlite.org/)

## Contribution

Contributions are more than welcome! Thank you!

## Donations

Donations are more than welcome! Thank you!

1. BTC ``
1. ETH ``
1. TON ``
1. USDT (BEP20) ``
1. USDT (ETH20) ``