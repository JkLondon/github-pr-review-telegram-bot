Below is an example of a README file written in English, along with choosing the MIT License for the project.

---

**README.md**

```markdown
# GitHub PR Review Telegram Bot

A Go bot that monitors open GitHub pull requests requiring review and sends notifications via Telegram. Configurable through environment variables and supports GitHub token authentication. Ideal for developers who want to stay updated on PR reviews.

## Features

- **Automatic Monitoring:** Periodically polls the GitHub API for open pull requests that require your review.
- **Telegram Notifications:** Sends real-time notifications to your Telegram chat.
- **Easy Configuration:** Uses environment variables for quick and easy setup.
- **Token Support:** Optionally supports a GitHub personal access token for accessing private repositories.

## Requirements

- **Go:** Version 1.16 or later.
- **GitHub Account:** For accessing the GitHub API.
- **Telegram Bot:** A bot token from BotFather.
- **Optional:** GitHub personal access token for private repository access.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/github-pr-review-telegram-bot.git
   cd github-pr-review-telegram-bot
   ```

2. **Set up the necessary environment variables:**

   ```bash
   export GITHUB_USERNAME="your_github_username"
   export GITHUB_TOKEN="your_github_token" # Optional: Required for private repositories.
   export TELEGRAM_TOKEN="your_telegram_bot_token"
   export TELEGRAM_CHAT_ID="your_telegram_chat_id"
   ```

3. **Build the bot:**

   ```bash
   go build -o github-pr-review-telegram-bot
   ```

4. **Run the bot:**

   ```bash
   ./github-pr-review-telegram-bot
   ```

## Configuration

The bot uses the following environment variables:

- `GITHUB_USERNAME`: Your GitHub username used to filter pull requests requiring review.
- `GITHUB_TOKEN`: Your GitHub personal access token (optional, but recommended for accessing private repos).
- `TELEGRAM_TOKEN`: Your Telegram bot token.
- `TELEGRAM_CHAT_ID`: The Telegram chat ID where notifications will be sent.


## Contributing

Contributions are welcome! If you have ideas, improvements, or bug fixes, feel free to fork the repository and submit a pull request. For major changes, please open an issue first to discuss your proposed changes.
