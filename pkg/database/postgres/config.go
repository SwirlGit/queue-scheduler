package postgres

import (
	"fmt"
)

type Config struct {
	URL      string
	Username string
	Password string
}

func (c *Config) fullURL(appName string) string {
	return fmt.Sprintf("%s://%s:%s@%s?application_name=%s", "postgres", c.Username, c.Password, c.URL, appName)
}
