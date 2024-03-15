package httptransport

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GreetingRequest struct {
	Name string `json:"name"`
}
type GreetingResponse struct {
	Message string `json:"message"`
}

// PostGreet returns the greeting
func (h *Handler) Greet(w http.ResponseWriter, r *http.Request) {
	var requestData GreetingRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := GreetingResponse{
		Message: fmt.Sprintf("greetings %s", requestData.Name),
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
