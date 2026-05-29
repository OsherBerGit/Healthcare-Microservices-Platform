package com.healthcare.billingservice.adapter.out.persistence;

import com.healthcare.billingservice.domain.model.Invoice;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface InvoiceJpaRepository extends JpaRepository<Invoice, UUID> {
}
