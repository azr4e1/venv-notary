# linux

rm -rf ./build

echo "compiling for linux"
env GOOS=linux GOARCH=amd64 go build -o build/vn-linux-amd64 vn/main.go
env GOOS=linux GOARCH=arm64 go build -o build/vn-linux-arm64 vn/main.go

echo "compiling for macos"
env GOOS=darwin GOARCH=amd64 go build -o build/vn-darwin-amd64 vn/main.go
env GOOS=darwin GOARCH=arm64 go build -o build/vn-darwin-arm64 vn/main.go

# windows
GOOS=windows

echo "compiling for windows"
env GOOS=windows GOARCH=amd64 go build -o build/vn-windows-amd64.exe vn/main.go
env GOOS=windows GOARCH=arm64 go build -o build/vn-windows-arm64.exe vn/main.go
