# kubectl-sandbox

Kubectl Sandbox plugin gives you sandbox environment for kubectl. It will download k3s (k3s.io) and register as a service. Lightweight Kubernetes will start in your environment and you could try whatever you want.

#### spoiler alert

 ```                  
              This plugin needs ROOT access to create|start|stop|delete k3s service.
 ```

- `kubectl sandbox` configure environment and start your sandbox
- `kubectl sandbox load` load [sample app (guestbook)](https://raw.githubusercontent.com/kubernetes/examples/master/guestbook/all-in-one/guestbook-all-in-one.yaml) for your k3s instance 
- `kubectl sandbox delete` delete your k3s instance
- `kubectl sandbox reset` reset your k3s instance
