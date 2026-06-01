package com.healthcare.appointmentservice.domain.repository;

import com.healthcare.appointmentservice.domain.model.Appointment;

import java.util.Optional;
import java.util.UUID;

public interface AppointmentRepository {
    Appointment save(Appointment appointment);
    Optional<Appointment> findById(UUID id);
    void deleteByPatientId(String patientId);
}
