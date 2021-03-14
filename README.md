## docker, k8s

1) docker basic
```
$ make build

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

2) k8s Deployment
```
$ make build
$ ncc build -t reg.xxxxcorp.com/solver/myfirst:0.1 .

$ kubectl create -f deployment.yaml

$ kubectl get pods

$ kubectl exec -it myfirst-deployment-65d7f884d6-tffss /bin/bash

$ kubectl delete -f Deployment.yaml

# deployment를 사용하면 리비전 관리, 다양한 pod의 배포 정책을 지정할 수 있음

# 롤링업데이트로 이미지 변경
#    - Deployment.yaml 에서 image 버전 수정 후 kubectl apply -f deployment 호출하면 동일한 동작함
$ kubectl set image deployment myfirst-deployment myfirst:reg.xxxxcorp.com/solver/myfirst:0.2 --record
deployment.extensions/myfirst-deployment image updated

$ kubectl get pods
NAME                                  READY   STATUS              RESTARTS   AGE
myfirst-deployment-65d7f884d6-89sv6   1/1     Running             0          8m46s
myfirst-deployment-65d7f884d6-khfkt   1/1     Running             0          8m46s
myfirst-deployment-7464cbd84f-97d2c   0/1     ContainerCreating   0          3s

# 이전 리비전 정보 확인
$ kubectl rollout history deployment myfirst-deployment
deployment.extensions/myfirst-deployment
REVISION  CHANGE-CAUSE
1         kubectl apply --filename=Deployment.yaml --record=true
2         kubectl set image deployment myfirst-deployment myfirst=reg.xxxxcorp.com/solver/myfirst:0.2 --record=true

# revision 1번으로 롤백
$ kubectl rollout undo deployment myfirst-deployment --to-revision=1
deployment.extensions/myfirst-deployment rolled back

# 롤백 후 history 확인 결과 현재 리비전이 3으로 지정된 것을 알 수 있음
$ kubectl rollout history deployment myfirst-deployment
deployment.extensions/myfirst-deployment
REVISION  CHANGE-CAUSE
2         kubectl set image deployment myfirst-deployment myfirst=reg.xxxxcorp.com/solver/myfirst:0.2 --record=true
3         kubectl apply --filename=Deployment.yaml --record=true

# replica factor 변경
$ kubectl scale --replicas=1 deployment  myfirst-deployment
deployment.extensions/myfirst-deployment scaled

# 서비스 디스커버리 전략
* LoadBalancer VIP
    * [service-name].[namespace].svc.[cluster-name].io.xxxxcorp.com

```
$ curl myfirst-svc-lb.devmail.svc.xd1.io.navercorp.com:3000/ping
{"message":"pong22"}
```

* 내부 ClusterIP
    * 같은 namepsace에서 [service-name]
    * 다른 namespace에서 [service-name].[namespace]

# 배포 전략
* 롤링 업데이트: 하나씩 순차적으로 배포(default)
```

3) k8s Service
```
$ kubectl apply -f Deployment2.yaml

# pod에 할당된 IP 확인
$ kubectl get pods -o wide
NAME                                   READY   STATUS    RESTARTS   AGE     IP               NODE           NOMINATED NODE   READINESS GATES
hostname-deployment-7b46bfbbb8-6cfwg   1/1     Running   0          2m32s   10.171.143.68    ad1x0131.xxx   <none>           <none>
hostname-deployment-7b46bfbbb8-h7xxf   1/1     Running   0          2m32s   10.171.188.178   ad1x0426.xxx   <none>           <none>
hostname-deployment-7b46bfbbb8-kpq4n   1/1     Running   0          2m32s   10.171.219.66    ad1x0625.xxx   <none>           <none>

* ClusterIP 타입: k8s 내부에서만 포트들에 접근할때 사용(외부 노출x)
* NodePort 타입: 포드에 접근할 수 있는 포트를 외부, 내부 모두 가능하게함
    * -> 단, 각 노드의 포트번호를 알아야함... '$ kubectl get nodes'
