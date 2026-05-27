package com.healthcare.patientservice.domain;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDate;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Patient {
    private UUID id;

    private String firstName;
    private String lastName;
    private LocalDate dateOfBirth;

    private String medicalHistorySummary;

    public boolean isMinor() {
        return LocalDate.now().minusYears(18).isBefore(this.dateOfBirth);
    }
}