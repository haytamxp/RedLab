package audit

import "time"

type Record struct {

	Event Event `json:"event"`

	User string `json:"user"`

	Resource string `json:"resource"`

	Action string `json:"action"`

	Time time.Time `json:"time"`

	IPAddress string `json:"ip_address"`
}