package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateMenuService(
	ctx context.Context,
	request menuPresenter.UpdateMenuRequest,
	userID string,
	menuID int,
) (err error) {

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			err = errors.WithStack(constant.ErrUnknownSource)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	q := repository.NewRepository(tx)

	existingMenu, err := q.ReadMenuByIDQuery(ctx, menuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found for update")
			return constant.ErrMenuIdNotFound
		}
		s.log.Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	params := request.ToUpdateParams(userID, menuID)

	// PROTEKSI: Cek apakah menu memiliki child
	childCount, err := q.CountMenuChildren(ctx, menuID)
	if err != nil {
		s.log.Error().Err(err).Int("menu_id", menuID).Msg("error counting menu children")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	// HANYA override field lain yang tidak diset di request
	if request.Name == "" {
		params.Name = existingMenu.Name
	}
	if request.Description == nil {
		params.Description = existingMenu.Description
	}
	if request.URL == nil {
		params.URL = existingMenu.URL
	}
	if request.Icon == nil {
		params.Icon = existingMenu.Icon
	}
	if request.Active == nil {
		params.Active = existingMenu.Active
	}
	if request.Display == nil {
		params.Display = existingMenu.Display
	}

	// TIDAK ADA OVERRIDE UNTUK ParentID LAGI!
	// Biarkan nilai dari ToUpdateParams (yang sudah handle null sebagai root)

	parentChanged := hasParentChanged(existingMenu.ParentID, params.ParentID)
	groupChanged := existingMenu.Group != params.Group
	// PROTEKSI: Jika menu memiliki child, tidak boleh pindah parent
	if childCount > 0 && parentChanged {
		return errors.WithStack(constant.ErrMenuHasChildren)
	}

	// LOGIKA GROUP: Jika menu memiliki parent, group harus mengikuti parent
	if params.ParentID.Valid {
		parentMenu, err := q.ReadMenuByIDQuery(ctx, int(params.ParentID.Int32))
		if err != nil {
			s.log.Error().Err(err).Int32("parent_id", params.ParentID.Int32).Msg("error reading parent menu")
			return errors.WithStack(constant.ErrUnknownSource)
		}
		// Override group dengan group parent
		params.Group = parentMenu.Group

	}

	// LOGIKA SORT: Hitung ulang sort berdasarkan perubahan
	if parentChanged {
		var newSort int

		if params.ParentID.Valid {
			// Menu dipindahkan ke parent baru
			newSort, err = q.ReadNextSortForParentAndGroup(ctx, params.ParentID.Int32, params.Group)
			if err != nil {
				s.log.Error().Err(err).Msg("error getting new sort value for new parent")
				return errors.WithStack(constant.ErrUnknownSource)
			}
		} else {
			// Menu dipindahkan ke root
			newSort, err = q.ReadSortForGroup(ctx, params.Group)
			if err != nil {
				s.log.Error().Err(err).Msg("error getting new sort value for root")
				return errors.WithStack(constant.ErrUnknownSource)
			}
		}

		params.Sort = newSort

	} else if groupChanged && !params.ParentID.Valid {
		// Hanya group yang berubah dan menu adalah root
		newSort, err := q.ReadSortForGroup(ctx, params.Group)
		if err != nil {
			s.log.Error().Err(err).Msg("error getting new sort value for group change")
			return errors.WithStack(constant.ErrUnknownSource)
		}
		params.Sort = newSort
	} else {
		// Tidak ada perubahan yang mempengaruhi sort
		params.Sort = int(existingMenu.Sort)
	}

	// Update menu utama
	err = q.UpdateMenuQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", request).Msg("error to update menu")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	// FITUR BARU: Jika menu memiliki child dan group berubah, update semua child
	if childCount > 0 && groupChanged && !parentChanged {

		err = q.UpdateChildrenGroup(ctx, menuID, params.Group)
		if err != nil {
			s.log.Error().Err(err).
				Int("parent_id", menuID).
				Str("new_group", params.Group).
				Msg("error updating children group")
			return errors.WithStack(constant.ErrUnknownSource)
		}

	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return
}

func hasParentChanged(oldParent, newParent sql.NullInt32) bool {
	if !oldParent.Valid && !newParent.Valid {
		return false
	}

	if oldParent.Valid != newParent.Valid {
		return true
	}

	return oldParent.Int32 != newParent.Int32
}
