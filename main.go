package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/KunalDuran/duranzcli/module/data"
	"github.com/KunalDuran/duranzcli/module/process"
)

func main() {

	startTime := time.Now()
	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := "password"

	// Connect the database
	_, err := data.InitDuranzDB(dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		log.Panic(err)
	}

	var processName, format, fileName string
	if len(os.Args) < 2 {
		// process_name: all, players, teams, match, matchstats, playerstats, delete
		// format: odi, test, t20, ipl
		// file_name: 123.json
		log.Fatal("Usage: duranzcli <process_name> <format> <file_name(when action=onefile)>")
	}

	// "args": ["all", "odi", "433606.json"]
	// "args": ["venues", "odi", "433606.json"]
	// "args": ["all", "ipl"]
	// "args": ["delete"]
	// "args": ["initial"]
	processName = strings.ToLower(os.Args[1])
	if len(os.Args) == 3 {
		format = strings.ToLower(os.Args[2])
	} else if len(os.Args) == 4 {
		format = strings.ToLower(os.Args[2])
		fileName = strings.ToLower(os.Args[3])
		if fileName == "" || !strings.Contains(fileName, ".json") {
			log.Fatal("Enter Valid JSON file")
		}
	}
	folderPath, ok := data.GamePath[strings.ToLower(format)]
	if !ok && processName != "delete" && processName != "initial" {
		log.Fatal("Format not correct: only odi, test, t20, ipl allowed")
	}
	// load cache layer
	data.PseudoCacheLayer(strings.ToLower(format))

	fmt.Println("Taking database from : ", process.DATASET_BASE)
	process.DATASET_BASE += (folderPath + process.SLASH)

	process.Activate(processName, fileName)

	fmt.Println(time.Since(startTime))

}
