# devctl

I want a developer interface that codifies rote logic execution. 

## Getting Started

### Run Tests

To run all tests:

```
go test ./...
```

### Run Locally

To run the CLI locally (from the project root):

```
go run main.go [command]
```

For example, to see available commands:

```
go run main.go --help
```

### Build

To build the project for your current OS/architecture:

```
make build
```

This will produce a binary named `devctl` in the `bin/` directory for your current system.

To cross-compile for another system (e.g., Linux amd64):

```
GOOS=linux GOARCH=amd64 make build
```

This will also produce a binary named `devctl` in the `bin/` directory, but for the specified target system.

---

### Todo

[ ] JWT generator 
[ ] persistent configuration
[ ] persistent secrets
[ ] plugin capability 

## Install from Pre-built Binary

You can install `devctl` by downloading a pre-built binary from the [Releases](https://github.com/danlafeir/devctl/releases) page of this repository.

1. Go to the [Releases](https://github.com/danlafeir/devctl/releases) page.
2. Download the binary for your operating system and architecture (e.g., `devctl-linux-amd64`, `devctl-darwin-arm64`, etc.).
3. (If the binary is compressed, extract it first.)
4. Make the binary executable:
   ```sh
   chmod +x devctl
   ```
5. Move it to a directory in your `PATH`, for example:
   ```sh
   sudo mv devctl /usr/local/bin/
   ```
6. Verify the installation:
   ```sh
   devctl --help
   ```

You can now use `devctl` from anywhere in your terminal. 

## Quick Install (Online Script)

You can install the latest `devctl` release automatically with a one-liner (Linux/macOS):

```sh
curl -sSL https://github.com/danlafeir/devctl/releases/latest/download/install.sh | sh
```

This script will detect your OS and architecture, download the correct binary, and install it to `/usr/local/bin` (you may be prompted for your password).

**Security tip:** Always review install scripts before piping to `sh`.

For more details or manual installation, see the [Releases](https://github.com/danlafeir/devctl/releases) page. 