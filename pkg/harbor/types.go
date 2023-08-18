package harbor

import (
	"time"
)

const (
	AccessOperationDelete = "delete"
)

type Tag struct {
	Digest       string    `json:"digest"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	Architecture string    `json:"architecture"`
	OS           string    `json:"os"`
	Created      time.Time `json:"created"`
	PullTime     time.Time `json:"pull_time"`
}

type AccessLog struct {
	LogID     int64  `json:"log_id"`
	ProjectID int64  `json:"project_id"`
	RepoName  string `json:"repo_name"`
	Tag       string `json:"repo_tag"`
	Operation string `json:"operation"`
}
