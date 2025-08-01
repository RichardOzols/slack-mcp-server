# Heroku Deployment Guide

This document provides instructions for deploying the Slack MCP Server to Heroku.

## Prerequisites

- Heroku account
- Heroku CLI installed
- Git repository access

## Quick Deploy

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## Manual Deployment

### 1. Clone the repository

```bash
git clone https://github.com/korotovsky/slack-mcp-server.git
cd slack-mcp-server
```

### 2. Create a Heroku app

```bash
heroku create your-app-name
```

### 3. Set environment variables

Set the required Slack credentials (you need either xoxp OR both xoxc/xoxd):

```bash
# Option 1: User OAuth token
heroku config:set SLACK_MCP_XOXP_TOKEN=xoxp-your-token

# Option 2: Browser tokens (alternative to xoxp)
heroku config:set SLACK_MCP_XOXC_TOKEN=xoxc-your-token
heroku config:set SLACK_MCP_XOXD_TOKEN=xoxd-your-token

# Optional: SSE API key for authentication
heroku config:set SLACK_MCP_SSE_API_KEY=your-api-key

# Optional: Log level (default: info)
heroku config:set SLACK_MCP_LOG_LEVEL=info

# Optional: Enable message posting (disabled by default for safety)
heroku config:set SLACK_MCP_ADD_MESSAGE_TOOL=true
```

### 4. Deploy

```bash
git push heroku main
```

### 5. Access your server

Your Slack MCP Server will be available at:
```
https://your-app-name.herokuapp.com/sse
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SLACK_MCP_XOXC_TOKEN` | Yes* | Slack browser token (xoxc-...) |
| `SLACK_MCP_XOXD_TOKEN` | Yes* | Slack browser cookie d (xoxd-...) |
| `SLACK_MCP_XOXP_TOKEN` | Yes* | User OAuth token (xoxp-...) - alternative to xoxc/xoxd |
| `SLACK_MCP_SSE_API_KEY` | No | Bearer token for SSE transport authentication |
| `SLACK_MCP_LOG_LEVEL` | No | Log level: debug, info, warn, error, panic, fatal (default: info) |
| `SLACK_MCP_ADD_MESSAGE_TOOL` | No | Enable message posting. Set to `true` for all channels, comma-separated channel IDs for specific channels, or leave empty to disable (default: disabled) |

*You need either `SLACK_MCP_XOXP_TOKEN` **or** both `SLACK_MCP_XOXC_TOKEN` and `SLACK_MCP_XOXD_TOKEN`.

## Configuration Notes

- The server automatically runs in SSE mode on Heroku (required for HTTP endpoints)
- Host is automatically set to `0.0.0.0` for Heroku compatibility
- Port is automatically configured using Heroku's `PORT` environment variable
- Logs are output in JSON format for better Heroku log aggregation

## Authentication Setup

For detailed authentication setup instructions, see:
- [Authentication Setup](docs/01-authentication-setup.md)

## Security Considerations

- Never commit tokens to your repository
- Use Heroku's config vars to set sensitive environment variables
- Consider enabling message posting only for specific channels if needed
- Review the [Security](../SECURITY.md) documentation

## Troubleshooting

### Application not starting
- Check that you have set the required Slack tokens
- Verify your tokens are valid using the authentication guide
- Check the Heroku logs: `heroku logs --tail`

### Cache warnings
- The server will warn about caches during startup - this is normal
- Users and channels caches will be populated on first API calls

### Connection issues
- Ensure your Heroku app is not sleeping (consider using Heroku's always-on dyno)
- Check that your SSE client is connecting to the correct endpoint: `https://your-app.herokuapp.com/sse`