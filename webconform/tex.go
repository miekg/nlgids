package webconform

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

const Template = "conform.tex.tmpl"

var funcMap = template.FuncMap{
	"half":            Half,
	"divide":          Divide,
	"euro":            Euro,
	"divideTimesRate": DivideTimesRate,
}

func Half(a float64) float64                     { return a / 2 }
func Euro(a, rate float64) float64               { return a * rate }
func DivideTimesRate(a, b, rate float64) float64 { return a / b * rate }
func Divide(a float64, b int) float64            { return a / float64(b) }
func DivideFloat(a, b float64) float64           { return a / b }

// Create parses the templates and runs pdflatex on the resulting tex file. It returns generated PDF.
func (c *Conform) Create(tmplDir, tmpl string) ([]byte, error) {
	c.escapeTex()
	err := c.FillOut()
	if err != nil {
		return nil, err
	}

	dst, err := ioutil.TempDir("/tmp", "pdflatex")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dst)

	sources, err := TeXFiles(tmplDir)
	if err != nil {
		return nil, err
	}
	for _, src := range sources {
		basename := path.Base(src)
		if err := os.Symlink(src, path.Join(dst, basename)); err != nil {
			return nil, err
		}
	}

	conformTex := path.Join(dst, "conform.tex")

	// We use @@ in the tmpl, because {{ is in heavy use by TeX.
	// This also creates a dummy template called "", that we don't use.
	tmplName := path.Base(tmpl)
	t, err := template.New(tmplName).Delims("@@", "@@").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		return nil, err
	}
	if err := c.ExecuteTemplateAndWrite(t, tmplName, conformTex); err != nil {
		return nil, err
	}

	pdf, err := pdfLateX(dst, conformTex)
	if err != nil {
		return nil, err
	}
	return pdf, nil
}

// TeXFiles returns the tex and jpg files found in the directory dir.
func TeXFiles(dir string) ([]string, error) {
	tex := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		switch path.Ext(f.Name()) {
		case ".tex", ".jpg":
			tex = append(tex, path.Join(dir, f.Name()))
		}
	}
	return tex, nil
}

func pdfLateX(dir, texfile string) ([]byte, error) {
	cmd := exec.Command("pdflatex", texfile)
	cmd.Dir = dir

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err := cmd.Wait()

	pdf := strings.Replace(texfile, ".tex", ".pdf", 1)
	data, err := ioutil.ReadFile(pdf)
	return data, err
}
