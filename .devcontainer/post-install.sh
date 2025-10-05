#!/usr/bin/env bash
set -e

# Detect architecture (arm64 for Apple Silicon, amd64 for Intel)
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
  ARCH_NAME="arm64"
else
  ARCH_NAME="amd64"
fi

# --- KIND ---
echo "Installing kind for Darwin ($ARCH_NAME)..."
curl -Lo ./kind "https://kind.sigs.k8s.io/dl/latest/kind-darwin-${ARCH_NAME}"
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# --- Kubebuilder ---
echo "Installing kubebuilder for Darwin ($ARCH_NAME)..."
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/darwin/${ARCH_NAME}"
chmod +x kubebuilder
sudo mv kubebuilder /usr/local/bin/

# --- kubectl ---
echo "Installing kubectl for Darwin ($ARCH_NAME)..."
# Get the latest stable version
KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)

curl -Lo kubectl "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/darwin/${ARCH_NAME}/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/kubectl

# --- Docker network ---
docker network create -d=bridge --subnet=172.19.0.0/24 kind || true

# --- Verify installs ---
kind version
kubebuilder version
docker --version
go version
kubectl version --client
