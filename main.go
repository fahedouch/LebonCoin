package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Stats struct {
	sync.Mutex
	Requests map[string]int
}

var stats Stats

func fizzBuzzHandler(w http.ResponseWriter, r *http.Request) {
	str1 := r.FormValue("str1")
	str2 := r.FormValue("str2")
	int1Str := r.FormValue("int1")
	int2Str := r.FormValue("int2")
	limitStr := r.FormValue("limit")

	int1, err := strconv.Atoi(int1Str)
	if err != nil {
		http.Error(w, "Invalid int1", http.StatusBadRequest)
		return
	}

	int2, err := strconv.Atoi(int2Str)
	if err != nil {
		http.Error(w, "Invalid int2", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	params := fmt.Sprintf("int1=%s&int2=%s&limit=%s&str1=%s&str2=%s", int1Str, int2Str, limitStr, str1, str2)

	// use lock to ensure safe write race
	stats.Lock()
	stats.Requests[params]++
	stats.Unlock()

	var result []string
	for i := 1; i <= limit; i++ {
		if i%int1 == 0 && i%int2 == 0 {
			result = append(result, str1+str2)
		} else if i%int1 == 0 {
			result = append(result, str1)
		} else if i%int2 == 0 {
			result = append(result, str2)
		} else {
			result = append(result, strconv.Itoa(i))
		}
	}

	response := strings.Join(result, ",")

	fmt.Fprint(w, response)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// use lock to ensure safe read race
	stats.Lock()
	defer stats.Unlock()

	maxHits := 0
	var maxParams string

	for params, hits := range stats.Requests {
		if hits > maxHits {
			maxHits = hits
			maxParams = params
		}
	}

	fmt.Fprintf(w, "Max hitted request: %s\nHits: %d", maxParams, maxHits)
}

func main() {
	stats = Stats{Requests: make(map[string]int)}

	// handlers
	http.HandleFunc("/fizzbuzz", fizzBuzzHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
