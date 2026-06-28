package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/codenomdev/viona/pkg/util"
	"github.com/codenomdev/viona/templates"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

const goModTpl = `module viona
go 1.25.1
`

type codenomBuilder struct {
	buildingMaterial *buildingMaterial
	BuildError       error
}

type buildingMaterial struct {
	moduleReplacement   string
	plugins             []*pluginInfo
	outputPath          string
	tmpDir              string
	originalCodenomInfo OriginalCodenomInfo
}

type OriginalCodenomInfo struct {
	Version  string
	Revision string
	Time     string
}

type pluginInfo struct {
	// Name of the plugin e.g. github.com/codenomdev/viona-plugins/xxx
	Name string
	// Path to the plugin. If path exist, read plugin from local filesystem
	Path string
	// Version of the plugin
	Version string
}

func newCodenomBuilder(buildDir, outputPath string, plugins []string, originalCodenomInfo OriginalCodenomInfo) *codenomBuilder {
	material := &buildingMaterial{originalCodenomInfo: originalCodenomInfo}
	parentDir, _ := filepath.Abs(".")
	if buildDir != "" {
		material.tmpDir = filepath.Join(parentDir, buildDir)
	} else {
		material.tmpDir, _ = os.MkdirTemp(parentDir, "codenom_build")
	}
	if len(outputPath) == 0 {
		outputPath = filepath.Join(parentDir, "new_codenom")
	}
	material.outputPath, _ = filepath.Abs(outputPath)
	material.plugins = formatPlugins(plugins)
	material.moduleReplacement = os.Getenv("CODENOM_MODULE")
	return &codenomBuilder{
		buildingMaterial: material,
	}
}

func (a *codenomBuilder) DoTask(task func(b *buildingMaterial) error) {
	if a.BuildError != nil {
		return
	}
	a.BuildError = task(a.buildingMaterial)
}

// BuildNewCodenom builds a new codenom with specified plugins
func BuildNewCodenom(buildDir, outputPath string, plugins []string, originalCodenomInfo OriginalCodenomInfo) (err error) {
	builder := newCodenomBuilder(buildDir, outputPath, plugins, originalCodenomInfo)
	builder.DoTask(createMainGoFile)
	builder.DoTask(downloadGoModFile)
	builder.DoTask(movePluginToVendor)
	builder.DoTask(copyUIFiles)
	builder.DoTask(buildUI)
	builder.DoTask(mergeI18nFiles)
	builder.DoTask(buildBinary)
	builder.DoTask(cleanByproduct)
	return builder.BuildError
}

func formatPlugins(plugins []string) (formatted []*pluginInfo) {
	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		// plugin description like this 'github.com/codenomdev/viona-plugin/github-connector@latest=/local/path'
		info := &pluginInfo{}
		plugin, info.Path, _ = strings.Cut(plugin, "=")
		info.Name, info.Version, _ = strings.Cut(plugin, "@")
		// Resolve local path to absolute since build runs in a temp directory
		if len(info.Path) > 0 {
			if absPath, err := filepath.Abs(info.Path); err == nil {
				info.Path = absPath
			}
		}
		formatted = append(formatted, info)
	}
	return formatted
}

