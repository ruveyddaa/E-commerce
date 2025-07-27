package pkg

import "log"

func LogErrorWithCorrelation(err error, correlationID string) {

	log.Printf("[CorrelationID: %s] %s", correlationID, err.Error())

}
func LogInfoWithCorrelation(message string, correlationID string) {
	log.Printf("[CorrelationID: %s] %s", correlationID, message)
}
