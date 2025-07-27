package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

// Part is a simplified representation of a bike part for the frontend
type Part struct {
	ID          string   `json:"id"`
	Brand       string   `json:"brand"`
	Model       string   `json:"model"`
	Category    string   `json:"category"`
	SubCategory string   `json:"sub_category"`
	Price       float64  `json:"price"`
	MSRP        float64  `json:"msrp,omitempty"`
	Discount    float64  `json:"discount,omitempty"`
	Currency    string   `json:"currency"`
	InStock     bool     `json:"in_stock"`
	Description string   `json:"description"`
	Images      []string `json:"images,omitempty"`
	URL         string   `json:"url"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Global variables
var (
	// API base URL
	apiBaseURL = "http://localhost:8080/api/v1"

	// DOM elements
	searchButton     js.Value
	searchInput      js.Value
	categoryFilter   js.Value
	resultsContainer js.Value
)

func main() {
	// Set up channel to keep the program running
	c := make(chan struct{})

	fmt.Println("WebAssembly module initialized")

	// Get DOM elements
	document := js.Global().Get("document")
	searchButton = document.Call("getElementById", "search-button")
	searchInput = document.Call("getElementById", "search-input")
	categoryFilter = document.Call("getElementById", "category-filter")
	resultsContainer = document.Call("getElementById", "results-container")

	// Register event listeners
	searchButton.Call("addEventListener", "click", js.FuncOf(handleSearch))
	searchInput.Call("addEventListener", "keypress", js.FuncOf(handleSearchKeyPress))

	// Load initial results
	go fetchParts("", "")

	// Keep the program running
	<-c
}

// handleSearch is called when the search button is clicked
func handleSearch(this js.Value, args []js.Value) interface{} {
	// Prevent default form submission
	if len(args) > 0 {
		args[0].Call("preventDefault")
	}

	// Get search query and category
	query := searchInput.Get("value").String()
	category := categoryFilter.Get("value").String()

	// Fetch parts
	go fetchParts(query, category)

	return nil
}

// handleSearchKeyPress is called when a key is pressed in the search input
func handleSearchKeyPress(this js.Value, args []js.Value) interface{} {
	// Check if Enter key was pressed
	event := args[0]
	if event.Get("key").String() == "Enter" {
		// Prevent default form submission
		event.Call("preventDefault")

		// Trigger search
		handleSearch(this, args)
	}

	return nil
}

// fetchParts fetches parts from the API
func fetchParts(query, category string) {
	// Show loading indicator
	showLoading()

	// Build API URL
	url := fmt.Sprintf("%s/parts/search?q=%s", apiBaseURL, query)
	if category != "" {
		url += fmt.Sprintf("&category=%s", category)
	}

	// Create fetch promise
	promise := js.Global().Call("fetch", url)

	// Handle response
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		response := args[0]

		// Check if response is ok
		if !response.Get("ok").Bool() {
			// Handle HTTP error
			errorMsg := fmt.Sprintf("Error: %s (%d)", response.Get("statusText").String(), response.Get("status").Int())
			showError(errorMsg)
			return nil
		}

		// Parse response as JSON
		jsonPromise := response.Call("json")
		jsonPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Handle parsed JSON
			result := args[0]

			// Check for API error response
			if !result.Get("error").IsUndefined() {
				errorMsg := result.Get("message").String()
				showError(errorMsg)
				return nil
			}

			// Convert JS array to Go slice
			length := result.Length()
			parts := make([]Part, length)
			for i := 0; i < length; i++ {
				item := result.Index(i)
				partJSON := js.Global().Get("JSON").Call("stringify", item).String()
				var part Part
				json.Unmarshal([]byte(partJSON), &part)
				parts[i] = part
			}

			// Display parts
			displayParts(parts)
			return nil
		})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Handle JSON parsing error
			err := args[0]
			errorMsg := fmt.Sprintf("Error parsing response: %s", err.Get("message").String())
			showError(errorMsg)
			return nil
		}))

		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Handle network error
		err := args[0]
		errorMsg := fmt.Sprintf("Network error: %s", err.Get("message").String())
		showError(errorMsg)
		return nil
	}))
}

// showLoading displays a loading indicator
func showLoading() {
	resultsContainer.Set("innerHTML", `
		<div class="col-span-full flex justify-center items-center py-8">
			<div class="loader"></div>
			<p class="ml-4">Loading parts...</p>
		</div>
	`)
}

// showError displays an error message
func showError(message string) {
	resultsContainer.Set("innerHTML", fmt.Sprintf(`
		<div class="col-span-full bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
			<p>%s</p>
		</div>
	`, message))
}

// displayParts displays the parts in the results container
func displayParts(parts []Part) {
	// Check if there are any parts
	if len(parts) == 0 {
		resultsContainer.Set("innerHTML", `
			<div class="col-span-full bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded">
				<p>No parts found matching your search criteria.</p>
			</div>
		`)
		return
	}

	// Build HTML for each part
	html := ""
	for _, part := range parts {
		// Get main image or placeholder
		imageURL := "https://via.placeholder.com/400x300"
		if len(part.Images) > 0 {
			imageURL = part.Images[0]
		}

		// Format price
		price := fmt.Sprintf("%.2f", part.Price)
		msrp := ""
		if part.MSRP > 0 && part.MSRP > part.Price {
			msrp = fmt.Sprintf("%.2f", part.MSRP)
		}

		// Stock status
		stockStatus := `<span class="inline-block bg-green-100 text-green-800 px-2 py-1 rounded text-xs">In Stock</span>`
		if !part.InStock {
			stockStatus = `<span class="inline-block bg-red-100 text-red-800 px-2 py-1 rounded text-xs">Out of Stock</span>`
		}

		// Build part card HTML
		partHTML := fmt.Sprintf(`
			<div class="part-card bg-white rounded-lg shadow-md overflow-hidden">
				<div class="h-48 bg-gray-200">
					<img src="%s" alt="%s %s" class="w-full h-full object-cover">
				</div>
				<div class="p-4">
					<div class="flex justify-between items-start">
						<div>
							<h3 class="font-bold text-lg">%s %s</h3>
							<p class="text-sm text-gray-600">%s</p>
						</div>
						<div class="text-right">
							<p class="font-bold text-lg">$%s</p>
							%s
						</div>
					</div>
					<p class="text-sm mt-2">%s</p>
					<div class="mt-4 flex justify-between">
						%s
						<a href="#" class="text-blue-600 hover:text-blue-800 text-sm" data-part-id="%s">View Details</a>
					</div>
				</div>
			</div>
		`,
			imageURL,
			part.Brand,
			part.Model,
			part.Brand,
			part.Model,
			part.Category,
			price,
			msrpHTML(msrp),
			part.Description,
			stockStatus,
			part.ID,
		)

		html += partHTML
	}

	// Set the HTML
	resultsContainer.Set("innerHTML", html)

	// Add event listeners to "View Details" links
	addViewDetailsListeners()
}

// msrpHTML returns HTML for the MSRP if available
func msrpHTML(msrp string) string {
	if msrp != "" {
		return fmt.Sprintf(`<p class="text-sm line-through text-gray-500">$%s</p>`, msrp)
	}
	return ""
}

// addViewDetailsListeners adds event listeners to "View Details" links
func addViewDetailsListeners() {
	document := js.Global().Get("document")
	links := document.Call("querySelectorAll", ".part-card a[data-part-id]")

	length := links.Length()
	for i := 0; i < length; i++ {
		link := links.Index(i)
		partID := link.Get("dataset").Get("partId").String()

		link.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Prevent default link behavior
			args[0].Call("preventDefault")

			// Show part details
			showPartDetails(partID)

			return nil
		}))
	}
}

// showPartDetails displays details for a specific part
func showPartDetails(partID string) {
	// Show loading indicator
	showLoading()

	// Build API URL
	url := fmt.Sprintf("%s/parts/%s", apiBaseURL, partID)

	// Create fetch promise
	promise := js.Global().Call("fetch", url)

	// Handle response
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		response := args[0]

		// Check if response is ok
		if !response.Get("ok").Bool() {
			// Handle HTTP error
			errorMsg := fmt.Sprintf("Error: %s (%d)", response.Get("statusText").String(), response.Get("status").Int())
			showError(errorMsg)
			return nil
		}

		// Parse response as JSON
		jsonPromise := response.Call("json")
		jsonPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Handle parsed JSON
			result := args[0]

			// Check for API error response
			if !result.Get("error").IsUndefined() {
				errorMsg := result.Get("message").String()
				showError(errorMsg)
				return nil
			}

			// Convert JS object to Go struct
			partJSON := js.Global().Get("JSON").Call("stringify", result).String()
			var part Part
			json.Unmarshal([]byte(partJSON), &part)

			// Display part details
			displayPartDetails(part)
			return nil
		})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Handle JSON parsing error
			err := args[0]
			errorMsg := fmt.Sprintf("Error parsing response: %s", err.Get("message").String())
			showError(errorMsg)
			return nil
		}))

		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Handle network error
		err := args[0]
		errorMsg := fmt.Sprintf("Network error: %s", err.Get("message").String())
		showError(errorMsg)
		return nil
	}))
}

// displayPartDetails shows detailed information for a single part
func displayPartDetails(part Part) {
	// Build image carousel HTML
	imagesHTML := ""
	if len(part.Images) > 0 {
		imagesHTML = `<div class="flex overflow-x-auto space-x-4 mb-4 p-2">`
		for _, img := range part.Images {
			imagesHTML += fmt.Sprintf(`
				<div class="flex-shrink-0 w-64 h-64">
					<img src="%s" alt="%s %s" class="w-full h-full object-contain">
				</div>
			`, img, part.Brand, part.Model)
		}
		imagesHTML += `</div>`
	} else {
		imagesHTML = fmt.Sprintf(`
			<div class="w-full h-64 bg-gray-200 mb-4">
				<img src="https://via.placeholder.com/640x480" alt="%s %s" class="w-full h-full object-contain">
			</div>
		`, part.Brand, part.Model)
	}

	// Stock status
	stockStatus := `<span class="badge badge-success">In Stock</span>`
	if !part.InStock {
		stockStatus = `<span class="badge badge-danger">Out of Stock</span>`
	}

	// Price information
	priceHTML := fmt.Sprintf(`<span class="text-2xl font-bold">$%.2f</span>`, part.Price)
	if part.MSRP > 0 && part.MSRP > part.Price {
		priceHTML += fmt.Sprintf(`
			<div class="flex items-center">
				<span class="line-through text-gray-500 mr-2">$%.2f</span>
				<span class="badge badge-success">Save %.0f%%</span>
			</div>
		`, part.MSRP, part.Discount)
	}

	// Build the part details HTML
	html := fmt.Sprintf(`
		<div class="bg-white rounded-lg shadow-md overflow-hidden col-span-full">
			<div class="p-6">
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-2xl font-bold">%s %s</h2>
					<a href="#" id="back-to-results" class="text-blue-600 hover:text-blue-800">Back to results</a>
				</div>

				<div class="flex flex-col md:flex-row">
					<div class="md:w-1/2 mb-4 md:mb-0 md:pr-6">
						%s
					</div>

					<div class="md:w-1/2">
						<div class="mb-4">
							<p class="text-sm text-gray-600">%s > %s</p>
							<div class="flex items-center mt-2">
								%s
								%s
							</div>
						</div>

						<hr class="my-4">

						<div class="mb-4">
							<h3 class="font-bold mb-2">Description</h3>
							<p>%s</p>
						</div>

						<div class="mb-4">
							<button id="view-on-site" class="btn-primary mr-2" data-url="%s">
								View on retailer site
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	`,
		part.Brand, part.Model,
		imagesHTML,
		part.Category, part.SubCategory,
		stockStatus,
		priceHTML,
		part.Description,
		part.URL,
	)

	// Set the HTML
	resultsContainer.Set("innerHTML", html)

	// Add event listener to "Back to results" button
	document := js.Global().Get("document")
	backButton := document.Call("getElementById", "back-to-results")
	backButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Prevent default link behavior
		args[0].Call("preventDefault")

		// Go back to results
		query := searchInput.Get("value").String()
		category := categoryFilter.Get("value").String()
		go fetchParts(query, category)

		return nil
	}))

	// Add event listener to "View on retailer site" button
	viewOnSiteButton := document.Call("getElementById", "view-on-site")
	viewOnSiteButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Open the URL in a new tab
		url := viewOnSiteButton.Get("dataset").Get("url").String()
		js.Global().Get("window").Call("open", url, "_blank")

		return nil
	}))
}
