package com.healthcare.appointmentservice.adapter.out.persistence;

import com.healthcare.appointmentservice.domain.model.Appointment;
import com.healthcare.appointmentservice.domain.repository.AppointmentRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
@RequiredArgsConstructor
public class AppointmentPersistenceAdapter implements AppointmentRepository {

    private final AppointmentJpaRepository jpaRepository;

    @Override
    public Appointment save(Appointment appointment) {
        return jpaRepository.save(appointment);
    }

    @Override
    public Optional<Appointment> findById(UUID id) {
        return jpaRepository.findById(id);
    }

    @Override
    public void deleteByPatientId(String patientId) { jpaRepository.deleteByPatientId(patientId); }
}
