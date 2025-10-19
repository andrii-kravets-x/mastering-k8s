# Task
## List
1. Початківці: використовуючи гайд, створити, збілдати та запустити свій перший контролер.
2. Просунутий рівень: розібратися з тестами, отримати метрики контролера
3. Мах: запустити контролер в контрл плаін (в кубері) з leader election: true

Examples:
- репо з прикладом що розбирали на занняті https://github.com/den-vasyliev/mastering-k8s/tree/main/new-controller
- Офіційно і коротко про контролери https://kubernetes.io/docs/concepts/architecture/controller/

## Solution for task_1(beginner level)

We will use **kubebuilder**

See other options [here](./ways-to-build-controllers.md)

### 1. Download kubebuilder and install locally.
```bash
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```

### 2. Init project
`kubebuilder init --domain my.domain --repo my.domain/guestbook`

### 3. Controller-code scaffold
**Warning! this command was used to auto-generate, but code was changed afterwards, use files in controller-root-directory**
```
kubebuilder create api --group web --version v1 --kind App --resource --controller
```

### 4. Generate openapi (kustomize)
```bash
make generate
make manifests
```

### 5. Generate and install CRD to cluster
```bash
kustomize build config/crd > CustomResourceDefinition.yaml

kubectl apply -f CustomResourceDefinition.yaml
customresourcedefinition.apiextensions.k8s.io/apps.web.my.domain created
```

### 6. Check CRDs in cluster
```bash
kubectl get crd
NAME                 CREATED AT
apps.web.my.domain   2025-10-19T19:53:36Z
```

### 7. Create CR using generated sample manifest
```bash
kubectl apply -f config/samples/web_v1_app.yaml
app.web.my.domain/app-sample created


kubectl get apps.web.my.domain 
NAME         AGE
app-sample   2m27s


# this describe output is taken after controller was launched
kubectl describe apps.web.my.domain
Name:         app-sample
Namespace:    default
Labels:       app.kubernetes.io/managed-by=kustomize
              app.kubernetes.io/name=controller-root-directory
Annotations:  <none>
API Version:  web.my.domain/v1
Kind:         App
Metadata:
  Creation Timestamp:  2025-10-19T21:34:00Z
  Generation:          1
  Resource Version:    81882
  UID:                 feac91f8-a18f-439f-9d13-0645ec6b384e
Spec:
  Foo:  this is string for my first controller
Status:
  Ready:  true
Events:   <none>
```

### 8. Build and launch controller
```bash
go build -o bin/manager cmd/main.go
./bin/manager
```

See logs and status change on [screenshot](./2025-10-19_23-37.png)




### In case you fail - delete CRD (this will also delete CR defined by CRD) and retry with fixes
```
kubectl delete --ignore-not-found=true -f CustomResourceDefinition.yaml
```
