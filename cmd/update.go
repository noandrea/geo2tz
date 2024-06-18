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
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/noandrea/geo2tz/v2/helpers"
	"github.com/noandrea/geo2tz/v2/web"
	"github.com/spf13/cobra"
)

const (
	// geo data url
	LatestReleaseURL = "https://github.com/evansiroky/timezone-boundary-builder/releases/latest"
)

// updateCmd represents the build command
var updateCmd = &cobra.Command{
	Use:   "update VERSION",
	Short: "Download the timezone data from the latest release or a specific version",
	Example: `To build from the latest version:
geo2tz update latest

To build from a specific version:
geo2tz update 2023d
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		versionName := args[0]
		return update(versionName, web.Settings.Tz.DatabaseName)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&web.Settings.Tz.DatabaseName, "db", web.TZDBFile, "Destination database filename")
	updateCmd.Flags().StringVar(&web.Settings.Tz.VersionFile, "version-file", web.TZVersionFile, "Version file")
}

func update(versionName, targetFile string) (err error) {
	var release = web.NewTzRelease(versionName)
	// do we need the latest version?
	if versionName == "latest" {
		release, err = getLatest()
		if err != nil {
			return
		}
	}
	if err := fetchAndCacheFile(targetFile, release.GeoDataURL); err != nil {
		return err
	}
	helpers.SaveJSON(release, web.Settings.Tz.VersionFile)
	return
}

func fetchAndCacheFile(filename string, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	if n != resp.ContentLength {
		fmt.Println(n, resp.ContentLength)
	}
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