// mergeI18nFiles merge i18n files
func mergeI18nFiles(b *buildingMaterial) (err error) {
	fmt.Printf("try to merge i18n files\n")

	type YamlPluginContent struct {
		Plugin map[string]any `yaml:"plugin"`
	}

	pluginAllTranslations := make(map[string]*YamlPluginContent)
	for _, plugin := range b.plugins {
		i18nDir := filepath.Join(b.tmpDir, fmt.Sprintf("vendor/%s/i18n", plugin.Name))
		fmt.Println("i18n dir: ", i18nDir)
		if !util.CheckDirExist(i18nDir) {
			continue
		}

		entries, err := os.ReadDir(i18nDir)
		if err != nil {
			return err
		}

		for _, file := range entries {
			// ignore directory
			if file.IsDir() {
				continue
			}
			// ignore non-YAML file
			if filepath.Ext(file.Name()) != ".yaml" {
				continue
			}
			buf, err := os.ReadFile(filepath.Join(i18nDir, file.Name()))
			if err != nil {
				log.Debugf("read translation file failed: %s %s", file.Name(), err)
				continue
			}

			translation := &YamlPluginContent{}
			if err = yaml.Unmarshal(buf, translation); err != nil {
				log.Debugf("unmarshal translation file failed: %s %s", file.Name(), err)
				continue
			}

			if pluginAllTranslations[file.Name()] == nil {
				pluginAllTranslations[file.Name()] = &YamlPluginContent{Plugin: make(map[string]any)}
			}
			for k, v := range translation.Plugin {
				pluginAllTranslations[file.Name()].Plugin[k] = v
			}
		}
	}

	originalI18nDir := filepath.Join(b.tmpDir, "vendor/github.com/codenomdev/viona/i18n")
	entries, err := os.ReadDir(originalI18nDir)
	if err != nil {
		return err
	}

	for _, file := range entries {
		// ignore directory
		if file.IsDir() {
			continue
		}
		// ignore non-YAML file
		filename := file.Name()
		if filepath.Ext(filename) != ".yaml" && filename != "i18n.yaml" {
			continue
		}

		// if plugin don't have this translation file, ignore it
		if pluginAllTranslations[filename] == nil {
			continue
		}

		out, _ := yaml.Marshal(pluginAllTranslations[filename])

		buf, err := os.OpenFile(filepath.Join(originalI18nDir, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Debugf("read translation file failed: %s %s", filename, err)
			continue
		}

		_, _ = buf.WriteString("\n")
		_, _ = buf.Write(out)
		_ = buf.Close()
	}
	return err
}

// buildUI run pnpm install and pnpm build commands to build ui
func buildUI(b *buildingMaterial) (err error) {
	localUIBuildDir := filepath.Join(b.tmpDir, "vendor/github.com/codenomdev/viona/ui")

	pnpmInstallCmd := b.newExecCmd("pnpm", "pre-install")
	pnpmInstallCmd.Dir = localUIBuildDir
	if err = pnpmInstallCmd.Run(); err != nil {
		return err
	}

	pnpmBuildCmd := b.newExecCmd("pnpm", "build")
	pnpmBuildCmd.Dir = localUIBuildDir
	if err = pnpmBuildCmd.Run(); err != nil {
		return err
	}
	return nil
}

// copyUIFiles copy ui files from codenom module to tmp dir
func copyUIFiles(b *buildingMaterial) (err error) {
	goListCmd := b.newExecCmd("go", "list", "-mod=mod", "-m", "-f", "{{.Dir}}", "github.com/codenomdev/viona")
	buf := new(bytes.Buffer)
	goListCmd.Stdout = buf
	if err = goListCmd.Run(); err != nil {
		return fmt.Errorf("failed to run go list: %w", err)
	}

	vionaDir := strings.TrimSpace(buf.String())
	goModUIDir := filepath.Join(vionaDir, "ui")
	localUIBuildDir := filepath.Join(b.tmpDir, "vendor/github.com/codenomdev/viona/ui/")
	// The node_modules folder generated during development will interfere packaging, so it needs to be ignored.
	if err = copyDirEntries(os.DirFS(goModUIDir), ".", localUIBuildDir, "node_modules"); err != nil {
		return fmt.Errorf("failed to copy ui files: %w", err)
	}

	pluginsDir := filepath.Join(b.tmpDir, "vendor/github.com/codenomdev/viona-plugins/")
	localUIPluginDir := filepath.Join(localUIBuildDir, "src/plugins/")

	// copy plugins dir
	fmt.Printf("try to copy dir from %s to %s\n", pluginsDir, localUIPluginDir)

	// if plugins dir not exist means no plugins
	if !util.CheckDirExist(pluginsDir) {
		return nil
	}

	pluginsDirEntries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins dir: %w", err)
	}
	for _, entry := range pluginsDirEntries {
		if !entry.IsDir() {
			continue
		}
		sourcePluginDir := filepath.Join(pluginsDir, entry.Name())
		// check if plugin is a ui plugin
		packageJsonPath := filepath.Join(sourcePluginDir, "package.json")
		fmt.Printf("check if %s is a ui plugin\n", packageJsonPath)
		if !util.CheckFileExist(packageJsonPath) {
			continue
		}

		pnpmInstallCmd := b.newExecCmd("pnpm", "install")
		pnpmInstallCmd.Dir = sourcePluginDir
		if err = pnpmInstallCmd.Run(); err != nil {
			return fmt.Errorf("failed to install plugin dependencies: %w", err)
		}

		localPluginDir := filepath.Join(localUIPluginDir, entry.Name())
		fmt.Printf("try to copy dir from %s to %s\n", sourcePluginDir, localPluginDir)
		if err = copyDirEntries(os.DirFS(sourcePluginDir), ".", localPluginDir, "node_modules"); err != nil {
			return fmt.Errorf("failed to copy ui files: %w", err)
		}
	}
	formatUIPluginsDirName(localUIPluginDir)
	return nil
}

