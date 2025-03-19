package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// GitHubSearchResult describes the response from the GitHub API search.
type GitHubSearchResult struct {
	TotalCount int           `json:"total_count"`
	Items      []GitHubIssue `json:"items"`
}

// GitHubIssue describes the required fields from a pull request.
// Other fields can be added as needed.
type GitHubIssue struct {
	HTMLURL string `json:"html_url"`
	Title   string `json:"title"`
}

// fetchPullRequests retrieves a list of pull requests requiring review for the specified user.
func fetchPullRequests(username, githubToken string) ([]GitHubIssue, error) {
	// Form the search query: open pull requests where a review has been requested.
	query := fmt.Sprintf("is:open is:pr review-requested:%s", username)
	apiURL := fmt.Sprintf("https://api.github.com/search/issues?q=%s", url.QueryEscape(query))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	// If a token is provided, add authentication.
	if githubToken != "" {
		req.Header.Set("Authorization", "token "+githubToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(bodyBytes))
	}

	var result GitHubSearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// sendTelegramMessage sends a message to Telegram via the Bot API.
func sendTelegramMessage(telegramToken, chatID, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API error: %s", string(bodyBytes))
	}
	return nil
}

func main() {
	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error reading it, relying on environment variables")
	}

	// Retrieve the necessary data from environment variables.
	githubUsername := os.Getenv("GITHUB_USERNAME")
	githubToken := os.Getenv("GITHUB_TOKEN")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if githubUsername == "" || telegramToken == "" || chatID == "" {
		log.Fatal("Please set the environment variables: GITHUB_USERNAME, TELEGRAM_TOKEN, and TELEGRAM_CHAT_ID")
	}

	// Load Spain timezone (Europe/Madrid).
	spainLocation, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Printf("Error loading Spain timezone: %v", err)
		spainLocation = time.Local
	}

	// Start a loop that checks for pull requests every 4 hours.
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		now := time.Now().In(spainLocation)
		// Check if current time is between 9:00 and 21:00 Spain time.
		if now.Hour() < 9 || now.Hour() >= 21 {
			log.Printf("Skipping execution. Current time %v is outside allowed hours (09:00 - 21:00 Spain time).", now.Format("15:04"))
		} else {
			log.Println("Checking pull requests for review...")
			prs, err := fetchPullRequests(githubUsername, githubToken)
			if err != nil {
				log.Println("Error fetching pull requests:", err)
			} else if len(prs) > 0 {
				message := "Pull requests requiring review:\n"
				for _, pr := range prs {
					message += fmt.Sprintf("- %s\n%s\n", pr.Title, pr.HTMLURL)
				}
				err = sendTelegramMessage(telegramToken, chatID, message)
				if err != nil {
					log.Println("Error sending Telegram message:", err)
				} else {
					log.Println("Message sent successfully!")
				}
			} else {
				log.Println("No pull requests requiring review found.")
			}
		}
		<-ticker.C
	}
}
