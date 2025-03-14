package cmd

import (
	"fmt"
	"log"
	"os"

	serp "github.com/dneil5648/serpcli/serphandler"
	"github.com/spf13/cobra"
)

var searchEngine, ouputFile, query string

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&searchEngine, "engine", "google", "Search Engine of choice")
	getCmd.Flags().StringVar(&ouputFile, "output", "output.csv", "Output file path")
	getCmd.Flags().StringVar(&query, "query", "", "Insert your query here")
}

var getCmd = &cobra.Command{
	Use:   "get [query]",
	Short: "CLI Tool for retrieving SerpAPI results for a Specfied Search Engine",
	Long: `This tool will allow you to get all of the results from a Google search and write it to a CSV file. 
	
	Usage:
	  gdorks get [flags]
	
	Flags:
		  --apikey 	string		[Required] SerpAPI Api Key if Env variable not set
		  --engine 	string		[Optional] Search Engine to use
		  --output 	string		[Required] Output Filepath with file name
		  --query  	string		[Required] The query you are looking to run 
		  

	Example:
		gdorks get "example.com AND site:github.com" -s google.com -o "hello.csv"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check required parameters

		if query == "" {
			// If no --query flag, check for positional argument
			if len(args) > 0 {
				query = args[0]

			} else {
				log.Fatal("Search query is required. Use --query flag or provide as positional argument")
			}
		}
		key := os.Getenv("SERP_API_KEY")
		if key == "" {
			log.Fatal("Please store api key: export SERP_API_KEY=<your-api-key>")
		}
		fmt.Println(key)
		fmt.Printf("Running Query: %s\n", query)
		fmt.Printf("Using search engine: %s\n", searchEngine)
		fmt.Printf("Results will be saved to: %s\n", ouputFile)

		// Create handler
		handler := serp.CreateSerpHandler(
			key,
			searchEngine,
			ouputFile,
		)

		// Execute query and stream to CSV
		err := handler.Query(query)
		if err != nil {
			log.Fatalf("Error executing query: %v", err)
		}

		fmt.Println("Search completed successfully")
	},
}
