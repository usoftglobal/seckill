package libs

type Config struct {
	MySQL struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Password string `yaml:"password"`
		DB string `yaml:"db"`
	}
	Redis struct {
		Host string `yaml:"host"`
	}
}