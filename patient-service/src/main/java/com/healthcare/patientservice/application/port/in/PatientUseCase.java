package com.healthcare.patientservice.application.port.in;

import com.healthcare.patientservice.domain.Patient;

import java.util.UUID;

public interface PatientUseCase {
    Patient createPatient(CreatePatientCommand command);
    Patient getPatient(UUID id);
}
