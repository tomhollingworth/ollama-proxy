# Ollama Proxy for Vast.ai

This is a simple Go-based HTTP proxy for forwarding all your Ollama API requests to a Vast.ai instance. It is designed to help you securely and conveniently access your Ollama server running on Vast.ai, handling authentication tokens and cookies automatically.

## Features
- Proxies all HTTP requests to your Ollama server on Vast.ai
- Handles authentication token exchange and cookies
- Simple to deploy and run anywhere you have Go installed

## Usage

### 1. Build the Proxy

```
go build -o ollama-proxy main.go
```

### 2. Run the Proxy

You must provide the `-address` flag with your Vast.ai Ollama endpoint, including the `token` query parameter. For example:

```
./ollama-proxy -address "http://vast-ai-ollama-host:11434?token=YOUR_TOKEN" -listen ":11434"
```

- `-address` (required): The Vast.ai Ollama endpoint, including the `token` parameter.
- `-listen` (optional): The local address to listen on (default is `:11434`).

### 3. Point Your Ollama Client

Configure your Ollama client or API consumer to use `http://localhost:11434` (or your chosen listen address) as the API endpoint. All requests will be securely proxied to your Vast.ai instance.

## Example

```
./ollama-proxy -address "http://123.45.67.89:11434?token=abcdef123456" -listen ":11434"
```

Then, use your Ollama client as usual, pointing it to `http://localhost:11434`.

## How It Works
- The proxy exchanges your token for a session cookie with the Vast.ai Ollama server.
- All incoming requests are forwarded to the remote Ollama server, preserving headers and paths.
- Responses are streamed back to your client.

## Requirements
- Go 1.18 or newer
- Access to a Vast.ai instance running Ollama

## Security Note
- Keep your Vast.ai token secure. Do not share it or commit it to version control.

## License
MIT
