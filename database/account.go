package database

import (
	"login-server/api/login"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"time"
)

type Account struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PremDays int    `json:"premdays"`
	LastDay  int    `json:"lastday"`
}

func (acc *Account) Authenticate(db *sql.DB) error {
	h := sha1.New()
	h.Write([]byte(acc.Password))

	p := h.Sum(nil)

	statement := fmt.Sprintf(
		"SELECT id, premdays, lastday FROM accounts WHERE email = '%s' AND password = '%x'",
		acc.Email,
		p,
	)

	return db.QueryRow(statement).Scan(&acc.ID, &acc.PremDays, &acc.LastDay)
}

func (acc *Account) GetSession() login.Session {
	return login.Session{
		IsPremium:      acc.PremDays > 0,
		PremiumUntil:   acc.getPremiumTime(),
		SessionKey:     fmt.Sprintf("%s\n%s", acc.Email, acc.Password),
		ShowRewardNews: true,
		Status:         "active",
	}
}

func (acc *Account) getPremiumTime() int {
	if acc.PremDays > 0 {
		return int(time.Now().UnixNano()/1e6) + acc.PremDays*86400
	}
	return 0
}
