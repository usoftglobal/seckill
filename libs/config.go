package libs

type Config struct {
	MySQL struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Password string `yaml:"password"`
		DB string `yaml:"db"`
	}
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}
