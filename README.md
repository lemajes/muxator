```text
███    ███ ██    ██ ██   ██  █████  ████████  ██████  ██████  
████  ████ ██    ██  ██ ██  ██   ██    ██    ██    ██ ██   ██ 
██ ████ ██ ██    ██   ███   ███████    ██    ██    ██ ██████  
██  ██  ██ ██    ██  ██ ██  ██   ██    ██    ██    ██ ██   ██ 
██      ██  ██████  ██   ██ ██   ██    ██     ██████  ██   ██
```
---
### What is it ?
* It serves a port on local machine as a socks proxy
* It create multiple Tor connections.
* It routes each request through a different Tor connection.

---

## How to use
```bash
git clone https://github.com/lemajes/muxator.git
cd muxator
go get muxator/tor
go get muxator/socks
go run .
##
```
