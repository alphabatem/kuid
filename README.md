# KUID (Kompressed UUID)

KUID is a Go package that provides a compressed universally unique identifier implementation. It offers a more compact representation of UUIDs by using base62 encoding, making them URL-safe and more space-efficient while maintaining compatibility with standard UUIDs.

## Features

- Generate random KUIDs
- Convert between UUID and KUID formats
- Base62 encoding for URL-safe identifiers
- Compact 22-character representation (compared to 36 characters for UUID)
- Binary encoding/decoding support
- Thread-safe implementation

## Installation

```bash
go get github.com/alphabatem/kuid
```

## Usage

### Generate a New KUID

```go
kuid, err := kuid.NewKUID()
if err != nil {
    log.Fatal(err)
}
fmt.Println(kuid.String()) // Outputs a 22-character base62 string
```

### Convert UUID to KUID

```go
uuid := "550e8400-e29b-41d4-a716-446655440000"
kuid, err := kuid.FromUUID(uuid)
if err != nil {
    log.Fatal(err)
}
fmt.Println(kuid.String())
```

### Convert KUID to UUID

```go
kuid, _ := kuid.NewKUID()
uuid := kuid.ToUUID()
fmt.Println(uuid) // Standard UUID format
```

### Create KUID from String

```go
str := "4B1FkY2xHJRF6PTYR8Xj2Z"
kuid, err := kuid.FromString(str)
if err != nil {
    log.Fatal(err)
}
```

### Compare KUIDs

```go
kuid1, _ := kuid.NewKUID()
kuid2, _ := kuid.NewKUID()
if kuid1.Equal(kuid2) {
    fmt.Println("KUIDs are equal")
}
```

## Technical Details

KUID internally stores the identifier as two uint64 values (most significant bits and least significant bits). The string representation uses base62 encoding (0-9, A-Z, a-z) to achieve a compact 22-character format:

- Each uint64 is encoded into 11 characters
- Total length is 22 characters
- Maintains full UUID compatibility
- URL-safe characters only

## Performance

The base62 encoding/decoding operations are optimized for performance. The package uses minimal memory allocations and efficient algorithms for conversions.

## Limitations

- Base62 encoding of uint64 values must fit within 11 characters
- String representations are always 22 characters
- UUID conversions must be valid UUID format

## Thread Safety

All KUID operations are thread-safe. The package can be safely used in concurrent applications.

## Error Handling

The package provides specific error types for different failure cases:

- `ErrInvalidLength`: Input string has incorrect length
- `ErrInvalidChar`: Invalid character in input string
- `ErrInvalidUUID`: Malformed UUID string

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
