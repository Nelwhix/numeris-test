package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/thanhpk/randstr"
	"hash/crc32"
	"time"
)

func generateTokenString() string {
	tokenEntropy := randstr.String(40)
	crc32bHash := crc32.ChecksumIEEE([]byte(tokenEntropy))

	return fmt.Sprintf("%s%x", tokenEntropy, crc32bHash)
}

func GetOrCreateToken(m models.Model, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	cToken, err := m.FindValidTokenForUser(ctx, userID)
	if err == nil {
		return cToken.Token, nil
	}

	expires := time.Now().Add(24 * time.Hour * 7)
	tokenString := generateTokenString()
	request := models.CreateTokenRequest{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: expires,
	}

	err = m.InsertIntoTokens(ctx, request)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckTokenValidity(ctx context.Context, m models.Model, tokenString string) error {
	cToken, err := m.FindToken(ctx, tokenString)
	if err != nil {
		return err
	}

	if time.Now().After(cToken.ExpiresAt) {
		return errors.New("expired token")
	}

	lastUsed := time.Now()
	cToken.LastUsedAt = &lastUsed
	cToken.UpdatedAt = time.Now()
	err = m.UpdateToken(ctx, cToken)
	if err != nil {
		return err
	}

	return nil
}
