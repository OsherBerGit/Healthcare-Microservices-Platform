package com.healthcare.appointmentservice.adapter.out.persistence;

import com.healthcare.appointmentservice.domain.model.Appointment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface AppointmentJpaRepository extends JpaRepository<Appointment, UUID> {
}
