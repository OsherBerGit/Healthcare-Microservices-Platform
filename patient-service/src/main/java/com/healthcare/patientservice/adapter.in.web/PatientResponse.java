package com.healthcare.patientservice.adapter.in.web;

import java.time.LocalDate;
import java.util.UUID;

public record PatientResponse(
        UUID id,
        String firstName,
        String lastName,
        LocalDate dateOfBirth,
        String medicalHistorySummary
) { }
