# 🐾 Apexclaw

<p align="center">
  <img src="https://img.shields.io/github/stars/amarnathcjd/apexclaw?style=flat-square&color=ff69b4" alt="stars">
  <img src="https://img.shields.io/github/forks/amarnathcjd/apexclaw?style=flat-square&color=9370db" alt="forks">
  <img src="https://img.shields.io/github/license/amarnathcjd/apexclaw?style=flat-square&color=00bfff" alt="license">
</p>

**Your personal AI assistant that lives in Telegram.** Not just a chatbot—a capable agent that thinks, uses tools, and actually gets things done.

Apexclaw is powered by GLM and comes with **100+ built-in tools** covering everything from file management and web automation to email, calendar, WhatsApp, and beyond. Ask it anything, and it will figure out how to help.

---

## ✨ What it can do

- 🖼️ **Analyze images** — Send photos and get detailed descriptions
- 🎙️ **Transcribe voice** — Reply with voice notes and it transcribes + acts
- 📱 **WhatsApp integration** — Read and send WhatsApp messages directly
- 🌐 **Browse the web** — Use a real headless browser to navigate, click, and read
- 📧 **Email & Calendar** — Read Gmail, send emails, manage Google Calendar events
- 🧠 **Long-term memory** — Save facts, notes, and a searchable knowledge base
- 🎬 **Media tools** — Search IMDB, fetch movie details, download YouTube videos
- 🐍 **Run code** — Execute Python scripts, shell commands, system operations
- 🗺️ **Travel & navigation** — Flight searches, geocoding, directions, timezone conversion
- 📊 **Data & finance** — Live stock prices, currency conversion, weather, news
- 💬 **Interactive buttons** — Send messages with inline buttons for quick actions
- 🎯 **Task automation** — Schedule tasks, Pomodoro timers, daily digests

---

## 🚀 Quick Start

### One-line install (Linux/macOS)

```bash
curl -fsSL https://claw.gogram.fun | bash
apexclaw
```

This automatically:
1. Downloads the latest binary for your OS
2. Installs to `/usr/local/bin` or `~/.local/bin`
3. Launches an interactive setup wizard
4. Prompts for Telegram credentials
5. Starts the agent

### Manual build

```bash
# Requirements: Go 1.22+, ffmpeg (for voice)
git clone https://github.com/amarnathcjd/apexclaw
cd apexclaw
go build -o apexclaw .
./apexclaw
```

