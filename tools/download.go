package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var DownloadYtdlp = &ToolDef{
	Name:        "download_ytdlp",
	Description: "Download video or audio using yt-dlp if it is installed on the system map.",
	Args: []ToolArg{
		{Name: "url", Description: "URL to download", Required: true},
		{Name: "audio_only", Description: "Set to 'true' to extract audio only", Required: false},
		{Name: "options", Description: "Extra command line flags (e.g. '-f best')", Required: false},
	},
	Execute: func(args map[string]string) string {
		url := strings.TrimSpace(args["url"])
		if url == "" {
			return "Error: url is required"
		}

		if _, err := exec.LookPath("yt-dlp"); err != nil {
			return "Error: yt-dlp is not installed or not in PATH."
		}

		var cmdArgs []string
		if args["audio_only"] == "true" {
			cmdArgs = append(cmdArgs, "-x", "--audio-format", "mp3")
		}
		if opts := strings.TrimSpace(args["options"]); opts != "" {
			cmdArgs = append(cmdArgs, strings.Split(opts, " ")...)
		}
		cmdArgs = append(cmdArgs, url)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		cmd := exec.CommandContext(ctx, "yt-dlp", cmdArgs...)
		out, err := cmd.CombinedOutput()

		res := string(out)
		if len(res) > 4000 {
			res = res[len(res)-4000:]
		}

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Sprintf("Timeout (5m).\n...%s", res)
			}
			return fmt.Sprintf("Error: %v\n...%s", err, res)
		}
		return fmt.Sprintf("Success:\n...%s", res)
	},
}

var DownloadAria2c = &ToolDef{
	Name:        "download_aria2c",
	Description: "Download files using aria2c if it is installed on the system map.",
	Args: []ToolArg{
		{Name: "url", Description: "URL to download", Required: true},
		{Name: "options", Description: "Extra command line flags (e.g. '-x 16')", Required: false},
	},
	Execute: func(args map[string]string) string {
		url := strings.TrimSpace(args["url"])
		if url == "" {
			return "Error: url is required"
		}

		if _, err := exec.LookPath("aria2c"); err != nil {
			return "Error: aria2c is not installed or not in PATH."
		}

		var cmdArgs []string
		if opts := strings.TrimSpace(args["options"]); opts != "" {
			cmdArgs = append(cmdArgs, strings.Split(opts, " ")...)
		}
		cmdArgs = append(cmdArgs, url)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		cmd := exec.CommandContext(ctx, "aria2c", cmdArgs...)
		out, err := cmd.CombinedOutput()

		res := string(out)
		if len(res) > 4000 {
			res = res[len(res)-4000:]
		}

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Sprintf("Timeout (5m).\n...%s", res)
			}
			return fmt.Sprintf("Error: %v\n...%s", err, res)
		}
		return fmt.Sprintf("Success:\n...%s", res)
	},
}

var YouTubeTranscript = &ToolDef{
	Name:        "youtube_transcript",
	Description: "Fetch transcripts from YouTube videos for summarization, QA, and content extraction",
	Args: []ToolArg{
		{Name: "url", Description: "YouTube video URL", Required: true},
		{Name: "language", Description: "Subtitle language: 'en' (default), 'es', 'fr', 'de', etc.", Required: false},
		{Name: "format", Description: "Output format: 'text' (default), 'markdown', 'json'", Required: false},
	},
	Execute: func(args map[string]string) string {
		url := args["url"]
		if url == "" {
			return "Error: url is required"
		}

		language := args["language"]
		if language == "" {
			language = "en"
		}

		format := args["format"]
		if format == "" {
			format = "text"
		}

		return getYouTubeTranscript(url, language, format)
	},
}

func getYouTubeTranscript(url string, language string, format string) string {
	if !commandExists("yt-dlp") {
		return "Error: yt-dlp not found. Install it with: pip install yt-dlp"
	}
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("yt_%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Sprintf("Error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "subs")
	cmd := exec.Command(
		"yt-dlp",
		"--write-subs",
		"--write-auto-subs",
		"--skip-download",
		"--sub-lang", language,
		"--output", outputPath,
		url,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Error fetching transcript: %v", err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Sprintf("Error reading temp directory: %v", err)
	}

	var vttFile string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".vtt") {
			vttFile = filepath.Join(tempDir, f.Name())
			break
		}
	}

	if vttFile == "" {
		return "Error: No subtitles found for this video. Try a different language or ensure the video has captions available."
	}

	content, err := os.ReadFile(vttFile)
	if err != nil {
		return fmt.Sprintf("Error reading transcript file: %v", err)
	}

	cleanText := cleanVTT(string(content))

	switch format {
	case "markdown":
		return formatMarkdown(url, cleanText)
	case "json":
		return formatJSON(url, cleanText)
	default:
		return cleanText
	}
}

func cleanVTT(content string) string {
	lines := strings.Split(content, "\n")
	var textLines []string
	seen := make(map[string]bool)

	timestampPattern := regexp.MustCompile(`\d{2}:\d{2}:\d{2}\.\d{3}\s-->\s\d{2}:\d{2}:\d{2}\.\d{3}`)
	htmlTagPattern := regexp.MustCompile(`<[^>]+>`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || line == "WEBVTT" || line == "Kind: captions" || line == "Language: en" {
			continue
		}
		if regexp.MustCompile(`^\d+$`).MatchString(line) {
			continue
		}
		if timestampPattern.MatchString(line) {
			continue
		}
		if strings.HasPrefix(line, "NOTE") || strings.HasPrefix(line, "STYLE") {
			continue
		}

		line = htmlTagPattern.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if seen[line] {
			continue
		}
		seen[line] = true

		textLines = append(textLines, line)
	}

	return strings.Join(textLines, "\n")
}

func formatMarkdown(url string, text string) string {
	var sb strings.Builder
	sb.WriteString("# YouTube Transcript\n\n")
	fmt.Fprintf(&sb, "**Source:** %s\n\n", url)
	sb.WriteString("## Content\n\n")

	sentences := strings.Split(text, "\n")
	paragraph := ""

	for _, sentence := range sentences {
		paragraph += sentence + " "

		if len(paragraph) > 300 {
			sb.WriteString(paragraph + "\n\n")
			paragraph = ""
		}
	}

	if paragraph != "" {
		sb.WriteString(paragraph + "\n")
	}

	return sb.String()
}

func formatJSON(url string, text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "\"", "\\\"")
	text = strings.ReplaceAll(text, "\n", "\\n")
	text = strings.ReplaceAll(text, "\r", "\\r")
	text = strings.ReplaceAll(text, "\t", "\\t")

	return fmt.Sprintf(`{
  "source": "%s",
  "content": "%s",
  "type": "youtube_transcript"
}`, url, text)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
