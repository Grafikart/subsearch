package main

import (
	"errors"
	"github.com/Grafikart/subsearch/opensubtitle"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/urfave/cli.v1"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "subsearch"
	app.Usage = "Search subtitles for your files (using opensubtitles.org)"
	app.Action = search
	err := app.Run(os.Args)
	if err != nil {
		color.Red("Error: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func search(cli *cli.Context) (err error) {
	args := cli.Args()
	if len(args) < 1 {
		return errors.New("no paths found, use \"subsearch <path1> <path2>\"")
	}
	for _, file := range args {
		c, err := opensubtitle.NewClient()
		if err != nil {
			return err
		}
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		if err := s.Color("blue"); err != nil {
			return err
		}
		s.Start()
		subtitles, err := c.Search(file)
		s.Stop()
		if err != nil {
			return err
		}
		survey.SelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
	{{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}{{ SelectFocusIcon }} {{else}}{{color "default+hb"}}  {{end}}
	{{- $choice}}
	{{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`
		options := subtitles.ToMap()
		prompt := &survey.Select{
			Message: "Choose a file to download :",
			Options: getKeys(options),
		}
		v := ""
		if err := survey.AskOne(prompt, &v, nil); err != nil {
			return err
		}
		s.Start()
		err = options[v].Download(file + ".srt")
		s.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

func getKeys(m map[string]opensubtitle.Subtitle) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
