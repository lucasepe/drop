```
┌┬┐┬─┐┌─┐┌─┐
 ││├┬┘│ │├─┘
─┴┘┴└─└─┘┴  
```

> Lightweight and secure HTTP server for hosting static files from a specified directory

Overview
========

Drop is a lightweight and secure HTTP server for hosting static files from a specified directory.

This project is useful for various scenarios, including:

- Testing WebAssembly (WASM) applications - without the need for a complex web server
- Sharing files between machines - over a local network
- Hosting simple static websites - for development purposes
- Providing a lightweight file access point - for devices in an IoT network

Features
========

- 📂 Serves static files from a specified directory
- 📑 Automatically generates a stylish index if index.html is missing
- 📜 Consistent MIME type resolution across different environments
- 👀 Access Log
- 🔒 Basic Authentication for access
- 🧩 Customizable HTTP response headers for specific file
- 🔥 Dynamic HTTP response headers for specific file
- 🔐 HTTPS/TLS support for encrypted communication
- 👮‍♀️ Prevent Dot Files Access (e.g., .env, .gitignore)
- 👮‍♀️ Prevent Symlink Access
- 📡 Support for OPTIONS requests, returning allowed HTTP methods
- ⚡ Proper handling of HEAD requests (returns headers like Content-Type and Content-Length plus your custom headers)
- ⛔ Blocks unsupported HTTP methods (POST, PUT, DELETE, etc.) with 405 Method Not Allowed
- 🚀 Graceful shutdown on termination signals

# How To

## Installation

You can install `drop` using different methods depending on your operating system.  

### macOS (Homebrew)

The easiest way to install `drop` on macOS is via Homebrew:  

```sh
brew install lucasepe/cli-tools/drop
```

### Windows & Linux 

#### Download Prebuilt Binaries

- Go to the [Releases](https://github.com/lucasepe/drop/releases) section.  
- Download the latest binary for your OS and Arch.  
- Add it to your system's `$PATH` if necessary.  

#### Alternative: Install via `go install` (requires Go installed)

If you have Go installed, you can install `drop` directly with:  

```sh
go install github.com/lucasepe/drop@latest
```  

This will place the binary in `$GOPATH/bin` (or `$HOME/go/bin` if `$GOPATH` is not set).

Ensure this path is in your `$PATH` to use `drop` globally.  


## Basic Authentication

To enable Basic Authentication put into the serving folder an `.users` file.

This is a flat file that contains the user name and the SHA-256 crypt hashed password for each user.

`.users` file sample:

```
admin:$5$azZ$NH//nNpYkwzlwe03A4ZmLxZz0lQTmJ0Ongj9KIfC6o6
```

You can generate each row using openssl:

```
printf "admin:$(openssl passwd -5 -salt 'azZ' '12345')\n"
```

## Custom HTTP Response Headers

`drop` allows you to define custom HTTP response headers based on file request patterns.  

Headers are defined in a custom `.headers` file (similar to an `.ini` format). 

The general structure follows:  

- **global headers** (applied to all responses)
- **pattern-based headers** (applied only to matching file paths)


**How it works**

- patterns follow **glob-style matching** (e.g., `*.js` matches all JavaScript files)
- if a request matches multiple patterns, only the first match is applied
- global headers are **always applied first**, followed by any matching pattern-specific headers


**Example**  

```yaml
# Global headers (applied to all responses)
X-Greeting: Hello World!

# Pattern-based headers
[*.mod]
X-Type: Go Module File
```

## Custom dynamic HTTP Headers 

`drop` allows also dynamic values based on the request and server state. 

Here’s a list of possible variables:  

| **Variable** | **Description** | **Example Value** | **Use Case**           |
|--------------|-----------------|-------------------|------------------------|
| `{{ SERVER_ADDR }}` | Server IP and port | `192.168.1.10:8080` | CSP, CORS  |
| `{{ SERVER_NAME }}` | Server hostname | `example.com` | CSP, logging |
| `{{ REMOTE_ADDR }}` | Client IP address | `203.0.113.45` | CORS, security         |
| `{{ REQUEST_URI }}` | Full request URI (path + query string) | `/static/libs/hello.wasm?a=xxxx`| Debugging |
| `{{ REQUEST_PATH }}` | Request path only | `/libs/hello.wasm` | CSP, CORS |
| `{{ USER_AGENT }}`   | Client's User-Agent | `Mozilla/5.0 (Windows NT 10.0...)` | Security, analytics |
| `{{ REFERER }}`      | Referrer URL of the request | `https://google.com/` | Security, analytics |

**Example**

```yaml
# https://content-security-policy.com/examples/
Content-Security-Policy: default-src 'self'; img-src 'self'; style-src 'unsafe-inline' http://${SERVER_ADDR} https://${SERVER_ADDR} https://cdnjs.cloudflare.com;
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block

# No Cache (for debug)
[*.wasm]
Cache-Control: no-store, no-cache, must-revalidate, max-age=0
Pragma: no-cache
Expires: 0
```