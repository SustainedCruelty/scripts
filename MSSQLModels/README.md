# MSSQLModels

Converts the schema of a mssql database to csharp classes

## Usage
```
$ classgen.exe --help
Usage of classgen.exe:
  -all
        retrieve the classes from all tables
  -d string
        the database to connect to 
  -l string
        (optional) reads table names from a text file
  -n string
        what namespace the csharp class should belong to (default "Models")
  -o string
        (optional) where to save the model(s)
  -p string
        the database password
  -s string
        the database server
  -t string
        the table to generate the class for
  -u string
        the database user
```
## Building from source
Compile to an executable file using 'go build'
```
$ go build main.go
```

## TODO

- Support more than one dbms (e.g. mysql, postgres, oracle etc.)
- Add concurrency