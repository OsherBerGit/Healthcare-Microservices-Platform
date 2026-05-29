package com.healthcare.appointmentservice.adapter.in.web;

import java.time.LocalDateTime;

public record AppointmentRequest(
        String patientId,
        String doctorId,
        LocalDateTime appointmentTime,
        String reasonForVisit
) { }
