## SQLfuzz

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