//go:build cgo && sqlite_vec

package sqlite3

/*
#cgo linux LDFLAGS: -lm
#include "sqlite3-binding.h"

int sqlite3_vec_init(sqlite3 *db, char **pzErrMsg, const sqlite3_api_routines *pApi);
*/
import "C"

import "fmt"

// sqlite-vec (sqlite3_opt_vec.c) loads into every connection: the vec0 virtual table and the
// vec_* functions.
func init() {
	if rc := C.sqlite3_auto_extension((*[0]byte)(C.sqlite3_vec_init)); rc != C.SQLITE_OK {
		panic(fmt.Sprintf("sqlite3: registering sqlite-vec failed with SQLite error %d", int(rc)))
	}
}
