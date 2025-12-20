package migrations

import "pintuniv-go/internal/database"

func CreateTryoutIndexes() {
	// index untuk leaderboard / resume / limit
	database.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tryout_attempt_user
		ON tryout_attempts(user_id, tryout_id, finished_at)
	`)

	// unique constraint biar draft jawaban tidak dobel
	database.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS uniq_tryout_answer
		ON tryout_answers(attempt_id, question_id)
	`)
}
