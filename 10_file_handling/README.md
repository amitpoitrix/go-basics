# Files in Go

## 1. Reading files

Read a small file in one call with `os.ReadFile`, and stream a large one line by line with `bufio.Scanner`.

```go
// Whole file at once: simplest, fine for small files
data, err := os.ReadFile("config.json")
if err != nil {
    return err
}
fmt.Println(string(data))

// Line by line: constant memory, right for large files
f, err := os.Open("app.log") // read-only
if err != nil {
    return err
}
defer f.Close()

scanner := bufio.NewScanner(f)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}
if err := scanner.Err(); err != nil {
    return err
}

// Fixed-size chunks: for binary data
buf := make([]byte, 4096)
for {
    n, err := f.Read(buf)
    if n > 0 {
        process(buf[:n]) // use only the bytes actually read
    }
    if err == io.EOF {
        break // normal end of file
    }
    if err != nil {
        return err
    }
}
```

- Check the error from `os.Open` before the `defer f.Close()`.
- `io.EOF` is not a failure. It is how a reader says there is nothing left.
- A `Scanner` line is limited to 64 KB by default. Call `scanner.Buffer` to raise it, or use `bufio.Reader`.

## 2. Writing files

Write a whole file with `os.WriteFile`, append with `os.OpenFile`, and wrap the file in a `bufio.Writer` when making many small writes.

```go
// Whole file: creates it, or replaces the contents if it exists
err := os.WriteFile("out.txt", []byte("hello\n"), 0644)

// Append: open with flags
f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    return err
}
defer f.Close()

if _, err := f.WriteString("new entry\n"); err != nil {
    return err
}

// Buffered: far fewer system calls for many small writes
out, err := os.Create("big.txt") // create the file, or empty it if it exists
if err != nil {
    return err
}
defer out.Close()

w := bufio.NewWriter(out)
for i := 0; i < 1000; i++ {
    fmt.Fprintln(w, "line", i)
}
if err := w.Flush(); err != nil { // without Flush, the last buffered data is lost
    return err
}
```

| Flag | Meaning |
| --- | --- |
| `os.O_RDONLY`, `os.O_WRONLY`, `os.O_RDWR` | Read only, write only, or both |
| `os.O_CREATE` | Create the file if it does not exist |
| `os.O_APPEND` | Every write goes to the end |
| `os.O_TRUNC` | Empty the file on open |
| `os.O_EXCL` | With `O_CREATE`, fail if the file already exists |

- The permission `0644` means the owner can read and write, and everyone else can read. Use `0600` for secrets.
- `os.Create` is shorthand for opening with create, truncate and read-write, at mode `0666`.
- For data that must survive a crash, call `f.Sync()` and check the error from `Close`.

## 3. Paths, directories and file info

Use `os` for operations on files and directories, and `path/filepath` to build and take apart paths in a portable way.

```go
// Does it exist?
if _, err := os.Stat("data.txt"); errors.Is(err, os.ErrNotExist) {
    fmt.Println("no such file")
}

info, _ := os.Stat("data.txt")
fmt.Println(info.Size(), info.IsDir(), info.ModTime())

// Create, rename, delete
os.MkdirAll("logs/2026/10", 0755) // creates parents too; no error if it exists
os.Rename("old.txt", "new.txt")
os.Remove("new.txt")  // one file or an empty directory
os.RemoveAll("logs") // a directory and everything in it

// List one directory
entries, _ := os.ReadDir(".")
for _, e := range entries {
    fmt.Println(e.Name(), e.IsDir())
}

// Walk a whole tree
filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    if !d.IsDir() && filepath.Ext(path) == ".go" {
        fmt.Println(path)
    }
    return nil
})

// Build and split paths
p := filepath.Join("data", "2026", "report.csv") // data/2026/report.csv on Linux
fmt.Println(filepath.Dir(p), filepath.Base(p), filepath.Ext(p)) // data/2026 report.csv .csv

// Temporary file
tmp, _ := os.CreateTemp("", "upload-*.txt")
defer os.Remove(tmp.Name())
```

