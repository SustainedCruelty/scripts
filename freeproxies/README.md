### freeproxies
web scrapes a list of http and https proxies and saves them to a textfile
### Usage
```bash
Usage of freeproxies.exe:
  -o string
        what file to save the proxies to (default "proxies.txt")
  -p string
        whether to fetch http or https proxies (default "http")
```
### Compile
Build your own version from source
```bash
go build freeproxies.go freeproxylist.go geonode.go proxyscrape.go
```
### Current sources
- https://free-proxy-list.net/
- https://geonode.com/free-proxy-list
- https://proxyscrape.com/free-proxy-list
