package audit

type Event string

const (

	EventLogin Event = "LOGIN"

	EventLogout Event = "LOGOUT"

	EventRegister Event = "REGISTER"

	EventCreateUser Event = "CREATE_USER"

	EventUpdateUser Event = "UPDATE_USER"

	EventDeleteUser Event = "DELETE_USER"

	EventLDAPSync Event = "LDAP_SYNC"

	EventCreateAssessment Event = "CREATE_ASSESSMENT"

	EventUpdateAssessment Event = "UPDATE_ASSESSMENT"

	EventDeleteAssessment Event = "DELETE_ASSESSMENT"

	EventGenerateReport Event = "GENERATE_REPORT"
)