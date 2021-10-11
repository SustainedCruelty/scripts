### mentibot
spam mentimeter wordclouds with words from a textfile
### Usage
```
Usage of mentibot.exe:
  -c string
        the code youre supposed to enter on menti.com
  -r int
        how often to loop through the file (default 1)
  -w string
        file that contains the words to be added (default "words.txt")
```
### Example
```
> go run mentibot.go -c 47313558 -r 2 -w menti_words.txt
```
### Compiling
Compile your own version from source using go build
```
go build mentibot.go
```
### Python Version
This repo also contains a python version of the script doing the exact same thing. It doesn't have any command line paramaters so you just have to follow the instructions

