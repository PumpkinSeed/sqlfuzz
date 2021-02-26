.PHONY: all
all: slqfuzz_darwin_amd64 slqfuzz_windows_amd64 slqfuzz_linux_amd64 slqfuzz_linux_arm64

slqfuzz_darwin_amd64:
	env GOOS=darwin GOARCH=amd64 go build -o sqlfuzz_darwin_amd64 main.go

# Temporary removed
slqfuzz_darwin_arm64:
	env GOOS=darwin GOARCH=arm64 go build -o sqlfuzz_darwin_amd64 main.go

slqfuzz_windows_amd64:
	env GOOS=windows GOARCH=amd64 go build -o sqlfuzz_windows_amd64.exe main.go

slqfuzz_linux_amd64:
	env GOOS=linux GOARCH=amd64 go build -o sqlfuzz_linux_amd64 main.go

slqfuzz_linux_arm64:
	env GOOS=linux GOARCH=arm64 go build -o sqlfuzz_linux_arm64 main.go

clean:
	rm -rf sqlfuzz*