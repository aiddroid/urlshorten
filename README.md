### URL Shorten
A golang URL shorten service with redis.

### features
- URL shorten
- URL unshorten
- URL info query
- Web API(JSON) 

#### usage
- start service
```
go run main.go
```

create a short URL:
```
curl -i -XPOST "http://localhost:8089/api/shorten" --data '{"url":"https://www.qq.com","expire":600}'

=== response ===
{"code":0,"data":"http://localhost:8089/Q","message":"OK"}
```

visit a short URL:
```
curl -i "http://localhost:8089/Q"
```

query short URL info:
```
curl -i "http://0.0.0.0:8089/api/info?link=http://localhost:8089/P"
```


