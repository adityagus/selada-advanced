package main

import (
	"database/sql"
	"strings"
	"testing"
)

func TestColumnConversion(t *testing.T) {
	for _, tc := range []struct{ kind, full, want string }{
		{"smallint", "smallint unsigned", "integer"},
		{"int", "int unsigned", "bigint"},
		{"bigint", "bigint unsigned", "numeric(20,0)"},
		{"tinyint", "tinyint(1)", "smallint"},
		{"decimal", "decimal(18,2) unsigned", "decimal(18,2)"},
		{"datetime", "datetime", "timestamp without time zone"},
	} {
		got, err := pgType(column{kind: tc.kind, full: tc.full})
		if err != nil || got != tc.want {
			t.Errorf("%s: got %s, %v; want %s", tc.full, got, err, tc.want)
		}
	}
	if _, err := definition(column{kind: "bigint", full: "bigint unsigned", extra: "auto_increment"}); err == nil {
		t.Fatal("numeric identity must be rejected")
	}
	got, err := definition(column{name: "note", kind: "varchar", full: "varchar(30)", nullable: "NO", def: sql.NullString{String: "O'Brien", Valid: true}})
	if err != nil || !strings.Contains(got, "DEFAULT 'O''Brien'") {
		t.Fatalf("literal default: %s, %v", got, err)
	}
	if _, err := definition(column{kind: "timestamp", full: "timestamp", extra: "on update CURRENT_TIMESTAMP"}); err == nil {
		t.Fatal("unsupported extras must be rejected")
	}
}
