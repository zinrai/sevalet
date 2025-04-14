# Sevalet

Sevalet ( Server valet ) is a Go application for securely executing commands on a server via HTTP API. It only executes allowed commands and is useful when you need to remotely execute commands on a server via API.

## Features

- Execute commands on server via HTTP API
- Only allowed commands and arguments can be executed
- Restrict commands and arguments through configuration files
- Timeout functionality

## Installation

```bash
$ go install github.com/zinrai/sevalet/cmd/server@latest
```

## Configuration

Configure executable commands and arguments in the `configs/commands.yaml` file.

### Starting

```bash
$ ./sevalet -config configs/app.yaml
```

## API Endpoints

### Execute Command

**Request:**

```http
POST /execute HTTP/1.1
Content-Type: application/json

{
  "command": "ls",
  "args": ["-la", "/tmp"],
  "timeout": 30
}
```

**Success Response:**

```json
{
  "api": {
    "status": "success",
    "code": "SUCCESS",
    "message": "Request processed successfully"
  },
  "command": {
    "executed": true,
    "exit_code": 0,
    "stdout": "output result",
    "stderr": "",
    "execution_time": "0.234s"
  }
}
```

### List Allowed Commands

**Request:**

```http
GET /commands HTTP/1.1
```

**Response:**

```json
{
  "api": {
    "status": "success",
    "code": "SUCCESS",
    "message": "Request processed successfully"
  },
  "commands": [
    {
      "name": "ls",
      "description": "List directory contents",
      "allowed_args": ["-l", "-a", "-h", "/tmp", "/var/log"]
    },
    {
      "name": "cat",
      "description": "Display file contents",
      "allowed_args": ["/var/log/syslog", "/etc/hostname"]
    }
  ]
}
```

## Security Notes

- Sevalet does not provide authentication. Implement authentication using a reverse proxy or similar.
- The application runs with the permissions of the user executing it, so apply appropriate permission restrictions.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
