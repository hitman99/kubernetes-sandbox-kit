package kubernetes

const kubeconfig = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %s
    server: %s
  name: sandbox
contexts:
- context:
    cluster: sandbox
    user: sandbox-user
    namespace: %s
  name: sandbox
current-context: sandbox
kind: Config
preferences: {}
users:
- name: sandbox-user
  user:
    token: %s
`
