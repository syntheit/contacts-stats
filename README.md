# Contacts Stats Visualization

A Go-based tool to analyze your VCF contact files or CardDAV addressbooks and visualize the geographic distribution of your contacts.

## Features

- **Sources**:
  - **CardDAV**: Fetches contacts directly from your Nextcloud/WebDAV server.
  - **VCF File**: Parses a local `.vcf` file.
- **Geographic Analysis**: Identifies countries and regions based on phone numbers.
- **Interactive Map**: Web-based visualization.
- **Dockerized**: Easy to deploy with Docker and Docker Compose.

## Usage

### Docker (Recommended)

1. **Configure**: Create a `.env` file based on the example below.

   ```bash
   CARDDAV_URL=https://nextcloud.example.com/remote.php/dav/addressbooks/users/daniel/contacts/
   CARDDAV_USER=your_username
   CARDDAV_PASSWORD=your_password
   ```

2. **Update Stats**: Run the update command to fetch contacts and generate `stats.json`.

   ```bash
   docker-compose run contacts-stats update
   ```

3. **Serve**: Start the web server.

   ```bash
   docker-compose up
   ```

   Visit [http://localhost:8080](http://localhost:8080).

### Local Development

Prerequisites: Go 1.25+ or Nix.

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run Update**:
   ```bash
   # From VCF
   go run main.go update --vcf contacts.vcf

   # From CardDAV (using env vars)
   export CARDDAV_URL=...
   export CARDDAV_USER=...
   export CARDDAV_PASSWORD=...
   go run main.go update
   ```

3. **Run Server**:
   ```bash
   go run main.go serve
   ```

## Configuration

| Environment Variable | Description |
|----------------------|-------------|
| `CARDDAV_URL`        | URL to the CardDAV addressbook collection (should end with `/`) |
| `CARDDAV_USER`       | Username for CardDAV auth |
| `CARDDAV_PASSWORD`   | Password/App Password for CardDAV auth |
| `VCF_PATH`           | Path to fallback VCF file if CardDAV is not used |
| `PORT`               | Web server port (default 8080) |

## Command Line Flags

Legacy flags (`--vcf`, `--serve`) are still supported but using the subcommands `update` and `serve` is recommended.
