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

==========

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

==========

# 서비스 디스커버리 전략
* LoadBalancer VIP
    * [service-name].[namespace].svc.[cluster-name].io.xxxxcorp.com

$ curl myfirst-svc-lb.devmail.svc.xd1.io.navercorp.com:3000/ping
{"message":"pong22"}

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
```

4) k8s ConfigMap, Secret, kustomize
```
// configmap -> cm 도 가능
$ kubectl create configmap myfirst-configmap --from-literal name=tuyy --from-literal age=31
configmap/myfirst-configmap created

$ kubectl get cm myfirst-configmap
NAME                DATA   AGE
myfirst-configmap   2      10s

$ kubectl describe cm myfirst-configmap
Name:         myfirst-configmap
Namespace:    devmail
Labels:       <none>
Annotations:  <none>

Data
====
age:
----
31
name:
----
tuyy
Events:  <none>

$ kubectl get configmap myfirst-configmap -o yaml
apiVersion: v1
data:
  age: "31"
  name: tuyy
kind: ConfigMap
metadata:
  creationTimestamp: "2021-03-14T02:11:36Z"
  name: myfirst-configmap
  namespace: devmail
  resourceVersion: "4482254835"
  selfLink: /api/v1/namespaces/devmail/configmaps/myfirst-configmap
  uid: eb4de920-4953-4b2a-9b48-a693693ea52f

==========

# 컨피그맵 사용 방법 2가지
* 컨피그맵의 값을 컨테이너의 환경변수로 사용하기

$ cat myfirst-pod-for-cm.yaml
apiVersion: v1
kind: Pod
metadata:
    name: myfirst-pod-example
spec:
    containers:
    - name: myfirst-pod
      image: busybox
      args: ['tail', '-f', '/dev/null']
      envFrom:
      - configMapRef:
          name: myfirst-configmap

$ kubectl exec myfirst-pod-example env |grep 'name\|age'
age=31
name=tuyy

-> 참고로 envFrom 대신 env로 변경하면 특정 cm에 특정 key를 환경변수로 지정할 수 있다.

==========

* 컨피그맵의 값을 포드 내부의 파일로 마운트해서 사용하기

$ cat myfirst-pod-for-cm2.yaml
apiVersion: v1
kind: Pod
metadata:
    name: myfirst-pod-example2
spec:
    containers:
        - name: myfirst-pod
          image: busybox
          args: ['tail', '-f', '/dev/null']
          volumeMounts:
              - name: cm-volume
                mountPath: /etc/config

    volumes:
        - name: cm-volume
          configMap:
              name: myfirst-configmap
              #items:
              #    - key: age
              #      page: age_custom  # 특정 key를 특정 파일명으로 생성할 수 있음

$ kubectl exec myfirst-pod-example2 ls /etc/config
age
name

$ kubectl exec myfirst-pod-example2 cat /etc/config/name
tuyy

$ echo hi my name is tuyy >> test.txt

## 파일로 cm 생성. 아래 cmd에서 mytest= 을 test.txt로 가능하며 key가 test.txt 가 된다.
$ kubectl create cm test-file-cm --from-file mytest=test.txt
configmap/test-file-cm created

$ kubectl describe cm test-file-cm
Name:         test-file-cm
Namespace:    devmail
Labels:       <none>
Annotations:  <none>

Data
====
mytest:
----
hi my name is tuyy

Events:  <none>

## 여러 key,value가 포함된 파일로 cm 생성
$ cat multiple-kv.env
NAME1=TUYY1
NAME2=TUYY2
NAME3=TUYY3
NAME4=TUYY4

$ kubectl create cm multiple-kv-cm --from-env-file multiple-kv.env
configmap/multiple-kv-cm created

$ kubectl describe cm multiple-kv-cm
Name:         multiple-kv-cm
Namespace:    devmail
Labels:       <none>
Annotations:  <none>

Data
====
NAME1:
----
TUYY1
NAME2:
----
TUYY2
NAME3:
----
TUYY3
NAME4:
----
TUYY4
Events:  <none>

## yaml 파일로 cm 생성

$ kubectl create cm yaml-cm --from-literal name=tuyy0 --from-literal age=31 --dry-run -o yaml
apiVersion: v1
data:
  age: "31"
  name: tuyy0
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: yaml-cm

$ kubectl create cm yaml-cm --from-literal name=tuyy0 --from-literal age=31 --dry-run -o yaml > yaml-cm.ymal

$ kubectl get cm
NAME                DATA   AGE

$ kubectl apply -f yaml-cm.ymal
configmap/yaml-cm created

$ kubectl get cm
NAME                DATA   AGE
yaml-cm             2      10s

$ kubectl describe cm yaml-cm
Name:         yaml-cm
Namespace:    devmail
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"v1","data":{"age":"31","name":"tuyy0"},"kind":"ConfigMap","metadata":{"annotations":{},"creationTimestamp":null,"name":"yam...

Data
====
age:
----
31
name:
----
tuyy0
Events:  <none>

==========

# Secret은 ConfigMap 과 비슷하다. 암호화가 필요한 값들을 지정해야한다.
* 참고로 yaml 파일에서 configmap 을 secret으로만 치환하면 그대로 사용 가능하다.

$ kubectl create secret generic my-password --from-literal password=1q2w3e4r
secret/my-password created

$ kubectl get secret
NAME                            TYPE                                  DATA   AGE
my-password                     Opaque                                1      7s

$ kubectl describe secret my-password
Name:         my-password
Namespace:    devmail
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
password:  8 bytes

$ kubectl get secret my-password -o yaml
apiVersion: v1
data:
  password: MXEydzNlNHI=
kind: Secret
metadata:
  creationTimestamp: "2021-03-14T02:59:43Z"
  name: my-password
  namespace: devmail
  resourceVersion: "4483816789"
  selfLink: /api/v1/namespaces/devmail/secrets/my-password
  uid: fb573e6d-0095-4f75-8cf1-8bac89e95260
type: Opaque

==========

# kustomize는 kubectl 명령어 1.14 버전부터 사용가능하며, YAML 파일의 속성을 별도로 정의해 사용하거나 여러 YAML을 묶는 등 다양한 용도로 사용할 수 있는 기능이다.

$ cat kustomization.yaml
configMapGenerator:
    - name: kustomize-cm
      files:
          - test1=test.txt
          - test2=test.txt

$ kubectl kustomize .
apiVersion: v1
data:
  test1: |
    hi my name is tuyy
  test2: |
    hi my name is tuyy
kind: ConfigMap
metadata:
  name: kustomize-cm-dmd7fbtd6g

$ kubectl get cm
No resources found.

$ kubectl apply -k .
configmap/kustomize-cm-dmd7fbtd6g created

$ kubectl get cm
NAME                      DATA   AGE
kustomize-cm-dmd7fbtd6g   2      5s

```