- Build paths with `filepath.Join`, never by adding strings with a slash. It uses the right separator for the operating system.
- `path/filepath` is for the file system. The plain `path` package is for slash-separated paths such as URLs.
- In tests, `t.TempDir()` gives a directory that is removed automatically.

## 4. Copying, JSON files and embedded files

A file is both an `io.Reader` and an `io.Writer`, so it plugs straight into `io.Copy`, the JSON encoder and every other streaming API.

```go
// Copy a file without loading it into memory
func copyFile(dst, src string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    if _, err := io.Copy(out, in); err != nil { // streams in small chunks
        return err
    }
    return out.Sync() // make sure the bytes reached the disk
}

// Read JSON from a file
type Config struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}

func loadConfig(path string) (Config, error) {
    var cfg Config
    f, err := os.Open(path)
    if err != nil {
        return cfg, err
    }
    defer f.Close()
    err = json.NewDecoder(f).Decode(&cfg)
    return cfg, err
}

// Write JSON to a file
func saveConfig(path string, cfg Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

// Embed a file into the binary at compile time (Go 1.16)
import _ "embed"

//go:embed schema.sql
var schema string
```

- `io.Copy` works between any reader and writer: file to file, file to HTTP response, network to file.
- An embedded file ships inside the executable, so there is nothing extra to deploy. Use `embed.FS` for whole directories.
- `encoding/csv` reads and writes CSV through the same reader and writer interfaces.

## 5. Gotchas

Most file bugs come from a forgotten close, an unflushed buffer, or reading a huge file into memory.

```go
// 1. Deferring Close before checking the error
f, err := os.Open(path)
defer f.Close() // wrong place: check err first
if err != nil {
    return err
}

// 2. defer inside a loop keeps every file open until the function returns
for _, p := range paths {
    f, _ := os.Open(p)
    defer f.Close()
}

// 3. Ignoring the error from Close on a file you wrote
func save(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, f.Close()) // a failed Close can mean lost data
    }()
    _, err = f.Write(data)
    return err
}
```

- Every opened file uses a file descriptor. Leaking them eventually fails with "too many open files".
- `os.ReadFile` on a multi-gigabyte file can exhaust memory. Stream it instead.
- A forgotten `Flush` on a `bufio.Writer` silently truncates the output.
- Never join a user-supplied name into a path unchecked; `../` can escape the intended directory. Go 1.24 added `os.Root`, which confines access to one directory.
- A write to a file is not atomic. For safe replacement, write to a temporary file and `os.Rename` it over the original.
- `io/ioutil` is deprecated since Go 1.16. Its functions now live in `os` and `io`.

## 6. Rapid-fire questions

One-line answers for a quick revision.

| Question | Short answer |
| --- | --- |
| Simplest way to read a file? | `os.ReadFile(path)`. |
| How do you read a large file? | Line by line with `bufio.Scanner`, or in chunks with `Read`. |
| `os.Open` vs `os.Create` vs `os.OpenFile`? | Read-only; create or truncate; full control through flags. |
| How do you append to a file? | `os.OpenFile` with `O_APPEND`, `O_CREATE` and `O_WRONLY`. |
| Why `defer f.Close()`? | To release the file descriptor on every return path. |
| How do you check that a file exists? | `os.Stat` and `errors.Is(err, os.ErrNotExist)`. |
| What does `io.EOF` mean? | The normal end of input, not an error. |
| Why use `bufio`? | Buffering cuts the number of system calls. |
| What must follow writes to a `bufio.Writer`? | `Flush`. |
| How do you build a portable path? | `filepath.Join`. |
| How do you copy a file efficiently? | `io.Copy` from a reader to a writer. |
| What does `//go:embed` do? | Compiles a file's contents into the binary. |
| What does mode `0644` mean? | Owner reads and writes; group and others read. |
| What replaced `ioutil.ReadFile`? | `os.ReadFile`, since Go 1.16. |
