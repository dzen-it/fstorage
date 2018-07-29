# FStorage

  

## Overview

FStorage is embeded file storage with access via http.

## Installing

```go get -u github.com/dzen-it/fstorage```

  

## Using

  

#### Base usage:

```go
limit:=1<<30 // set limit of memory 1GB for data dir.
s, err := storage.NewFileStorage("./data",limit) // more see into fstorage/storage
if err != nil {
	panic(err)
}
server := fstorage.NewServer(s, nil)
server.Start(":8080")
```
#### Configuring:
```go
server.RPS = 9.42 // requests per second per one IP
server.CS = 3 // number of concurent sessions per one IP
server.MaxFilesize = 10<<20*100 // maximum size of the uploaded file
```

#### Set pre and  post processing via http middlewares:
```go
type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w customResponseWriter) WriteHeader(statusCode int) {
	w.status  = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

middleware:=func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp  := customResponseWriter{
			status: 0,
			ResponseWriter: w,
		}

		_, file  := path.Split(r.URL.Path)

		fmt.Println("Pre processing of file", file)
		next.ServeHTTP(resp, r)
		fmt.Println("Post processing: done with code", resp.status)
	})
}

server.AddMiddleware(middleware)
```
**Note:**  To access the file before writing to the storage, you must use the `io.ReadeCloser` from `r.Body`. To avoid trouble use `io.TeeReader()`.

#### Add custom http handler 
```go
// add endpoint GET /greet?name=
r := chi.NewRouter()
r.Get("/greeter", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", r.FormValue("name"))
})
server.MountHandler("/", r)
```

## REST API   
At the moment there is a fixed minimum request interval of 1 minute.

| Endpoint | Method | Body | Succes response | Description | 
|---|---|---|---|---|
| /files/{name_of_file} | PUT | Bytes of file  | Code: 201<br> Body: c1133418aba4ed90f78881498fc1a1ce68870f569489a661d89d89eb3416a7f4 | Upload file |
| /files/{hash_of_file} | GET || Code: 200 | Download file |
| /files/{hash_of_file} | DELETE | | Code: 204 | Delete file | 

### Headers hash control
If a header is exists when the file is uploading, the hash will be calculated, if it does not match the hash from the header, then will return the error.

| Header | Hash type |
|---|---|
| X-FStorage-Hash-Control-MD5 | MD5|
| X-FStorage-Hash-Control-SHA1 | SHA1|
| X-FStorage-Hash-Control-SHA256 | SHA2 256 Bit|
| X-FStorage-Hash-Control-SHA512 | SHA2 512 Bit|
| X-FStorage-Hash-Control-Keccak256 | SHA3 256 Bit|
| X-FStorage-Hash-Control-Keccak512 | SHA3 512 Bit|

#### Example:
```bash
curl --header "X-FStorage-Hash-Control-sha256: 832e4ba158a563a0b2eae3c010033229984f96cf1a4ad8d5c5c3226ff2d4daf6" \
-X PUT http://localhost:8080/files/something.txt
```

## TODO:

 - [ ] Redis: for storing metadata
