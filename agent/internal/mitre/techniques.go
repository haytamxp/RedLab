package mitre

type Technique struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Techniques = map[string]Technique{
	"T1482": {
		ID:   "T1482",
		Name: "Domain Trust Discovery",
	},
	"T1087.002": {
		ID:   "T1087.002",
		Name: "Domain Account Discovery",
	},
	"T1069.002": {
		ID:   "T1069.002",
		Name: "Permission Groups Discovery: Domain Groups",
	},
}
