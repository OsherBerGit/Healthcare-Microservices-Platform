package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	patientServiceURLs := []string{"http://localhost:8081"}
	appointmentServiceURLs := []string{"http://localhost:8082"}
	triageServiceURLs := []string{"http://127.0.0.1:8000"}

	mux := http.NewServeMux()

	patientLB := NewLoadBalancer(patientServiceURLs)
	appointmentLB := NewLoadBalancer(appointmentServiceURLs)
	triageLB := NewLoadBalancer(triageServiceURLs)

	mux.Handle("/api/patients", verifyJWTMiddleware(patientLB))
	mux.Handle("/api/patients/", verifyJWTMiddleware(patientLB))

	mux.Handle("/api/appointments", verifyJWTMiddleware(idempotencyMiddleware(appointmentLB)))
	mux.Handle("/api/appointments/", verifyJWTMiddleware(idempotencyMiddleware(appointmentLB)))

	mux.Handle("/api/triage", verifyJWTMiddleware(triageLB))
	mux.Handle("/api/triage/", verifyJWTMiddleware(triageLB))

	mux.Handle("/api/admission/full-process", verifyJWTMiddleware(http.HandlerFunc(sagaAdmissionHandler)))

	mux.HandleFunc("/health", healthCheckHandler)

	log.Println("Starting Healthcare API Gateway with Load Balancing on port 8080...")
	if err := http.ListenAndServe(":8080", securityHeadersMiddleware(loggingMiddleware(mux))); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is Up and Running\n"))
}

func sagaAdmissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")

	patientReq, _ := json.Marshal(map[string]interface{}{"name": reqData["name"]})
	resp1, err := makeSagaRequest("POST", "http://localhost:8080/api/patients", patientReq, token)
	if err != nil || (resp1.StatusCode != http.StatusOK && resp1.StatusCode != http.StatusCreated) {
		http.Error(w, "Saga failed at Patient creation", http.StatusInternalServerError)
		return
	}

	parsedTime, err := time.Parse("2006-01-02", reqData["date"].(string))
	if err != nil {
		parsedTime = time.Now().Add(24 * time.Hour)
	}
	formattedTime := parsedTime.Format("2006-01-02T15:04:05")

	appointmentReq, _ := json.Marshal(map[string]interface{}{
		"patientId":       reqData["name"],
		"doctorId":        "doc-default-001",
		"appointmentTime": formattedTime,
		"reasonForVisit":  "Triage Admission Process",
	})

	resp2, err := makeSagaRequest("POST", "http://localhost:8080/api/appointments", appointmentReq, token)
	if err != nil || (resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusCreated) {
		makeSagaRequest("DELETE", "http://localhost:8080/api/patients/rollback", patientReq, token)
		http.Error(w, "Saga failed at Appointment creation. Patient rolled back.", http.StatusInternalServerError)
		return
	}

	triageReq, _ := json.Marshal(map[string]interface{}{
		"age":                     reqData["age"],
		"heart_rate":              reqData["heart_rate"],
		"blood_pressure_systolic": reqData["blood_pressure_systolic"],
		"temperature":             reqData["temperature"],
		"symptoms_summary":        "Patient admitted via Saga orchestration",
	})

	resp3, err := makeSagaRequest("POST", "http://localhost:8080/api/triage/predict", triageReq, token)
	if err != nil || resp3.StatusCode != http.StatusOK {
		makeSagaRequest("DELETE", "http://localhost:8080/api/appointments/rollback", appointmentReq, token)
		makeSagaRequest("DELETE", "http://localhost:8080/api/patients/rollback", patientReq, token)
		http.Error(w, "Saga failed at Triage prediction. Appointment & Patient rolled back.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Full admission process completed successfully via Saga"}`))
}

func makeSagaRequest(method, url string, body []byte, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	idempotencyKey := fmt.Sprintf("saga-%d", time.Now().UnixNano())
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("X-Idempotency-Key", idempotencyKey)

	client := &http.Client{}
	return client.Do(req)
}
