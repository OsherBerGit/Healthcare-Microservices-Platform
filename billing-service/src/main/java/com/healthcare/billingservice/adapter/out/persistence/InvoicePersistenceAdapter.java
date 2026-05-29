package com.healthcare.billingservice.adapter.out.persistence;

import com.healthcare.billingservice.domain.model.Invoice;
import com.healthcare.billingservice.domain.repository.InvoiceRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
@RequiredArgsConstructor
public class InvoicePersistenceAdapter implements InvoiceRepository {

    private final InvoiceJpaRepository jpaRepository;

    @Override
    public Invoice save(Invoice invoice) {
        return jpaRepository.save(invoice);
    }

    @Override
    public Optional<Invoice> findById(UUID id) {
        return jpaRepository.findById(id);
    }
}
