package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Google Calendar API tools using Maton gateway

func calendarAPIRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	apiKey := os.Getenv("MATON_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MATON_API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://gateway.maton.ai/google-calendar/calendar/v3%s", endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Google Calendar API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

var CalendarListEvents = &ToolDef{
	Name:        "calendar_list_events",
	Description: "List upcoming events from Google Calendar. Requires MATON_API_KEY env var.",
	Secure:      true,
	Args: []ToolArg{
		{Name: "calendar_id", Description: "Calendar ID (default 'primary' for main calendar)", Required: false},
		{Name: "time_min", Description: "Start time (RFC 3339 format, e.g. '2024-01-15T10:00:00Z')", Required: false},
		{Name: "time_max", Description: "End time (RFC 3339 format, e.g. '2024-01-20T23:59:59Z')", Required: false},
		{Name: "max_results", Description: "Maximum events to return (default 10, max 250)", Required: false},
	},
	Execute: func(args map[string]string) string {
		calendarID := strings.TrimSpace(args["calendar_id"])
		if calendarID == "" {
			calendarID = "primary"
		}

		maxResults := "10"
		if v := strings.TrimSpace(args["max_results"]); v != "" {
			if _, err := strconv.Atoi(v); err == nil {
				maxResults = v
			}
		}

		endpoint := fmt.Sprintf("/calendars/%s/events?maxResults=%s&orderBy=startTime&singleEvents=true", calendarID, maxResults)

		if timeMin := strings.TrimSpace(args["time_min"]); timeMin != "" {
			endpoint += fmt.Sprintf("&timeMin=%s", timeMin)
		} else {
			now := time.Now().Format(time.RFC3339)
			endpoint += fmt.Sprintf("&timeMin=%s", now)
		}

		if timeMax := strings.TrimSpace(args["time_max"]); timeMax != "" {
			endpoint += fmt.Sprintf("&timeMax=%s", timeMax)
		}

		respBody, err := calendarAPIRequest("GET", endpoint, nil)
		if err != nil {
			return fmt.Sprintf("Error fetching events: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Sprintf("Error parsing response: %v", err)
		}

		items, ok := result["items"].([]interface{})
		if !ok || len(items) == 0 {
			return "No upcoming events found."
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📅 Found %d event(s):\n\n", len(items)))
		for i, item := range items {
			if eventMap, ok := item.(map[string]interface{}); ok {
				summary, _ := eventMap["summary"].(string)
				if summary == "" {
					summary = "(no title)"
				}

				var startTime, endTime string
				if start, ok := eventMap["start"].(map[string]interface{}); ok {
					if dt, ok := start["dateTime"].(string); ok {
						startTime = dt
					} else if d, ok := start["date"].(string); ok {
						startTime = d
					}
				}
				if end, ok := eventMap["end"].(map[string]interface{}); ok {
					if dt, ok := end["dateTime"].(string); ok {
						endTime = dt
					} else if d, ok := end["date"].(string); ok {
						endTime = d
					}
				}

				description := ""
				if desc, ok := eventMap["description"].(string); ok && desc != "" {
					if len(desc) > 100 {
						description = desc[:100] + "..."
					} else {
						description = desc
					}
				}

				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, summary))
				sb.WriteString(fmt.Sprintf("   Start: %s\n", startTime))
				sb.WriteString(fmt.Sprintf("   End: %s\n", endTime))
				if description != "" {
					sb.WriteString(fmt.Sprintf("   Note: %s\n", description))
				}
				sb.WriteString("\n")
			}
		}
		return strings.TrimRight(sb.String(), "\n")
	},
}

var CalendarCreateEvent = &ToolDef{
	Name:        "calendar_create_event",
	Description: "Create a new event in Google Calendar. Requires MATON_API_KEY env var.",
	Secure:      true,
	Args: []ToolArg{
		{Name: "summary", Description: "Event title (required)", Required: true},
		{Name: "start_time", Description: "Start time (RFC 3339 format, e.g. '2024-01-15T10:00:00Z')", Required: true},
		{Name: "end_time", Description: "End time (RFC 3339 format, e.g. '2024-01-15T11:00:00Z')", Required: true},
		{Name: "description", Description: "Event description (optional)", Required: false},
		{Name: "location", Description: "Event location (optional)", Required: false},
		{Name: "calendar_id", Description: "Calendar ID (default 'primary')", Required: false},
		{Name: "attendees", Description: "Attendee emails (comma-separated, optional)", Required: false},
	},
	Execute: func(args map[string]string) string {
		summary := strings.TrimSpace(args["summary"])
		startTime := strings.TrimSpace(args["start_time"])
		endTime := strings.TrimSpace(args["end_time"])
		description := strings.TrimSpace(args["description"])
		location := strings.TrimSpace(args["location"])
		calendarID := strings.TrimSpace(args["calendar_id"])
		attendeesStr := strings.TrimSpace(args["attendees"])

		if calendarID == "" {
			calendarID = "primary"
		}

		if summary == "" || startTime == "" || endTime == "" {
			return "Error: summary, start_time, and end_time are required"
		}

		event := map[string]interface{}{
			"summary": summary,
			"start": map[string]string{
				"dateTime": startTime,
			},
			"end": map[string]string{
				"dateTime": endTime,
			},
		}

		if description != "" {
			event["description"] = description
		}
		if location != "" {
			event["location"] = location
		}

		if attendeesStr != "" {
			var attendees []map[string]string
			for _, email := range strings.Split(attendeesStr, ",") {
				if email = strings.TrimSpace(email); email != "" {
					attendees = append(attendees, map[string]string{"email": email})
				}
			}
			if len(attendees) > 0 {
				event["attendees"] = attendees
			}
		}

		eventJSON, _ := json.Marshal(event)
		endpoint := fmt.Sprintf("/calendars/%s/events", calendarID)

		respBody, err := calendarAPIRequest("POST", endpoint, strings.NewReader(string(eventJSON)))
		if err != nil {
			return fmt.Sprintf("Error creating event: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Sprintf("Error parsing response: %v", err)
		}

		if id, ok := result["id"].(string); ok {
			return fmt.Sprintf("✅ Event created successfully — ID: %s", id)
		}

		return "✅ Event created successfully"
	},
}

var CalendarDeleteEvent = &ToolDef{
	Name:        "calendar_delete_event",
	Description: "Delete an event from Google Calendar by ID.",
	Secure:      true,
	Args: []ToolArg{
		{Name: "event_id", Description: "Event ID to delete (required)", Required: true},
		{Name: "calendar_id", Description: "Calendar ID (default 'primary')", Required: false},
	},
	Execute: func(args map[string]string) string {
		eventID := strings.TrimSpace(args["event_id"])
		calendarID := strings.TrimSpace(args["calendar_id"])

		if eventID == "" {
			return "Error: event_id is required"
		}

		if calendarID == "" {
			calendarID = "primary"
		}

		endpoint := fmt.Sprintf("/calendars/%s/events/%s", calendarID, eventID)

		_, err := calendarAPIRequest("DELETE", endpoint, nil)
		if err != nil {
			return fmt.Sprintf("Error deleting event: %v", err)
		}

		return fmt.Sprintf("✅ Event %s deleted successfully", eventID)
	},
}

var CalendarUpdateEvent = &ToolDef{
	Name:        "calendar_update_event",
	Description: "Update an existing event in Google Calendar.",
	Secure:      true,
	Args: []ToolArg{
		{Name: "event_id", Description: "Event ID to update (required)", Required: true},
		{Name: "summary", Description: "Event title (optional)", Required: false},
		{Name: "start_time", Description: "Start time (RFC 3339 format, optional)", Required: false},
		{Name: "end_time", Description: "End time (RFC 3339 format, optional)", Required: false},
		{Name: "description", Description: "Event description (optional)", Required: false},
		{Name: "location", Description: "Event location (optional)", Required: false},
		{Name: "calendar_id", Description: "Calendar ID (default 'primary')", Required: false},
	},
	Execute: func(args map[string]string) string {
		eventID := strings.TrimSpace(args["event_id"])
		calendarID := strings.TrimSpace(args["calendar_id"])

		if eventID == "" {
			return "Error: event_id is required"
		}

		if calendarID == "" {
			calendarID = "primary"
		}

		endpoint := fmt.Sprintf("/calendars/%s/events/%s", calendarID, eventID)

		respBody, err := calendarAPIRequest("GET", endpoint, nil)
		if err != nil {
			return fmt.Sprintf("Error fetching event: %v", err)
		}

		var event map[string]interface{}
		if err := json.Unmarshal(respBody, &event); err != nil {
			return fmt.Sprintf("Error parsing event: %v", err)
		}

		if summary := strings.TrimSpace(args["summary"]); summary != "" {
			event["summary"] = summary
		}
		if startTime := strings.TrimSpace(args["start_time"]); startTime != "" {
			event["start"] = map[string]string{"dateTime": startTime}
		}
		if endTime := strings.TrimSpace(args["end_time"]); endTime != "" {
			event["end"] = map[string]string{"dateTime": endTime}
		}
		if description := strings.TrimSpace(args["description"]); description != "" {
			event["description"] = description
		}
		if location := strings.TrimSpace(args["location"]); location != "" {
			event["location"] = location
		}

		eventJSON, _ := json.Marshal(event)

		_, err = calendarAPIRequest("PUT", endpoint, strings.NewReader(string(eventJSON)))
		if err != nil {
			return fmt.Sprintf("Error updating event: %v", err)
		}

		return fmt.Sprintf("✅ Event %s updated successfully", eventID)
	},
}
