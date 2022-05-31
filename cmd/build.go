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
	"log"

	"github.com/evanoberholster/timezoneLookup"
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
)

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVar(&snappy, "snappy", true, "Use Snappy compression (true/false)")
	buildCmd.Flags().StringVar(&jsonFilename, "json", "combined-with-oceans.json", "GEOJSON Filename")
}

func build(cmd *cobra.Command, args []string) {
	if dbFilename == "" || jsonFilename == "" {
		log.Printf("Options:\n\t -snappy=true\t Use Snappy compression\n\t -json=filename\t GEOJSON filename \n\t -db=filename\t Database destination\n\t -type=boltdb\t Type of Storage (boltdb or memory) ")
	} else {
		tz := timezoneLookup.MemoryStorage(snappy, dbFilename)
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
}
