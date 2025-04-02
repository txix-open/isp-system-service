package repository

const (
	pgUniqueViolationErrorCode = "23505"
	pgFkViolationErrorCode     = "23503"

	applicationPkConstrainName          = "application_pkey"
	applicationUniqueNameConstrainName  = "uq_name_application_group_id"
	applicationFkAppGroupConstraintName = "fk_application_group_id"
)
