package main

import (
	"github.com/xtdlib/log"
)

func main() {
	for i := 0; i < 10; i++ {
		log.Info("Hello, World!" + string(i+'0'))
	}

	// Ensure all logs are sent before exiting

	// To verify logs in Victoria Logs:
	// curl http://oci-aca-001:9428/select/logsql/query -d 'query=*'
}
