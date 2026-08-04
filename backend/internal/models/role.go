package models

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleStudent Role = "student"
	RoleTrainer Role = "trainer"
	RoleViewer  Role = "viewer"
)