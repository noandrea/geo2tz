/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/evanoberholster/timezoneLookup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the location database",
	Long:  `The commands replicates the functionality of the evanoberholster/timezoneLookup timezone command`,
	Run:   build,
}

var (
	snappy       bool
	jsonFilename string
	dbFilename   string
	storageType  string
	encoding     string
)

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVar(&snappy, "snappy", true, "Use Snappy compression (true/false)")
	buildCmd.Flags().StringVar(&jsonFilename, "json", "combined-with-oceans.json", "GEOJSON Filename")
	buildCmd.Flags().StringVar(&dbFilename, "db", "timezone", "Destination database filename")
	buildCmd.Flags().StringVar(&storageType, "type", "boltdb", "Storage: boltdb or memory")
	buildCmd.Flags().StringVar(&encoding, "encoding", "msgpack", "BoltDB encoding type: json or msgpack")
}

func build(cmd *cobra.Command, args []string) {
	if dbFilename == "" || jsonFilename == "" {
		log.Println("Options:\n\t -snappy=true\t Use Snappy compression\n\t -json=filename\t GEOJSON filename \n\t -db=filename\t Database destination\n\t -type=boltdb\t Type of Storage (boltdb or memory) ")
		return
	}
	var tz timezoneLookup.TimezoneInterface
	if storageType == "memory" {
		tz = timezoneLookup.MemoryStorage(snappy, dbFilename)
	} else if storageType == "boltdb" {
		encodingTz, err := timezoneLookup.EncodingFromString(encoding)
		if err != nil {
			log.Errorln("invalid encoding", err)
			return
		}
		tz = timezoneLookup.BoltdbStorage(snappy, dbFilename, encodingTz)
	} else {
		log.Println("\"-db\" No database type specified")
		return
	}

	if jsonFilename != "" {
		err := tz.CreateTimezones(jsonFilename)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println("\"-json\" No GeoJSON source file specified")
		return
	}
	tz.Close()
}
