## SQLfuzz

[![Go Report Card](https://goreportcard.com/badge/github.com/PumpkinSeed/sqlfuzz)](https://goreportcard.com/report/github.com/PumpkinSeed/sqlfuzz) [![GoDoc](https://godoc.org/github.com/PumpkinSeed/sqlfuzz?status.svg)](https://godoc.org/github.com/PumpkinSeed/sqlfuzz) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org) ![sqlfuzz test workflow](https://github.com/PumpkinSeed/sqlfuzz/actions/workflows/test.yml/badge.svg)

Load random data into SQL tables for testing purposes. The tool can get the layout of the SQL table and fill it up with random data. 

- [Installation](#installation)
- [Usage](#usage)
- [Flags](#flags)
- [Package usage](#package-usage)

### Installation

#### MacOS

```
wget https://github.com/PumpkinSeed/sqlfuzz/releases/download/{RELEASE}/sqlfuzz_darwin_amd64 -O /usr/local/bin/sqlfuzz
chmod +x /usr/local/bin/sqlfuzz
```

#### Linux

```
# amd64 build
wget https://github.com/PumpkinSeed/sqlfuzz/releases/download/{RELEASE}/sqlfuzz_linux_amd64 -O /usr/local/bin/sqlfuzz
chmod +x /usr/local/bin/sqlfuzz

# arm64 build
wget https://github.com/PumpkinSeed/sqlfuzz/releases/download/{RELEASE}/sqlfuzz_linux_arm64 -O /usr/local/bin/sqlfuzz
chmod +x /usr/local/bin/sqlfuzz
```

#### Windows

You can download the Windows build [here](https://github.com/PumpkinSeed/sqlfuzz/releases/download/v0.3.0/sqlfuzz_windows_amd64.exe)

#### Build from source

```
wget https://github.com/PumpkinSeed/sqlfuzz/archive/{RELEASE}.zip
# unzip
# cd into dir
go install main.go
```

### Usage

```
# MySQL
sqlfuzz -u username -p password -d database -h 127.0.0.1 -t table -n 100000 -w 100

# Postgres
sqlfuzz -u username -p password -d database -h 127.0.0.1 -t table -n 100000 -w 100 -P 5432 -D postgres
```

#### Flags

- `u`: User for database connection
- `p`: Password for database connection
- `d`: Database name for database connection
- `h`: Host for database connection
- `P`: Port for database connection
- `D`: Driver for database connection (supported: `mysql`, `postgres`)
- `t`: Table for fuzzing
- `n`: Number of rows to fuzz
- `w`: Concurrent workers to work on fuzzing
- 's': Seed value for reproducibility of data

### Package usage

TODO: Write package 