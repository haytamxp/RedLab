package permissions

var RolePermissions = map[string][]string{

	RoleAdministrator: {

		CreateUser,

		DeleteUser,

		UpdateUser,

		ViewUsers,

		CreateAssessment,

		ViewAssessment,

		GenerateReport,

		LDAPSync,

		ManageAgents,
	},

	RoleManager: {

		ViewUsers,

		CreateAssessment,

		ViewAssessment,

		GenerateReport,
	},

	RoleAnalyst: {

		CreateAssessment,

		ViewAssessment,

		GenerateReport,
	},

	RoleStudent: {

		ViewAssessment,
	},

	RoleViewer: {

		ViewAssessment,
	},
}

func HasPermission(role string, permission string) bool {

	permissions, ok := RolePermissions[role]

	if !ok {

		return false
	}

	for _, p := range permissions {

		if p == permission {

			return true
		}
	}

	return false
}