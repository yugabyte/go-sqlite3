//go:build cgo && sqlite_vec

package sqlite3

import (
	"database/sql"
	"os"
	"strings"
	"testing"
)

func TestVecVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var got string
	if err := db.QueryRow("SELECT vec_version()").Scan(&got); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile("sqlite-vec/VERSION")
	if err != nil {
		t.Fatal(err)
	}
	first, _, _ := strings.Cut(string(b), "\n")
	if want := "v" + strings.TrimPrefix(first, "sqlite-vec "); got != want {
		t.Fatalf("vec_version() = %s, but sqlite-vec/VERSION says %s", got, want)
	}
}

func TestVec0(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)
	db, err := sql.Open("sqlite3", tempFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// A partition key and a metadata TEXT column: deleting a row whose text is longer than 12
	// characters failed with a spurious SQLITE_DONE before sqlite-vec 0.1.9.
	if _, err := db.Exec(`CREATE VIRTUAL TABLE v USING vec0(
		part TEXT PARTITION KEY, label TEXT, embedding float[2] distance_metric=cosine)`); err != nil {
		t.Fatal(err)
	}
	rows := []struct {
		id    int
		part  string
		label string
		vec   string
	}{
		{1, "a", "east, longer than twelve", "[1, 0]"},
		{2, "a", "north", "[0, 1]"},
		{3, "a", "north-east", "[1, 1]"},
		{4, "b", "east in b", "[1, 0.1]"},
	}
	for _, r := range rows {
		if _, err := db.Exec("INSERT INTO v(rowid, part, label, embedding) VALUES (?, ?, ?, ?)",
			r.id, r.part, r.label, r.vec); err != nil {
			t.Fatal(err)
		}
	}

	knn := func(part string) []int {
		t.Helper()
		rs, err := db.Query(`SELECT rowid FROM v WHERE embedding MATCH '[1, 0]' AND part = ? AND k = 2
			ORDER BY distance`, part)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Close()
		var ids []int
		for rs.Next() {
			var id int
			if err := rs.Scan(&id); err != nil {
				t.Fatal(err)
			}
			ids = append(ids, id)
		}
		if err := rs.Err(); err != nil {
			t.Fatal(err)
		}
		return ids
	}
	if got := knn("a"); len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Fatalf("KNN in partition a = %v, want [1 3]", got)
	}
	if got := knn("b"); len(got) != 1 || got[0] != 4 {
		t.Fatalf("KNN in partition b = %v, want [4]", got)
	}

	if _, err := db.Exec("DELETE FROM v WHERE rowid = 1"); err != nil {
		t.Fatal(err)
	}
	if got := knn("a"); len(got) != 2 || got[0] != 3 || got[1] != 2 {
		t.Fatalf("KNN in partition a after DELETE = %v, want [3 2]", got)
	}
}
