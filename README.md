# goCsInspect

A high-performance, extensible tool for inspecting Counter-Strike (CS:GO/CS2) skin/item data via the Steam Game Coordinator. goCsInspect provides a REST API and command-line utilities to fetch, parse, and detail item information, supporting multi-client management, caching, and automated game file updates.

---

## Features

- **REST API** for item/skin inspection via inspect links.
- **Multi-client management** for parallel Steam account usage.
- **Automated game file updates** (items_game.txt, csgo_english.txt) from upstream sources.
- **Detailed item parsing** using protobuf and CS:GO item parser.
- **Persistent storage** of inspection results (SQLite).
- **Configurable request TTLs, cooldowns, and logging.**
- **Extensible architecture** for adding new features or storage backends.

---

## Installation

### Prerequisites

- Go 1.24+
- (Optional) SQLite3 for persistent storage

### Build

Clone the repository and build the binaries:

```sh
git clone https://github.com/volodymyrzuyev/goCsInspect.git
cd goCsInspect
make build_all
```

This will produce binaries in the `bin/` directory.

---

## Usage

### 1. Configure

Copy the example config and edit as needed:

```sh
cp config.yaml.example config.yaml
```

Edit `config.yaml` to add your Steam accounts and adjust settings (see below).

### 2. Run the REST API

```sh
bin/cmd/goCsInspectAPI/goCsInspect --config config.yaml
```

The API will bind to the address specified in your config (default: `0.0.0.0:8080`).

### 3. Fetch Protobuf Data (Optional)

To fetch and cache protobuf data for test or batch processing:

```sh
bin/cmd/protoFetcher/goCsInspect --config config.yaml
```

---

## Configuration

All settings are managed via `config.yaml`. Key options:

- `accounts`: List of Steam accounts (username, password, 2FA/shared secret).
- `requestttl`: Timeout for unresolved requests.
- `clientcooldown`: Minimum delay between requests per client.
- `gameitemslocation`, `gamelanguagelocation`: Paths to CS game files.
- `autoupdategamefiles`: Enable/disable auto-updating of game files.
- `gamefilesautoupdateinverval`: Interval for auto-updates.
- `databasestring`: Path to SQLite database.
- `loglevel`: Logging verbosity (`DEBUG`, `INFO`, `WARN`, `ERROR`).
- `bindip`: IP and port for REST API.

See `config.yaml.example` for a full annotated example.

---

## API

The REST API exposes endpoints for inspecting items via inspect links. Example usage:

```
GET /?url=<inspect_link>
```

Returns detailed item information in JSON.

---

## Development & Testing

- Run all tests:  
  ```sh
  make test
  ```

- Regenerate SQL code (if schema changes):  
  ```sh
  make gen_sql
  ```

---

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or documentation improvements.

---

## License

MIT License. See [LICENSE](LICENSE) for details.
