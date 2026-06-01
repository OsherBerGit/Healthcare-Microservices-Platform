package com.healthcare.appointmentservice.application;

import com.healthcare.appointmentservice.domain.model.Appointment;
import com.healthcare.appointmentservice.domain.model.AppointmentStatus;
import com.healthcare.appointmentservice.domain.repository.AppointmentRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
public class AppointmentCommandService {

    private final AppointmentRepository appointmentRepository;
    private final AppointmentEventPublisher eventPublisher;

    @Transactional
    public Appointment scheduleAppointment(Appointment appointment) {
        appointment.setStatus(AppointmentStatus.SCHEDULED);

        var savedAppointment = appointmentRepository.save(appointment);

        eventPublisher.publishAppointmentCreatedEvent(savedAppointment);

        return savedAppointment;
    }

    @Transactional
    public void deleteAppointmentByPatientId(String patientId) {
        appointmentRepository.deleteByPatientId(patientId);
    }
}
