package serphandler

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	g "github.com/serpapi/google-search-results-golang"
)

type SerpHandler struct {
	ApiKey       string
	SearchEngine string
	OutputFile   string
	MaxPages     int
	OutputExt    string
	Concurrency  int // Number of concurrent requests
}

// Create a new SerpHandler
func CreateSerpHandler(apiKey string, searchEngine string, outputFile string) *SerpHandler {
	ext := filepath.Ext(outputFile)
	return &SerpHandler{
		ApiKey:       apiKey,
		SearchEngine: searchEngine,
		OutputFile:   outputFile,
		MaxPages:     100, // Default to 100 pages max
		OutputExt:    ext,
		Concurrency:  3, // Default concurrency
	}
}

// SerpAPIResponse represents the top-level response from SerpAPI
type SerpAPIResponse struct {
	SearchMetadata    SearchMetadata             `json:"search_metadata"`
	SearchParameters  SearchParameters           `json:"search_parameters"`
	SearchInformation SearchInformation          `json:"search_information"`
	OrganicResults    []OrganicResult            `json:"organic_results"`
	Pagination        Pagination                 `json:"pagination"`
	SerpAPIPagination map[string]json.RawMessage `json:"serpapi_pagination"`
}

// SearchMetadata contains metadata about the search
type SearchMetadata struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	JsonEndpoint   string  `json:"json_endpoint"`
	CreatedAt      string  `json:"created_at"`
	ProcessedAt    string  `json:"processed_at"`
	TotalTimeTaken float64 `json:"total_time_taken"`
}

// SearchParameters contains the parameters used in the search
type SearchParameters struct {
	Engine            string `json:"engine"`
	Query             string `json:"q"`
	LocationRequested string `json:"location_requested,omitempty"`
	LocationUsed      string `json:"location_used,omitempty"`
	GoogleDomain      string `json:"google_domain"`
}

// SearchInformation contains summary information about the search
type SearchInformation struct {
	OrganicResultsState string  `json:"organic_results_state"`
	QueryDisplayed      string  `json:"query_displayed"`
	TotalResults        int     `json:"total_results"`
	TimeTakenDisplayed  float64 `json:"time_taken_displayed"`
}

// OrganicResult represents a single search result
type OrganicResult struct {
	Position         int    `json:"position"`
	Title            string `json:"title"`
	Link             string `json:"link"`
	DisplayedLink    string `json:"displayed_link"`
	Snippet          string `json:"snippet"`
	CachedPageLink   string `json:"cached_page_link,omitempty"`
	RelatedPagesLink string `json:"related_pages_link,omitempty"`
}

// SimplifiedResult contains just the fields we want to extract
type SimplifiedResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// Pagination contains pagination information
type Pagination struct {
	Current    int               `json:"current"`
	Next       string            `json:"next,omitempty"`
	OtherPages map[string]string `json:"other_pages,omitempty"`
}

