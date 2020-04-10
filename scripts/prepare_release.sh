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

go build -buildmode=plugin -o bin/stdlib/debug/calls.so bin/stdlib/debug/calls.go
rm bin/stdlib/debug/calls.go

go build -buildmode=plugin -o bin/stdlib/base/io/calls.so bin/stdlib/base/io/calls.go
rm bin/stdlib/base/io/calls.go

go build -buildmode=plugin -o bin/stdlib/base/conversion/calls.so bin/stdlib/base/conversion/calls.go
rm bin/stdlib/base/conversion/calls.go

go build -buildmode=plugin -o bin/stdlib/base/collections/calls.so bin/stdlib/base/collections/calls.go
rm bin/stdlib/base/collections/calls.go

go build -buildmode=plugin -o bin/stdlib/base/enumerable/calls.so bin/stdlib/base/enumerable/calls.go
rm bin/stdlib/base/enumerable/calls.go

go build -buildmode=plugin -o bin/stdlib/net/http/calls.so bin/stdlib/net/http/calls.go
rm bin/stdlib/net/http/calls.go

go build -buildmode=plugin -o bin/stdlib/os/calls.so bin/stdlib/os/calls.go
rm bin/stdlib/os/calls.go

go build -buildmode=plugin -o bin/stdlib/regex/calls.so bin/stdlib/regex/calls.go
rm bin/stdlib/regex/calls.go

go build -buildmode=plugin -o bin/stdlib/log/calls.so bin/stdlib/log/calls.go bin/stdlib/log/std_hook.go
rm bin/stdlib/log/calls.go bin/stdlib/log/std_hook.go