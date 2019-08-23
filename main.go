package main

import (
	_ "github.com/iikira/lfsutil/initcli"
	"github.com/iikira/lfsutil/internal/lfscommand"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	commentFlag = cli.StringFlag{
		Name: "comment",
	}
)

func main() {
	app := cli.NewApp()
	app.Version = "v1.0.0"
	app.Usage = "A simple tool for accessing git lfs server"
	app.Author = "iikira"
	app.Copyright = "(c) 2016-2019 iikira."
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "getobject",
			Aliases: []string{"go"},
			Action: func(c *cli.Context) (err error) {
				if len(c.Args()) == 0 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				opt := lfscommand.GetObjectOption{
					Args:      c.Args(),
					By:        c.String("by"),
					NoTest:    c.Bool("no-test"),
					PtrDir:    c.String("ptr_dir"),
					PtrSuffix: c.String("ptr_suffix"),
				}
				lfscommand.GetObject(&opt)
				return
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "by",
					Usage: "oid (oid:size or oid), file, ptr (filename or url)",
					Value: "oid",
				},
				cli.BoolFlag{
					Name: "no-test",
				},
				cli.StringFlag{
					Name:  "ptr_dir",
					Usage: "path to save ptr (TODO)",
					Value: "",
				},
				cli.StringFlag{
					Name:  "ptr_suffix",
					Usage: "suffix (TODO)",
					Value: "",
				},
				commentFlag,
			},
		},
		cli.Command{
			Name:    "upobject",
			Aliases: []string{"uo"},
			Action: func(c *cli.Context) (err error) {
				if len(c.Args()) == 0 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				opt := lfscommand.UpObjectOption{
					Args:      c.Args(),
					PtrDir:    c.String("ptr_dir"),
					PtrSuffix: c.String("ptr_suffix"),
					NoVerify:  c.Bool("no-verify"),
				}
				lfscommand.UpObject(&opt)
				return
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ptr_dir",
					Usage: "path to save ptr",
					Value: "",
				},
				cli.StringFlag{
					Name:  "ptr_suffix",
					Usage: "suffix",
					Value: "",
				},
				cli.BoolFlag{
					Name:  "no-verify",
					Usage: "skip verify",
				},
				commentFlag,
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	err := app.Run(os.Args[0:])
	if err != nil {
		log.Fatalln(err)
	}
}
