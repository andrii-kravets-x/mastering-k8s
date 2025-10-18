#!/bin/bash

# Exit on error
set -e

setup_configs() {
  # Get the host IP address
HOST_IP=$(hostname -I | awk '{print $1}')
POD_CIDR=192.168.0.0/24
SAN=10.10.0.1
SVC_CIDR=10.10.0.0/24
DNS_IP=10.10.0.10

    # Generate certificates and tokens if they don't exist
    openssl genrsa -out /tmp/sa.key 2048
    openssl rsa -in /tmp/sa.key -pubout -out /tmp/sa.pub
    TOKEN="1234567890"
    echo "${TOKEN},admin,admin,system:masters" > /tmp/token.csv

    # Always regenerate and copy CA certificate to ensure it exists
    echo "Generating CA certificate..."
    openssl genrsa -out /tmp/ca.key 2048
    openssl req -x509 -new -nodes -key /tmp/ca.key -subj "/CN=kubelet-ca" -days 365 -out /tmp/ca.crt
    
    sudo mkdir -p /etc/kubernetes/pki
    sudo cp /tmp/ca.crt /etc/kubernetes/pki/ca.crt
    sudo cp /tmp/ca.key /etc/kubernetes/pki/ca.key
    
    sudo mkdir -p /var/lib/kubelet/pki
    sudo cp /tmp/ca.crt /var/lib/kubelet/ca.crt
    sudo cp /tmp/ca.crt /var/lib/kubelet/pki/ca.crt


### Generate apiserver serving cert signed by CA
sudo cat > /tmp/apiserver-openssl.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
req_extensions     = req_ext
distinguished_name = dn

[ dn ]
CN = kube-apiserver

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
IP.1 = ${SAN}
IP.2 = ${HOST_IP}
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster.local
EOF

sudo openssl genrsa -out /etc/kubernetes/pki/apiserver.key 2048
sudo openssl req -new -key /etc/kubernetes/pki/apiserver.key -out /tmp/apiserver.csr -config /tmp/apiserver-openssl.cnf
sudo openssl x509 -req -in /tmp/apiserver.csr \
  -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial \
  -out /etc/kubernetes/pki/apiserver.crt -days 365 \
  -extensions req_ext -extfile /tmp/apiserver-openssl.cnf

    # Set up kubeconfig if not already configured
        sudo kubebuilder/bin/kubectl config set-credentials test-user --token=1234567890
        sudo kubebuilder/bin/kubectl config set-cluster test-env --server=https://${HOST_IP}:6443 --certificate-authority=/etc/kubernetes/pki/ca.crt
        sudo kubebuilder/bin/kubectl config set-context test-context --cluster=test-env --user=test-user --namespace=default 
        sudo kubebuilder/bin/kubectl config use-context test-context


    # Ensure proper permissions
    sudo chmod 644 /var/lib/kubelet/ca.crt
    sudo chmod 644 /var/lib/kubelet/config.yaml

# Generate kubelet serving certificate signed by cluster CA if not present
  echo "Generating kubelet serving certificate signed by cluster CA..."

  sudo mkdir -p /var/lib/kubelet/pki

  # Create CSR config
sudo  cat > /tmp/kubelet-openssl.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
req_extensions     = req_ext
distinguished_name = dn

[ dn ]
CN = system:node:$(hostname)

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $(hostname)
IP.1  = ${HOST_IP}
EOF

  # Generate key + CSR
  sudo openssl genrsa -out /var/lib/kubelet/pki/kubelet.key 2048
  sudo openssl req -new -key /var/lib/kubelet/pki/kubelet.key \
    -out /tmp/kubelet.csr -config /tmp/kubelet-openssl.cnf

  # Sign with cluster CA
  sudo openssl x509 -req -in /tmp/kubelet.csr \
    -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial \
    -out /var/lib/kubelet/pki/kubelet.crt -days 365 \
    -extensions req_ext -extfile /tmp/kubelet-openssl.cnf

  sudo chmod 600 /var/lib/kubelet/pki/kubelet.key
  sudo chmod 644 /var/lib/kubelet/pki/kubelet.crt

# Set up kubelet kubeconfig
sudo cp /root/.kube/config /var/lib/kubelet/kubeconfig
export KUBECONFIG=~/.kube/config
cp /tmp/sa.pub /tmp/ca.crt


# Create service account and configmap if they don't exist
sudo kubebuilder/bin/kubectl create sa default 2>/dev/null || true
sudo kubebuilder/bin/kubectl create configmap kube-root-ca.crt --from-file=ca.crt=/tmp/ca.crt -n default 2>/dev/null || true

}

setup_configs
