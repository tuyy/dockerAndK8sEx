## docker, k8s

1. docker
```
# Dockerfile을 이용하여 이미지 생성
$ docker build -t myfirst:0.1 .

# 이미지 생성 확인
$ docker images

# 컨테이너 실행
$ docker run -d --name myfirstct -p 80:8080 myfirst:0.1

# 컨테이너 실행 확인
$ docker ps -a

# 컨테이너에서 /sbin/ifconfig 실행
$ docker exec myfirstct /sbin/ifconfig

# 컨테이너 내부의 /bin/bash 창으로 이동
$ docker exec -it myfirstct /bin/bash
```

2. k8s
```
TODO
```