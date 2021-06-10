#!bin/sh

# https://cloud.google.com/spanner/docs/emulator#docker
gcloud config configurations create emulator
gcloud config set auth/disable_credentials true
gcloud config set project noted-episode-316407
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/

# Creating an instance
gcloud spanner instances create test-instance --config=emulator --description="Test instance" --nodes=1
gcloud spanner databases create test-database --instance test-instance

# Creating the schema
gcloud spanner databases ddl update test-database --instance=test-instance --ddl="
CREATE TABLE Users (
  Id    BYTES(16) NOT NULL,
  Name  STRING(MAX),
  Money INT64,
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
CREATE TABLE Offer (
  Id      BYTES(16) NOT NULL,
  Price   INT64,
  ItemId  BYTES(16) NOT NULL,

  CONSTRAINT FK_ItemOffer FOREIGN KEY (ItemId) REFERENCES Items(Id),
) PRIMARY KEY (Id)
"
