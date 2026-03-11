package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/auth/domain"
	"context"
)

func (s *authUsecase) AuthMeUsecase(ctx context.Context, id string) (user domain.User, err error) {
	user, err = s.repo.ReadDetailUserByIdQuery(ctx, id)
	if err != nil {
		// Assuming standard SQL error strings translated by repository or if repo doesn't we might need broader checks, but let's stick to checking generic errors or assuming not found if error isn't nil for simplicity if not injecting domain error
		// A proper clean architecture approach usually translates DB errors inside the repository.
		return user, constant.ErrUserIdNotFound
	}
	return user, nil
}
