package sessions

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/generate"
)

type SessionRepo struct {
	MainDB    *sqlx.DB
	ReplicaDB *sqlx.DB
}

func NewSessionRepo(mainDB *sqlx.DB, replicaDB *sqlx.DB) *SessionRepo {
	return &SessionRepo{
		MainDB:    mainDB,
		ReplicaDB: replicaDB,
	}
}

func (sr *SessionRepo) Create(ctx context.Context, req *CreateSession) (*SessionsData, error) {
	sessionID := generate.GenerateRandomID(&generate.IDOpts{
		Amount: 10,
	})

	var session SessionsData
	err := sr.MainDB.QueryRowContext(ctx, QueryInsert, sessionID, req.UserID, req.AccessToken, time.Now(), time.Now()).Scan(
		&session.SessionID,
		&session.UserID,
		&session.AccessToken,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (sr *SessionRepo) CheckValidatyUser(ctx context.Context, req *CheckValidaty) (*SessionsData, error) {
	var session SessionsData
	err := sr.ReplicaDB.QueryRowContext(ctx, QueryCheckValidity, req.SessionID, req.UserID).Scan(
		&session.SessionID,
		&session.UserID,
		&session.AccessToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.DeviceInfo,
		&session.LastActivity,
		&session.IsAccess,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (sr *SessionRepo) GetBySessionID(ctx context.Context, sessionID string) (*SessionsData, error) {
	var session SessionsData
	err := sr.ReplicaDB.QueryRowContext(ctx, QueryGetBySessionID, sessionID).Scan(
		&session.SessionID,
		&session.UserID,
		&session.AccessToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.DeviceInfo,
		&session.LastActivity,
		&session.IsAccess,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (sr *SessionRepo) GetByUserID(ctx context.Context, userID string) ([]*SessionsData, error) {
	rows, err := sr.ReplicaDB.QueryContext(ctx, QueryGetByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*SessionsData
	for rows.Next() {
		var session SessionsData
		err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&session.AccessToken,
			&session.IPAddress,
			&session.UserAgent,
			&session.DeviceInfo,
			&session.LastActivity,
			&session.IsAccess,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (sr *SessionRepo) InvalidateSession(ctx context.Context, sessionID string) error {
	_, err := sr.MainDB.ExecContext(ctx, QueryInvalidateSession, sessionID)
	return err
}

func (sr *SessionRepo) InvalidateAllUserSessions(ctx context.Context, userID string) error {

	_, err := sr.MainDB.ExecContext(ctx, QueryInvalidateAllUserSessions, userID)
	return err
}

func (sr *SessionRepo) UpdateLastActivity(ctx context.Context, sessionID string) error {
	now := time.Now()
	_, err := sr.MainDB.ExecContext(ctx, QueryUpdateLastActivity, now, sessionID)
	return err
}

func (sr *SessionRepo) GetByAccessToken(ctx context.Context, accessToken string) (*SessionsData, error) {
	var session SessionsData
	err := sr.ReplicaDB.QueryRowContext(ctx, QueryGetByAccessToken, accessToken).Scan(
		&session.SessionID,
		&session.UserID,
		&session.AccessToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.DeviceInfo,
		&session.LastActivity,
		&session.IsAccess,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}
