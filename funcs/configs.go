package funcs

import (
	"deeployer/tables"
	"fmt"
)

// config functionalities

func CreateConfig(userId int, sshKey string, githubKey string, projectName string, repourl string, host string, user string) tables.Config {
	c := tables.Config{
		UserId:      userId,
		SshKey:      sshKey,
		GitKey:      githubKey,
		ProjectName: projectName,
		RepoUrl:     repourl,
		Host:        host,
		User:        user,
	}
	configs := tables.Configs{
		c,
	}

	err := tables.InsertQuery(configs)
	if err != nil {
		fmt.Println("Error inserting configs", err)
	}
	return c
}
