# gonetcat

goalng nc

## how to install
```
go get -u github.com/umaumax/gonetcat
```

## how to use
### original netcat
```
$ nc localhost 7083 << !
GET / HTTP/1.0

!
```

### my command
```
$ gonetcat localhost:8080 << !
GET / HTTP/1.0

!
```

## NOTE
* 受信したContent-Lengthをパースして指定長さとなった場合に閉じる。
* 本来は送るリクエストのwriteを閉じて判断するべきなのだろうが、サーバの実装によって異なるため、以上のような実装となっている。
