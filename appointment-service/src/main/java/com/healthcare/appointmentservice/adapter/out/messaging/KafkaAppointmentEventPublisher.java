package com.healthcare.appointmentservice.adapter.out.messaging;

import com.healthcare.appointmentservice.application.AppointmentEventPublisher;
import com.healthcare.appointmentservice.domain.model.Appointment;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;
import tools.jackson.databind.ObjectMapper;

@Component
@RequiredArgsConstructor
@Slf4j
public class KafkaAppointmentEventPublisher implements AppointmentEventPublisher {

    private static final String TOPIC = "appointment-created";

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    @Override
    public void publishAppointmentCreatedEvent(Appointment appointment) {
        try {
            var eventPayload = objectMapper.writeValueAsString(appointment);

            kafkaTemplate.send(TOPIC, appointment.getId().toString(), eventPayload);

            log.info("Successfully published AppointmentCreatedEvent for ID: {}", appointment.getId());
        } catch (Exception e) {
            log.error("Failed to publish event to Kafka for Appointment ID: {}", appointment.getId(), e);
        }
    }
}
