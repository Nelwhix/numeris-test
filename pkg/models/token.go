package models

import (
	"context"
	"time"
)

type Token struct {
	UserID     string
	Token      string
	LastUsedAt *time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateTokenRequest struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

func (m *Model) InsertIntoTokens(ctx context.Context, request CreateTokenRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := "insert into personal_access_tokens(user_id, token, expires_at) values ($1, $2, $3)"
	_, err := m.Conn.Exec(ctx, sql, request.UserID, request.Token, request.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) FindToken(ctx context.Context, token string) (Token, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var cToken Token
	row := m.Conn.QueryRow(ctx, "select user_id, expires_at, last_used_at, updated_at FROM personal_access_tokens WHERE token = $1 AND expires_at > NOW()", token)
	err := row.Scan(&cToken.UserID, &cToken.ExpiresAt, &cToken.LastUsedAt, &cToken.UpdatedAt)
	if err != nil {
		return Token{}, err
	}

	return cToken, nil
}

func (m *Model) UpdateToken(ctx context.Context, token Token) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := "update personal_access_tokens set last_used_at = $1, updated_at = $2 where token = $3"
	_, err := m.Conn.Exec(ctx, sql, token.LastUsedAt, token.UpdatedAt, token.Token)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) FindValidTokenForUser(ctx context.Context, userID string) (Token, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var cToken Token
	row := m.Conn.QueryRow(ctx, "SELECT token FROM personal_access_tokens WHERE user_id = $1 AND expires_at > NOW()", userID)
	err := row.Scan(&cToken.Token)
	if err != nil {
		return Token{}, err
	}

	return cToken, nil
}
