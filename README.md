# ganef

**ganef**: noun *|ˈgä-nəf|* THIEF, RASCAL variant *goniff*

ganef is a simple program to query [sqlite](https://www.sqlite.org/) databases using [Go's templates](https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html).  ganef takes a sqlite database and a template as input and returns the output of executing the template against the database.  ganef adds the `query` function to the template.  `query` takes an [sqlite select statement](http://www.sqlitetutorial.net/sqlite-select/) and returns an array.  The array contains a map of `{column: value}` for each row from the query.

# Install

Go 1.8 or higher required.

``` bash
go get github.com/blezek/ganef
```

# Usage

## Download the example DB

```
# Get an example database
wget -O chinook.zip 'http://www.sqlitetutorial.net/download/sqlite-sample-database/?wpdmdl=94'
unzip chinook.zip
```

## List all the tracks

Note, the SQLite driver returns the column names in the *actual case* used to create the table, not what from the query string.

```
# Create a template
cat <<EOF > template.txt
All tracks
{{ range query "select TrackId, Name, Composer, unitprice from tracks" }}
TrackID: {{.TrackId}}
Name: {{.Name}}
Composer: {{.Composer}}
UnitPrice: {{.UnitPrice}}
{{end}}
EOF

ganef chinook.db template.txt -

```

## Pass variables from the command line

Variables can be passed using multiple `-v key=value` flags on the command line.  They are passed into the template as `{{.key}}`.

```
# Create a template
cat <<EOF > album_detail.txt
{{ range printf "select * from albums where albumid = %v" .album | query -}}
Detail for album "{{.Title}}"
{{end }}
{{ range printf "SELECT name, milliseconds, bytes, albumid FROM tracks WHERE albumid = %v;" .album | query -}}
Name: {{.Name}}
Length(ms): {{.Milliseconds}}
Size(bytes): {{.Bytes}}

{{end}}
EOF

ganef -v album=279 chinook.db album_detail.txt -
ganef -v album=275 chinook.db album_detail.txt -

```


# Building

In a clone of the repo (kudos to the fine [Hellogopher](https://github.com/cloudflare/hellogopher)):

``` bash
make
```

