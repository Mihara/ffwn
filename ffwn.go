package main

import (
	"encoding/xml"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gocarina/gocsv"
	"github.com/integrii/flaggy"
	"github.com/jhillyerd/enmime"
	"github.com/la5nta/wl2k-go/mailbox"
)

type Checkin struct {
	Callsign   string `csv:"callsign"`
	Name       string `csv:"name"`
	City       string `csv:"city"`
	Region     string `csv:"region"`
	Country    string `csv:"country"`
	Username   string `csv:"username"`
	Connection string `csv:"connection"`
	Question   string `csv:"question"`
}

var rmsCmd *flaggy.Subcommand
var patCmd *flaggy.Subcommand
var exportFile string
var outputFile string
var patCallsign string
var patDir string

func init() {

	outputFile = "output.csv"
	patDir = filepath.Join(xdg.DataHome, "pat", "mailbox")

	flaggy.SetName("ffwn-checkout")
	flaggy.SetDescription("Process incoming Winlink mail and collate data from FFWN check-in mails.")
	flaggy.DefaultParser.AdditionalHelpAppend = `
Copyright © 2023 Eugene Medvedev (R2AZE).
See the source code at: https://github.com/Mihara/ffwn
Released under the terms of WTFPL2 license.`

	flaggy.String(&outputFile, "o", "output", "Name of output file.")

	rmsCmd = flaggy.NewSubcommand("rms")
	rmsCmd.Description = "Process an XML message export from Winlink Express."
	rmsCmd.AddPositionalValue(&exportFile, "FILE.xml", 1, true,
		"A message export from your Winlink Express mailbox")

	patCmd = flaggy.NewSubcommand("pat")
	patCmd.Description = "Process a Pat mailbox."
	patCmd.AddPositionalValue(&patCallsign, "CALLSIGN", 1, true,
		"The callsign you were using to receive check-in messages.")

	patCmd.String(&patDir, "p", "pat_dir", "Location of Pat mail directory.")

	flaggy.AttachSubcommand(rmsCmd, 1)
	flaggy.AttachSubcommand(patCmd, 1)
	flaggy.Parse()

}

func MessageBody(txt string) (Checkin, error) {
	var result Checkin

	lines := strings.Split(txt, "\n")
	var nonEmptyLines []string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim != "" {
			nonEmptyLines = append(nonEmptyLines, trim)
		}
	}
	if len(nonEmptyLines) < 2 {
		return result, errors.New("message does not contain two lines")
	}
	cols := strings.Split(nonEmptyLines[0], ",")
	if len(cols) != 7 {
		return result, errors.New("missing columns in first line")
	}
	result = Checkin{
		Callsign:   strings.TrimSpace(cols[0]),
		Name:       strings.TrimSpace(cols[1]),
		City:       strings.TrimSpace(cols[2]),
		Region:     strings.TrimSpace(cols[3]),
		Country:    strings.TrimSpace(cols[4]),
		Username:   strings.TrimSpace(cols[5]),
		Connection: strings.TrimSpace(cols[6]),
		Question:   strings.TrimSpace(nonEmptyLines[1]),
	}

	return result, nil
}

func saveCSV(filename string, records []Checkin) {
	output, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("Could not open output file ", filename)
	}
	defer output.Close()

	if err := gocsv.MarshalFile(&records, output); err != nil {
		log.Fatal("Could not save data to output file ", filename)
	}
	log.Print("Done, see ", filename)
}

func main() {

	var records []Checkin

	if rmsCmd.Used {

		exportData, err := os.ReadFile(exportFile)
		if err != nil {
			log.Fatal("Could not open file ", exportFile)
		}

		type Message struct {
			XMLName  xml.Name `xml:"message"`
			Id       string   `xml:"id"`
			Subject  string   `xml:"subject"`
			MimeData string   `xml:"mime"`
		}

		var exportTree struct {
			XMLName  xml.Name  `xml:"Winlink_Express_message_export"`
			Messages []Message `xml:"message_list>message"`
		}

		if err := xml.Unmarshal(exportData, &exportTree); err != nil {
			log.Fatal(err)
		}

		for _, message := range exportTree.Messages {
			if strings.TrimSpace(message.Subject) == "FFWN" {
				env, _ := enmime.ReadEnvelope(strings.NewReader(message.MimeData))

				if record, err := MessageBody(env.Text); err == nil {
					records = append(records, record)
				} else {
					log.Print("Could not make sense of message ", message.Id)
				}
			}
		}

		saveCSV(outputFile, records)

	} else if patCmd.Used {
		box := filepath.Join(patDir, patCallsign, mailbox.DIR_INBOX)
		mail, err := mailbox.LoadMessageDir(box)
		if err != nil {
			log.Fatal("Could not open mailbox at ", box)
		}
		for _, message := range mail {
			if strings.TrimSpace(message.Subject()) == "FFWN" {
				body, err := message.Body()
				if err != nil {
					log.Print("Disembodied message ", message.MID())
					continue
				}
				if record, err := MessageBody(body); err == nil {
					records = append(records, record)
				} else {
					log.Print("Could not make sense of message ", message.MID())
				}
			}
		}

		saveCSV(outputFile, records)

	} else {
		flaggy.ShowHelp("")
	}

}
