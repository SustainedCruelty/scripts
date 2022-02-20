package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unicode"

	_ "github.com/denisenkom/go-mssqldb"
)

type Column struct {
	ColumnName string
	DataType   string
}

var (
	selectColumns = "SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ?;"
	selectTables  = "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES GO;"
)

func main() {

	var password string
	flag.StringVar(&password, "p", "", "the database password")

	var server string
	flag.StringVar(&server, "s", "", "the database server")

	var user string
	flag.StringVar(&user, "u", "", "the database user")

	var dbName string
	flag.StringVar(&dbName, "d", "", "the database to connect to")

	var table string
	flag.StringVar(&table, "t", "", "the table to generate the class for")

	var fetchAll bool
	flag.BoolVar(&fetchAll, "all", false, "retrieve the classes from all tables")

	var output string
	flag.StringVar(&output, "o", "", "(optional) where to save the model(s)")

	var tableList string
	flag.StringVar(&tableList, "l", "", "(optional) reads table names from a text file")

	var namespace string
	flag.StringVar(&namespace, "n", "Models", "what namespace the csharp class should belong to")

	flag.Parse()

	// append a slash to the folder path if it doesnt have one already
	if output != "" && output[len(output)-1] != '/' {
		output += "/"
	}

	// create the folders
	if output != "" {
		if err := os.MkdirAll(output, os.ModeDir); err != nil {
			log.Fatalf("failed to create folder(s): %s", err.Error())
		}
	}

	tables := make([]string, 0)

	if table != "" {
		tables = append(tables, table)
	}

	// read the text file if it was supplied
	if tableList != "" {
		file, err := os.Open(tableList)
		if err != nil {
			log.Fatalf("Failed to open %s: %s", tableList, err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tables = append(tables, scanner.Text())
		}
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, 1433, dbName)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("failed to establish a connection to the database: %s", err.Error())
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Error with the database connection: %s", err.Error())
	}
	// setting up the database connection

	log.Println("successfully connected to the database")

	// read all tables from the database if the -all flag is set
	if fetchAll {
		allTables, err := getAllTables(db)
		if err != nil {
			log.Fatalf("failed to retrieve all table names from the database: %s", err.Error())
		}
		tables = append(tables, allTables...)
	}

	for _, t := range tables {
		log.Printf("generating class for table %s", t)
		cols, err := getColumnsFromTable(t, db)
		if err != nil {
			log.Fatalf("Failed to retrieve column information from the database: %s", err.Error())
		}
		class := generateTable(namespace, t, cols)
		ioutil.WriteFile(output+t+".cs", []byte(class), 0644)
	}

}

func getColumnsFromTable(tableName string, db *sql.DB) ([]Column, error) {
	var res []Column

	stmt, err := db.Prepare(selectColumns)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tableName)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var col Column
		if err := rows.Scan(&col.ColumnName, &col.DataType); err != nil {
			return res, err
		}
		res = append(res, col)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func generateTable(namespace string, tableName string, columns []Column) string {

	// add imports
	res := "using Newtonsoft.Json;\nusing System;\n\n"
	// add namespace + class declaration
	res += fmt.Sprintf("namespace %s\n{\n\t[Serializable]\n\tpublic class %s\n\t{", namespace, tableName)

	// add properties
	for _, col := range columns {
		// data annotation
		res += fmt.Sprintf("\n\t\t[JsonProperty(\"%s\")]", firstCharLower(col.ColumnName))
		// property
		res += fmt.Sprintf("\n\t\tpublic %s %s { get; set; }", convertTypes(col.DataType), col.ColumnName)
	}
	// closing brackets
	res += "\n\t}\n}\n"

	return res
}

func convertTypes(dtype string) string {
	switch typeConv := dtype; typeConv {
	case "int":
		return "int"
	case "bigint":
		return "long"
	case "datetime":
		return "DateTime"
	case "bit":
		return "bool"
	default:
		return "string"
	}
}

func firstCharLower(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func getAllTables(db *sql.DB) ([]string, error) {
	var res []string

	stmt, err := db.Prepare(selectTables)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return res, err
		}
		res = append(res, table)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}
