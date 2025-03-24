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
- 🔐 HTTPS/TLS support for encrypted communication
- 👮‍♀️ Prevent Dot Files Access (e.g., .env, .gitignore)
- 👮‍♀️ Prevent Symlink Access
- 🚀 Graceful shutdown on termination signals


How To 
======

Ecco la sezione **Installation** in inglese per il README del tuo progetto:  

---

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


Basic Authentication
--------------------

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

Custom HTTP Response Headers
----------------------------  

`drop` allows you to define custom HTTP response headers based on file request patterns.  


Headers are defined in a custom `.headers` file (similar to an `.ini` format). 

The general structure follows:  

- **global headers** (applied to all responses)
- **pattern-based headers** (applied only to matching file paths)


**How it works**

- patterns follow **glob-style matching** (e.g., `*.js` matches all JavaScript files)
- if a request matches multiple patterns, only the first match is applied
- global headers are **always applied first**, followed by any matching pattern-specific headers


**Example Configuration**  

```ini
# Global headers (applied to all responses)
X-Greeting: Hello World!

# Pattern-based headers
[*.mod]
X-Type: Go Module File
```
