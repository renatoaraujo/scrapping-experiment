# Web Scraping Experiment

This repository contains a simple script used for web scraping various pages to extract company information. It's designed as an experiment to understand and handle the scraping process, including challenges like handling redirects, page interaction, and data extraction.

## Overview

The script targets specific URLs, initiating multiple scraping tasks concurrently. Each task navigates to a web page, interacts with elements to reveal content (if necessary), and extracts specific data, which is then displayed in the console.

Here's what the script does:

1. **Concurrent Scraping:** Launches several scraping tasks in parallel, each for a different URL.
2. **Error Handling:** Logs and retries failed scraping tasks.
3. **Page Interaction:** Uses a headless browser to interact with pages just like a human user might, including following redirects and clicking buttons.
4. **Data Extraction:** Parses the final page content to extract specific pieces of company information.

## Technology

This script is written in Go and uses the following key packages:

- `chromedp`: A Go library for driving browsers (specifically Chrome) using the DevTools Protocol, used for navigation, page interaction, and content retrieval.
- `goquery`: A Go library that provides jQuery-like functionality for parsing and querying HTML content.

## Running the Script

To run the script, you need to have Go installed and configured on your machine. Then, you can clone this repository and run the script using the Go command line.

```bash
# Clone the repository
git clone https://github.com/renatoaraujo/scrapping-experiment.git

# Navigate to the repository directory
cd scrapping-experiment

# Run the script
go run main.go
```

Please note that this script is an experiment and is not designed for production use. It serves as a demonstration of web scraping concepts and techniques.