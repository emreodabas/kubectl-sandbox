# kubectl-sandbox

Kubectl Sandbox plugin gives you sandbox environment for kubectl. It will download k3s (k3s.io) and register as a service. Lightweight Kubernetes will start in your environment and you could try whatever you want.

#### spoiler alert

 ```                  
              Plugin works only for linux distributions
              This plugin needs ROOT access to create|start|stop|delete k3s service.
              
 ```

- `kubectl sandbox` configure environment and start your sandbox
- `kubectl sandbox load` load sample app [(guestbook)](https://raw.githubusercontent.com/kubernetes/examples/master/guestbook/all-in-one/guestbook-all-in-one.yaml) for your k3s instance 
- `kubectl sandbox delete` delete your k3s instance
- `kubectl sandbox reset` reset your k3s instance


## Installation

 Download binary
 https://github.com/emreodabas/kubectl-sandbox/releases/latest
 
 ```
mkdir -p /opt/kubectl-plugins
tar xzvf kubectl-sandbox_{version}.tar.gz -C /opt/kubectl-plugins
sudo ln -s /opt/kubectl-plugins/kubectl-sandbox /usr/local/bin/kubectl-sandbox

 ```
