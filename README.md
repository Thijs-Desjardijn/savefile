# Save File Retention Manager

This tool automatically manages the number of saved files in a given directory, ensuring that only the most recent saves are kept. When the maximum allowed number of save files is exceeded, the oldest one(s) are deleted based on timestamps embedded in their filenames.

## Features

- **Automatic Cleanup** — Deletes old save files when the limit is reached.
- **Timestamp-Based Ordering** — Finds the oldest save using a `YYYYMMDD_HHMMSS` timestamp in the filename.
- **Configurable Limits** — Set the maximum number of stored save files.
- **Safe Filtering** — Ignores non-regular files and files without a valid timestamp.

## Filename Format

The cleanup logic expects filenames to contain a timestamp in the following format:
PREFIX_YYYYMMDD_HHMMSS.ext

- The timestamp must start at character index **5** in the filename and be exactly 15 characters long (`YYYYMMDD_HHMMSS`) **Note: this is automatically handled** .
- Files with invalid or missing timestamps are ignored.

## How It Works

1. Reads the list of files in the save directory.
2. Filters out:
   - Non-regular files
   - Filenames shorter than 20 characters
   - Files without a valid timestamp
3. Tracks the number of valid save files.
4. If the count exceeds `maxFiles`, deletes the oldest file(s).

## Functions and Methods

### `New(path string, codec EncoderDecoder) (*Saver, error)`
Creates a new `Saver` instance that manages save files in a given directory and handles encoding/decoding using the provided codec.                               
**Note that you are allowed to create multiple savers in 1 directory but this will cause unexpected behavior**.

**Parameters:**
- `path` — The filesystem path to the directory where save files will be stored.
- `codec` — An implementation of the `EncoderDecoder` interface.  
  Common options:
  - `GobCodec{}` for [gob](https://pkg.go.dev/encoding/gob) binary encoding.
  - `JSONCodec{}` for [json](https://pkg.go.dev/encoding/json) text encoding.

**Returns:**
- `Saver` — A new Saver instance.
- `error` — Non-nil if the Saver could not be created.

**Example:**
```go
saver, err := New("./saves", GobCodec{})
if err != nil {
    log.Fatal(err)
}
```

---

### `NewLimit(path string, codec EncoderDecoder, maxFiles int) (*Saver, error)`
Like `New`, but also enforces a maximum number of stored save files. When the limit is exceeded, the oldest file(s) are automatically deleted.                                          
**Note that you are allowed to create multiple savers in 1 directory but this will cause unexpected behavior**.

**Parameters:**
- `path` — Directory where save files are stored.
- `codec` — Encoder/decoder to use for file serialization.                                                                                  
  Common options:
  - `GobCodec{}` for [gob](https://pkg.go.dev/encoding/gob) binary encoding.
  - `JSONCodec{}` for [json](https://pkg.go.dev/encoding/json) text encoding.
- `maxFiles` — Maximum number of save files to keep. Must be greater than zero.

**Returns:**
- `Saver` — New Saver instance with file limit enabled.
- `error` — Non-nil if initialization fails.

**Example:**
```go
saver, err := NewLimit("./saves", JSONCodec{}, 10)
if err != nil {
    log.Fatal(err)
}
```

---

### `(s *Saver) Save(data any) error`
Saves the given data to a new file in the Saver’s directory using the configured codec.
After saving the data it will also run DeleteOld() if a maxFiles had been set using NewLimit().

**Parameters:**
- `data` — Any Go value to be saved.

**Returns:**
- `error` — Non-nil if encoding or writing fails.

**Example:**
```go
user := User{Name: "Alice", Score: 42}
if err := saver.Save(user); err != nil {
    log.Fatal(err)
}
```

---

### `(s *Saver) LoadLatest(target any) error`
Loads the most recent save file into the provided target variable.

**Parameters:**
- `target` — Pointer to the variable where the decoded data will be stored.

**Returns:**
- `error` — Non-nil if no valid save files are found or decoding fails.

**Example:**
```go
var user User
if err := saver.LoadLatest(&user); err != nil {
    log.Fatal(err)
}
```

---

### `(s *Saver) Load(file string, target any) error`
Loads a specific save file into the provided target variable.

**Parameters:**
- `file` — Filename of the save file to load (not the full path).
- `target` — Pointer to the variable where the decoded data will be stored.

**Returns:**
- `error` — Non-nil if the file cannot be read, decoded, or found.

**Example:**
```go
var user User
if err := saver.Load("save_20250813_123456.dat", &user); err != nil {
    log.Fatal(err)
}
```

---

### `(s *Saver) DeleteOld() error`
Deletes the oldest save files in the Saver’s directory based on its timestamp until the number of save files is <= s.maxFiles. If you didn't specify a s.maxFiles using NewLimit it will Delete the oldest file. Run the DeleteFile method if you do want to delete a specific file.

**Returns:**
- `error` — Non-nil if no valid save files are found or if deletion fails.

**Example:**
```go
if err := saver.DeleteOld(); err != nil {
    log.Fatal(err)
}
```

---

### `(s *Saver) DeleteFile(fileName string) error`
Deletes a specific save file from the Saver’s directory.

**Parameters:**
- `fileName` — The name of the file to delete (not the full path).

**Returns:**
- `error` — Non-nil if the file cannot be found or deletion fails.

**Example:**
```go
if err := saver.Delete("examplefile"); err != nil {
    log.Fatal(err)
}
```

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).
