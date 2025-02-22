// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package database

import (
	"fmt"
	"strings"

	"github.com/canonical/sqlair"
	"github.com/juju/collections/transform"
)

// SliceToPlaceholder returns a string that can be used in a SQL/DML
// statement as a parameter list for a [NOT] IN clause.
// For example, passing []int{1, 2, 3} would return "?,?,?".
// It also returns a suitable transformed slice of the input values to type any.
func SliceToPlaceholder[T any](in []T) (string, []any) {
	vals := make([]any, 0, len(in))
	return strings.Join(transform.Slice(in, func(item T) string {
		vals = append(vals, item)
		return "?"
	}), ","), vals
}

// SliceToPlaceholderTransform returns a string that can be used in SQL/DML
// statement as a parameter list for a [NOT] IN clause.
// For example, passing []int{1, 2, 3} would return "?,?,?".
// Also takes a transform function to alter the type and meaning of the in slice
// into a new slice that can be used with the parameters.
func SliceToPlaceholderTransform[T any](in []T, trans func(T) any) (string, []any) {
	vals := make([]any, 0, len(in))
	return strings.Join(transform.Slice(in, func(item T) string {
		vals = append(vals, trans(item))
		return "?"
	}), ","), vals
}

// MakeBindArgs returns a string of bind args for a given number of columns and
// rows.
func MakeBindArgs(columns, rows int) string {
	var r []string

	c := strings.Repeat("?, ", columns)
	c = c[:len(c)-2]
	for i := 0; i < rows; i++ {
		r = append(r, fmt.Sprintf("(%s)", c))
	}
	return strings.Join(r, ", ")
}

// MapToMultiPlaceholder returns a string of bind args for map key value inserts
// into a table and a flattened args slice.
// Example bind string (?, ?), (?, ?), (?, ?)
func MapToMultiPlaceholder[K comparable, V any](m map[K]V) (string, []any) {
	binds := make([]string, 0, len(m))
	vals := make([]any, 0, len(m)*2)
	for k, v := range m {
		binds = append(binds, "(?, ?)")
		vals = append(vals, k, v)
	}

	return strings.Join(binds, ","), vals
}

// SqlairClauseAnd creates a sqlair query condition where each
// of the non-empty map values becomes an AND operator.
func SqlairClauseAnd(columnValues map[string]any) (string, sqlair.M) {
	var terms []string
	args := sqlair.M{}
	for tableCol, value := range columnValues {
		if value == "" {
			continue
		}
		col := strings.ReplaceAll(tableCol, ".", "_")
		terms = append(terms, tableCol+" = $M."+col)
		args[col] = value
	}
	return strings.Join(terms, " AND "), args
}
