package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"mapleafgo.cn/sensitive/server"
)

func configStart(c *cli.Context) error {
	viper.SetConfigFile(c.Path("path"))
	viper.ReadInConfig()
	return server.Start()
}

func configFlag(c *cli.Context) error {
	viper.Set("port", c.Int("port"))
	viper.Set("path", c.String("path"))
	return server.Start()
}

func main() {
	c := cli.App{
		Name:  "sensitive",
		Usage: "敏感词过滤服务",
		Authors: []*cli.Author{
			{Name: "mapleafgo", Email: "mapleafgo@163.com"},
		},
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:   "config",
				Usage:  "配置文件启动",
				Flags:  []cli.Flag{&cli.StringFlag{Name: "path", Aliases: []string{"c"}, Usage: "*配置文件路径", Required: true}},
				Action: configStart,
			},
			{
				Name:  "flag",
				Usage: "命令启动",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "port", Aliases: []string{"o"}, Usage: "*服务端口", Required: true},
					&cli.StringFlag{Name: "path", Aliases: []string{"p"}, Usage: "词典路径"},
				},
				Action: configFlag,
			},
		},
	}
	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
