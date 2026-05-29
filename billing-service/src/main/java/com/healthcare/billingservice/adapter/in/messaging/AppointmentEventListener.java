package com.healthcare.billingservice.adapter.in.messaging;

import com.healthcare.billingservice.application.BillingCommandService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import tools.jackson.databind.JsonNode;
import tools.jackson.databind.ObjectMapper;

@Component
@RequiredArgsConstructor
@Slf4j
public class AppointmentEventListener {

    private final BillingCommandService billingService;
    private final ObjectMapper objectMapper;

    @KafkaListener(topics = "appointment-created", groupId = "billing-group")
    public void handleAppointmentCreated(String eventPayload) {
        try {
            log.info("[BILLING SERVICE] Received new appointment event from Kafka...");

            var eventData = objectMapper.readTree(eventPayload);
            var appointmentId = eventData.get("id").asText();
            var patientId = eventData.get("patientId").asText();

            billingService.generateInvoiceForAppointment(appointmentId, patientId);
        } catch (Exception e) {
            log.error("Failed to process billing event: {}", e.getMessage(), e);
        }
    }
}
