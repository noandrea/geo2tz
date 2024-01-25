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
	"net/http"
	"strings"
	"time"

	"github.com/evanoberholster/timezoneLookup/v2"
	"github.com/noandrea/geo2tz/v2/helpers"
	"github.com/noandrea/geo2tz/v2/web"
	"github.com/spf13/cobra"
)

const (
	// geo data url
	LatestReleaseURL = "https://github.com/evansiroky/timezone-boundary-builder/releases/latest"
	TZZipFile        = "tzdata/timezone.zip"
)

var (
	// cli parameters.
	cacheFile  string
	geoDataURL string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Build the location database for a specific version",
	Example: `geo2tz build 2023d`,
	Long:    `The commands replicates the functionality of the evanoberholster/timezoneLookup timezone command`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tzVersion := web.NewTzRelease(args[0])
		return update(tzVersion, cacheFile, web.Settings.Tz.DatabaseName)
	},
}

// buildCmd represents the build command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the location database by downloading the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		tzVersion, err := getLatest()
		if err != nil {
			return err
		}
		return update(tzVersion, cacheFile, web.Settings.Tz.DatabaseName)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&cacheFile, "cache", TZZipFile, "Temporary cache filename")
	buildCmd.Flags().StringVar(&web.Settings.Tz.DatabaseName, "db", web.TZDBFile, "Destination database filename")
	buildCmd.Flags().StringVar(&web.Settings.Tz.VersionFile, "version-file", web.TZVersionFile, "Version file")
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&geoDataURL, "geo-data-url", "", "URL to download geo data from")
	updateCmd.Flags().StringVar(&cacheFile, "cache", TZZipFile, "Temporary cache filename")
	updateCmd.Flags().StringVar(&web.Settings.Tz.DatabaseName, "db", web.TZDBFile, "Destination database filename")
	updateCmd.Flags().StringVar(&web.Settings.Tz.VersionFile, "version-file", web.TZVersionFile, "Version file")
}

func update(release web.TzRelease, zipFile, dbFile string) (err error) {
	// remove old file
	if err = helpers.DeleteQuietly(zipFile, dbFile); err != nil {
		return
	}

	var (
		tzc   timezoneLookup.Timezonecache
		total int
	)
	fmt.Printf("building database %s v%s from %s\n", dbFile, release.Version, release.GeoDataURL)
	if err = timezoneLookup.ImportZipFile(zipFile, release.GeoDataURL, func(tz timezoneLookup.Timezone) error {
		total += len(tz.Polygons)
		tzc.AddTimezone(tz)
		return nil
	}); err != nil {
		return
	}
	if err = tzc.Save(dbFile); err != nil {
		return
	}
	tzc.Close()
	fmt.Println("polygons added:", total)
	fmt.Println("saved timezone data to:", dbFile)

	// remove tmp file
	if err = helpers.DeleteQuietly(cacheFile); err != nil {
		return
	}
	err = helpers.SaveJSON(release, web.Settings.Tz.VersionFile)
	return
}

func getLatest() (web.TzRelease, error) {
	// create http client
	client := &http.Client{
		Timeout: 1 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// don't follow redirects
			return http.ErrUseLastResponse
		},
	}
	r, err := client.Head(LatestReleaseURL)
	if err != nil {
		err = fmt.Errorf("failed to get release url: %w", err)
		return web.TzRelease{}, err
	}
	defer r.Body.Close()
	// get the tag name
	releaseURL, err := r.Location()
	if err != nil {
		err = fmt.Errorf("failed to get release url: %w", err)
		return web.TzRelease{}, err
	}
	v := web.NewTzRelease(releaseURL.Path[strings.LastIndex(releaseURL.Path, "/")+1:])
	return v, nil
}
