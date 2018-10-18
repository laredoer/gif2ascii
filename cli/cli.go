package cli

import (
	"errors"
	"log"
	"os"

	"github.com/wule61/gif2ascii/gif"

	"github.com/urfave/cli"
)

type Cli struct {
	*cli.App
	FilePath string
}

func New() *Cli {
	c := &Cli{cli.NewApp(), ""}
	c.Name = "gif2ascii"
	c.Version = "1.0.0"
	c.Usage = "make gif to ascii gif"

	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "git path",
			Destination: &c.FilePath, //取到的FLAG值，赋值到这个变量
		},
	}

	c.Action = func(c *cli.Context) error {
		if c.String("f") != "" {
			images := gif.GetEach(c.String("f"))
			for k, v := range images {
				str := gif.Ascllimage(v)
				gif.DrawImg(str, k)
			}
			gif.CreateGif()
			os.RemoveAll("./tmp")
			log.Fatal("Wrote output.png OK.")
		} else {
			return errors.New("没有指定gif文件路径")
		}
		return nil
	}

	return c
}
