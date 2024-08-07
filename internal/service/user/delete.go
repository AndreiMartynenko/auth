package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndreiMartynenko/auth/internal/model"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Deleted user with id: %d", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}
