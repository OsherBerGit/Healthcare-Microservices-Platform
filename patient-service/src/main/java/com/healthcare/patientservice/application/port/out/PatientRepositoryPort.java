package com.healthcare.patientservice.application.port.out;

import com.healthcare.patientservice.domain.Patient;

import java.util.Optional;
import java.util.UUID;

public interface PatientRepositoryPort {
    Patient save(Patient patient);
    Optional<Patient> findById(UUID id);
    void deleteByFirstName(String firstName);
}
