package com.healthcare.auditservice.domain;

import jakarta.persistence.*;
import lombok.*;
import java.time.LocalDateTime;

import java.util.UUID;

@Entity
@Table(name = "audit_logs")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class AuditLog {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    private String eventType;

    @Column(columnDefinition = "TEXT")
    private String eventPayload;

    private LocalDateTime timestamp;

}
