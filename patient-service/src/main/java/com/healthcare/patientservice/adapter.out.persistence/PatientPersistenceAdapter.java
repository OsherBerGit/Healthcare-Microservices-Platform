package com.healthcare.patientservice.adapter.out.persistence;

import com.healthcare.patientservice.application.port.out.PatientRepositoryPort;
import com.healthcare.patientservice.domain.Patient;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public class PatientPersistenceAdapter implements PatientRepositoryPort {

    private final PatientMongoRepository mongoRepository;

    public PatientPersistenceAdapter(PatientMongoRepository mongoRepository) {
        this.mongoRepository = mongoRepository;
    }

    @Override
    public Patient save(Patient patient) {
        var document = PatientDocument.builder()
                .id(patient.getId())
                .firstName(patient.getFirstName())
                .lastName(patient.getLastName())
                .dateOfBirth(patient.getDateOfBirth())
                .medicalHistorySummary(patient.getMedicalHistorySummary())
                .build();

        var savedDocument = mongoRepository.save(document);
        return patient;
    }

    @Override
    public Optional<Patient> findById(UUID id) {
        return mongoRepository.findById(id).map(doc ->
                Patient.builder()
                        .id(doc.getId())
                        .firstName(doc.getFirstName())
                        .lastName(doc.getLastName())
                        .dateOfBirth(doc.getDateOfBirth())
                        .medicalHistorySummary(doc.getMedicalHistorySummary())
                        .build()
        );
    }

    @Override
    public void deleteByFirstName(String firstName) {
        mongoRepository.deleteByFirstName(firstName);
    }
}
