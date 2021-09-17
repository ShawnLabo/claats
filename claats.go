package claats

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/googlecodelabs/tools/claat/cmd"
	"github.com/haya14busa/ghglob"
)

type Option struct {
	In      string
	Out     string
	Pattern string
	GA      string
	Title   string
}

func Generate(opt Option) error {
	files, err := ghglob.GlobList([]string{opt.Pattern}, ghglob.Option{Root: opt.In})
	if err != nil {
		return fmt.Errorf("ghglob.GlobList: %w", err)
	}

	for _, f := range files {
		if err := doClaat(f, opt); err != nil {
			return fmt.Errorf("doClaat: %w", err)
		}
	}

	if err := createIndex(files, opt); err != nil {
		return fmt.Errorf("createIndex: %w", err)
	}

	return nil
}

func doClaat(path string, opt Option) error {
	relPath, err := filepath.Rel(opt.In, path)
	if err != nil {
		return fmt.Errorf("filepath.Rel: %w", err)
	}

	eo := cmd.CmdExportOptions{
		Expenv:   "web",
		GlobalGA: opt.GA,
		Output:   filepath.Dir(filepath.Dir(filepath.Join(opt.Out, relPath))),
		Prefix:   "https://storage.googleapis.com",
		Srcs:     []string{path},
		Tmplout:  "html",
	}

	if err := os.MkdirAll(eo.Output, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	code := cmd.CmdExport(eo)

	if code == 0 {
		return nil
	}

	return fmt.Errorf("cmd.CmdExport: exited status %d", code)
}

type lab struct {
	Path     string    `json:"-"`
	Updated  time.Time `json:"updated"`
	Duration uint32    `json:"duration"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Tags     []string  `json:"tags"`
}

type templateData struct {
	Title string
	GA    string
	Labs  []*lab
}

func createIndex(files []string, opt Option) error {
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return fmt.Errorf("template.New.Parse: %w", err)
	}

	labs := []*lab{}

	for _, f := range files {
		l, err := getLab(f, opt)
		if err != nil {
			return fmt.Errorf("getLab: %w", err)
		}

		labs = append(labs, l)
	}

	data := templateData{
		Title: opt.Title,
		GA:    opt.GA,
		Labs:  labs,
	}

	if err := t.Execute(os.Stdout, data); err != nil {
		return fmt.Errorf("t.Execute: %w", err)
	}

	return nil
}

func getLab(path string, opt Option) (*lab, error) {
	relPath, err := filepath.Rel(opt.In, path)
	if err != nil {
		return nil, fmt.Errorf("filepath.Rel: %w", err)
	}

	dir := filepath.Dir(relPath)

	f, err := os.Open(filepath.Join(opt.Out, dir, "codelab.json"))
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	l := &lab{}
	if err := json.NewDecoder(f).Decode(l); err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	l.Path = filepath.Join(dir, "index.html")

	return l, nil
}
