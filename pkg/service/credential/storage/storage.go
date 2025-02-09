package storage

import (
	"github.com/TBD54566975/ssi-sdk/credential"
	"github.com/tbd54566975/ssi-service/internal/util"
	"github.com/tbd54566975/ssi-service/pkg/storage"
)

type StoredCredential struct {
	ID           string                          `json:"id"`
	Credential   credential.VerifiableCredential `json:"credential"`
	Issuer       string                          `json:"issuer"`
	Subject      string                          `json:"subject"`
	Schema       string                          `json:"schema"`
	IssuanceDate string                          `json:"issuanceDate"`
}

type Storage interface {
	StoreCredential(credential StoredCredential) error
	GetCredential(id string) (*StoredCredential, error)
	GetCredentialsByIssuer(issuer string) ([]StoredCredential, error)
	GetCredentialsBySubject(subject string) ([]StoredCredential, error)
	GetCredentialsBySchema(schema string) ([]StoredCredential, error)
	DeleteCredential(id string) error
}

func NewCredentialStorage(s storage.ServiceStorage) (Storage, error) {
	gotBolt, ok := s.(*storage.BoltDB)
	if !ok {
		return nil, util.LoggingNewError("unsupported storage type")
	}
	boltStorage, err := NewBoltCredentialStorage(gotBolt)
	if err != nil {
		return nil, util.LoggingErrorMsg(err, "could not instantiate credential bolt storage")
	}
	return boltStorage, err
}
