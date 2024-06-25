package config

type AzureResourceConfig struct {
	Name          *string `yaml:"Name"`
	ResourceGroup *string `yaml:"ResourceGroup"`
	Region        *string `yaml:"Region"`
}
