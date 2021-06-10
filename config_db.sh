#!/bin/sh

# Creating the schema
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="DROP TABLE Offer" # TODO : Remove this
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="DROP TABLE Offers"
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="DROP TABLE Items"
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="DROP TABLE Users"

gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
CREATE TABLE Users (
  Id    BYTES(16) NOT NULL,
  Name  STRING(MAX),
  Money INT64,

  CONSTRAINT CS_PositiveMoney CHECK(Money >= 0),
) PRIMARY KEY (Id)
"

gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
CREATE TABLE Items (
  Id          BYTES(16) NOT NULL,
  Description STRING(MAX),

  UserId BYTES(16) NOT NULL,

  CONSTRAINT FK_UserItem FOREIGN KEY (UserId) REFERENCES Users(Id),
) PRIMARY KEY (Id)
"

gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
CREATE TABLE Offers (
  Id      BYTES(16) NOT NULL,
  Price   INT64,
  ItemId  BYTES(16) NOT NULL,

  CONSTRAINT FK_ItemOffer FOREIGN KEY (ItemId) REFERENCES Items(Id),
  CONSTRAINT CS_PositivePrice CHECK(Price >= 0),
) PRIMARY KEY (Id)
"
