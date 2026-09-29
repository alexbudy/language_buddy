package store

import (
	"context"
	"fmt"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
)

// Profile DB model and queries

type Profile struct {
    ID                 int64
    Name               string
    EnableSpeech       bool
    DefaultNumQuestions int
    DefaultNumAnswers   int
    Lang1              string
    Lang2              string
}

// GetProfiles returns all profiles.
func (s *Store) GetProfiles(ctx context.Context) ([]Profile, error) {
	var profiles []Profile
	getProfilesQry := `SELECT id, name, enable_speech, default_num_questions, 
						default_num_answers, lang1, lang2 FROM profiles 
						WHERE deleted_at IS NULL ORDER BY created_at ASC`
	rows, err := s.db.QueryContext(ctx, getProfilesQry)
	if err != nil {
		return nil, fmt.Errorf("store: get profiles: %w", err)
	}
	defer rows.Close()
	log.DebugSQL("Successfully queried %s", getProfilesQry)

	for rows.Next() {
		var p Profile
		if err := rows.Scan(&p.ID, &p.Name, &p.EnableSpeech, &p.DefaultNumQuestions, &p.DefaultNumAnswers, &p.Lang1, &p.Lang2); err != nil {
			return nil, fmt.Errorf("store: get profiles: %w", err)
		}
		profiles = append(profiles, p)
	}
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("store: get profiles: rows: %w", err)
    }

	return profiles, nil
}