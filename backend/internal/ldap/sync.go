package ldap

import (
	"fmt"
	"strings"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type SyncUser struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	IsActive  bool
	LDAPUser  bool
	Role      models.Role
}

func ConvertUser(user *User) *SyncUser {
	return &SyncUser{
		Username:  user.SAMAccountName,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.Enabled,
		LDAPUser:  true,
		Role:      MapGroupsToRole(user.Groups),
	}
}

func MapGroupsToRole(groups []string) models.Role {
	for _, group := range groups {
		group = strings.ToLower(group)

		switch {
		case strings.Contains(group, "redlab-admin"):
			return models.RoleAdmin

		case strings.Contains(group, "redlab-trainer"):
			return models.RoleTrainer

		case strings.Contains(group, "redlab-student"):
			return models.RoleStudent

		case strings.Contains(group, "redlab-viewer"):
			return models.RoleViewer
		}
	}

	return models.RoleViewer
}

func ValidateUser(user *User) error {
	if user == nil {
		return fmt.Errorf("LDAP user is nil")
	}

	if user.SAMAccountName == "" {
		return fmt.Errorf("LDAP user has no sAMAccountName")
	}

	if user.DN == "" {
		return fmt.Errorf("LDAP user has no distinguishedName")
	}

	return nil
}