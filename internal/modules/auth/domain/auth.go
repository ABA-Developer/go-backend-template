package domain

import (
	"context"
	"database/sql"
	"time"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/jwt"
)

type Session struct {
	ID                    string     `json:"id"`
	UserID                string     `json:"user_id"`
	AccessToken           string     `json:"access_token"`
	AccessTokenExpiredAt  time.Time  `json:"access_token_expired_at"`
	RefreshToken          string     `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time  `json:"refresh_token_expired_at"`
	IPAddress             string     `json:"ip_address"`
	UserAgent             string     `json:"user_agent"`
	CreatedAt             *time.Time `json:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at"`
}

type LoginAttemp struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	IPAddress string     `json:"ip_address"`
	CreatedAt *time.Time `json:"created_at"`
}

type LoginRecord struct {
	ID          int        `json:"id"`
	UserID      string     `json:"user_id"`
	AccessToken string     `json:"access_token"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
	Browser     string     `json:"browser"`
	OS          string     `json:"os"`
	IPAddress   string     `json:"ip_address"`
	Action      string     `json:"action"` // "login" or "logout"
	CreatedAt   *time.Time `json:"created_at"`
}

type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	FullName  string     `json:"full_name"`
	Role      string     `json:"role"`
	RoleID    int        `json:"role_id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Active    bool       `json:"active"`
	Phone     *string    `json:"phone"`
	ImgPath   *string    `json:"img_path"`
	ImgName   *string    `json:"img_name"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *string    `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type AuthRepository interface {
	WithTx(tx *sql.Tx) AuthRepository
	CreateSessionQuery(ctx context.Context, args Session) (err error)
	UpdateSessionQuery(ctx context.Context, args Session) (err error)
	DeleteSessionQuery(ctx context.Context, id string) (err error)
	ReadDetailSessionQuery(ctx context.Context, id string) (data Session, err error)
	ReadDetailUserByEmailQuery(ctx context.Context, email string) (data User, err error)
	CreateLoginAttemp(ctx context.Context, args LoginAttemp) (err error)
	CreateLoginRecord(ctx context.Context, args LoginRecord) (err error)
	CheckPermissionQuery(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error)
	ReadDetailUserByIdQuery(ctx context.Context, id string) (data User, err error)
}

type AuthUsecase interface {
	CheckPermissionUsecase(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error)
	AuthMeUsecase(ctx context.Context, id string) (data User, err error)
	LogoutUsecase(ctx context.Context, claims *jwt.AccessTokenPayload, iPAddress string) (err error)
	LoginUsecase(ctx context.Context, email, password, userAgent, iPAddress string) (data Session, user User, err error)
}
