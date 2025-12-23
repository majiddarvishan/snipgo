# BiMap - Thread-Safe Bidirectional Map

A high-performance, thread-safe bidirectional map implementation in Go that allows efficient lookups in both directions (key→value and value→key).

## Features

- ✅ **Thread-Safe**: Built-in mutex protection for concurrent access
- ✅ **Bidirectional Lookups**: O(1) lookups by key or by value
- ✅ **Deterministic Iteration**: Ordered key storage for predictable pagination
- ✅ **Automatic Cleanup**: Removes old mappings when keys/values are reused
- ✅ **Zero Dependencies**: Pure Go implementation

## Installation

```bash
go get github.com/majiddarvishan/snipgo
```

## Usage

### Basic Operations

```go
package main

import (
    "fmt"
    "github.com/majiddarvishan/snipgo"
)

func main() {
    // Create a new BiMap
    bm := snipgo.NewBiMap()

    // Add key-value pairs
    bm.Set("user1", "email1@example.com")
    bm.Set("user2", "email2@example.com")
    bm.Set("user3", "email3@example.com")

    // Lookup by key
    if value, exists := bm.Get("user1"); exists {
        fmt.Println("Found:", value) // Output: email1@example.com
    }

    // Reverse lookup by value
    if key, exists := bm.GetByValue("email2@example.com"); exists {
        fmt.Println("Key:", key) // Output: user2
    }

    // Get size
    fmt.Println("Size:", bm.Len()) // Output: 3

    // Delete a mapping
    bm.Delete("user1")
}
```

### Pagination

BiMap maintains an ordered list of keys for deterministic pagination:

```go
bm := snipgo.NewBiMap()

// Add multiple entries
for i := 0; i < 100; i++ {
    bm.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
}

// Get first 10 entries
page1 := bm.GetValuesWithRange(0, 10)

// Get next 10 entries
page2 := bm.GetValuesWithRange(10, 10)

// Get entries 20-29
page3 := bm.GetValuesWithRange(20, 10)
```

### Concurrent Access

BiMap is thread-safe and can be used safely from multiple goroutines:

```go
bm := snipgo.NewBiMap()

// Multiple goroutines can safely access the BiMap
go func() {
    for i := 0; i < 1000; i++ {
        bm.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
    }
}()

go func() {
    for i := 0; i < 1000; i++ {
        if value, exists := bm.Get(fmt.Sprintf("key%d", i)); exists {
            fmt.Println(value)
        }
    }
}()

go func() {
    for i := 0; i < 1000; i++ {
        bm.Delete(fmt.Sprintf("key%d", i))
    }
}()
```

## API Reference

### Constructor

#### `NewBiMap() *BiMap`
Creates and returns a new empty BiMap instance.

```go
bm := snipgo.NewBiMap()
```

### Methods

#### `Set(key string, value string)`
Adds or updates a key-value pair. If the key already exists, the old value mapping is removed. If the value already exists with a different key, that old key mapping is removed.

```go
bm.Set("user1", "email1@example.com")
bm.Set("user1", "newemail@example.com") // Updates existing key
```

**Thread-safe**: Yes

---

#### `Get(key string) (string, bool)`
Returns the value associated with the given key and a boolean indicating whether the key exists.

```go
value, exists := bm.Get("user1")
if exists {
    fmt.Println("Value:", value)
}
```

**Thread-safe**: Yes
**Time Complexity**: O(1)

---

#### `GetByValue(value string) (string, bool)`
Returns the key associated with the given value and a boolean indicating whether the value exists.

```go
key, exists := bm.GetByValue("email1@example.com")
if exists {
    fmt.Println("Key:", key)
}
```

**Thread-safe**: Yes
**Time Complexity**: O(1)

---

#### `GetValuesWithRange(start, limit int) map[string]string`
Returns a subset of key-value pairs starting at the specified index with the given limit. Keys are returned in sorted order for deterministic pagination.

```go
// Get entries 10-19
page := bm.GetValuesWithRange(10, 10)
```

**Thread-safe**: Yes
**Time Complexity**: O(limit)

---

#### `Delete(key string)`
Removes the key-value pair associated with the given key. No-op if the key doesn't exist.

```go
bm.Delete("user1")
```

