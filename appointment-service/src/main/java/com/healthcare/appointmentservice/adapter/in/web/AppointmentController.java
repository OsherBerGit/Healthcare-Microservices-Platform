package com.healthcare.appointmentservice.adapter.in.web;

import com.healthcare.appointmentservice.application.AppointmentCommandService;
import com.healthcare.appointmentservice.domain.model.Appointment;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/appointments")
@RequiredArgsConstructor
public class AppointmentController {

    private final AppointmentCommandService commandService;

    @PostMapping
    public ResponseEntity<Appointment> scheduleAppointment(@RequestBody AppointmentRequest request) {
        var appointmentToSchedule = Appointment.builder()
                .patientId(request.patientId())
                .doctorId(request.doctorId())
                .appointmentTime(request.appointmentTime())
                .reasonForVisit(request.reasonForVisit())
                .build();

        var scheduledAppointment = commandService.scheduleAppointment(appointmentToSchedule);

        return ResponseEntity.status(HttpStatus.CREATED).body(scheduledAppointment);
    }
}
