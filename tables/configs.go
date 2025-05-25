package tables

import (
	"database/sql"
	"deeployer/db"
	"fmt"
)

type Config struct {
	Id            int
	UserId        int
	SshKey        string
	GitKey        string
	ProjectName   string
	RepoUrl       string
	Host          string
	User          string
	IsDockerised  bool
	BuildCommands sql.NullString
}

type Configs []Config

func GetProjectConfig(id int) (Config, error) {
	var c Config
	query := "select * from configs where id = ?;"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		return c, fmt.Errorf("error getting config %v", err)
	}
	defer rows.Close()
	if rows.Next() {

		err := rows.Scan(&c.Id, &c.UserId, &c.SshKey, &c.GitKey, &c.ProjectName, &c.RepoUrl, &c.Host, &c.User, &c.IsDockerised, &c.BuildCommands)
		if err != nil {
			return c, fmt.Errorf("Error getting config values, %v", err)
		}
	} else {
		return c, fmt.Errorf("No config found for this id= %d", id)
	}

	if err := rows.Err(); err != nil {
		return c, fmt.Errorf("Row iteration error, %v", err)
	}
	return c, nil
}

func (c Configs) insertQuery() error {
	str := "insert into configs (`user_id`, `ssh_key`, `github_key`, `project_name`, `repourl`, `host`, `user`) values"
	for i := 0; i < len(c); {
		s := "(%d, '%s', '%s', '%s', '%s', '%s', '%s')"
		values := fmt.Sprintf(s, c[i].UserId, c[i].SshKey, c[i].GitKey, c[i].ProjectName, c[i].RepoUrl, c[i].Host, c[i].User)
		str += values
		i += 1
	}

	_, err := db.DB.Query(str)
	if err != nil {
		return fmt.Errorf("Error creating a config %v", err)
	}
	return nil
}

func GetAllConfigs() ([]Config, error) {
	query := "select project_name, is_dockerised, repourl, build_commands, host, user from configs;"

	rows, err := db.DB.Query(query)
	var projects []Config

	if err != nil {
		return projects, fmt.Errorf("error in query, err: %v", err)
	}

	for rows.Next() {
		var c Config
		rows.Scan(&c.ProjectName, &c.IsDockerised, &c.RepoUrl, &c.BuildCommands, &c.Host, &c.User)
		projects = append(projects, c)
	}

	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("err while iteration, err %v", err)
	}
	return projects, nil
}
