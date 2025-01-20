# HTTP File Server

## Overview

The HTTP File Server is a Go-based file server designed to handle file storage and retrieval over HTTP. It uses `gin-gonic/gin` for routing and middleware, enabling efficient API development. The server includes a robust storage layer with concurrency handling, custom path transformations, and structured logging.

---

## Features

- **File Operations**: Supports reading, writing, updating, and deleting files via HTTP endpoints.
- **Concurrency Control**: Mutex-based locking for file-specific operations ensures safe concurrent access.
- **Customizable Storage Paths**: Flexible path transformation functions allow for content-addressable storage or other custom schemes.
- **Structured Logging**: Built-in logging mechanism tracks key operations and errors.
- **API Versioning**: Endpoints are organized under `/v1/fileserver`.

---

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/your-repository/http-file-server.git
   cd http-file-server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

---

## API Endpoints

### Base URL
```
http://localhost:4000/v1/fileserver
```

### Endpoints

1. **Read File**
   - `GET /:key`
   - Returns the file associated with the given key.
   - Example:
     ```bash
     curl -X GET http://localhost:4000/v1/fileserver/example-key
     ```

2. **Write File**
   - `POST /:key`
   - Uploads a file using a multipart form with the field name `file`.
   - Example:
     ```bash
     curl -X POST http://localhost:4000/v1/fileserver/example-key -F "file=@path/to/your/file.txt"
     ```

3. **Update File**
   - `PATCH /:key`
   - [Not Implemented]

4. **Delete File**
   - `DELETE /:key`
   - Deletes the file associated with the given key.
   - Example:
     ```bash
     curl -X DELETE http://localhost:4000/v1/fileserver/example-key
     ```

---

## Storage Path Transformation

The storage system supports customizable path transformation functions to determine how files are organized on disk.

### DefaultPathTransformFunc

- Splits the key into user details and file name using `^` as a delimiter.
- Returns the user details as the directory path and the second part as the file name.

Example:
```go
var DefaultPathTransformFunc = func(storageRoot string, key string) (string, string) {
    keyContains := strings.Split(key, "^")
    userDetails := keyContains[0]
    fileName := keyContains[1]
    return storageRoot + "/" + userDetails, fileName
}
```

### CASPathTransformFunc

- Implements a content-addressable storage (CAS) path transformation.
- Splits the key into user details and file name using `^` as a delimiter.
- Computes:
  - SHA-1 hash for user details to construct a directory path.
  - MD5 hash for the file name to create a unique identifier with its original extension.

Example:
```go
var CASPathTransformFunc = func(storageRoot string, key string) (string, string) {
    keyContains := strings.Split(key, "^")
    userDetails := keyContains[0]
    fileExt := filepath.Ext(keyContains[1])
    fileName := filepath.Base(keyContains[1])

    hash := sha1.Sum([]byte(userDetails))
    hashStr := hex.EncodeToString(hash[:])
    blockSize := 8
    sliceLen := len(hashStr) / blockSize
    paths := make([]string, sliceLen)

    for i := 0; i < sliceLen; i++ {
        from, to := i*blockSize, (i*blockSize)+blockSize
        paths[i] = hashStr[from:to]
    }

    fileNameBytes := md5.Sum([]byte(fileName))
    fileName = hex.EncodeToString(fileNameBytes[:])

    return storageRoot + "/" + strings.Join(paths, "/"), fileName + fileExt
}
```

---

## Configuration

The server uses the following configuration options:
- **`StorageRoot`**: Root directory for file storage.
- **`PathTransformFunc`**: Function for customizing file storage paths (default: `CASPathTransformFunc`).
- **`Logger`**: Instance of `goLogger` for logging.

---

## Project Structure

- `main.go`: Entry point of the application, sets up routing and dependencies.
- `controller/controller.go`: Contains the `StorageController` for handling API routes and HTTP methods.
- `storage/storage.go`: Implements the storage backend, including file operations, concurrency management, and path transformations.

---

## Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin): Web framework for Go.
- [goLogger](https://github.com/sumit-behera-in/goLogger): Structured logging for Go applications.

---

## Contributing

Contributions are welcome! Feel free to submit a pull request or open an issue.

---

Let me know if you need further adjustments or additional details!