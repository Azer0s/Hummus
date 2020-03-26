rm -rf ./bin || true
mkdir bin || true
cp -r ./stdlib/ ./bin/stdlib/ || true

# Build stdlib native calls

go build -buildmode=plugin -o bin/stdlib/str/calls.so bin/stdlib/str/calls.go
rm bin/stdlib/str/calls.go

go build -buildmode=plugin -o bin/stdlib/pipe/calls.so bin/stdlib/pipe/calls.go
rm bin/stdlib/pipe/calls.go

go build -buildmode=plugin -o bin/stdlib/sync/calls.so bin/stdlib/sync/calls.go
rm bin/stdlib/sync/calls.go