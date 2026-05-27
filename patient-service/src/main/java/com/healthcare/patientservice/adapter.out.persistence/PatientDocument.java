package com.healthcare.patientservice.adapter.out.persistence;

import lombok.Builder;
import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDate;
import java.util.UUID;

@Data
@Builder
@Document(collection = "patients")
public class PatientDocument {
    @Id
    private UUID id;
    private String firstName;
    private String lastName;
    private LocalDate dateOfBirth;
    private String medicalHistorySummary;
}
