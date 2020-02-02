package cmd

import (
	"fmt"
	"github.com/l3uddz/wantarr/build"
	"github.com/l3uddz/wantarr/config"
	"github.com/l3uddz/wantarr/logger"
	pvrObj "github.com/l3uddz/wantarr/pvr"
	"github.com/l3uddz/wantarr/utils/paths"
	stringutils "github.com/l3uddz/wantarr/utils/strings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	// Global flags
	flagLogLevel     = 0
	flagConfigFolder = paths.GetCurrentBinaryPath()
	flagConfigFile   = "config.yaml"
	flagDatabaseFile = "vault.db"
	flagLogFile      = "activity.log"
	flagRefreshCache = false

	// Global vars
	pvrName      string
	lowerPvrName string
	pvrConfig    *config.Pvr
	pvr          pvrObj.Interface
	log          *logrus.Entry
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wantarr",
	Short: "A CLI application to search for wanted media files in the arr suite",
	Long: `A CLI application that can be used to search for wanted media files in the arr suite.

Allows searching for missing / wanted media files (episodes/movies).
It will monitor the queue and respect any limits set via the configuration file.
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Parse persistent flags
	rootCmd.PersistentFlags().StringVar(&flagConfigFolder, "config-dir", flagConfigFolder, "Config folder")
	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", flagConfigFile, "Config file")
	rootCmd.PersistentFlags().StringVarP(&flagDatabaseFile, "database", "d", flagDatabaseFile, "Database file")
	rootCmd.PersistentFlags().StringVarP(&flagLogFile, "log", "l", flagLogFile, "Log file")
	rootCmd.PersistentFlags().CountVarP(&flagLogLevel, "verbose", "v", "Verbose level")

}

func initConfig() {
	// Set core variables
	if !rootCmd.PersistentFlags().Changed("config") {
		flagConfigFile = filepath.Join(flagConfigFolder, flagConfigFile)
	}
	if !rootCmd.PersistentFlags().Changed("database") {
		flagDatabaseFile = filepath.Join(flagConfigFolder, flagDatabaseFile)
	}
	if !rootCmd.PersistentFlags().Changed("log") {
		flagLogFile = filepath.Join(flagConfigFolder, flagLogFile)
	}

	// Init Logging
	if err := logger.Init(flagLogLevel, flagLogFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("app")

	log.Infof("Using %s = %s (%s@%s)", stringutils.StringLeftJust("VERSION", " ", 10),
		build.Version, build.GitCommit, build.Timestamp)
	logger.ShowUsing()

	// Init Config
	if err := config.Init(flagConfigFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}
}

/* Private Helpers */

func pluckMediaItemIds(mediaItems []pvrObj.MediaItem) []int {
	var mediaItemIds []int

	for _, mediaItem := range mediaItems {
		mediaItemIds = append(mediaItemIds, mediaItem.ItemId)
	}

	return mediaItemIds
}
