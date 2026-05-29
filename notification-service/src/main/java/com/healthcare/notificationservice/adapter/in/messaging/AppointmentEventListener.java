package com.healthcare.notificationservice.adapter.in.messaging;

import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@Slf4j
public class AppointmentEventListener {
    @KafkaListener(topics = "appointment-created", groupId = "notification-group")
    public void handleAppointmentCreated(String eventPayload) {
        log.info("[NOTIFICATION SERVICE] Received new event from Kafka!");
        log.info("Sending confirmation email to patient...");
        log.info("Event Details: {}", eventPayload);
        log.info("Email sent successfully!\n");
    }
}
