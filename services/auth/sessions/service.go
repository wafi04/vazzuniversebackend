package sessions

import "context"

type SessionService struct {
	sessionRepo *SessionRepo
}

func NewSessionsService(sessionRepo *SessionRepo) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (ss *SessionService) Create(ctx context.Context, req *CreateSession) (*SessionsData, error) {
	return ss.sessionRepo.Create(ctx, req)
}

func (ss *SessionService) CheckValidatyUser(ctx context.Context, req *CheckValidaty) (*SessionsData, error) {
	return ss.sessionRepo.CheckValidatyUser(ctx, req)
}

func (ss *SessionService) GetBySessionID(ctx context.Context, sessionID string) (*SessionsData, error) {
	return ss.sessionRepo.GetBySessionID(ctx, sessionID)
}

func (ss *SessionService) GetByUserID(ctx context.Context, userID string) ([]*SessionsData, error) {
	return ss.sessionRepo.GetByUserID(ctx, userID)
}

func (ss *SessionService) InvalidateSession(ctx context.Context, sessionID string) error {
	return ss.sessionRepo.InvalidateSession(ctx, sessionID)
}

func (ss *SessionService) InvalidateAllUserSessions(ctx context.Context, userID string) error {
	return ss.sessionRepo.InvalidateAllUserSessions(ctx, userID)
}

func (ss *SessionService) UpdateLastActivity(ctx context.Context, sessionID string) error {
	return ss.sessionRepo.UpdateLastActivity(ctx, sessionID)
}

func (ss *SessionService) GetByAccessToken(ctx context.Context, accessToken string) (*SessionsData, error) {
	return ss.sessionRepo.GetByAccessToken(ctx, accessToken)
}
