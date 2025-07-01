package xxl_job_logger

import (
	"testing"
)

func Test_info(t *testing.T) {
	GetLogHandler().SetLogId(1).Info("Test")
}
