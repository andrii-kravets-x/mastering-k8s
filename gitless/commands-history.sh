# Install fluxcd
helm upgrade --install flux-operator oci://ghcr.io/controlplaneio-fluxcd/charts/flux-operator --namespace flux-system --create-namespace
# flux instance
k apply -f flux-instance.yml
# envoy gateway
helm install envoy-gateway oci://docker.io/envoyproxy/gateway-helm --version v1.3.2 --namespace envoy-gateway-system --create-namespace
# envoy gateway config
k apply -f gateway.yml 
# flux github token secret
read -s -GITHUB_TOKEN?"Enter GitHub token: "
# 
kubectl create secret generic github-auth \
  --from-literal=username=git \
  --from-literal=password=${GITHUB_TOKEN} \
  -n flux-system
# flux image registry secret
kubectl create secret docker-registry ghcr-auth \
  --docker-server=ghcr.io \
  --docker-username=den-vasyliev \
  --docker-password=${GITHUB_TOKEN} \
  -n flux-system
# Install gitops stack
k apply -f gitops.yml
# secret in demo needs to be created manually
#gitrepo
# flux install
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
yes|brew install fluxcd/tap/flux
. <(flux completion zsh)
flux reconcile image repository -n demo kbot
# imageautomation
# imagepolicy
# imageupdateautomation
# image update
# apply gitless
k apply -f gitless.yml 
