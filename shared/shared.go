package shared

import "fmt"

type User struct {
	Id    []byte `spanner:"Id"`
	Name  string `spanner:"Name"`
	Money int64  `spanner:"Money"`
}
type Item struct {
	Id          []byte `spanner:"Id"`
	Description string `spanner:"Description"`
	UserId      []byte `spanner:"UserId"`
}

type Offer struct {
	Id     []byte `spanner:"Id"`
	Price  int64  `spanner:"Price"`
	ItemId []byte `spanner:"ItemId"`
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
