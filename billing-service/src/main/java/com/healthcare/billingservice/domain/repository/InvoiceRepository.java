package com.healthcare.billingservice.domain.repository;

import com.healthcare.billingservice.domain.model.Invoice;

import java.util.Optional;
import java.util.UUID;

public interface InvoiceRepository {
    Invoice save(Invoice invoice);
    Optional<Invoice> findById(UUID id);
}
