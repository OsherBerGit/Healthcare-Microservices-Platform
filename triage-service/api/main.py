import os
import joblib
import pandas as pd
from fastapi import FastAPI, HTTPException
from schemas import PatientTriageRequest, TriageResponse
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor

resource = Resource.create({"service.name": "triage-service"})
provider = TracerProvider(resource=resource)
processor = BatchSpanProcessor(OTLPSpanExporter(endpoint="http://jaeger:4318/v1/traces"))
provider.add_span_processor(processor)
trace.set_tracer_provider(provider)

app = FastAPI(
    title="Triage AI Service",
    description="Predicts patient urgency using Real-World ICU Data Artifacts",
    version="2.0.0"
)

FastAPIInstrumentor.instrument_app(app)

script_dir = os.path.dirname(os.path.abspath(__file__))
ARTIFACTS_PATH = os.path.join(script_dir, "..", "data", "models", "triage_artifacts.pkl")

if os.path.exists(ARTIFACTS_PATH):
    artifacts = joblib.load(ARTIFACTS_PATH)
    model = artifacts['model']
    imputer = artifacts['imputer']
    scaler = artifacts['scaler']
    print("ML Artifacts (Model, Imputer, Scaler) loaded successfully!")
else:
    artifacts = None
    print(f"Warning: Artifacts not found at {ARTIFACTS_PATH}! Run train_model.py first.")

@app.get("/health")
async def health_check():
    return {"status": "ok", "model_loaded": artifacts is not None}

@app.post("/api/triage/predict", response_model=TriageResponse)
async def predict_urgency(request: PatientTriageRequest):
    if artifacts is None:
        raise HTTPException(status_code=500, detail="ML Artifacts are not loaded on the server.")

    raw_input = pd.DataFrame([{
        'age': request.age,
        'heart_rate': request.heart_rate,
        'blood_pressure_systolic': request.blood_pressure_systolic,
        'temperature': request.temperature
    }])

    try:
        bp = request.blood_pressure_systolic if request.blood_pressure_systolic > 0 else 1
        raw_input['shock_index'] = request.heart_rate / bp
        raw_input['high_risk_elderly'] = int((request.age > 65) and ((request.heart_rate > 110) or (request.temperature > 38.5)))

        imputed_data = imputer.transform(raw_input)
        scaled_data = scaler.transform(imputed_data)
        
        prediction = model.predict(scaled_data)[0]
        probabilities = model.predict_proba(scaled_data)[0]
        confidence = max(probabilities)
        
        level_map = {
            2: "Emergent (High Urgency) - Immediate Attention",
            3: "Urgent (Medium Urgency) - See within 30 mins",
            4: "Less Urgent (Routine) - Safe to wait"
        }

        return TriageResponse(
            acuity_score=int(prediction),
            urgency_level=level_map.get(int(prediction), "Unknown"),
            confidence=float(confidence)
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Prediction error: {str(e)}")