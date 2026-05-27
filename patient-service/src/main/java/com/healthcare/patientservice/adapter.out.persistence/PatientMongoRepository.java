package com.healthcare.patientservice.adapter.out.persistence;

import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.UUID;

public interface PatientMongoRepository extends MongoRepository<PatientDocument, UUID> {
}
