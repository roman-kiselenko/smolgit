## smolgit

**smolgit** offers a minimalist interface with essential [git](https://git-scm.com/) functionalities, making it perfect for small teams or individual developers. There is no fancy javascript frontend or super duper features.

### Motivation
**smolgit** is perfect for those who value simplicity and efficiency in their workflow. The goal is to eliminate the bloat and offer a tool that is easy to set up, use, and maintain, without sacrificing the critical capabilities that developers rely on.

### Usage

1. **git operations**: Easily perform `pull`, `push`, and `fetch` operations from the main interface.
2. **repository visualization**: Browse the file tree, view the latest commit, and explore the commit and branch and tag lists.
3. **file content**: See the content of non-binary files.


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

### Prerequisites

- go (>= 1.22+)
- git
- sqlite


## Built with

- [gin](https://github.com/gin-gonic/gin)
- [hack](https://hackcss.egoist.dev/)
- [sqlite](https://www.sqlite.org/)

## Contact

For questions or feedback, please open an issue or contact us at [roman.kiselenko.dev@gmail.com](roman.kiselenko.dev@gmail.com).
