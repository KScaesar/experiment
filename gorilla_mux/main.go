package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products/{p_id}", ProductHandler).Methods(http.MethodGet)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)
	qs := r.URL.Query()

	// path= map[p_id:aa1231]
	// query= map[page:[2] size:[10]]
	// 2
	fmt.Println("path=", path)
	fmt.Println("query=", qs)
	fmt.Println(qs.Get("page"))

	var err error
	if err != nil {
		// WriteHeader 與 Write 的呼叫順序
		// https://tachingchen.com/tw/blog/pitfall-of-golang-header-operation/?fbclid=IwAR3eGt8aYoTkKakOjxNRmfhq5lJL7vDj7HnT3MwyiseW_tl5lWtbKakPa7E

		ReplyErrorResponse(w, fmt.Errorf("get fail"))
		return
	}

	payload := map[string]string{
		"id": path["p_id"],
	}
	ReplySuccessResponse(w, http.StatusOK, payload)
}

func ReplyErrorResponse(w http.ResponseWriter, err error) {
	resp := &HttpResponse{
		Message: err.Error(),
		Payload: struct{}{},
	}
	data, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(data)
}

func ReplySuccessResponse(w http.ResponseWriter, httpCode int, payload any) {
	if payload == nil {
		payload = struct{}{}
	}

	resp := &HttpResponse{
		Message: "ok",
		Payload: payload,
	}
	data, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(data)
}

type HttpResponse struct {
	Message string `json:"message,omitempty"`
	Payload any    `json:"payload,omitempty"`
}
