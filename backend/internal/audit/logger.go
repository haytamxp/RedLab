package audit

import (
	"log"
	"time"
)

func Log(
	event Event,
	user string,
	resource string,
	action string,
	ip string,
) {

	record := Record{

		Event: event,

		User: user,

		Resource: resource,

		Action: action,

		IPAddress: ip,

		Time: time.Now(),
	}

	log.Printf(
		"[AUDIT] %s | %s | %s | %s | %s",
		record.Time.Format(time.RFC3339),
		record.Event,
		record.User,
		record.Resource,
		record.Action,
	)
}