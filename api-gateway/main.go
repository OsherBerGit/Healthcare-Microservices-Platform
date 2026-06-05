package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() *trace.TracerProvider {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("jaeger:4318"),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "api-gateway"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	patientServiceURLs := []string{"http://patient:8081"}
	appointmentServiceURLs := []string{"http://appointment:8082"}
	triageServiceURLs := []string{"http://triage:8000"}

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

	wrappedMux := otelhttp.NewHandler(securityHeadersMiddleware(loggingMiddleware(mux)), "api-gateway")

	log.Println("Starting Healthcare API Gateway with Load Balancing on port 8080...")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
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
	ctx := r.Context()

	patientReq, _ := json.Marshal(map[string]interface{}{"name": reqData["name"]})
	resp1, err := makeSagaRequest(ctx, "POST", "http://localhost:8080/api/patients", patientReq, token)
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

	resp2, err := makeSagaRequest(ctx, "POST", "http://localhost:8080/api/appointments", appointmentReq, token)
	if err != nil || (resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusCreated) {
		makeSagaRequest(ctx, "DELETE", "http://localhost:8080/api/patients/rollback", patientReq, token)
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

	resp3, err := makeSagaRequest(ctx, "POST", "http://localhost:8080/api/triage/predict", triageReq, token)
	if err != nil || resp3.StatusCode != http.StatusOK {
		makeSagaRequest(ctx, "DELETE", "http://localhost:8080/api/appointments/rollback", appointmentReq, token)
		makeSagaRequest(ctx, "DELETE", "http://localhost:8080/api/patients/rollback", patientReq, token)
		http.Error(w, "Saga failed at Triage prediction. Appointment & Patient rolled back.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Full admission process completed successfully via Saga"}`))
}

func makeSagaRequest(ctx context.Context, method, url string, body []byte, token string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	idempotencyKey := fmt.Sprintf("saga-%d", time.Now().UnixNano())
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("X-Idempotency-Key", idempotencyKey)

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return client.Do(req)
}