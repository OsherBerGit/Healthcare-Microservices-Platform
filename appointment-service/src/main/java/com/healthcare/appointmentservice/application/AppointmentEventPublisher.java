package com.healthcare.appointmentservice.application;

import com.healthcare.appointmentservice.domain.model.Appointment;

public interface AppointmentEventPublisher {
    void publishAppointmentCreatedEvent(Appointment appointment);
}
