# entry-server

## 开发
```sh
# 启动mysql、redis等镜像
docker-compose up
```

## 环境变量

SERVER_ENV
- prod
- dev

## 启动容器
```sh
docker run -dp 8081:8080 --env SERVER_ENV=dev --name entry-server mirrors.xx.com/jht/entry-server:${tag}
```

## 测试
```sh
# 主版本
curl localhost:8080 -H "host:jht.xx.com"

# 人员灰度
curl localhost:8080 -H "host:jht.xx.com" -H "staffname:jht"

# 百分比灰度
curl localhost:8080 -H "host:jht.xx.com" -H "staffid:666"

```