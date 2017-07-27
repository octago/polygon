## How to test

Set `$GOPATH` to proper value.

Run server:
```
$ make server
```

Run client:
```
$ make client
```

Check app logs from the client:
```
# cat /var/log/app.log
...
2017/07/27 17:29:17 Listening on [::]:46199
2017/07/27 17:29:17 Listening on [::]:40323
2017/07/27 17:29:17 listen tcp :0: socket: too many open files
```

Up ulimit:
```
ulimit -n 256
```

Try to run app again:
```
# app
...
2017/07/27 17:31:20 Listening on [::]:39817
2017/07/27 17:31:20 Listening on [::]:37351
2017/07/27 17:31:20 Finished OK
```
