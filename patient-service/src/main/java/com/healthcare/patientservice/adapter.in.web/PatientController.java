package com.healthcare.patientservice.adapter.in.web;

import com.healthcare.patientservice.application.port.in.CreatePatientCommand;
import com.healthcare.patientservice.application.port.in.PatientUseCase;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/patients")
public class PatientController {

    private final PatientUseCase patientUseCase;

    public PatientController(PatientUseCase patientUseCase) {
        this.patientUseCase = patientUseCase;
    }

    @PostMapping({"", "/"})
    public ResponseEntity<PatientResponse> createPatient(@RequestBody CreatePatientRequest request) {
        var command = new CreatePatientCommand(
                request.firstName(),
                request.lastName(),
                request.dateOfBirth(),
                request.medicalHistorySummary()
        );

        var patient = patientUseCase.createPatient(command);

        var response = new PatientResponse(
                patient.getId(),
                patient.getFirstName(),
                patient.getLastName(),
                patient.getDateOfBirth(),
                patient.getMedicalHistorySummary()
        );

        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<PatientResponse> getPatient(@PathVariable UUID id) {
        var patient = patientUseCase.getPatient(id);

        var response = new PatientResponse(
                patient.getId(),
                patient.getFirstName(),
                patient.getLastName(),
                patient.getDateOfBirth(),
                patient.getMedicalHistorySummary()
        );

        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/rollback")
    public ResponseEntity<Void> rollbackPatient(@RequestBody CreatePatientRequest request) {
        patientUseCase.deletePatientByName(request.firstName());
        return ResponseEntity.noContent().build();
    }
}
