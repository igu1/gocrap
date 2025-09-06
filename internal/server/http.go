package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/igu1/gocrap/internal/flow"
)

func RegisterRoutes() {
	http.HandleFunc("/run", flowHandler)
}

func flowHandler(w http.ResponseWriter, r *http.Request) {
	f := flow.Flow{Mem: make(map[string]interface{})}
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := f.Run(f)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res.Mem); err != nil {
		fmt.Println("encode error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
