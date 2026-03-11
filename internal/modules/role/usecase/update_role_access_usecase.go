package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/pkg/errors"
)

func (s *roleUsecase) UpdateRoleAccessUsecase(ctx context.Context, roleID int, requestPayloads []domain.UpdateRoleMenuPermission) (err error) {

	_, err = s.roleRepo.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	// Rely on the repository's native transaction to apply all updates safely.
	err = s.roleRepo.UpdateRoleAccessTx(ctx, roleID, requestPayloads)
	if err != nil {
		if errors.Is(err, constant.ErrMenuPermissionIdNotFound) {
			s.log.Warn().Err(err).Msg("validation error on update role access")
			return err
		}
		s.log.Error().Err(err).Msg("error updating role access via transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