func (s *SerpHandler) fetchPage(ctx context.Context, query string, ch chan<- map[string]interface{}, errCh chan<- error, doneCh chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		// Signal that we're done fetching
		// fmt.Println("Signaling fetch is done...")
		doneCh <- struct{}{}
		// fmt.Println("FETCH DONE signal sent successfully")
	}()

	parameter := map[string]string{
		"engine": s.SearchEngine,
		"q":      query,
	}

	consecutiveEmptyResults := 0
	maxConsecutiveEmptyResults := 3
	fmt.Println("Fetching Pages...")

	for page := 0; page < s.MaxPages; page++ {
		select {
		case <-ctx.Done():
			return
		default:
			// Only add start parameter for pages after the first one
			if page > 0 {
				parameter["start"] = fmt.Sprintf("%d", page*10)
			}

			search := g.NewGoogleSearch(parameter, s.ApiKey)
			result, err := search.GetJSON()
			if err != nil {
				errCh <- fmt.Errorf("error fetching page %d: %w", page, err)
				return
			}

			// Check if organic results are empty
			organicResults, hasOrganic := result["organic_results"].([]interface{})
			if !hasOrganic || len(organicResults) == 0 {
				consecutiveEmptyResults++
				if consecutiveEmptyResults >= maxConsecutiveEmptyResults {
					errCh <- fmt.Errorf("reached maximum consecutive empty results (%d). stopping additional requests", maxConsecutiveEmptyResults)
					return
				}
			} else {
				// Reset counter if we found results
				consecutiveEmptyResults = 0
			}

			// Send the result to the channel
			ch <- result

			// Check if we've reached the last page
			serpPagination, ok := result["serpapi_pagination"].(map[string]interface{})
			if !ok {
				return // Cannot get pagination, stop here
			}

			// Check if there's a next page
			next, hasNext := serpPagination["next"]
			if !hasNext || next == nil || next == "" {
				return // No more pages
			}

			// Add a small delay to avoid rate limiting
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (s *SerpHandler) writeData(ctx context.Context, ch <-chan map[string]interface{}, errCh chan<- error, doneCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Writing Pages...")
	f, err := os.OpenFile(s.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errCh <- fmt.Errorf("failed to open output file: %w", err)
		return
	}
	defer f.Close()

	// Create CSV writer
	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Check if file is empty, write header if needed
	fileInfo, err := f.Stat()
	if err != nil {
		errCh <- fmt.Errorf("failed to get file stats: %w", err)
		return
	}

	if fileInfo.Size() == 0 {
		// Write headers
		if err := writer.Write([]string{"title", "link", "snippet", "date"}); err != nil {
			errCh <- fmt.Errorf("error writing CSV header: %w", err)
			return
		}
		writer.Flush()
	}

	// Process results from the channel
	fetchDone := false
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done in writeData")
			return
		case <-doneCh:
			// All fetches are done
			fmt.Println("Received FETCH DONE signal in writeData")
			fetchDone = true
			// Don't return yet, still need to process any remaining results in the channel
		case result, ok := <-ch:
			if !ok {
				// Channel is closed, we're done
				fmt.Println("Result channel closed in writeData")
				return
			}
			s.processResult(result, writer, errCh)
			writer.Flush() // Flush after each result for safety
		default:
			// No more results in channel right now
			if fetchDone {
				// If fetching is done and no more results, we can exit
				fmt.Println("No more results and fetch is done, exiting writeData")
				return
			}
			// Small sleep to avoid busy waiting
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Helper method to process a single result
func (s *SerpHandler) processResult(result map[string]interface{}, writer *csv.Writer, errCh chan<- error) {
	// Process organic results
	organicResults, ok := result["organic_results"].([]interface{})
	if !ok {
		errCh <- fmt.Errorf("failed to parse organic_results")
		return
	}

	for _, r := range organicResults {
		org, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		title, _ := org["title"].(string)
		link, _ := org["link"].(string)
		date, _ := org["date"].(string)
		snippet, _ := org["snippet"].(string)

		record := []string{title, link, snippet, date}
		if err := writer.Write(record); err != nil {
			errCh <- fmt.Errorf("error writing to CSV: %w", err)
			continue
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		errCh <- fmt.Errorf("error flushing CSV data: %w", err)
	}
}

// Query performs a search query and writes results to the output file
func (s *SerpHandler) Query(q string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	wg := &sync.WaitGroup{}
	resultCh := make(chan map[string]interface{}, s.Concurrency)
	errCh := make(chan error, s.Concurrency)
	doneCh := make(chan struct{}, 1) // Buffer of 1 so fetch can send without blocking

	fmt.Println("Starting search query:", q)

	// Start fetcher(s)
	wg.Add(1)
	go s.fetchPage(ctx, q, resultCh, errCh, doneCh, wg)

	// Start the writer
	wg.Add(1)
	go s.writeData(ctx, resultCh, errCh, doneCh, wg)

	// Error handling goroutine
	var queryErr error
	errorDone := make(chan struct{})

	go func() {
		select {
		case err := <-errCh:
			queryErr = err
			fmt.Println("Error received:", err)
			cancel() // Cancel operation on first error
		case <-ctx.Done():
			fmt.Println("Context done in error handler")
			// Context timeout or cancel
		}
		close(errorDone)
	}()

	// Wait for everything to complete
	// fmt.Println("Waiting for all goroutines to complete...")
	wg.Wait()
	// fmt.Println("All goroutines completed")

	// Clean up channels - important order to avoid hanging
	close(resultCh)
	// fmt.Println("Result channel closed")

	// Cancel context to ensure error goroutine completes
	cancel()
	// fmt.Println("Context cancelled")

	// Now safe to wait for error handler
	select {
	case <-errorDone:
		// fmt.Println("Error handler completed")
	case <-time.After(500 * time.Millisecond):
		// fmt.Println("Timed out waiting for error handler")
	}

	close(errCh)
	// fmt.Println("Error channel closed")
	// fmt.Println("Query completed")

	return queryErr
}
