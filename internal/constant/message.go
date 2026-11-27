package constant

const (
	MsgInternalError    = "Terjadi kesalahan di server."
	MsgServerInitFailed = "Server gagal dijalankan."
	MsgRouteNotFound    = "Endpoint tidak ditemukan."
	MsgMethodNotAllowed = "Method tidak didukung untuk endpoint ini."
	MsgSongInvalid      = "Data lagu tidak valid."
	MsgSlugTakenClient  = "Slug sudah dipakai lagu lain."
	MsgSlugTaken        = "slug sudah dipakai"

	MsgFieldRequired = "wajib diisi"
	MsgSlugFormat    = "hanya huruf kecil, angka, dan tanda hubung (contoh: izayoi-lionheart)"
	MsgInvalidStatus = "harus draft atau published"
	MsgMaxLenFmt     = "maksimal %d karakter"

	MsgSchemaVersionRequired = "schemaVersion wajib diisi"
	MsgSpanMin               = "span minimal 1"
	MsgSpanMismatchFmt       = "total span %d tidak sama dengan jumlah baris %d"
	MsgRubyUnbalanced        = "kurung ruby { } tidak seimbang"
	MsgFootnoteMissingFmt    = "penanda [^%d] dipakai tapi tidak ada di footnotes"
	MsgFootnoteUnusedFmt     = "catatan kaki %d ada di footnotes tapi tidak dipakai"

	MsgQueryInvalid    = "Parameter query tidak valid."
	MsgLimitInvalidFmt = "harus berupa angka bulat antara 1 dan %d"
	MsgCursorInvalid   = "cursor tidak valid"

	MsgBodyEmpty        = "body request kosong"
	MsgBodyTooLargeFmt  = "body request melebihi %d byte"
	MsgBodyInvalidJSON  = "JSON tidak valid: %w"
	MsgBodySingleObject = "body request hanya boleh berisi satu objek JSON"

	MsgDatabaseURLRequired = "config: DATABASE_URL wajib diisi"
	MsgDatabaseURLInvalid  = "server: DATABASE_URL tidak bisa dibaca"
	MsgPoolCreateFailed    = "server: gagal membuat pool database: %w"
)
