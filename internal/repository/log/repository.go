package log

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/AndreiMartynenko/auth/internal/model"
	"github.com/AndreiMartynenko/auth/internal/repository"
	"github.com/AndreiMartynenko/common/pkg/db"
)

const (
	tableName = "transaction_log"

	idColumn        = "id"
	timestampColumn = "timestamp"
	logColumn       = "log"
)

type repo struct {
	db db.Client
}

// NewRepository creates new object of repository layer.
func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Log(ctx context.Context, text *model.Log) error {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(logColumn).
		Values(text.Text).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "log_repository.Log",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