// downloadGoModFile run go mod commands to download dependencies
func downloadGoModFile(b *buildingMaterial) (err error) {
	codenomReplacement := b.moduleReplacement

	// If no replacement specified and current binary is v2+, auto-determine replacement.
	// This is needed because go mod tidy would otherwise resolve github.com/codenomdev/viona
	// to the latest v1.x version, causing v2+ features (e.g. AI/MCP) to disappear.
	if len(codenomReplacement) == 0 && b.originalCodenomInfo.Version != "" {
		ver, verErr := semver.NewVersion(strings.TrimPrefix(b.originalCodenomInfo.Version, "v"))
		if verErr == nil && ver.Major() >= 2 {
			codenomReplacement = fmt.Sprintf("github.com/codenomdev/viona@%s", b.originalCodenomInfo.Version)
		}
	}

	if len(codenomReplacement) > 0 {
		// For v2+ versioned module paths (e.g. github.com/codenomdev/viona@v2.0.0),
		// go mod tidy rejects the version because the module path lacks a /v2 suffix.
		// Work around this by cloning the repo locally and using a local path replacement.
		localPath, resolveErr := resolveModuleReplacement(codenomReplacement, b.tmpDir)
		if resolveErr != nil {
			return resolveErr
		}
		replacement := fmt.Sprintf("%s=%s", "github.com/codenomdev/viona", localPath)
		err = b.newExecCmd("go", "mod", "edit", "-replace", replacement).Run()
		if err != nil {
			return err
		}
	}

	err = b.newExecCmd("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}

	err = b.newExecCmd("go", "mod", "vendor").Run()
	if err != nil {
		return err
	}
	return
}

// createMainGoFile creates main.go file in tmp dir that content is mainGoTpl
func createMainGoFile(b *buildingMaterial) (err error) {
	fmt.Printf("[build] build dir: %s\n", b.tmpDir)
	err = util.CreateDirIfNotExist(b.tmpDir)
	if err != nil {
		return err
	}

	var (
		remotePlugins []string
	)
	for _, p := range b.plugins {
		remotePlugins = append(remotePlugins, versionedModulePath(p.Name, p.Version))
	}

	mainGoFile := &bytes.Buffer{}
	// tmpl, err := template.New("main").Parse(mainGoTpl)
	// if err != nil {
	// 	return err
	// }
	tmpl := template.Must(template.ParseFS(templates.TemplatesFS, "main_go.tmpl"))
	err = tmpl.Execute(mainGoFile, map[string]any{
		"remote_plugins": remotePlugins,
	})
	if err != nil {
		return err
	}

	err = util.WriteFile(filepath.Join(b.tmpDir, "main.go"), mainGoFile.String())
	if err != nil {
		return err
	}

	err = util.WriteFile(filepath.Join(b.tmpDir, "go.mod"), goModTpl)
	if err != nil {
		return err
	}

	for _, p := range b.plugins {
		// If user set a path, use it to replace the module with local path
		if len(p.Path) > 0 {
			var replacement string
			if len(p.Version) > 0 {
				replacement = fmt.Sprintf("%s@%s=%s", p.Name, p.Version, p.Path)
			} else {
				replacement = fmt.Sprintf("%s=%s", p.Name, p.Path)
			}
			err = b.newExecCmd("go", "mod", "edit", "-replace", replacement).Run()
		} else if len(p.Version) > 0 {
			// If user specify a version, use it to get specific version of the module
			err = b.newExecCmd("go", "get", fmt.Sprintf("%s@%s", p.Name, p.Version)).Run()
		}
		if err != nil {
			return err
		}
	}
	return
}

// resolveModuleReplacement resolves the CODENOM_MODULE value to a usable local path or
// remote replacement string. For v2+ versioned module paths (e.g. github.com/codenomdev/viona@v2.0.0),
// Go module system rejects the version because the module path has no /v2 suffix. In that case
// the repository is cloned locally and the local path is returned instead.
func resolveModuleReplacement(replacement, tmpDir string) (string, error) {
	// Local paths can be used as-is.
	if strings.HasPrefix(replacement, "/") || strings.HasPrefix(replacement, "./") || strings.HasPrefix(replacement, "../") {
		return replacement, nil
	}

	// Parse module@version format.
	moduleName, version, hasVersion := strings.Cut(replacement, "@")
	if !hasVersion {
		return replacement, nil
	}

	// Only handle v2+ versions on module paths without the /vN suffix.
	ver, err := semver.StrictNewVersion(strings.TrimPrefix(version, "v"))
	if err != nil || ver.Major() < 2 {
		return replacement, nil
	}
	if strings.HasSuffix(moduleName, fmt.Sprintf("/v%d", ver.Major())) {
		return replacement, nil
	}

	// Clone the repo to a local directory and return its path.
	gitURL := "https://" + moduleName
	tag := "v" + strings.TrimPrefix(version, "v")
	localPath := filepath.Join(filepath.Dir(tmpDir), fmt.Sprintf("codenom_src%s", strings.ReplaceAll(version, ".", "_")))

	if _, statErr := os.Stat(localPath); statErr == nil {
		fmt.Printf("[build] using cached local clone at %s\n", localPath)
		return localPath, nil
	}

	fmt.Printf("[build] v2+ module detected, cloning %s@%s to local path %s...\n", moduleName, version, localPath)
	cloneCmd := exec.Command("git", "clone", "--depth=1", "--branch="+tag, gitURL, localPath)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err = cloneCmd.Run(); err != nil {
		return "", fmt.Errorf(
			"failed to clone %s@%s: %w\nTip: set CODENOM_MODULE to a local checkout path instead, e.g. CODENOM_MODULE=/path/to/codenom",
			moduleName, version, err,
		)
	}

	fmt.Printf("[build] successfully cloned to %s\n", localPath)
	return localPath, nil
}

// movePluginToVendor move plugin to vendor dir
// Traverse the plugins, and if the plugin path is not github.com/codenomdev/viona-plugins, move the contents of the current plugin to the vendor/github.com/codenomdev/viona-plugins/ directory.
func movePluginToVendor(b *buildingMaterial) (err error) {
	pluginsDir := filepath.Join(b.tmpDir, "vendor/github.com/codenomdev/viona-plugins/")
	for _, p := range b.plugins {
		pluginDir := filepath.Join(b.tmpDir, "vendor/", p.Name)
		pluginName := filepath.Base(p.Name)
		if !strings.HasPrefix(p.Name, "github.com/codenomdev/viona-plugins/") {
			fmt.Printf("try to copy dir from %s to %s\n", pluginDir, filepath.Join(pluginsDir, pluginName))
			err = copyDirEntries(os.DirFS(pluginDir), ".", filepath.Join(pluginsDir, pluginName), "node_modules")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyDirEntries(sourceFs fs.FS, sourceDir, targetDir string, ignoreDir ...string) (err error) {
	err = util.CreateDirIfNotExist(targetDir)
	if err != nil {
		return err
	}
	ignoreThisDir := func(path string) bool {
		for _, s := range ignoreDir {
			if strings.HasPrefix(path, s) {
				return true
			}
			// Also ignore nested occurrences, e.g. src/plugins/foo/node_modules
			if strings.Contains(path, string(filepath.Separator)+s) {
				return true
			}
			if strings.Contains(path, "/"+s) {
				return true
			}
		}
		return false
	}

	err = fs.WalkDir(sourceFs, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignoreThisDir(path) {
			return nil
		}

		// Convert the path to use forward slashes, important because we use embedded FS which always uses forward slashes
		path = filepath.ToSlash(path)

		// Construct the absolute path for the source file/directory
		srcPath := filepath.Join(sourceDir, path)
		srcPath = filepath.ToSlash(srcPath)

		// Construct the absolute path for the destination file/directory
		dstPath := filepath.Join(targetDir, path)

		if d.IsDir() {
			// Create the directory in the destination
			err := os.MkdirAll(dstPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dstPath, err)
			}
		} else {
			// Open the source file
			srcFile, err := sourceFs.Open(srcPath)
			if err != nil {
				return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
			}
			defer srcFile.Close()

			// Create the destination file
			dstFile, err := os.Create(dstPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file %s: %w", dstPath, err)
			}
			defer dstFile.Close()

			// Copy the file contents
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return fmt.Errorf("failed to copy file contents from %s to %s: %w", srcPath, dstPath, err)
			}
		}

		return nil
	})

	return err
}

// format plugins dir name from dash to underline
func formatUIPluginsDirName(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("read ui plugins dir failed: [%s] %s\n", dirPath, err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() || !strings.Contains(entry.Name(), "-") {
			continue
		}
		newName := strings.ReplaceAll(entry.Name(), "-", "_")
		if err := os.Rename(filepath.Join(dirPath, entry.Name()), filepath.Join(dirPath, newName)); err != nil {
			fmt.Printf("rename ui plugins dir failed: [%s] %s\n", dirPath, err)
		} else {
			fmt.Printf("rename ui plugins dir success: [%s] -> [%s]\n", entry.Name(), newName)
		}
	}
}

// buildBinary build binary file
func buildBinary(b *buildingMaterial) (err error) {
	versionInfo := b.originalCodenomInfo
	cmdPkg := "github.com/codenomdev/viona/cmd"
	ldflags := fmt.Sprintf("-X %s.Version=%s -X %s.Revision=%s -X %s.Time=%s",
		cmdPkg, versionInfo.Version, cmdPkg, versionInfo.Revision, cmdPkg, versionInfo.Time)
	err = b.newExecCmd("go", "build",
		"-ldflags", ldflags, "-o", b.outputPath, ".").Run()
	if err != nil {
		return err
	}
	return
}

// cleanByproduct delete tmp dir
func cleanByproduct(b *buildingMaterial) (err error) {
	return os.RemoveAll(b.tmpDir)
}

func (b *buildingMaterial) newExecCmd(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	fmt.Println(cmd.Args)
	cmd.Dir = b.tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func versionedModulePath(modulePath, moduleVersion string) string {
	if moduleVersion == "" {
		return modulePath
	}
	ver, err := semver.StrictNewVersion(strings.TrimPrefix(moduleVersion, "v"))
	if err != nil {
		return modulePath
	}
	major := ver.Major()
	if major > 1 {
		modulePath += fmt.Sprintf("/v%d", major)
	}
	return path.Clean(modulePath)
}
