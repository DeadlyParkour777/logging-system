package generator

import (
	"fmt"
	"math/rand"

	"github.com/DeadlyParkour777/logging-system/pkg/logs"
)

func GenerateRandomLog(serviceName string) *logs.SendLogRequest {
	randNum := rand.Intn(100)
	req := &logs.SendLogRequest{ServiceName: serviceName}

	if randNum < 70 {
		req.Level = "INFO"
		req.Message = fmt.Sprintf("User %d successfully logged in.", rand.Intn(1000))
		req.Metadata = map[string]string{"user_id": fmt.Sprintf("%d", rand.Intn(1000))}
	} else if randNum < 90 {
		req.Level = "WARN"
		req.Message = "Warn message"
	} else {
		req.Level = "ERROR"
		req.Message = "Failed to process payment for order #12345"
		req.Metadata = map[string]string{"order_id": "12345", "error_code": "5003"}
	}
	return req
}
