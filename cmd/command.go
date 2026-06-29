package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/codenomdev/viona/internal/apps/cli"
	"github.com/codenomdev/viona/plugins"
	"github.com/spf13/cobra"
)

var (
	// Root command
	rootCmd = &cobra.Command{
		Use:   "codenom",
		Short: "The root base command(s).",
		Long:  "The main command(s). Use --help for list available commands.",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start API server",
		Run: func(cmd *cobra.Command, _ []string) {
			res, err := bootstrap(cmd.Context(), true)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer res.Cleanup()

			runApp(res.Context)
		},
	}

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build Codenom with plugins",
		Long:  `Build a new Codenom with plugins that you need`,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("try to build a new codenom with plugins:\n%s\n", strings.Join(buildWithPlugins, "\n"))
			err := cli.BuildNewCodenom(buildDir, buildOutput, buildWithPlugins, cli.OriginalCodenomInfo{
				Version:  Version,
				Revision: Revision,
				Time:     Time,
			})
			if err != nil {
				fmt.Printf("build failed %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("build new codenom successfully %s\n", buildOutput)
		},
	}

	migrateCmd = &cobra.Command{
		Use:   "db:migrate",
		Short: "The Database migration(s)",
		Run: func(cmd *cobra.Command, _ []string) {
			res, err := bootstrap(cmd.Context(), true)
			if err != nil {
				panic(err)
			}
			defer res.Cleanup()

			runDBMigrations(res.Context)
		},
	}

	seederCmd = &cobra.Command{
		Use:   "db:seed",
		Short: "The database seed",
		Run: func(cmd *cobra.Command, _ []string) {
			res, err := bootstrap(cmd.Context(), true)
			if err != nil {
				panic(err)
			}
			defer res.Cleanup()

			RunSeeders(
				res.Context,
				TableName,
				ViewAll,
			)
		},
	}

	makeMigrationCmd = cli.MakeGeneratorCmd(cli.GeneratorConfig{
		CommandName: "make:migration",
		Short:       "Make database migration",
		Long:        "Make new database migration. Use --help for list available commands.",
		Dir:         "internal/database/migrations",
		FileSuffix:  "table",
		Template:    "make_migration.tmpl",
		SuccessMsg:  "Database migration created:",
	})

	makeSeederCmd = cli.MakeGeneratorCmd(cli.GeneratorConfig{
		CommandName: "make:seed",
		Short:       "Make database seed",
		Long:        "Make new database seeder. Use --help for list available commands.",
		Dir:         "internal/database/seeders",
		FileSuffix:  "seed",
		Template:    "make_seeder.tmpl",
		SuccessMsg:  "Database seed created:",
	})

	pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "Print all plugins packed in the binary",
		Long:  `Print all plugins packed in the binary`,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Get all plugin:\n")
			_ = plugins.CallBase(func(base plugins.Base) error {
				info := base.Info()
				fmt.Printf("%s[%s] made by %s\n", info.SlugName, info.Version, info.Author)
				return nil
			})
		},
	}
)

func init() {
	rootCmd.Version = fmt.Sprintf("%s\nrevision: %s\nbuild time: %s", Version, Revision, Time)
	// rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "config.yaml", "config path, eg: c name_config.yaml.")
	buildCmd.Flags().StringSliceVarP(&buildWithPlugins, "with", "w", []string{}, "plugins needed to build")

	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "build output path")

	buildCmd.Flags().StringVarP(&buildDir, "build-dir", "b", "", "dir for build process")

	// Flag migration db
	migrateCmd.Flags().StringVarP(&TableName, "table", "t", "", "tablename for migration, eg: 2024_11_28_000000_user_table. For list tablename migration, please run with flag --view-all or -v.")
	migrateCmd.Flags().BoolVarP(&Rollback, "rollback", "r", false, "rollback database migration, undo the last migration.")
	// migrateCmd.Flags().BoolVarP(&ViewAll, "view-all", "v", false, "view all table availables. eg: 2024_11_28_000000_user_table.")
	// migrateCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "DRY-RUN mode — no changes will be applied.")

	seederCmd.Flags().StringVarP(&TableName, "table", "t", "", "tablename for seeder, eg: 2024_11_28_000101_user_roles_seed. For list tablename seed, please run with flag --view-all or -v.")
	seederCmd.Flags().BoolVarP(&ViewAll, "view-all", "v", false, "view all tablename availables. eg: 2024_11_28_000101_user_roles_seed.")

	rootCmd.PersistentFlags().StringVarP(
		&ConfigPath,
		"config",
		"c",
		"config.yaml",
		"config path",
	)

	// Register parent cmd
	for _, cmd := range []*cobra.Command{serveCmd, migrateCmd, makeMigrationCmd, makeSeederCmd, seederCmd, pluginCmd, buildCmd} {
		rootCmd.AddCommand(cmd)
	}
}
