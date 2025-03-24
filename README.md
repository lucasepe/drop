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

- [x] 📂 Serves static files from a specified directory
- [x] 📑 Automatically generates a stylish index if index.html is missing
- [x] 📜 Consistent MIME type resolution across different environments
- [x] 👀 Access Log
- [x] 🔒 Basic Authentication for access
- [ ] 🧩 Customizable HTTP response headers for specific file
- [x] 🔐 HTTPS/TLS support for encrypted communication
- [x] 👮‍♀️ Prevent Dot Files Access (e.g., .env, .gitignore)
- [x] 👮‍♀️ Prevent Symlink Access
- [x] 🚀 Graceful shutdown on termination signals


How To 
======

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
