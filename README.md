# ganef

**ganef**: noun *\ˈgä-nəf\* THIEF, RASCAL variant *goniff*

ganef is a simple program to query [sqlite](https://www.sqlite.org/) databases using [Go's templates](https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html).  ganef takes a sqlite database and a template as input and returns the output of executing the template against the database.  ganef adds the `query` function to the template.  `query` takes an [sqlite select statement](http://www.sqlitetutorial.net/sqlite-select/) and returns an array.  The array contains a map of `{column: value}` for each row from the query.



## Build

In the wild use:

``` bash
go get github.com/blezek/gonof
```

In a clone of the repo (kudos to the fine [Hellogopher](https://github.com/cloudflare/hellogopher)):

``` bash
make
```

