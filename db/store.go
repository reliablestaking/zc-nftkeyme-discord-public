package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type (
	// Store struct to store Db
	Store struct {
		Db *sqlx.DB
	}

	// DiscordUser struct to store
	DiscordUser struct {
		ID                   int            `db:"id"`
		DiscordUserID        string         `db:"discord_user_id"`
		DiscordUsername      string         `db:"discord_username"`
		DiscordEmail         string         `db:"discord_email"`
		NftkeymeID           sql.NullString `db:"nftkeyme_id"`
		NftkeymeEmail        sql.NullString `db:"nftkeyme_email"`
		NftkeymeAccessToken  sql.NullString `db:"nftkeyme_access_token"`
		NftkeymeRefreshToken sql.NullString `db:"nftkeyme_refresh_token"`
		NumAssets            sql.NullInt64  `db:"num_assets"`
	}
)

// GetUserByDiscordID Gets a user using their discord id
func (s Store) GetUserByDiscordID(discordUserID string) (*DiscordUser, error) {
	discordUser := DiscordUser{}
	err := s.Db.Get(&discordUser, "SELECT * FROM discord_user where discord_user_id = $1", discordUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &discordUser, nil
}

// GetAllDiscordUsers Gets all users
func (s Store) GetAllDiscordUsers() ([]DiscordUser, error) {
	discordUsers := []DiscordUser{}
	err := s.Db.Select(&discordUsers, "SELECT * FROM discord_user")
	if err != nil {
		return nil, err
	}

	return discordUsers, nil
}

// InsertDiscordUser inserts a new user into the db
func (s Store) InsertDiscordUser(discordUserID, discordUsername, discordEmail string) error {
	insertUserQuery := `INSERT INTO discord_user (discord_user_id,discord_username,discord_email) VALUES($1, $2, $3)`

	rows, err := s.Db.Query(insertUserQuery, discordUserID, discordUsername, discordEmail)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// UpdateDiscordUser updates a new user in the db
func (s Store) UpdateDiscordUser(discordUserID, nftkeymeAccessToken, nftkeymeRefreshToken string) error {
	insertUserQuery := `UPDATE discord_user SET nftkeyme_access_token = $1, nftkeyme_refresh_token = $2 WHERE discord_user_id = $3`

	rows, err := s.Db.Query(insertUserQuery, nftkeymeAccessToken, nftkeymeRefreshToken, discordUserID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// UpdateDiscordUserNftkeyInfo updates nftkey me user info
func (s Store) UpdateDiscordUserNftkeyInfo(discordUserID, nftkeymeID, nftkeymeEmail string) error {
	insertUserQuery := `UPDATE discord_user SET nftkeyme_id = $1, nftkeyme_email = $2 WHERE discord_user_id = $3`

	rows, err := s.Db.Query(insertUserQuery, nftkeymeID, nftkeymeEmail, discordUserID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// UpdateDiscordUserNumAssets updates user with new asset count
func (s Store) UpdateDiscordUserNumAssets(discordUserID string, numAssets int) error {
	insertUserQuery := `UPDATE discord_user SET num_assets = $1 WHERE discord_user_id = $2`

	rows, err := s.Db.Query(insertUserQuery, numAssets, discordUserID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
