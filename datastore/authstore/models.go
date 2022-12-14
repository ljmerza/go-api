// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package authstore

import (
	"time"

	"github.com/google/uuid"
)

// The permission table stores an approval of a mode of access to a resource.
type Permission struct {
	// The unique ID for the table.
	PermissionID uuid.UUID
	// Unique External ID to be given to outside callers.
	PermissionExtlID string
	// A human-readable string which represents a resource (e.g. an HTTP route or document, etc.).
	Resource string
	// A string representing the action taken on the resource (e.g. POST, GET, edit, etc.)
	Operation string
	// A description of what the permission is granting, e.g. "grants ability to edit a billing document".
	PermissionDescription string
	// A boolean denoting whether the permission is active (true) or not (false).
	Active bool
	// The application which created this record.
	CreateAppID uuid.UUID
	// The user which created this record.
	CreateUserID uuid.NullUUID
	// The timestamp when this record was created.
	CreateTimestamp time.Time
	// The application which performed the most recent update to this record.
	UpdateAppID uuid.UUID
	// The user which performed the most recent update to this record.
	UpdateUserID uuid.NullUUID
	// The timestamp when the record was updated most recently.
	UpdateTimestamp time.Time
}

// The role table stores a job function or title which defines an authority level.
type Role struct {
	// The unique ID for the table.
	RoleID uuid.UUID
	// Unique External ID to be given to outside callers.
	RoleExtlID string
	// A human-readable code which represents the role.
	RoleCd string
	// A longer description of the role.
	RoleDescription string
	// A boolean denoting whether the role is active (true) or not (false).
	Active bool
	// The application which created this record.
	CreateAppID uuid.UUID
	// The user which created this record.
	CreateUserID uuid.NullUUID
	// The timestamp when this record was created.
	CreateTimestamp time.Time
	// The application which performed the most recent update to this record.
	UpdateAppID uuid.UUID
	// The user which performed the most recent update to this record.
	UpdateUserID uuid.NullUUID
	// The timestamp when the record was updated most recently.
	UpdateTimestamp time.Time
}

// The role_permission table stores which roles have which permissions.
type RolePermission struct {
	// The unique role which can have 1 to many permissions set in this table.
	RoleID uuid.UUID
	// The unique permission that is being given to the role.
	PermissionID uuid.UUID
	// The application which created this record.
	CreateAppID uuid.UUID
	// The user which created this record.
	CreateUserID uuid.NullUUID
	// The timestamp when this record was created.
	CreateTimestamp time.Time
	// The application which performed the most recent update to this record.
	UpdateAppID uuid.UUID
	// The user which performed the most recent update to this record.
	UpdateUserID uuid.NullUUID
	// The timestamp when the record was updated most recently.
	UpdateTimestamp time.Time
}

// The role_user table stores which roles have which users.
type RoleUser struct {
	// The unique role which can have one to many users set in this table.
	RoleID uuid.UUID
	// The unique user that is being given the role.
	UserID uuid.UUID
	// The application which created this record.
	CreateAppID uuid.UUID
	// The user which created this record.
	CreateUserID uuid.NullUUID
	// The timestamp when this record was created.
	CreateTimestamp time.Time
	// The application which performed the most recent update to this record.
	UpdateAppID uuid.UUID
	// The user which performed the most recent update to this record.
	UpdateUserID uuid.NullUUID
	// The timestamp when the record was updated most recently.
	UpdateTimestamp time.Time
}
