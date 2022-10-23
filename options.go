package entgqlplus

import "os"

type (
	ExtensionOption = func(*Extension)
	Database        = string
)

var (
	SQLite     Database = "sqlite"
	MySQL      Database = "mysql"
	PostgreSQL Database = "postgres"
)

// WithMutation(b bool) enables entgqlplus to generate the Mutations.
// Default value is WithMutation(false).
func WithMutation(b bool) ExtensionOption {
	return func(e *Extension) {
		e.config.Mutation = b
	}
}

// WithSubscription(b bool) enables entgqlplus to generate the Subscriptions.
// Works only if WithMutation(true) is enabled.
// Default value is WithSubscription(false).
func WithSubscription(b bool) ExtensionOption {
	return func(e *Extension) {
		if e.config.Mutation {
			e.config.Subscription = b
		}
	}
}

// WithEchoServer(b bool) enables entgqlplus to generate the server, routes and the handlers.
// Default value is WithEchoServer(false).
func WithEchoServer(b bool) ExtensionOption {
	return func(e *Extension) {
		if b && !e.config.Echo {
			e.config.Echo = b
		}
	}
}

// WithJWTAuth(b bool) enables entgqlplus to generate the login route and the Protected middleware
// Works only if WithEcho(true) is enabled.
// Default value is WithJWTAuth(false).
func WithJWTAuth(b bool) ExtensionOption {
	return func(e *Extension) {
		if e.config.Echo {
			e.config.JWT = b
		}
	}
}

// WithDatabase(b Database) enables entgqlplus to generate the necessary code to connect to the database and migration.
// Default value is WithDatabase(entgql.SQLite).
func WithDatabase(d Database) ExtensionOption {
	return func(e *Extension) {
		e.config.Database = d
	}
}

// WithConfigPath(p string) enables entgqlplus locate the gqlgen.yml config file.
// Default value is With WithConfigPath("../gqlgen.yml").
func WithConfigPath(p string) ExtensionOption {
	return func(e *Extension) {
		_, err := os.Stat(p)
		catch(err)
		e.config.GqlGenPath = p
	}
}

// WithFileUpload(b bool) adds upload mutation.
// this only works if WithMutation(true) is enabled.
// Default is WithFileUpload(false).
func WithFileUpload(b bool) ExtensionOption {
	return func(e *Extension) {
		if e.config.Mutation {
			e.config.FileUpload = b
		}
	}
}
