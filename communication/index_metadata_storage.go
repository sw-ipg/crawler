package communication

import (
	"context"
	"crawler/app/model"
	"crawler/communication/sqlqueries"
	"time"

	"github.com/jmoiron/sqlx"
)

type IndexMetadataStorage struct {
	db *sqlx.DB
}

func NewIndexMetadataStorage(db *sqlx.DB) IndexMetadataStorage {
	return IndexMetadataStorage{
		db: db,
	}
}

func (s IndexMetadataStorage) GetLastDateOfDoc(ctx context.Context, path string) (t time.Time, err error) {
	err = s.db.GetContext(ctx, &t, sqlqueries.GetLastDateOfDocSql, path)
	return t, err
}

func (s IndexMetadataStorage) AddDoc(ctx context.Context, doc model.IndexDoc) error {
	_, err := s.db.NamedExecContext(ctx, sqlqueries.StoreIndexDocSql, doc)
	return err
}

func (s IndexMetadataStorage) GetChecksumsForPath(ctx context.Context, path string) (c []uint32, err error) {
	err = s.db.SelectContext(ctx, &c, sqlqueries.GetChecksumsForPathSql, path)
	return c, err
}
