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
	"fmt"
	"os"

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
	// geo data url
	GeoDataURL = "https://github.com/evansiroky/timezone-boundary-builder/releases/download/2022b/timezones-with-oceans.geojson.zip"
	// cli parameters.
	snappy       bool
	jsonFilename string
	dbFilename   string
)

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&dbFilename, "db", "timezone", "Destination database filename")
	buildCmd.Flags().BoolVar(&snappy, "snappy", true, "Use Snappy compression (true/false)")
	buildCmd.Flags().StringVar(&jsonFilename, "json", "combined-with-oceans.json", "GEOJSON Filename")
}

func build(*cobra.Command, []string) {
	if dbFilename == "" || jsonFilename == "" {
		fmt.Printf(`Options:
  -snappy=true   Use Snappy compression
  -json=filename GEOJSON filename
  -db=filename   Database destination
`)
		return
	}

	tz := timezoneLookup.MemoryStorage(snappy, dbFilename)

	if !fileExists(jsonFilename) {
		fmt.Printf("json file %v does not exists, will try to download from the source", jsonFilename)
		return
	}

	if jsonFilename != "" {
		err := tz.CreateTimezones(jsonFilename)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println(`"--json" No GeoJSON source file specified`)
		return
	}

	tz.Close()

}

func fileExists(filePath string) bool {
	f, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}
