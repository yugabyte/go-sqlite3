//go:build cgo && sqlite_vec

// cgo compiles only the C files in the package's own folder, so this pulls in sqlite-vec's
// source (sqlite-vec/, synced by upgrade/sqlite-vec.sh) as part of SQLite itself.
#ifdef USE_LIBSQLITE3
#error "the sqlite_vec tag compiles sqlite-vec against this package's SQLite, so it can't be used with libsqlite3"
#endif
#define SQLITE_CORE 1
#define SQLITE_VEC_STATIC 1
#include "sqlite-vec/sqlite-vec.c"
