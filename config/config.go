
package config


type Config struct {
	Network   Network
	ProjectID string
}

func NewConfig(network Network, projectID string) *Config {
	return &Config{network,projectID}
}


func (c *Config) GetURL() string {
	return c.Network.URL() + c.ProjectID
	
}
