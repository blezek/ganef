package main

import (
	"database/sql"
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var usage string = `usage: gonof [options] <sqlitedb> <template> <output>

options:
  -h help
  -n use nullable strings
`

var db *sql.DB
var indicateNull bool = false

type rowmap map[string]interface{}

func doQuery(q string) ([]rowmap, error) {

	r := make([]rowmap, 0)
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	typeNames := make([]string, len(cols))
	for i, _ := range types {
		typeNames[i] = strings.ToLower(types[i].DatabaseTypeName())
	}
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columnPointers {
			// log.Printf("Looking at %s - type is %v", cols[i], typeNames[i])
			t := typeNames[i]
			if strings.Contains(t, "text") || strings.Contains(t, "varchar") || strings.Contains(t, "char") || strings.Contains(t, "date") {
				columnPointers[i] = new(sql.NullString)
			} else if strings.Contains(t, "double") || strings.Contains(t, "float") {
				columnPointers[i] = new(sql.NullFloat64)
			} else if strings.Contains(t, "integer") {
				columnPointers[i] = new(sql.NullInt64)
			} else {
				columnPointers[i] = new(interface{})
			}
			// columnPointers[i] = reflect.New(types[i].ScanType())
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(rowmap)
		for i, colName := range cols {
			var err error
			var vv driver.Value
			switch columnPointers[i].(type) {
			case *sql.NullString:
				v := columnPointers[i].(*sql.NullString)
				if indicateNull {
					m[colName] = *(v)
				} else {
					vv, err = v.Value()
					m[colName] = vv
				}
			case *sql.NullFloat64:
				v := columnPointers[i].(*sql.NullFloat64)
				if indicateNull {
					m[colName] = *(v)
				} else {
					vv, err = v.Value()
					m[colName] = vv
				}
			case *sql.NullInt64:
				v := columnPointers[i].(*sql.NullInt64)
				if indicateNull {
					m[colName] = *(v)
				} else {
					vv, err = v.Value()
					m[colName] = vv
				}
			}
			// t := types[i].DatabaseTypeName()
			// if t == "text" {

			// m[colName] = *(columnPointers[i].(*string))
			// val := columnPointers[i].(*interface{})
			// m[colName] = *val
			// }
			if err != nil {
				log.Fatalf("could not convert column %s: %s", colName[i], err.Error())
			}
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		r = append(r, m)
	}
	return r, err
}

func main() {

	help := false
	flag.BoolVar(&help, "h", false, "get help for the application")
	flag.BoolVar(&indicateNull, "n", false, "return indicators of null to template (eg https://golang.org/pkg/database/sql/#NullString), default is to use zero values")
	flag.Parse()

	if help {
		fmt.Println(usage)
		os.Exit(0)
	}

	if len(flag.Args()) < 3 {
		log.Fatal(usage)
	}

	dbFilename := flag.Args()[0]
	templateFilename := flag.Args()[1]
	outputFilename := flag.Args()[2]
	var err error
	db, err = sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Fatalf("Failed to open sqlite DB %v, %v", dbFilename, err)
	}

	// Add some helper functions
	var funcs = template.FuncMap{
		"query": doQuery,
		// "json": func(v interface{}) (string, error) {
		// 	a, err := json.Marshal(v)
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	return string(a), nil
		// },
		// "humanizeTime": humanize.Time,
		"now": time.Now,
		// 	"markdown": func(s string) template.HTML {
		// 		return template.HTML(string(blackfriday.MarkdownCommon([]byte(s))))
		// 	},
	}

	var t = template.New(templateFilename).Funcs(funcs)
	t = template.Must(t.ParseFiles(templateFilename))
	out, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("Failed to open %v, %v", outputFilename, err)
	}

	data := map[string]interface{}{
		"db": db,
	}
	err = t.Execute(out, data)
	if err != nil {
		log.Fatalf("Failed to execute template %v", err)
	}

}
