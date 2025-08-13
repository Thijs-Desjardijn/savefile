# Save File Retention Manager

This tool automatically manages the number of saved files in a given directory, ensuring that only the most recent saves are kept. When the maximum allowed number of save files is exceeded, the oldest one(s) are deleted based on timestamps embedded in their filenames.

## Features

- **Automatic Cleanup** — Deletes old save files when the limit is reached.
- **Timestamp-Based Ordering** — Finds the oldest save using a `YYYYMMDD_HHMMSS` timestamp in the filename.
- **Configurable Limits** — Set the maximum number of stored save files.
- **Safe Filtering** — Ignores non-regular files and files without a valid timestamp.

## Filename Requirements

The cleanup logic expects filenames to contain a timestamp in the following format:
PREFIX_YYYYMMDD_HHMMSS.ext

- The timestamp must start at character index **5** in the filename and be exactly 15 characters long (`YYYYMMDD_HHMMSS`).
- Files with invalid or missing timestamps are ignored.

## How It Works

1. Reads the list of files in the save directory.
2. Filters out:
   - Non-regular files
   - Filenames shorter than 20 characters
   - Files without a valid timestamp
3. Tracks the number of valid save files.
4. If the count exceeds `maxStoredFiles`, deletes the oldest file(s).

## Functions and Methods

### `New(path string, codec EncoderDecoder) (*Saver, error)`
Creates a new `Saver` instance that manages save files in a given directory and handles encoding/decoding using the provided codec.

**Parameters:**
- `path` — The filesystem path to the directory where save files will be stored.
- `codec` — An implementation of the `EncoderDecoder` interface.  
  Common options:
  - `GobCodec{}` for [gob](https://pkg.go.dev/encoding/gob) binary encoding.
  - `JSONCodec{}` for [json](https://pkg.go.dev/encoding/json) text encoding.

**Returns:**
- `*Saver` — A pointer to the new Saver instance.
- `error` — Non-`nil` if the Saver could not be created (e.g., directory issues).

**Example:**
```go
saver, err := New("./saves", GobCodec{})
if err != nil {
    log.Fatal(err)
}
```

### NewLimit(path string, codec EncoderDecoder, maxFiles int) (*Saver, error)
Like New, but also enforces a maximum number of stored save files. When the limit is exceeded, the oldest file(s) are automatically deleted.

**Parameters:**
- `path` — Directory where save files are stored.
- `codec` — Encoder/decoder to use for file serialization.
- `maxFiles` — Maximum number of save files to keep. Must be greater than zero.

**Returns:**
- `Saver` — New Saver instance with file limit enabled.
- `error` — Non-nil if initialization fails.

**Example:**
```go
saver, err := NewLimit("./saves", JSONCodec{}, 10)
if err != nil {
    return err
}
```

*** (s *Saver) Save(data any) error
Saves the given data to a new file in the Saver’s directory using the configured codec.

**Parameters:**
- `data` — Any Go value (struct, map, slice, etc.) to be saved.

**Behavior:**
Creates a new filename with a timestamp (YYYYMMDD_HHMMSS) to ensure chronological ordering.
Encodes data using the Saver’s codec (GobCodec or JSONCodec).
Writes the encoded data to disk.
If maxStoredFiles is set, old saves are deleted as needed.

**Returns:**
- `error` — Non-nil if encoding or writing fails.

**Example:**
```go
player := Player{Name: "Alice", Score: 42}
if err := saver.Save(player); err != nil {
    return err
}
```

###(s *Saver) Load(file string, target any) error
Loads a specific save file into the provided target variable.

**Parameters:**
- `file` — Filename of the save file to load (not the full path).
target — Pointer to the variable where the decoded data will be stored.

Behavior:

Reads the specified file from the Saver’s directory.

Decodes it using the Saver’s codec.

Stores the result in target.

Returns:

error — Non-nil if the file cannot be read, decoded, or found.

Example:
```go
var player Player
if err := saver.Load("save_20250813_123456.dat", &player); err != nil {
    log.Fatal(err)
}
```
(s *Saver) LoadLatest(target any) error
Loads the most recent save file into the provided target variable.

Parameters:

target — Pointer to the variable where the decoded data will be stored.

Behavior:

Scans the Saver’s directory for valid save files.

Finds the file with the most recent timestamp in its name.

Decodes its contents into target.

Returns:

error — Non-nil if no valid save files are found or decoding fails.

Example:
```go
var latestPlayer Player
if err := saver.LoadLatest(&latestPlayer); err != nil {
    log.Fatal(err)
}
```


## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).