On first run, you'll be prompted for:
- **Telegram API ID** (from [my.telegram.org](https://my.telegram.org))
- **Telegram API Hash** (from [my.telegram.org](https://my.telegram.org))
- **Bot Token** (from [@BotFather](https://t.me/botfather))
- **Owner ID** (your Telegram Chat ID)

---

## 🔧 Setup Guides

### Gmail (Maton API)

For best email experience, use the new Maton API integration:

1. Get an API key from [maton.ai/settings](https://maton.ai/settings)
2. Add to `.env`:
```ini
MATON_API_KEY=your_api_key_here
```
3. Done! The following tools are now available:
   - `gmail_list_messages` — List/search emails with filters (unread, starred, from, subject, date range, attachments)
   - `gmail_get_message` — Get full message details by ID
   - `gmail_send_message` — Send emails (supports CC, BCC)
   - `gmail_modify_labels` — Add/remove labels (starred, unread, trash, important, etc.)

**Legacy email (IMAP/SMTP):**
```ini
EMAIL_ADDRESS=your.email@gmail.com
EMAIL_PASSWORD=your-16-char-app-password
EMAIL_IMAP_HOST=imap.gmail.com
EMAIL_IMAP_PORT=993
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
```

### Google Calendar (Maton API)

Same Maton API key unlocks calendar tools:
- `calendar_list_events` — List upcoming events (with date ranges, max 250)
- `calendar_create_event` — Create events with attendees, location, description
- `calendar_delete_event` — Delete events by ID
- `calendar_update_event` — Update existing events

Times use RFC 3339 format: `2024-01-15T10:00:00Z`

### WhatsApp Integration

WhatsApp support via official reverse engineering:

1. Start apexclaw and trigger WhatsApp connection
2. Scan the QR code with your WhatsApp phone
3. Connection saved to `~/.apexclaw/whatsapp/store.db`
4. Available tools:
   - Send/receive WhatsApp messages
   - Integrated with AI responses (replies in WhatsApp)

---

## 🧰 All Tools

### System & Execution
| Tool | Purpose |
|---|---|
| `exec` | Run shell commands with auto-detected timeout |
| `run_python` | Execute Python scripts |
| `system_info` | Get CPU, RAM, disk usage |
| `process_list` | List running processes |
| `kill_process` | Terminate a process by PID |
| `clipboard_get` | Read clipboard contents |
| `clipboard_set` | Write to clipboard |

### Files & Directory
| Tool | Purpose |
|---|---|
| `read_file` | Read file contents |
| `write_file` | Write or create files |
| `append_file` | Append content to file |
| `list_dir` | List directory contents |
| `create_dir` | Create directories |
| `delete_file` | Delete files |
| `move_file` | Move or rename files |
| `search_files` | Find files by pattern |

### Memory & Knowledge Base
| Tool | Purpose |
|---|---|
| `save_fact` | Persist a key-value fact |
| `recall_fact` | Retrieve saved facts |
| `list_facts` | List all facts |
| `delete_fact` | Delete a fact |
| `update_note` | Create/overwrite named notes |
| `kb_add` | Add to searchable knowledge base |
| `kb_search` | Search KB with TF-IDF ranking |
| `kb_list` | List KB entries |
| `kb_delete` | Remove from KB |

### Web & Search
| Tool | Purpose |
|---|---|
| `web_fetch` | Fetch and read webpages |
| `web_search` | Search the web with results |
| `http_request` | Make raw HTTP requests |
| `rss_feed` | Parse RSS/Atom feeds |
| `wikipedia` | Search and read Wikipedia |
| `news_headlines` | Get live news headlines |
| `reddit_feed` | Read top posts from subreddits |
| `youtube_search` | Search YouTube videos |

### Media & Entertainment
| Tool | Purpose |
|---|---|
| `imdb_search` | Search movies, shows, actors |
| `imdb_title` | Get details for IMDB ID |
| `pinterest_search` | Search Pinterest boards/pins |
| `pinterest_get_pin` | Get Pinterest pin details |
| `download_ytdlp` | Download videos/audio via yt-dlp |
| `download_aria2c` | Download files via aria2c |

### Browser Automation
| Tool | Purpose |
|---|---|
| `browser_open` | Open URL in headless Chrome |
| `browser_click` | Click elements by CSS selector |
| `browser_type` | Type into inputs |
| `browser_get_text` | Extract page text |
| `browser_eval` | Run JavaScript on page |
| `browser_screenshot` | Take page screenshots |

### Email & Communication
| Tool | Purpose |
|---|---|
| `gmail_list_messages` | List/search Gmail (with filters) |
| `gmail_get_message` | Get full message by ID |
| `gmail_send_message` | Send emails with CC/BCC |
| `gmail_modify_labels` | Add/remove Gmail labels |
| `read_email` | Read emails (IMAP legacy) |
| `send_email` | Send emails (SMTP legacy) |
| `text_to_speech` | Convert text to voice notes |

### Calendar & Scheduling
| Tool | Purpose |
|---|---|
| `calendar_list_events` | List calendar events |
| `calendar_create_event` | Create events |
| `calendar_update_event` | Update events |
| `calendar_delete_event` | Delete events |
| `schedule_task` | Schedule one-off or repeating tasks |
| `cancel_task` | Cancel scheduled tasks |
| `list_tasks` | List all scheduled tasks |
| `timer` | Set countdown timers |
| `pomodoro` | Start Pomodoro sessions |

### GitHub
| Tool | Purpose |
|---|---|
| `github_search` | Search repositories |
| `github_read_file` | Read files from repos |

### Travel & Navigation
| Tool | Purpose |
|---|---|
| `flight_airport_search` | Look up airport info |
| `flight_route_search` | Search flight routes |
| `flight_countries` | List supported countries |
| `nav_geocode` | Geocode addresses to coordinates |
| `nav_route` | Get directions between points |
| `nav_sunshade` | Calculate sun shading for drives |

### Data & Utilities
| Tool | Purpose |
|---|---|
| `weather` | Get live weather by location |
| `stock_price` | Get live stock prices |
| `currency_convert` | Convert between currencies |
| `unit_convert` | Convert units (length, weight, temp, etc.) |
| `timezone_convert` | Convert times between timezones |
| `translate` | Translate text to other languages |
| `ip_lookup` | Look up IP information |
| `dns_lookup` | Resolve DNS records |
| `calculate` | Evaluate math expressions |
| `hash_text` | Hash strings (MD5, SHA256, etc.) |
| `encode_decode` | Base64 encode/decode |
| `regex_match` | Test regex patterns |
| `color_info` | Get hex/RGB color info |
| `text_process` | Trim, split, replace text |
| `datetime` | Get current date/time |
| `random` | Generate random numbers |
| `echo` | Echo back messages |

### Productivity
| Tool | Purpose |
|---|---|
| `todo_add` | Add to-do items |
| `todo_list` | List all to-dos |
| `todo_done` | Mark to-do as done |
| `todo_delete` | Delete to-do items |
| `daily_digest` | Set up daily briefing |
| `cron_status` | Check scheduled task status |

### Documents
| Tool | Purpose |
|---|---|
| `read_document` | Read stored documents |
| `list_documents` | List all documents |
| `summarize_document` | Summarize documents |

### Telegram
| Tool | Purpose |
|---|---|
| `tg_send_message` | Send text messages |
| `tg_send_file` | Send files/photos |
| `tg_send_message_buttons` | Send messages with inline buttons |
| `tg_download` | Download Telegram media |
| `tg_get_chat_info` | Get chat/user info |
| `tg_forward` | Forward messages |
| `tg_delete_msg` | Delete messages |
| `tg_pin_msg` | Pin messages |
| `tg_react` | React with emojis |
| `tg_get_reply` | Get replied-to message |
| `set_bot_dp` | Update bot profile picture |

---

## 🔐 Web Dashboard (Optional)

Apexclaw includes a web interface for settings and monitoring:

- **Default login code**: `123456` (shown on first startup)
- **JWT authentication**: 1-hour sessions with refresh tokens
- **Change code anytime**: Via web UI or Telegram `/webcode` command
- **Telegram integration**: Use `/webcode show/set/random` to manage codes

To access the dashboard:
1. Start apexclaw
2. Open `http://localhost:8080`
3. Enter the 6-digit login code
4. Manage settings, view logs, and monitor activity

---

## 📝 Adding Custom Tools

Create a new `ToolDef` in any file under `tools/` directory and register it in `tools/tools.go`:

```go
var MyTool = &core.ToolDef{
	Name: "my_tool",
	Description: "What your tool does",
	Args: []core.ToolArg{
		{Name: "input", Description: "Input parameter", Required: true},
	},
	Execute: func(args map[string]string) string {
		// Your implementation here
		return result
	},
}
```

Apexclaw automatically picks it up—no restart needed (in debug mode).

---

## 📊 Logging & Debugging

Apexclaw logs to `~/.apexclaw/logs/` with rotating files. Log levels: DEBUG, INFO, WARN, ERROR.

Use the execution profiler to trace tool calls:
- `debug_trace` tool — Track timing, results, and context for each tool execution
- Helps identify slow operations or failures

---

## 📦 Dependencies

- **Go 1.22+**
- **ffmpeg** (for voice transcription)
- Built-in support for:
  - Telegram Bot API
  - Chrome/Chromium (for browser automation)

---

## 🎯 Planned Features

- [ ] Web plugin system for third-party tools
- [ ] Persistent refresh token storage
- [ ] Vector embeddings for smarter memory recall
- [ ] Mobile app companion
- [ ] Multi-language interface

---

## 📄 License

MIT License — See LICENSE file for details

---

## 🤝 Contributing

Found a bug or have a feature request? Open an issue or pull request on GitHub.

Questions? Contact [@amarnathcjd](https://t.me/amarnathcjd) on Telegram.

---

**Apexclaw** — *An AI assistant that actually gets things done.*
