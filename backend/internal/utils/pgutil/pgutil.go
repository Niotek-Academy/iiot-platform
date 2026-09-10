// Package pgutil holds tiny conversion helpers shared across services when
// talking to pgx-generated types.
package pgutil

import "github.com/jackc/pgx/v5/pgtype"

// ToText converts an optional *string (nil = "not provided") into the
// pgtype.Text that sqlc's COALESCE(sqlc.narg(...), ...) pattern expects.
func ToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}