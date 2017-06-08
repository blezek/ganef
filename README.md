# ganef

**ganef**: noun *|ˈgä-nəf|* THIEF, RASCAL variant *goniff*

ganef is a simple program to query [sqlite](https://www.sqlite.org/) databases using [Go's templates](https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html).  ganef takes a sqlite database and a template as input and returns the output of executing the template against the database.  ganef adds the `query` function to the template.  `query` takes an [sqlite select statement](http://www.sqlitetutorial.net/sqlite-select/) and returns an array.  The array contains a map of `{column: value}` for each row from the query.

# Install

``` bash
go get github.com/blezek/gonof
```

# Usage

```
# Get an example database
wget -O chinook.zip 'http://www.sqlitetutorial.net/download/sqlite-sample-database/?wpdmdl=94'
unzip chinook.zip

# Create a template
cat <<EOF > template.txt
All tracks
{{ range query "select TrackId, Name, Composer, unitprice from tracks limit 20" }}
TrackID: {{.TrackId}}
Name: {{.Name}}
Composer: {{.Composer}}
UnitPrice: {{.UnitPrice}}
{{end}}
EOF

ganif chinook.db template.txt -

```

# Building

In a clone of the repo (kudos to the fine [Hellogopher](https://github.com/cloudflare/hellogopher)):

``` bash
make
```

