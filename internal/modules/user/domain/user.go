package domain

import (
	"context"

	"be-dashboard-nba/api/presenter"
)

// User adalah entitas utama untuk domain User.
type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	FullName  string  `json:"full_name"`
	Role      string  `json:"role"`
	RoleID    int     `json:"role_id"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Active    bool    `json:"active"`
	Phone     *string `json:"phone"`
	ImgPath   *string `json:"img_path"`
	ImgName   *string `json:"img_name"`
	CreatedBy string  `json:"created_by"`
	CreatedAt string  `json:"created_at"`
	UpdatedBy *string `json:"updated_by"`
	UpdatedAt *string `json:"updated_at"`
}

// UserListRow adalah Read Model khusus untuk query list user.
type UserListRow struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	RoleID   int    `json:"role_id"`
	Active   bool   `json:"active"`
}

// UserDetailRow adalah Read Model untuk melihat detail spesifik user (tanpa password/timestamp).
type UserDetailRow struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Role     string  `json:"role"`
	RoleID   int     `json:"role_id"`
	Active   bool    `json:"active"`
}

// UserProfileRow adalah Read Model khusus untuk profil user yang memiliki info Avatar/Gambar.
type UserProfileRow struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Role     string  `json:"role"`
	RoleID   int     `json:"role_id"`
	Active   bool    `json:"active"`
	ImgPath  *string `json:"img_path"`
	ImgName  *string `json:"img_name"`
}

// UserPaginationResponse adalah format balikan untuk list data.
type UserPaginationResponse struct {
	Data       []UserListRow
	Pagination presenter.Pagination
}

// ==========================================
// KONTRAK INTERFACE UNTUK CLEAN ARCHITECTURE
// ==========================================

type UserFilter struct {
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
	Page      int
}

type UpdateUserPayload struct {
	ID        string
	Name      string
	FullName  string
	Email     string
	RoleID    int
	Phone     *string
	Active    *bool
	UpdatedBy *string
}

// UserRepository mendefinisikan operasi-operasi database (CRUD) yang bisa dilakukan pada entitas User.
type UserRepository interface {
	CreateUserWithRoleTx(ctx context.Context, user User) (err error)
	UpdateUserWithRoleTx(ctx context.Context, user User) (err error)
	DeleteUserQuery(ctx context.Context, id string) (err error)

	ReadUsersQuery(ctx context.Context, filter UserFilter) (data []UserListRow, err error)
	ReadCountUserQuery(ctx context.Context, filter UserFilter) (count int, err error)
	ReadUserByIDQuery(ctx context.Context, id string) (data UserDetailRow, err error)
	ReadUserProfileQuery(ctx context.Context, id string) (data UserProfileRow, err error)
	IsUserEmailExistsQuery(ctx context.Context, email string) (exists bool, err error)
	IsUpdateUserEmailExistsQuery(ctx context.Context, email, id string) (exists bool, err error)
}

// UserUsecase mendefinisikan logika bisnis (service) apa saja yang bisa dilakukan pada entitas User.
type UserUsecase interface {
	ReadDetailUserUsecase(ctx context.Context, id string) (data UserDetailRow, err error)
	UpdateUserUsecase(ctx context.Context, payload UpdateUserPayload) (err error)
	ReadUsersUsecase(ctx context.Context, filter UserFilter) (data UserPaginationResponse, err error)
	CreateUserUsecase(ctx context.Context, user User) (err error)
	DeleteUserUsecase(ctx context.Context, userID string, deletedBy string) (err error)
	ReadUserProfileUsecase(ctx context.Context, userID string) (data UserProfileRow, err error)
}
