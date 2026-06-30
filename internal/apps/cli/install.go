package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codenomdev/viona/i18n"
	"github.com/codenomdev/viona/pkg/util"
)

func InstallI18nBundle(baseDir string, replace bool) {
	fmt.Println("[i18n] try to install i18n bundle...")

	if len(os.Getenv("SKIP_REPLACE_I18N")) > 0 {
		replace = false
	}

	i18nDir := filepath.Join(baseDir, "i18n")

	if err := util.CreateDirIfNotExist(i18nDir); err != nil {
		fmt.Println(err.Error())
		return
	}

	i18nList, err := i18n.I18n.ReadDir(".")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("[i18n] find i18n bundle %d\n", len(i18nList))

	for _, item := range i18nList {
		targetPath := filepath.Join(i18nDir, item.Name())

		content, err := i18n.I18n.ReadFile(item.Name())
		if err != nil {
			continue
		}

		exist := util.CheckFileExist(targetPath)

		if exist && !replace {
			continue
		}

		if exist {
			fmt.Printf("[i18n] replacing %s\n", item.Name())
			_ = os.Remove(targetPath)
		}

		fmt.Printf("[i18n] install %s bundle...\n", item.Name())

		err = util.WriteFile(targetPath, string(content))
		if err != nil {
			fmt.Printf("[i18n] install %s failed: %s\n", item.Name(), err.Error())
		} else {
			fmt.Printf("[i18n] install %s success\n", item.Name())
		}
	}
}
