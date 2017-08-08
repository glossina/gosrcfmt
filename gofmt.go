package gosrcfmt

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/sirkon/message"
)

const (
	header = `
/* This file was autogenerated via
 {{ .Split }}
 {{ .Cmd }}
 {{ .Split }}
do not touch it with bare hands!
Although you perfectly can, just remove this warning first.
*/

{{ .Code }}

`
)

// Format formats data from src as a go code with gofmt utility
func Format(dest io.Writer, data []byte) {
	t := template.New("header")
	tmpl, err := t.Parse(header)
	if err != nil {
		panic(err)
	}

	var ctx struct {
		Split string
		Cmd   string
		Code  string
	}
	ctx.Cmd = strings.Join(os.Args, " ")
	length := len(ctx.Cmd)
	if length < 38 {
		length = 38
	}
	ctx.Split = strings.Repeat("-", length)
	ctx.Code = string(data)
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}
	program := buf.String()

	cmd := exec.Command("gofmt")
	cmd.Stdin = strings.NewReader(program)
	cmd.Stdout = dest
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		message.Criticalf("%s\n\n---------------------------------------\n%s", program, err)
	}
}

// FormatReader formats data from src as a go code with gofmt utility
func FormatReader(dest io.Writer, src io.Reader) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		panic(err)
	}
	Format(dest, data)
}
