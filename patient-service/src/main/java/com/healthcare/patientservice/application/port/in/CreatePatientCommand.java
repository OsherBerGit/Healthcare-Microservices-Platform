package com.healthcare.patientservice.application.port.in;

import java.time.LocalDate;

public record CreatePatientCommand(
        String firstName,
        String lastName,
        LocalDate dateOfBirth,
        String medicalHistorySummary
) { }
