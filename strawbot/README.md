### strawbot
simple voting bot for strawpoll.de using https proxies
### Usage
```
Usage of strawbot.exe:
  -c int
        what level of concurrency to use (default 10)
  -o string
        what option to vote for
  -p string
        file containing a list of https proxies (default "proxies.txt")
  -t int
        how many seconds until timeout for each request (default 20)
  -u string
        link to the strawpoll
```
### Example
```
> go run strawbot.go -u https://strawpoll.de/wc333bb -c 10 -t 20 -o opt2 -p httpsproxies.txt
[+] Successfully casted vote with proxy 103.119.55.216:8080
[+] Successfully casted vote with proxy 43.129.194.136:10008
[+] Successfully casted vote with proxy 20.81.106.180:8888
[+] Successfully casted vote with proxy 80.59.199.212:8080
[+] Successfully casted vote with proxy 149.28.91.128:1088
[+] Successfully casted vote with proxy 213.32.75.88:9300
                            ...
```
### Compiling
Compile your own version from source using go build
```
go build strawbot.go
```
### Python Version
This repo also contains a python version of the script doing the exact same thing. It doesn't have any command line paramaters so you just have to follow the instructions
```
> python strawbot.py
[?] Enter the poll's URL: https://strawpoll.de/wc333bb
[?] What option to vote for: opt3
[?] List of proxies to use: httpsproxies.txt
[?] How many seconds until timeout?: 20
[?] How many threads to use: 10

[*] Proceeding to vote for option 'opt3' with a total of 579 proxies, 10 threads and a timeout of 20 seconds
[?] Do you want to continue? [y/n]: y

[+] Successfully casted vote with proxy 31.170.175.192:53281 in 0.9 seconds
[+] Successfully casted vote with proxy 143.0.66.197:999 in 3.33 seconds
[+] Successfully casted vote with proxy 2.188.222.154:48562 in 4.03 seconds
[+] Successfully casted vote with proxy 82.99.232.18:58689 in 4.99 seconds
[-] Failed casting vote with proxy 83.171.103.67:8080 (proxy failure / timeout);
                                ...
```

