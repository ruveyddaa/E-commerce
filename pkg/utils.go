package pkg

import (
	"github.com/google/uuid"
	"log"
)

func LogErrorWithCorrelation(err error, correlationID string) {

	log.Printf("[CorrelationID: %s] %s", correlationID, err.Error())

}
func LogInfoWithCorrelation(message string, correlationID string) {
	log.Printf("[CorrelationID: %s] %s", correlationID, message)
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
