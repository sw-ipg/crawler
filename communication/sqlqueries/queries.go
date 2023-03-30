package sqlqueries

import _ "embed"

//go:embed get_last_date_of_doc.sql
var GetLastDateOfDocSql string

//go:embed store_index_doc.sql
var StoreIndexDocSql string

//go:embed get_checksums_for_path.sql
var GetChecksumsForPathSql string
