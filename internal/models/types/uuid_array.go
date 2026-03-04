package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (u *UUIDArray) Scan(value any) error {
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("Failed to parse UUIDArray: unsupport data type")
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	parts := strings.Split(str, ",")

	*u = make(UUIDArray, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(strings.Trim(s, `"`))

		if s == "" {
			continue
		}

		uu, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("Invalid UUID in Array: %v", err)
		}

		*u = append(*u, uu)
	}

	return nil
}

func (u UUIDArray) Value() (driver.Value, error) {
	if len(u) == 0 {
		return "{}", nil
	}

	dbFormat := make([]string, 0, len(u))

	for _, v := range u {
		dbFormat = append(dbFormat, fmt.Sprintf(`"%s"`, v.String()))
	}

	return "{" + strings.Join(dbFormat, ",") + "}", nil
}

func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
