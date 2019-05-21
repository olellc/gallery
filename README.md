# gallery

gallery is a hello-world style foto hosting service.
It stores uploaded fotos in the SQLite database.
It also generates previews for jpeg, png and gif file formats.

# Compiling

This project uses [go-sqlite3](https://github.com/mattn/go-sqlite3) which in turn uses `CGO`. So to compile this package you are required to set the environment variable `CGO_ENABLED=1` and have a `gcc` compile present within your path.

# Command line options

```
-h, --help     show help message

--addr=        TCP network address to listen on (default: :8080)
--dbpath=      path to the SQLite database file (default: temp file)
--keep-schema  keep existing schema
```

# REST API

- `POST /save` - post new foto. Returns an id of the newly created foto.
- `GET /fotos` -  get list of foto ids as json array
- `GET /raw_foto/{id}` - get an original foto
- `GET /preview/{id}` - get a preview
- `DELETE /foto/{id}` - delete a foto

# Usage example

```
$ http POST :8080/save < 1.jpg
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json; charset=utf-8
Date: Tue, 21 May 2019 07:19:20 GMT

1


$ http POST :8080/save < 2.png
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json; charset=utf-8
Date: Tue, 21 May 2019 07:19:40 GMT

2


$ http POST :8080/save < 3.gif
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json; charset=utf-8
Date: Tue, 21 May 2019 07:19:51 GMT

3


$ http :8080/fotos
HTTP/1.1 200 OK
Content-Length: 8
Content-Type: application/json; charset=utf-8
Date: Tue, 21 May 2019 07:19:59 GMT

[
    1, 
    2, 
    3
]


$ http :8080/raw_foto/1
HTTP/1.1 200 OK
Content-Length: 7403
Content-Type: image/jpeg
Date: Tue, 21 May 2019 07:20:17 GMT

+-----------------------------------------+
| NOTE: binary data not shown in terminal |
+-----------------------------------------+


$ http :8080/preview/1
HTTP/1.1 200 OK
Content-Length: 9418
Content-Type: image/png
Date: Tue, 21 May 2019 07:20:27 GMT

+-----------------------------------------+
| NOTE: binary data not shown in terminal |
+-----------------------------------------+


$ http DELETE :8080/foto/2
HTTP/1.1 200 OK
Content-Length: 0
Date: Tue, 21 May 2019 07:20:39 GMT


$ http :8080/fotos
HTTP/1.1 200 OK
Content-Length: 6
Content-Type: application/json; charset=utf-8
Date: Tue, 21 May 2019 07:20:51 GMT

[
    1, 
    3
]
```
