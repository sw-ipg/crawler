package model

import "time"

type IndexDoc struct {
	Id              int       `db:"id"`
	Domain          string    `db:"domain"`
	Path            string    `db:"path"`
	CRC32Checksum   uint32    `db:"crc32_checksum"`
	Date            time.Time `db:"date"`
	StorageFileName string    `db:"storage_file_name"`
}
