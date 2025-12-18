package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgDateToTime(d pgtype.Date) time.Time {
    t, _ := d.Value()
    return t.(time.Time)
}
