package db

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	zdgoutil "github.com/zerodoctor/zdgo-util"
)

const (
	ENV_SCHEMA string = `
CREATE TABLE IF NOT EXISTS envs (
	project_name TEXT NOT NULL,
	file_name    TEXT NOT NULL,
	file_content TEXT,
	created_at   INTEGER NOT NULL,
	PRIMARY KEY (project_name, file_name)
);`

	ENV_INSERT string = `
INSERT INTO envs (
	project_name, file_name, file_content, created_at
) VALUES (
	$1, $2, $3, $4
) ON CONFLICT (project_name, file_name) DO UPDATE SET
	project_name = excluded.project_name, 
	file_name    = excluded.file_name, 
	file_content = excluded.file_content, 
	created_at   = excluded.created_at;`

	ENV_QUERY_FILE string = `SELECT * FROM envs WHERE project_name = $1 AND file_name = $2;`
	ENV_QUERY_ALL  string = `SELECT * FROM envs;`
)

type Handler struct {
	*sqlx.DB
}

func NewHandler() (*Handler, error) {
	h := &Handler{}

	db, err := sqlx.Connect("sqlite3", util.EXEC_PATH+"/lite.db")
	if err != nil {
		return nil, err
	}
	h.DB = db
	logger.Info("connected to lite.db")

	return h, nil
}

func (h *Handler) CreateTables() {
	if err := h.Transact(true, func(tx *sqlx.Tx) error {
		_, err := tx.Exec(ENV_SCHEMA)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Printf("[error=%s]", err.Error())
	}
}

func (h *Handler) SaveEnvFile(file string, project string) error {
	all, err := util.GetAllFiles(file)
	if err != nil {
		return err
	}

	return h.Transact(false, func(tx *sqlx.Tx) error {
		for _, f := range all {
			fmt.Printf("saving [file=%s]...", f.Path+"/"+f.Name())
			data, err := ioutil.ReadFile(f.Path + "/" + f.Name())
			if err != nil {
				return err
			}

			h.Exec(ENV_INSERT, project, f.Name(), string(data), time.Now().Unix())
		}

		return nil
	})
}

type Env struct {
	ProjectName string    `db:"project_name" json:"project_name"`
	FileName    string    `db:"file_name" json:"file_name"`
	FileContent string    `db:"file_content" json:"file_content"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func (h *Handler) ReadEnvFile(project string, file string) ([]Env, error) {
	var result []Env
	err := h.Select(&result, ENV_QUERY_FILE, project, file)
	return result, err
}

func (h *Handler) ReadAllEnv() ([]Env, error) {
	var result []Env
	err := h.Select(&result, ENV_QUERY_FILE)
	return result, err
}

func (h *Handler) Transact(retry bool, fn func(*sqlx.Tx) error) error {
	tx, err := h.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !retry {
		return tx.Commit()
	}

	return zdgoutil.Retry(func() error {
		return tx.Commit()
	}, zdgoutil.RetryAmountOpt(3), zdgoutil.RetryDurationOpt(3*time.Second))
}
