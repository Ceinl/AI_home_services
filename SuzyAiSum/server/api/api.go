package api

import (
	"encoding/json"
	"log"
	"net/http"
)
type Api struct{
	API_KEY string
	SYSTEM_PROMPT string
}

func NewApi(API_KEY string) *Api {
	return &Api{
		API_KEY: API_KEY,
	}
}


func (a *Api) POSTSummarize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		User_input string `json:"user_input"`
	}
	var req RequestBody 

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid request.", http.StatusBadRequest)
		return
	}

	if  req.User_input == "" {
		http.Error(w, "Empty request.", http.StatusBadRequest)
		return
	}
	if a.API_KEY == "" {
		http.Error(w, "API key is not set.", http.StatusBadRequest)
		return
	}

	log.Println("User input:",req.User_input) 
	summary, err := Summarize(req.User_input,a.SYSTEM_PROMPT ,a.API_KEY)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to summarize.", http.StatusBadRequest)
		return
	}
	log.Println("Summary: ",summary)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"summary": summary,
	})
}

func (a *Api) POSTSystemPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		System_prompt string `json:"system_prompt"`
	}

	var req RequestBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid request.", http.StatusBadRequest)
		return
	}

	a.SYSTEM_PROMPT = req.System_prompt

	log.Println("Set system prompt:",a.SYSTEM_PROMPT)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"system_prompt": a.SYSTEM_PROMPT,
	})
}

