package com.healthcare.patientservice.application.service;

import com.healthcare.patientservice.application.port.in.CreatePatientCommand;
import com.healthcare.patientservice.application.port.in.PatientUseCase;
import com.healthcare.patientservice.application.port.out.PatientRepositoryPort;
import com.healthcare.patientservice.domain.Patient;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class PatientManagementService implements PatientUseCase {

    private final PatientRepositoryPort patientRepository;

    public PatientManagementService(PatientRepositoryPort patientRepository) {
        this.patientRepository = patientRepository;
    }

    @Override
    public Patient createPatient(CreatePatientCommand command) {
        var newPatient = Patient.builder()
                .id(UUID.randomUUID())
                .firstName(command.firstName())
                .lastName(command.lastName())
                .dateOfBirth(command.dateOfBirth())
                .medicalHistorySummary(command.medicalHistorySummary())
                .build();

        return patientRepository.save(newPatient);
    }

    @Override
    public Patient getPatient(UUID id) {
        return patientRepository.findById(id).orElseThrow();
    }
}
