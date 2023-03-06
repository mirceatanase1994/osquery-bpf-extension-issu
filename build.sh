echo "Building logging extension..."
cd log_ext/
CGO_ENABLED=0 go build -o bin/log.ext log_ext.go
cd ../

echo "Building table extension..."
cd table_ext/
CGO_ENABLED=0 go build -o bin/foobar.ext table_ext.go
cd ../

echo "Done!"