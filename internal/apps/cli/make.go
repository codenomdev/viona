package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/codenomdev/viona/pkg/util"
	"github.com/codenomdev/viona/templates"
	"github.com/gookit/goutil/strutil"
	"github.com/spf13/cobra"
)

type GeneratorConfig struct {
	CommandName string
	Short       string
	Long        string
	Dir         string
	FileSuffix  string
	Template    string
	SuccessMsg  string
}

func MakeGeneratorCmd(cfg GeneratorConfig) *cobra.Command {
	return &cobra.Command{
		Use:   cfg.CommandName,
		Short: cfg.Short,
		Long:  cfg.Long,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic("name is required")
			}

			name := strutil.Lowercase(args[0])
			name = strutil.SnakeCase(name, "_")

			id := util.TimeNow().Format("20060102150405")

			filename := fmt.Sprintf("%s_%s_%s.go", id, name, cfg.FileSuffix)
			path := filepath.Join(cfg.Dir, filename)

			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				panic(err)
			}

			file, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			tpl := template.Must(template.ParseFS(templates.TemplatesFS, cfg.Template))

			err = tpl.Execute(file, map[string]string{
				"ID":   id,
				"Name": fmt.Sprintf("%s_%s", name, cfg.FileSuffix),
			})
			if err != nil {
				panic(err)
			}

			fmt.Println(cfg.SuccessMsg, path)
		},
	}
}
