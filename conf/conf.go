package conf

type Config struct {
	Port   int    `toml:"port"`
	Cookie string `toml:"cookie"`
	Secret string `toml:"secret"`
}
