// Package taskconstants defines task alias constants used throughout the
// application for task registration and enqueueing via TaskDefinitionEnqueueByAlias.
//
// Using these constants instead of string literals provides compile-time safety,
// a single source of truth, and makes it easy to find all usages of a given task.
package taskconstants

const (
	// BlindIndexRebuildTaskAlias is the alias for the blind index rebuild task.
	BlindIndexRebuildTaskAlias = "BlindIndexUpdate"

	// CleanUpTaskAlias is the alias for the cleanup task.
	CleanUpTaskAlias = "CleanUpTask"

	// EmailTestTaskAlias is the alias for the email test task.
	EmailTestTaskAlias = "EmailTestTask"

	// EmailToAdminTaskAlias is the alias for the admin notification email task.
	EmailToAdminTaskAlias = "EmailToAdminTask"

	// EmailToAdminOnNewContactFormSubmittedTaskAlias is the alias for the
	// contact form submission admin notification task.
	EmailToAdminOnNewContactFormSubmittedTaskAlias = "email-to-admin-on-new-contact-form-submitted"

	// EmailToAdminOnNewUserRegisteredTaskAlias is the alias for the
	// new user registration admin notification task.
	EmailToAdminOnNewUserRegisteredTaskAlias = "email-to-admin-on-new-user-registered"

	// HelloWorldTaskAlias is the alias for the hello world task.
	HelloWorldTaskAlias = "HelloWorldTask"

	// StatsVisitorEnhanceTaskAlias is the alias for the stats visitor
	// enhancement task.
	StatsVisitorEnhanceTaskAlias = "StatsVisitorEnhanceTask"
)
