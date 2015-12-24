package actions

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jmoiron/jsonq"
	"github.com/parnurzeal/gorequest"
	"github.com/redbooth/rbtools/validations"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Backup struct {
	Id           string
	State        string
	Progress     string
	Name         string
	Md5          string
	Download_url string
	Message      string
	Location     string
}

var (
	BackupIdRegexp = regexp.MustCompile(`backups\/(\d*)$`)
)

func decodeJSON(reader io.Reader) (jq *jsonq.JsonQuery, err error) {
	data := map[string]interface{}{}
	err = json.NewDecoder(reader).Decode(&data)

	if err != nil {
		log.Fatal("Error decoding backup response: %s", err)
	}

	jq = jsonq.NewQuery(data)
	return
}

func parseBackupCreationResponse(reader io.Reader) (backup *Backup, err error) {
	jq, err := decodeJSON(reader)

	if err != nil {
		log.Fatal("Error parsing backup response: %s", err)
	} else {
		location, _ := jq.String("location")
		backup = new(Backup)
		backup.Id = BackupIdRegexp.FindStringSubmatch(location)[1]
	}

	return
}

func logError(message string, err error) {
	if err != nil {
		log.Fatal("Error parsing %s : %s", message, err)
	}
}

func parseBackup(reader io.Reader) (backup *Backup, err error) {
	response, err := decodeJSON(reader)

	if err != nil {
		log.Fatal("Error parsing backup response: %s", err)
	} else {
		bkup, _ := response.Object("backup")
		jq := jsonq.NewQuery(bkup)
		backup = new(Backup)
		id, err := jq.Int("id")
		backup.Id = strconv.Itoa(id)
		logError("backup.Id", err)
		backup.State, err = jq.String("state")
		logError("backup.State", err)
		progress, err := jq.Int("progress")
		backup.Progress = strconv.Itoa(progress)
		logError("backup.Progress", err)
		backup.Name, err = jq.String("name")
		logError("backup.Name", err)
		if backup.State == "processed" {
			backup.Md5, err = jq.String("md5")
			logError("backup.Md5", err)
			backup.Download_url, err = response.String("download_url")
			logError("backup.Download_url", err)
		}

		fmt.Printf("info: creating backup with progress: %s (state: %s) \n", backup.Progress, backup.State)
	}

	return
}

func createHandler(ctx *cli.Context, callback func(backup *Backup)) func(gorequest.Response, string, []error) {
	return func(response gorequest.Response, body string, errs []error) {
		backup, err := parseBackupCreationResponse(strings.NewReader(body))

		if err != nil {
			log.Fatal(" Error parsing created backup response: %s", err)
			return
		}

		pollBackupStatus(ctx, backup, callback)(nil, "", nil)
	}
}

func debugRequest(data []byte, err error) {
	debug, cerr := strconv.ParseBool(os.Getenv("DEBUG"))
	if !debug || cerr != nil {
		return
	}

	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

func downloadBackup(ctx *cli.Context, backup *Backup) {
	var (
		host                = ctx.GlobalString("host")
		username            = ctx.GlobalString("username")
		password            = ctx.GlobalString("password")
		download_backup_url = fmt.Sprintf("%s/manager/backups/%s/download.json", host, backup.Id)
	)

	fmt.Printf("info: downloading backup... \n")
	client := &http.Client{}
	request, err := http.NewRequest("GET", download_backup_url, nil)
	request.SetBasicAuth(username, password)
	debugRequest(httputil.DumpRequestOut(request, true))

	response, err := client.Do(request)
	defer response.Body.Close()

	content_dispostion := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content_dispostion)
	if err != nil {
		log.Fatal(" Error parsing content disposition: %s", err)
	}

	// TODO: Make download path customizable
	download_path := fmt.Sprintf("/tmp/%s", params["filename"])

	out, err := os.Create(download_path)
	defer out.Close()

	debugRequest(httputil.DumpResponse(response, true))
	written, err := io.Copy(out, response.Body)

	if err != nil {
		log.Fatal(" Error downloading backup: %s (%s)", err, written)
	}
	fmt.Printf("info: Downloaded backup with to %s \n", download_path)
}

func pollBackupStatus(ctx *cli.Context, backup *Backup, callback func(backup *Backup)) func(gorequest.Response, string, []error) {
	return func(response gorequest.Response, body string, errs []error) {
		var (
			host           = ctx.GlobalString("host")
			username       = ctx.GlobalString("username")
			password       = ctx.GlobalString("password")
			get_backup_url = fmt.Sprintf("%s/manager/backups/%s.json", host, backup.Id)
			polledBackup   *Backup
		)

		if response != nil {
			poll, err := parseBackup(strings.NewReader(body))
			polledBackup = poll

			if err != nil {
				log.Fatal(" Error parsing backup poll response: %s", err)
			}
		}

		if polledBackup != nil && polledBackup.State == "processed" {
			downloadBackup(ctx, polledBackup)
		} else {
			time.Sleep(2)
			gorequest.New().SetBasicAuth(username, password).Get(get_backup_url).End(pollBackupStatus(ctx, backup, callback))
		}
	}
}

func onBackupDownload(backup *Backup) {
	fmt.Printf("info: downloaded backup with id %s \n", backup.Id)
}

func CreateBackupAction(ctx *cli.Context) {

	var (
		host            = ctx.GlobalString("host")
		username        = ctx.GlobalString("username")
		password        = ctx.GlobalString("password")
		post_backup_url = fmt.Sprintf("%s/manager/backups.json", host)
	)

	validations.ValidateRequiredArguments(ctx)

	gorequest.New().SetBasicAuth(username, password).Post(post_backup_url).End(createHandler(ctx, onBackupDownload))
}

func DeleteBackupAction(ctx *cli.Context) {

	var (
		host              = ctx.GlobalString("host")
		username          = ctx.GlobalString("username")
		password          = ctx.GlobalString("password")
		backup_id         = ctx.Args().First()
		delete_backup_url = fmt.Sprintf("%s/manager/backups/%s.json", host, backup_id)
	)

	validations.ValidateRequiredArguments(ctx)

	response, _, err := gorequest.New().SetBasicAuth(username, password).Delete(delete_backup_url).End()

	if err != nil {
		log.Fatal("error: Deleting backup: %s (%s)", backup_id, err)
	}

	if response.StatusCode != 200 {
		log.Fatal("error: Deleting backup: %s (status: %s)", backup_id, response.Status)
	}

	fmt.Printf("Successfully deleted backup with id %s \n", backup_id)
}