**Thread-safe**: Yes
**Time Complexity**: O(n) where n is the number of keys (due to slice removal)

---

#### `Len() int`
Returns the number of key-value pairs in the BiMap.

```go
count := bm.Len()
fmt.Println("Total entries:", count)
```

**Thread-safe**: Yes
**Time Complexity**: O(1)

## Implementation Details

### Internal Structure

```go
type BiMap struct {
    mu         sync.RWMutex          // Protects all fields
    keyToValue map[string]string     // Key → Value lookup
    valueToKey map[string]string     // Value → Key lookup
    keys       []string              // Ordered keys (sorted)
}
```

### Key Characteristics

1. **Uniqueness**: Both keys and values must be unique across the BiMap
2. **Ordering**: Keys are maintained in sorted order for deterministic iteration
3. **Updates**: Setting a key that exists updates its value; setting a value that exists replaces its key
4. **Synchronization**: Uses RWMutex for efficient concurrent reads

### Performance

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Set       | O(log n)*      | O(1)             |
| Get       | O(1)           | O(1)             |
| GetByValue| O(1)           | O(1)             |
| Delete    | O(n)           | O(1)             |
| Len       | O(1)           | O(1)             |
| GetValuesWithRange | O(limit) | O(limit)   |

*O(log n) due to maintaining sorted key order

## Use Cases

BiMap is ideal for scenarios requiring bidirectional lookups:

- **Phone Number Mapping**: Original MSISDN ↔ Replacement MSISDN
- **User Sessions**: Session ID ↔ User ID
- **Translation Tables**: External ID ↔ Internal ID
- **Cache Systems**: Cache Key ↔ Resource Identifier
- **Routing Tables**: Source Address ↔ Destination Address

## Thread Safety Guarantees

All public methods are thread-safe and can be called concurrently from multiple goroutines. The implementation uses `sync.RWMutex` to allow:

- Multiple concurrent readers
- Exclusive access for writers
- Automatic deadlock prevention

## Limitations

1. **String Keys Only**: Currently supports only string keys and values
2. **Memory Overhead**: Stores each mapping twice (key→value and value→key)
3. **Delete Performance**: O(n) deletion due to ordered key maintenance
4. **No Iterators**: No direct iteration support (use `GetValuesWithRange` for pagination)

## Examples

### Example 1: Phone Number Replacement System

```go
// Map old phone numbers to new ones
phoneMap := snipgo.NewBiMap()
phoneMap.Set("+1234567890", "+9876543210")
phoneMap.Set("+1111111111", "+2222222222")

// Forward lookup: old → new
if newNumber, exists := phoneMap.Get("+1234567890"); exists {
    fmt.Println("Replace with:", newNumber)
}

// Reverse lookup: new → old
if oldNumber, exists := phoneMap.GetByValue("+9876543210"); exists {
    fmt.Println("Original was:", oldNumber)
}
```

### Example 2: User Session Management

```go
sessions := snipgo.NewBiMap()

// Map session IDs to user IDs
sessions.Set("session-abc123", "user-456")
sessions.Set("session-def789", "user-789")

// Find user by session
if userID, exists := sessions.Get("session-abc123"); exists {
    fmt.Println("User:", userID)
}

// Find session by user (single session per user)
if sessionID, exists := sessions.GetByValue("user-456"); exists {
    fmt.Println("Session:", sessionID)
}

// End session
sessions.Delete("session-abc123")
```

### Example 3: Paginated API Response

```go
data := snipgo.NewBiMap()

// Load data
for i := 0; i < 1000; i++ {
    data.Set(fmt.Sprintf("key%04d", i), fmt.Sprintf("value%04d", i))
}

// Implement pagination
func getPage(pageNum, pageSize int) map[string]string {
    start := pageNum * pageSize
    return data.GetValuesWithRange(start, pageSize)
}

// Get page 1 (items 0-9)
page1 := getPage(0, 10)

// Get page 2 (items 10-19)
page2 := getPage(1, 10)
```

## Testing

Run the test suite:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Guava's BiMap implementation
- Built for high-performance telecommunications systems
- Optimized for concurrent access patterns

## Support

For issues, questions, or contributions, please open an issue on GitHub.