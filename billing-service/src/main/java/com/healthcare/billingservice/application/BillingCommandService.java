package com.healthcare.billingservice.application;

import com.healthcare.billingservice.domain.model.Invoice;
import com.healthcare.billingservice.domain.model.PaymentStatus;
import com.healthcare.billingservice.domain.repository.InvoiceRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Service
@RequiredArgsConstructor
@Slf4j
public class BillingCommandService {

    private final InvoiceRepository invoiceRepository;

    @Transactional
    public void generateInvoiceForAppointment(String appointmentId, String patientId) {
        var newInvoice = Invoice.builder()
                .appointmentId(appointmentId)
                .patientId(patientId)
                .amount(new BigDecimal("150.00"))
                .status(PaymentStatus.PENDING)
                .createdAt(LocalDateTime.now())
                .build();

        invoiceRepository.save(newInvoice);

        log.info("Successfully generated new PENDING invoice for Appointment ID: {} | Patient: {}", appointmentId, patientId);
    }
}
