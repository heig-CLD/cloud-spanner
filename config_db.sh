#!/bin/bash

if [ "$1" == "--help" ]; then
  echo "Usage :"
  echo ""
  echo "--no-checks : Runs the script without the CHECK SQL constraints."
  echo "--help      : Displays the help page. This is how you got there !"
  exit 0
fi

# Creating the schema

# Please run this if you still use the old schema.
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="DROP TABLE Transfers"
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="DROP TABLE Users"

if [ "$1" == "--no-checks" ]; then # Run the commands without checks.
  gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
  CREATE TABLE Users (
    Id    BYTES(16)   NOT NULL,
    Name  STRING(MAX) NOT NULL,
    Money INT64       NOT NUL,
  ) PRIMARY KEY (Id)
  "
  gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
  CREATE TABLE Transfers (
    Id          BYTES(16) NOT NULL,
    Amount      INT64     NOT NULL,
    FromUserId  BYTES(16) NOT NULL,
    ToUserId    BYTES(16) NOT NULL,
    AtTimestamp TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    CONSTRAINT FK_FromUser FOREIGN KEY (FromUserId) REFERENCES Users(Id),
    CONSTRAINT FK_ToUser   FOREIGN KEY (ToUserId)   REFERENCES Users(Id),
  ) PRIMARY KEY (Id)
  "
# Run the commands with the checks (not on emulator).
else
  gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
  CREATE TABLE Users (
    Id    BYTES(16)   NOT NULL,
    Name  STRING(MAX) NOT NULL,
    Money INT64       NOT NULL,
    CONSTRAINT CS_PositiveMoney CHECK(Money >= 0),
  ) PRIMARY KEY (Id)
  "
  gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
  CREATE TABLE Transfers (
    Id          BYTES(16) NOT NULL,
    Amount      INT64     NOT NULL,
    FromUserId  BYTES(16) NOT NULL,
    ToUserId    BYTES(16) NOT NULL,
    AtTimestamp TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    CONSTRAINT FK_FromUser       FOREIGN KEY (FromUserId) REFERENCES Users(Id),
    CONSTRAINT FK_ToUser         FOREIGN KEY (ToUserId)   REFERENCES Users(Id),
    CONSTRAINT CS_PositiveAmount CHECK(Amount >= 0),
  ) PRIMARY KEY (Id)
  "
fi
