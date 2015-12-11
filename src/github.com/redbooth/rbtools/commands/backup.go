package commands

import (
	"github.com/codegangsta/cli"
	"github.com/redbooth/rbtools/actions"
)

const (
	DEFAULT_DOWNLOAD_PATH = "/tmp"
	DEFAULT_PRUNE_DAYS    = "3"
	DEFAULT_SCHEDULE      = "daily"
)

var Backup = cli.Command{
	Name:  "backup",
	Usage: "Perform backup-related actions against Redbooth instance",
	Subcommands: []cli.Command{
		backupCreate,
		backupDelete,
		backupPrune,
		backupSchedule,
	},
}

/********************
 *   CREATE         *
 ********************/
var download_path string
var backupCreate = cli.Command{
	Name:   "create",
	Usage:  "Create a backup and download it",
	Action: createAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "download-to",
			Value:       DEFAULT_DOWNLOAD_PATH,
			Usage:       "Path where backups will be downloaded to...",
			Destination: &download_path},
	},
}

func createAction(ctx *cli.Context) {
	var (
		download_path = ctx.String("download-to")
		host          = ctx.GlobalString("host")
		username      = ctx.GlobalString("username")
		password      = ctx.GlobalString("password")
	)

	actions.ValidateRequiredArguments(ctx)
	println("Created backup...downloading to:", download_path, host, username, password)
}

/********************
 *   DELETE         *
 ********************/
var backupDelete = cli.Command{
	Name:   "delete",
	Usage:  "Delete a backup",
	Action: deleteAction,
}

func deleteAction(ctx *cli.Context) {
	var (
		backup_id = ctx.Args().First()
	)

	actions.ValidateRequiredArguments(ctx)
	println("Deleting backup:", backup_id)
}

/********************
 *   PRUNE          *
 ********************/
var older_than string
var backupPrune = cli.Command{
	Name:   "prune",
	Usage:  "Prune backups older than `num_days`",
	Action: pruneAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "days",
			Value:       DEFAULT_PRUNE_DAYS,
			Usage:       "Prune backups older than `days`",
			Destination: &older_than},
	},
}

func pruneAction(ctx *cli.Context) {
	var (
		older_than = ctx.String("days")
	)

	actions.ValidateRequiredArguments(ctx)
	println("pruning all backups older than: ", older_than)
}

/********************
 *   PRUNE          *
 ********************/
var schedule string
var backupSchedule = cli.Command{
	Name:   "schedule",
	Usage:  "Schedule creation of backups regularly (implies daemon mode)",
	Action: scheduleAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "schedule",
			Value:       DEFAULT_SCHEDULE,
			Usage:       "Some schedule string",
			Destination: &schedule},
	},
}

func scheduleAction(ctx *cli.Context) {
	var (
		schedule = ctx.String("schedule")
	)

	actions.ValidateRequiredArguments(ctx)
	println("Scheduling backups ", schedule)
}
