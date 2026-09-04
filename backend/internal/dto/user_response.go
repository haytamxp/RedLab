package dto

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	LDAPUser  bool   `json:"ldap_user"`
}

type DirectoryUserResponse struct {
	DN                string   `json:"dn"`
	Username          string   `json:"username"`
	UserPrincipalName string   `json:"user_principal_name"`
	Email             string   `json:"email"`
	FirstName         string   `json:"first_name"`
	LastName          string   `json:"last_name"`
	DisplayName       string   `json:"display_name"`
	Enabled           bool     `json:"enabled"`
	Groups            []string `json:"groups"`
}
