package actions

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/parnurzeal/gorequest"
	"github.com/redbooth/rbtools/validations"
	"io"
	"regexp"
	"strings"
	// "os"
)

type BackupResponse struct {
	Message  string `json:"message"`
	Location string `json:"location"`
}

type Backup struct {
	Id    string
	State string
}

var (
	BackupIdRegexp = regexp.MustCompile(`backups\/(\d*)$`)
)

func decodeJSON(reader io.Reader) (response *BackupResponse, err error) {
	response = new(BackupResponse)
	err = json.NewDecoder(reader).Decode(response)
	return
}

func parseBackup(reader io.Reader) (backup *Backup, err error) {
	backupResponse, err := decodeJSON(reader)

	if err != nil {
		fmt.Printf(" boom ====> %s", err)
	} else {
		backup = new(Backup)
		backup.Id = BackupIdRegexp.FindStringSubmatch(backupResponse.Location)[1]
	}

	return
}

func createHandler(response gorequest.Response, body string, errs []error) {
	println("POST: ", body)
	backup, err := parseBackup(strings.NewReader(body))

	if err != nil {
		fmt.Printf(" boom ====> %s", err)
	} else {
		fmt.Printf("====> created backup with id %s \n", backup.Id)
	}
}

func CreateBackupAction(ctx *cli.Context) {

	var (
		// download_path   = ctx.String("download-to")
		host            = ctx.GlobalString("host")
		username        = ctx.GlobalString("username")
		password        = ctx.GlobalString("password")
		post_backup_url = fmt.Sprintf("%s/manager/backups.json", host)
		// logger          = log.New(os.Stdout, "logger: ", log.Lshortfile)
	)

	validations.ValidateRequiredArguments(ctx)

	gorequest.New().SetBasicAuth(username, password).Post(post_backup_url).End(createHandler)
}
