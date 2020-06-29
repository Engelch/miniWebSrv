# miniWebServer

This project implements a minimalistic web-server for educational and debugging purposes. It is written in go. 
It is under MIT license.

The server outputs HTML-header fields and the supplied HTML body.

## Example 

1. start the server in one terminal. `miniServer -p 8080`
2. call the server using `curl localhost:8080`

The output will look like:

```bash
$ curl localhost:8080
Hello, this is miniLogSrv at 2020-06-29T13:38:16Z
```

At the same time, the server will show:

```bash
$ ./build/debug/darwin-amd64/miniWebSrv-0.5.0 -p 8080
Info:2020-06-29T13:38:08Z:miniWebSrv::version:0.5.0:service starting at: 2020-06-29 15:38:08.688258 +0200 CEST m=+0.001760070
URL.Path: /
Host: localhost:8080
Header: map[Accept:[*/*] User-Agent:[curl/7.64.1]]
TrailingHeader: map[]
RemoteAddr: [::1]:62127
RequestURI: /
Content-Length: 0
Form: map[]
Body: vvvvvvvvvvvvvvvvvvvv

Body: ^^^^^^^^^^^^^^^^^^^^
```

## Versions

### 0.5

- change to default CLI parsing using github.com/urfave/cli
    - help mode 
    - port number now to be entered using -p
- Change to newest Makefile
- Using of bumpversion


    
    

