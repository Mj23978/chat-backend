package conf

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	cfg      = config{}
	Global   = &cfg.Global
	Log      = &cfg.Log
	Broker   = &cfg.Broker
	Registry = &cfg.Registry
	Database = &cfg.Database
	Cache    = &cfg.Cache
)

func init() {
	if !cfg.parse() {
		showHelp()
		os.Exit(-1)
	}
}

type global struct {
	Addr string `mapstructure:"addr"`
	Port string `mapstructure:"port"`
	Dc   string `mapstructure:"dc"`
}

type log struct {
	Level string `mapstructure:"level"`
}

type broker struct {
	URL string `mapstructure:"url"`
}

type registry struct {
	Addrs []string `mapstructure:"addrs"`
}

type cache struct {
	URL string `mapstructure:"url"`
}

type database struct {
	URL string `mapstructure:"url"`
}

type config struct {
	Global   global   `mapstructure:"global"`
	Log      log      `mapstructure:"log"`
	Broker   broker   `mapstructure:"broker"`
	Registry registry `mapstructure:"registry"`
	Cache    cache    `mapstructure:"cache"`
	Database database `mapstructure:"database"`
	CfgFile  string
}

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -h (show help info)")
}

func (c *config) load() bool {
	_, err := os.Stat(c.CfgFile)
	if err != nil {
		return false
	}

	viper.SetConfigFile(c.CfgFile)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file %s read failed. %v\n", c.CfgFile, err)
		return false
	}
	err = viper.GetViper().UnmarshalExact(c)
	if err != nil {
		fmt.Printf("config file %s loaded failed. %v\n", c.CfgFile, err)
		return false
	}
	fmt.Printf("config %s load ok!\n", c.CfgFile)
	return true
}

func (c *config) parse() bool {
	flag.StringVar(&c.CfgFile, "c", "conf/conf.toml", "config file")
	help := flag.Bool("h", false, "help info")
	flag.Parse()
	if !c.load() {
		return false
	}

	if *help {
		showHelp()
		return false
	}
	return true
}
