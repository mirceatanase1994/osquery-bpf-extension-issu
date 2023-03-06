This is an osquery-go table extension is based on the [osquery-go](https://github.com/osquery/osquery-go) table extension example.

The differences to the example are:
* it creates 3 tables instead of one (foobar, foobar2, foobar3)
* for each table, at each generate function call, a random number between 1 and 10 of new table entries are generated
* it adds a "timestamp" column which contains the timestamp that the table entry was generated, as a string