from pydantic import BaseModel, Field

class PatientTriageRequest(BaseModel):
    age: int = Field(..., description="Age of the patient")
    heart_rate: int = Field(..., description="Heart rate in BPM")
    blood_pressure_systolic: int = Field(..., description="Systolic blood pressure")
    temperature: float = Field(..., description="Body temperature in Celsius")
    symptoms_summary: str = Field(..., description="Short text describing symptoms")

class TriageResponse(BaseModel):
    acuity_score: int = Field(..., description="Urgency score from 1 (Resuscitation) to 5 (Non-urgent)")
    urgency_level: str = Field(..., description="Text description of urgency")
    confidence: float = Field(..., description="AI confidence in the prediction")