## smolgit

**smolgit** offers a minimalist [git](https://git-scm.com/) server, making it perfect for small teams or individual developers. Its minimal simple and just works. Small memory footprint, one binary to go.

### Motivation
**smolgit** is perfect for those who value simplicity and efficiency in their workflow. The goal is to eliminate the bloat and offer a tool that is easy to set up, use, and maintain.

### Usage

1. **git operations**: Easily perform `pull`, `push`, and `fetch` operations. Create repository just by `push`ing.
1. **repository visualization**: Browse the file tree, view the latest commit, and explore the commit and branch and tag lists.
1. **file content**: See the content of non-binary files.
1. **user management**: Simple user management, create users with `ssh-keys` to allow them `push\pull`


### Getting Started

#### Compiled

Download binary from release page.

#### Build

Build from source by cloning repository and run `make build` to build binary then `./bin/smolgit`

`./bin/smolgit --help`:

```shell
GLOBAL OPTIONS:
   --ssh_addr value                         (default: ":3081")
   --log_color                              (default: true)
   --log_json                               (default: false)
   --log_level value                        (default: "DEBUG")
   --server_addr value                      (default: ":3080")
   --server_brand value                     (default: "smolgit")
   --git_path value                         (default: "./tmp/smolgit")
   --db_path value                          (default: "./database.sqlite")
   --root_login value                       (default: "root")
   --root_password value                    (default: "password")
   --root_keys value [ --root_keys value ]
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

- [gin](https://github.com/gin-gonic/gin)
- [pico](https://picocss.com/docs)
- [sqlite](https://www.sqlite.org/)

## Contribution

Contributions are more than welcome!
