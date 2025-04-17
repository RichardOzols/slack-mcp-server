# Slack MCP Server

Model Context Protocol (MCP) server for Slack Workspaces. This integration supports both Stdio and SSE transports, proxy settings and does not require any permissions or bots being created or approved by Workspace admins 😏.

### Feature Demo

![ezgif-316311ee04f444](https://github.com/user-attachments/assets/35dc9895-e695-4e56-acdc-1a46d6520ba0)

## Setup Guide

### 1. Authentication Setup

Open up your Slack in your browser and login.

#### Lookup `SLACK_MCP_XOXC_TOKEN`

- Open your browser's Developer Console.
- In Firefox, under `Tools -> Browser Tools -> Web Developer tools` in the menu bar
- In Chrome, click the "three dots" button to the right of the URL Bar, then select
`More Tools -> Developer Tools`
- Switch to the console tab.
- Type "allow pasting" and press ENTER.
- Paste the following snippet and press ENTER to execute:
  `JSON.parse(localStorage.localConfig_v2).teams[document.location.pathname.match(/^\/client\/([A-Z0-9]+)/)[1]].token`

Token value is printed right after the executed command (it starts with
`xoxc-`), save it somewhere for now.

#### Lookup `SLACK_MCP_XOXD_TOKEN`

 - Switch to "Application" tab and select "Cookies" in the left navigation pane.
 - Find the cookie with the name `d`.  That's right, just the letter `d`.
 - Double-click the Value of this cookie.
 - Press Ctrl+C or Cmd+C to copy it's value to clipboard.
 - Save it for later.

### 2. Installation

Choose one of these installation methods:

#### 2.1. Docker

For detailed information about all environment variables, see [Environment Variables](https://github.com/korotovsky/slack-mcp-server?tab=readme-ov-file#environment-variables).

```bash
export SLACK_MCP_XOXC_TOKEN=xoxc-...
export SLACK_MCP_XOXD_TOKEN=xoxd-...

docker pull ghcr.io/korotovsky/slack-mcp-server:latest
docker run -i --rm \
  -e SLACK_MCP_XOXC_TOKEN \
  -e SLACK_MCP_XOXD_TOKEN \
  slack-mcp-server --transport stdio
```

#### 2.2. Docker Compose

```bash
wget -O docker-compose.yml https://github.com/korotovsky/slack-mcp-server/releases/latest/download/docker-compose.yml
wget -O .env https://github.com/korotovsky/slack-mcp-server/releases/latest/download/.env.dist
nano .env # Edit .env file with your tokens from step 1 of the setup guide
docker-compose up -d
```

#### 2.2.1 TLS and Exposing to the Internet

There are several reasons why you might need to setup HTTPS for your SSE.
 - `mcp-remote` is capable to handle only https schemes;
 - it is generally a good practice to use TLS for any service exposed to the internet;

You could use `ngrok`:

```bash
ngrok http 3001
```

and then use the endpoint `https://903d-xxx-xxxx-xxxx-10b4.ngrok-free.app` for your `mcp-remote` argument.

### 3. Configuration and Usage

You can configure the MCP server using command line arguments and environment variables.

Add the following to your `claude_desktop_config.json`:

#### Option 1 with `stdio` transport:
```json
{
  "mcpServers": {
    "slack": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "SLACK_MCP_XOXC_TOKEN=$SLACK_MCP_XOXC_TOKEN",
        "-e",
        "SLACK_MCP_XOXD_TOKEN=$SLACK_MCP_XOXD_TOKEN",
        "ghcr.io/korotovsky/slack-mcp-server",
        "mcp-server",
        "--transport",
        "stdio"
      ],
      "env": {
        "SLACK_MCP_XOXC_TOKEN": "xoxc-...",
        "SLACK_MCP_XOXD_TOKEN": "xoxd-..."
      }
    }
  }
}
```

#### Option 2 with `sse` transport:

Complete steps from 2.2 and run `docker compose up -d` to launch MCP server or with your preferred method and then configure it:

```json
{
  "mcpServers": {
    "slack": {
      "command": "npx",
      "args": [
        "-y",
        "mcp-remote",
        "https://x.y.z.q:3001/sse",
        "--header",
        "Authorization: Bearer ${SLACK_MCP_SSE_API_KEY}"
      ],
      "env": {
        "SLACK_MCP_SSE_API_KEY": "my-$$e-$ecret"
      }
    }
  }
}
```

#### Option 3 with `sse` transport on Windows:

Complete steps from 2.2 and run `docker compose up -d` to launch MCP server or with your preferred method and then configure it:

```json
{
  "mcpServers": {
    "slack": {
      "command": "C:\\Progra~1\\nodejs\\npx.cmd",
      "args": [
        "-y",
        "mcp-remote",
        "https://x.y.z.q:3001/sse",
        "--header",
        "Authorization: Bearer ${SLACK_MCP_SSE_API_KEY}"
      ],
      "env": {
        "SLACK_MCP_SSE_API_KEY": "my-$$e-$ecret"
      }
    }
  }
}
```

#### Console Arguments

| Argument              | Required ? | Description                                                              |
|-----------------------|------------|--------------------------------------------------------------------------|
| `--transport` or `-t` | Yes        | Select transport for the MCP Server, possible values are: `stdio`, `sse` |

#### Environment Variables

| Variable                       | Required ? | Default     | Description                                                                   |
|--------------------------------|------------|-------------|-------------------------------------------------------------------------------|
| `SLACK_MCP_XOXC_TOKEN`         | Yes        | `nil`       | Authentication data token field `token` from POST data field-set (`xoxc-...`) |
| `SLACK_MCP_XOXD_TOKEN`         | Yes        | `nil`       | Authentication data token from cookie `d` (`xoxd-...`)                        |
| `SLACK_MCP_SERVER_PORT`        | No         | `3001`      | Port for the MCP server to listen on                                          |
| `SLACK_MCP_SERVER_HOST`        | No         | `127.0.0.1` | Host for the MCP server to listen on                                          |
| `SLACK_MCP_SSE_API_KEY`        | No         | `nil`       | Authorization Bearer token when `transport` is `sse`                          |
| `SLACK_MCP_PROXY`              | No         | `nil`       | Proxy URL for the MCP server to use                                           |
| `SLACK_MCP_SERVER_CA`          | No         | `nil`       | Path to the CA certificate of the trust store                                 |
| `SLACK_MCP_SERVER_CA_INSECURE` | No         | `false`     | Trust all insecure requests (NOT RECOMMENDED)                                 |

## Available Tools

| Tool                   | Description                   |
|------------------------|-------------------------------|
| `conversationsHistory` | Get messages from the channel |

### Debugging Tools

```bash
# Run the inspector with stdio transport
npx @modelcontextprotocol/inspector go run mcp/mcp-server.go --transport stdio

# View logs
tail -n 20 -f ~/Library/Logs/Claude/mcp*.log
```

## Security

- Never share API tokens
- Keep .env files secure and private

## License

Licensed under MIT - see [LICENSE](LICENSE) file. This is not an official Slack product.
