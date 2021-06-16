package shared

import "fmt"

type User struct {
	Id    []byte `spanner:"Id"`
	Name  string `spanner:"Name"`
	Money int64  `spanner:"Money"`
}

type GcloudConfig struct {
	Project  string
	Instance string
	Database string
}

func (config GcloudConfig) Uri() string {
	return fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		config.Project,
		config.Instance,
		config.Database,
	)
}

func LocalConfig() GcloudConfig {
	return GcloudConfig{
		Project:  "noted-episode-316407",
		Instance: "test-instance",
		Database: "test-database",
	}
}
