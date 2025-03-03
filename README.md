# custom-load-balancer

### Step 1: Create the Go Application

Create a simple Go web application (`main.go`):

```
git mod init github/ckcowster/custom-load-balancer
```

```go
package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "Could not get hostname", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello from %s!", hostname)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
```

Create a `Dockerfile` to containerize the Go application:

```Dockerfile
FROM golang:1.16-alpine

WORKDIR /app

COPY main.go /app

RUN go build -o main main.go

CMD ["./main"]
```

### Step 2: Build and Push the Docker Image

Build the Docker image and push it to a container registry:

```sh
docker build -t localhost:5001/web-app:1 -f Dockerfile .
docker push localhost:5001/web-app:1
```

### Step 3: Create Kubernetes Manifests

Create a Deployment for the Go application:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
      - name: web-app
        image: localhost:5001/web-app:1
        ports:
        - containerPort: 80
```

Create a headless service for the Go application:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-app-headless
spec:
  clusterIP: None
  selector:
    app: web-app
  ports:
  - port: 80
    targetPort: 80
```

### Step 4: Implement the Custom Load Balancer in Go

We'll use the same custom load balancer code as before (`load_balancer.go`):

```go
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var podIPs = []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}
var weights = []float64{0.5, 0.3, 0.2}

func weightedChoice(choices []string, weights []float64) string {
	total := 0.0
	for _, weight := range weights {
		total += weight
	}
	r := rand.Float64() * total
	upto := 0.0
	for i, choice := range choices {
		if upto+weights[i] >= r {
			return choice
		}
		upto += weights[i]
	}
	return choices[len(choices)-1]
}

func loadBalance(w http.ResponseWriter, r *http.Request) {
	selectedPod := weightedChoice(podIPs, weights)
	resp, err := http.Get(fmt.Sprintf("http://%s", selectedPod))
	if err != nil {
		http.Error(w, "Failed to reach pod", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", loadBalance)
	http.ListenAndServe(":80", nil)
}
```

Create a `Dockerfile` for the custom load balancer:

```Dockerfile
FROM golang:1.16-alpine

WORKDIR /app

COPY load_balancer.go /app

RUN go build -o load_balancer load_balancer.go

CMD ["./load_balancer"]
```

### Step 5: Build and Push the Docker Image

Build and push the custom load balancer image:

```sh
docker build -t localhost:5001/clb-app:1 -f Dockerfile .
docker push localhost:5001/clb-app:1
```

### Step 6: Create Kubernetes Manifests for the Custom Load Balancer

Create a Deployment for the custom load balancer:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: custom-load-balancer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: custom-load-balancer
  template:
    metadata:
      labels:
        app: custom-load-balancer
    spec:
      containers:
      - name: custom-load-balancer
        image: your-dockerhub-username/custom-load-balancer-go:latest
        env:
        - name: POD_IPS
          value: "10.0.0.1,10.0.0.2,10.0.0.3"  # Replace with actual pod IPs
```

### Step 7: Deploy Everything to Kubernetes

Apply the manifests:

```sh
kubectl apply -f web-app-deployment.yaml
kubectl apply -f web-app-service.yaml
kubectl apply -f clb-app-deployment.yaml
```

### Step 8: Test the Setup

Expose the custom load balancer service using a NodePort or LoadBalancer service type:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: custom-load-balancer-service
spec:
  type: NodePort
  selector:
    app: custom-load-balancer
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30000  # Replace with an available port
```

Apply the service manifest:

```sh
kubectl apply -f clb-app-service.yaml
```

Access the custom load balancer service using the NodePort or LoadBalancer IP and port, and observe how it distributes traffic to the Go application pods.

This example provides a complete setup for deploying a simple web application with a custom load balancer written in Go in Kubernetes. If you have any questions or need further assistance, feel free to ask!

```
kubectl get nodes -o wide


~/github/yamls/hey -n 1000 -c 50 http://172.18.0.2:3000

./doit.sh

```