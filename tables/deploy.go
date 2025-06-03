package tables

import (
	"database/sql"
	"deeployer/db"
	"fmt"
	"log"
	"time"
)

type Step struct {
	Id           int
	BuildLogId   int
	Status       string
	StatusInt    int
	FailedLogs   sql.NullString
	FailedReason sql.NullString
}

type BuildLog struct {
	Id        int
	ConfigId  int
	Timestamp time.Time
}

type Builds []BuildLog
type Steps []Step

func (steps Steps) insertQuery() error {
	query := "insert into build_logs (`build_log_id`, `status`" +
		", `status_int`, `failed_logs`, `failed_reason`) values" +
		" (?, ?, ?, ?, ?);"

	for _, step := range steps {
		res, err := db.DB.Exec(
			query, step.BuildLogId, step.Status, step.StatusInt,
			step.FailedLogs, step.FailedReason,
		)
		if err != nil {
			return err
		}

		if rowCount, err := res.RowsAffected(); err == nil && rowCount > 0 {
			log.Println("step inserted")
		} else {
			return fmt.Errorf("insert failed without any error message, please check")
		}
	}
	return nil
}

func (builds Builds) insertQuery() error {
	query := "insert into config_logs (`config_id`, `timestamp`) values (?, ?);"
	for _, build := range builds {
		res, err := db.DB.Exec(query, build.ConfigId, build.Timestamp)
		if err != nil {
			return err
		}
		if rowCount, err := res.RowsAffected(); err == nil && rowCount > 0 {
			log.Println("config inserted")
		} else {
			return fmt.Errorf("Insert failed without any error message, please check")
		}
	}
	return nil
}
