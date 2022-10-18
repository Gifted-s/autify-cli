/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"autify/v1/models"
	"autify/v1/pkg"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Download web page locally",
	Long:  "With fetch, you can download web pages from command line and also view web page metadata",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if a metadata flag was passed
		pageUrl, err := cmd.Flags().GetString("metadata")
		if err != nil {
			color.Red("Flag Error:  %s", err.Error())
			os.Exit(1)
		}
		// If the --metadata flag was passed and the page url was passed along with it, then attempt to fetch metadata locally
		if pageUrl != "" {
			// validate URL
			_, err = url.ParseRequestURI(pageUrl)
			if err != nil {
				color.Red("Error: invalid url format:  %s", err.Error())
				os.Exit(1)
			}
			start := time.Now()
			// Fetch metadata
			meta, err := webpage.GetPageMeta(pageUrl)
			if err != nil {
				color.Red("fetch error: %s", err.Error())
				os.Exit(1)
			}
			// show metadata in console
			customMetaDisplay(meta)
			fmt.Printf("Execution time: %f seconds", time.Since(start).Seconds())
			return
		}
		// If the metadata flag was not passed, then attempt to download the pages from the internet
		urls := args
		if len(urls) == 0 {
			color.Red("Error: No Url passed")
			os.Exit(1)
		}
		// Validate URLs
		for _, u := range urls {
			_, err := url.ParseRequestURI(u)
			if err != nil {
				color.Red("Error: invalid url format:  %s", err.Error())
				os.Exit(1)
			}
		}
		start := time.Now()
		// Download starts
		err = webpage.DownloadPages(urls)
		if err != nil {
			color.Red("download error %s", err.Error())
		}
		color.Green("Download successful: %d page(s) downloaded", len(urls))
		fmt.Printf("Execution time: %f seconds", time.Since(start).Seconds())

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := webpage.Setup()
	if err != nil {
		os.Exit(1)
	}
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("metadata", "m", "", "URL of webpage to fetch")
}

func customMetaDisplay(meta models.WebPage) {
	color.Green("site: %s", meta.Site)
	color.Green("num_links: %d", meta.Num_Links)
	color.Green("images: %d", meta.Images)
	color.Green("last_fetch: %s", meta.Last_Fetch)
}
