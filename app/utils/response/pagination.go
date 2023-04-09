package response

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// Cursor is a structure maintaing information about the cursor
type Cursor struct {
	ID uint
}

// ToBase64String converts a cursor to a base64 encoded string
func (c *Cursor) ToBase64String() string {
	cursor := fmt.Sprintf("%d", c.ID)
	return b64.StdEncoding.EncodeToString([]byte(cursor))
}

// ParseCursor returns a cursor structure from a base64 encoded string
func ParseCursor(cursor string) (*Cursor, error) {
	fromID, err := b64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("error decoding pagination cursor")
	}

	fields := strings.Split(string(fromID), "/")
	if len(fields) < 1 {
		return nil, fmt.Errorf("invalid pagination cursor")
	}

	cID, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	return &Cursor{ID: uint(cID)}, nil
}
