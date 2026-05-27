package com.healthcare.patientservice.adapter.in.web;

import java.time.LocalDate;

public record CreatePatientRequest(
        String firstName,
        String lastName,
        LocalDate dateOfBirth,
        String medicalHistorySummary
) { }
