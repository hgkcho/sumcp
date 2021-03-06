package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

// ExitError handle error
type ExitError struct {
	exitCode int
	err      error
}

// NewExitError init ExitError
func NewExitError(exitCode int, err error) *ExitError {
	return &ExitError{
		exitCode: exitCode,
		err:      err,
	}
}

func (ee *ExitError) Error() string {
	if ee.err == nil {
		return ""
	}
	return fmt.Sprintf("%v", ee.err)
}

var version = "1.0.2"

func main() {
	name := "sumcp"

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "sumcp",
		Usage:                "copy multiple file contents to one file",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "force", Aliases: []string{"f"}},
		},
		UsageText: fmt.Sprintf(
			`%s

	  Version: %s

		copy data from multipul files to one file

		ex) %s a.txt b.txt c.txt

		expect
		cat a.txt > c.txt
		cat b.txt >> c.txt
	`, name, name, version,
		),
		Action: func(c *cli.Context) error {
			// -f ではじまるものはフラグ、なのでじょがい
			var args []string
			for _, v := range c.Args().Slice() {
				if !strings.HasPrefix(v, "-f") {
					args = append(args, v)
				}
			}

			argsLen := len(args)
			if argsLen < 2 {
				return fmt.Errorf("too few arguments. require more than 2 args, but get %d", argsLen)
			}

			sources := args[:(argsLen - 1)]
			target := args[(argsLen - 1)]
			force := c.Bool("force")

			if !force {
				fmt.Printf("Can I overwrite file %s? [y/N] \n", target)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				ans := scanner.Text()
				if ans != "y" && ans != "Y" {
					fmt.Fprintf(os.Stderr, "[Error] command %s is cancelled \n", name)
					os.Exit(exitCodeErr)
				}
			}

			start := time.Now()
			if err := run(sources, target); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(exitCodeErr)
			}
			end := time.Now()
			fmt.Fprintf(os.Stdout, "Done (time : %v)\n", end.Sub(start))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func run(sources []string, target string) error {

	for _, s := range sources {
		_, err := os.Lstat(s)
		if err != nil {
			os.Exit(exitCodeErr)
		}
	}

	file, err := os.Create(target)
	if err != nil {
		return err
	}

	defer file.Close()

	existed := make(map[string]bool, len(sources))
	wg := new(sync.WaitGroup)

	for _, source := range sources {
		source := source

		wg.Add(1)
		go func() {
			defer wg.Done()

			if existed[source] {
				fmt.Fprintf(os.Stderr, "duplicat source %s", source)
			}

			b, err := ioutil.ReadFile(source)
			info, err := os.Lstat(source)
			fmt.Fprintf(os.Stdout, "filename %s \n", info.Name())
			if err != nil {
				log.Fatal(err)
			}

			file.WriteString(string(b) + "\n")
		}()
		wg.Wait()
	}
	return nil
}
