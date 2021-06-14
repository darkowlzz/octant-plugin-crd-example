#!/usr/bin/env bash

set -eu

# This script queries kubernetes with the existing kubeconfig to gather
# information and create a new kubeconfig with a service account token.
# Refer: https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/#without-kubectl-proxy

namePrefix="box-user"
namespace="box"
kubeconfigFile="box-user.kubeconfig"

server=$(kubectl config view --minify | grep server | cut -f 2- -d ":" | tr -d " ")
secretName=$(kubectl -n ${namespace} get secrets | grep ^${namePrefix} | cut -f1 -d ' ')
ca=$(kubectl -n ${namespace} get secret/${secretName} -o jsonpath='{.data.ca\.crt}')
token=$(kubectl -n ${namespace} get secret/${secretName} -o jsonpath='{.data.token}' | base64 --decode)

echo "
apiVersion: v1
kind: Config
clusters:
- name: default-cluster
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: default-context
  context:
    cluster: default-cluster
    namespace: ${namespace}
    user: default-user
current-context: default-context
users:
- name: default-user
  user:
    token: ${token}
" > ${kubeconfigFile}

echo "New kubeconfig written into ${kubeconfigFile}"