* LoadBalancer 타입: 외부, 내부 접근 가능 + 로드벨런서 역할 지원함
    * 외부IP나 DNS가 할당됨!

==========

# ClusterIP로 서비스 생성 후 클러스터 내부에서 호출해보기
$ kubectl apply -f hostname-svc-clusterip.yaml
service/hostname-svc-clusterip created

$ kubectl get svc
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
hostname-svc-clusterip   ClusterIP   172.24.216.33   <none>        8080/TCP   24s

$ kubectl get svc,pods
NAME                             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/hostname-svc-clusterip   ClusterIP   172.24.216.33   <none>        8080/TCP   5m37s

NAME                                       READY   STATUS    RESTARTS   AGE
pod/hostname-deployment-7b46bfbbb8-6cfwg   1/1     Running   0          16m
pod/hostname-deployment-7b46bfbbb8-h7xxf   1/1     Running   0          16m
pod/hostname-deployment-7b46bfbbb8-kpq4n   1/1     Running   0          16m
pod/myfirst-deployment-65d7f884d6-ch6lh    1/1     Running   0          2m48s

$ kubectl exec -it myfirst-deployment-65d7f884d6-ch6lh bash

bash-5.0$ curl 172.24.216.33:8080
	<p>Hello,  hostname-deployment-7b46bfbbb8-h7xxf</p>	</blockquote>

bash-5.0$ curl 172.24.216.33:8080
	<p>Hello,  hostname-deployment-7b46bfbbb8-6cfwg</p>	</blockquote>


# 서비스의 라벨 셀렉터와 pod의 라벨이 매칭되었는지 엔드포인트 객체를 통해 확인
// $ kubectl get endpoints
$ kubectl get ep
NAME                     ENDPOINTS                                             AGE
hostname-svc-clusterip   10.171.143.68:80,10.171.188.178:80,10.171.219.66:80   9m46s

==========

# LoadBalancer 타입으로 서비스 생성 후 클러스터 외부에서 호출해보기
$ kubectl apply -f hostname-svc-loadbalancer.yaml
service/hostname-svc-loadbalancer created

$ kubectl get svc
NAME                        TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
hostname-svc-loadbalancer   LoadBalancer   172.24.164.3   10.108.227.101   80:39841/TCP   3s

$ kubectl get ep
NAME                        ENDPOINTS                                             AGE
hostname-svc-loadbalancer   10.171.143.68:80,10.171.188.178:80,10.171.219.66:80   9s

$ kubectl get pods -o wide
NAME                                   READY   STATUS    RESTARTS   AGE   IP               NODE           NOMINATED NODE   READINESS GATES
hostname-deployment-7b46bfbbb8-6cfwg   1/1     Running   0          50m   10.171.143.68    ad1k0131.xxx   <none>           <none>
hostname-deployment-7b46bfbbb8-h7xxf   1/1     Running   0          50m   10.171.188.178   ad1k0426.xxx   <none>           <none>
hostname-deployment-7b46bfbbb8-kpq4n   1/1     Running   0          50m   10.171.219.66    ad1k0625.xxx   <none>           <none>

$ curl 10.108.227.101
        <p>Hello,  hostname-deployment-7b46bfbbb8-kpq4n</p>     </blockquote>

# 주의할 점은 LoadBalancer type 서비스의 ClusterIP를 사용하기 위해서는 서비스의 targetPort 그리고 pod의 containerPort가 일치 해야함. ExternalIP를 사용하기 위해서는 서비스의 port와 targetPort 그리고 pod의 containerPort가 모두 일치 해야 함. 즉, 여러개의 포트를 사용하는 서비스의 경우 ClusterIP만 사용하는 Port 인 경우와 ExternalIP가 포함된 Port는 일치시켜야 조건에 차이가 있으니 주의해야함.

==========

```
