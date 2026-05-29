package com.healthcare.auditservice.listener;

import com.healthcare.auditservice.domain.AuditLog;
import com.healthcare.auditservice.repository.AuditLogRepository;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
public class AuditEventListener {

    private final AuditLogRepository auditLogRepository;

    public AuditEventListener(AuditLogRepository auditLogRepository) {
        this.auditLogRepository = auditLogRepository;
    }

    @KafkaListener(topics = "appointment-created", groupId = "audit-compliance-group")
    public void consumeEvent(String message) {
        System.out.println("Audit Service received event from Kafka! Documenting...");

        var logEntry = AuditLog.builder()
                .eventType("SYSTEM_EVENT")
                .eventPayload(message)
                .timestamp(LocalDateTime.now())
                .build();

        auditLogRepository.save(logEntry);

        System.out.println("Event successfully saved to audit_db");
    }
}
