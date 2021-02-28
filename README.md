## SQLfuzz

[![Go Report Card](https://goreportcard.com/badge/github.com/PumpkinSeed/sqlfuzz)](https://goreportcard.com/report/github.com/PumpkinSeed/sqlfuzz) [![GoDoc](https://godoc.org/github.com/PumpkinSeed/sqlfuzz?status.svg)](https://godoc.org/github.com/PumpkinSeed/sqlfuzz) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org) ![sqlfuzz test workflow](https://github.com/PumpkinSeed/sqlfuzz/actions/workflows/test.yml/badge.svg)

Load random data into SQL tables for testing purposes. The tool can get the layout of the SQL table and fill it up with random data. 

#### Usage

```
go install github.com/PumpkinSeed/sqlfuzz
```

```
sqlfuzz -u username -p password -d database -h 127.0.0.1 -t table -n 100000 -w 100
```

#### Flags

- `u`: User for database connection
- `p`: Password for database connection
- `d`: Database name for database connection
- `h`: Host for database connection
- `P`: Port for database connection
- `D`: Driver for database connection (currently only `mysql`)
- `t`: Table for fuzzing
- `n`: Number of rows to fuzz
- `w`: Concurrent workers to work on fuzzing