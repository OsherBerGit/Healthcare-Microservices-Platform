from fastapi import FastAPI, HTTPException
from schemas import PatientTriageRequest, TriageResponse
import joblib
import pandas as pd
import os

app = FastAPI(
    title="Triage AI Service",
    description="Predicts patient urgency using Machine Learning",
    version="1.0.0"
)

script_dir = os.path.dirname(os.path.abspath(__file__))
MODEL_PATH = os.path.join(script_dir, "triage_model.pkl")

if os.path.exists(MODEL_PATH):
    model = joblib.load(MODEL_PATH)
    print(f"ML Model loaded successfully from: {MODEL_PATH}")
else:
    model = None
    print(f"Warning: Model not found at {MODEL_PATH}! Please run train_model.py first.")

@app.get("/health")
async def health_check():
    return {"status": "ok", "model_loaded": model is not None}

@app.post("/api/triage/predict", response_model=TriageResponse)
async def predict_urgency(request: PatientTriageRequest):
    if model is None:
        raise HTTPException(status_code=500, detail="ML Model is not loaded on the server.")
    
    input_data = pd.DataFrame([{
        'age': request.age,
        'heart_rate': request.heart_rate,
        'blood_pressure_systolic': request.blood_pressure_systolic,
        'temperature': request.temperature
    }])
    
    prediction = model.predict(input_data)[0]
    probabilities = model.predict_proba(input_data)[0]

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