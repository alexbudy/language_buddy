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

// StoreProfile creates a new profile with the given name.
func (s *Store) StoreProfile(ctx context.Context, profName string) (Profile, error) {
	storeProfileQry := `INSERT INTO profiles (name) VALUES (?)`
	var profile Profile
	_, err := s.db.ExecContext(ctx, storeProfileQry, profName)
	if err != nil {
		return Profile{}, err
	}
	log.DebugSQL("Successfully inserted %s for profile %s", storeProfileQry, profName)

	profile, err = s.GetProfile(ctx, profName)
	if err != nil {
		return Profile{}, err
	}
	return profile, nil
}

// GetProfile returns the profile with the given name.
func (s *Store) GetProfile(ctx context.Context, profName string) (Profile, error) {
	var profile Profile
	getProfileQry := `SELECT id, name, enable_speech, default_num_questions, 
						default_num_answers, lang1, lang2 FROM profiles 
						WHERE deleted_at IS NULL AND name = ?`
	row := s.db.QueryRowContext(ctx, getProfileQry, profName)
	if err := row.Scan(&profile.ID, &profile.Name, &profile.EnableSpeech, &profile.DefaultNumQuestions, &profile.DefaultNumAnswers, &profile.Lang1, &profile.Lang2); err != nil {
		return Profile{}, fmt.Errorf("store: get profile: %w", err)
	}
	log.DebugSQL("Successfully retrieved %s for profile %s", getProfileQry, profName)
	return profile, nil

}

// DeleteProfile deletes the profile with the given name.
func (s *Store) DeleteProfile(ctx context.Context, profName string) error {
	deleteProfileQry := `DELETE FROM profiles WHERE name = ?`
	_, err := s.db.ExecContext(ctx, deleteProfileQry, profName)
	if err != nil {
		return fmt.Errorf("store: delete profile: %w", err)
	}
	log.DebugSQL("Successfully deleted %s for profile %s", deleteProfileQry, profName)
	return nil
}

// RenameProfile changes a profile's name and returns the renamed profile.
func (s *Store) RenameProfile(ctx context.Context, oldProfName string, newProfName string) (Profile, error) {
	renameProfileQry := `UPDATE profiles SET name = ? WHERE name = ?`
	_, err := s.db.ExecContext(ctx, renameProfileQry, newProfName, oldProfName)
	if err != nil {
		return Profile{}, fmt.Errorf("store: rename profile: %w", err)
	}
	log.DebugSQL("Successfully renamed %s for profile %s to %s", renameProfileQry, oldProfName, newProfName)
	
	return s.GetProfile(ctx, newProfName)
}

// UpdateProfile updates the settings for the profile with the given name.
func (s *Store) UpdateProfile(ctx context.Context, profile Profile) error {
	updateProfileQry := `UPDATE profiles SET enable_speech = ?, default_num_questions = ?, 
						default_num_answers = ?, lang1 = ?, lang2 = ? WHERE name = ?`
	_, err := s.db.ExecContext(ctx, updateProfileQry, profile.EnableSpeech, profile.DefaultNumQuestions, profile.DefaultNumAnswers, profile.Lang1, profile.Lang2, profile.Name)
	if err != nil {
		return fmt.Errorf("store: update profile: %w", err)
	}
	log.DebugSQL("Successfully updated %s for profile %s", updateProfileQry, profile.Name)
	return nil	
}