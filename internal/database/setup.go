package database

import (
	"fmt"
)

// SetupDatabase creates the necessary database tables and triggers.
func SetupDatabase() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// SQL for creating the users table
	usersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	// SQL for creating the tweets table
	tweetsTableSQL := `
	CREATE TABLE IF NOT EXISTS tweets (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		content VARCHAR(280) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	// SQL for the trigger function to update the updated_at timestamp
	updateTriggerSQL := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ language 'plpgsql';
	`

	// SQL to apply the trigger to the users table
	usersTriggerSQL := `
	DROP TRIGGER IF EXISTS update_users_updated_at ON users;
	CREATE TRIGGER update_users_updated_at
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE PROCEDURE update_updated_at_column();
	`

	// SQL to apply the trigger to the tweets table
	tweetsTriggerSQL := `
	DROP TRIGGER IF EXISTS update_tweets_updated_at ON tweets;
	CREATE TRIGGER update_tweets_updated_at
	BEFORE UPDATE ON tweets
	FOR EACH ROW
	EXECUTE PROCEDURE update_updated_at_column();
	`

	sqlCommands := []string{
		usersTableSQL,
		tweetsTableSQL,
		updateTriggerSQL,
		usersTriggerSQL,
		tweetsTriggerSQL,
	}

	// Execute all SQL commands
	for _, cmd := range sqlCommands {
		_, err := DB.Exec(cmd)
		if err != nil {
			return fmt.Errorf("error executing command: %v\n%s", err, cmd)
		}
	}

	fmt.Println("Database tables and triggers created successfully.")
	return nil
}
