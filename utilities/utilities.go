package utilities

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"
)

func MD5(data []byte) string {
	if data == nil {
		return "" // or return "NULL" or whatever you want to represent NULL
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}

func VersionHash(arr []sql.NullString) string {
	var strs []string = make([]string, 0)
	for _, v := range arr {
		if v.Valid {
			strs = append(strs, v.String)
		}
	}

	return MD5([]byte(strings.Join(strs, ",")))
}

func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}

	return sql.NullString{
		Valid:  true,
		String: *s,
	}
}
func NotNullString(s string) sql.NullString {
	return sql.NullString{
		Valid:  true,
		String: s,
	}
}
func NullStringPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}