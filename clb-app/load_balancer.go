package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func getPodIPs() []string {
	podIPsEnv := os.Getenv("POD_IPS")
	return strings.Split(podIPsEnv, ",")
}

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
	podIPs := getPodIPs()
	weights := []float64{0.5, 0.3, 0.2}
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
