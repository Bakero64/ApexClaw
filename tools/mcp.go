package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

var MCPCall = &ToolDef{
	Name:        "mcp_call",
	Description: "Call tools from MCP servers using mcporter. Requires mcporter to be installed.",
	Args: []ToolArg{
		{Name: "server", Description: "MCP server name or URL (e.g., 'linear', 'https://api.example.com/mcp.fetch')", Required: true},
		{Name: "tool", Description: "Tool name to call (e.g., 'list_issues', 'create_issue')", Required: true},
		{Name: "args", Description: "Tool arguments as JSON or key=value pairs (e.g., '{\"limit\":5}' or 'team=ENG limit:5')", Required: false},
		{Name: "output", Description: "Output format: 'text' (default), 'json'", Required: false},
	},
	Execute: func(args map[string]string) string {
		server := args["server"]
		tool := args["tool"]
		toolArgs := args["args"]
		output := args["output"]

		if server == "" || tool == "" {
			return "Error: server and tool are required"
		}

		if output == "" {
			output = "text"
		}

		return callMCP(server, tool, toolArgs, output)
	},
}

var MCPList = &ToolDef{
	Name:        "mcp_list",
	Description: "List available MCP servers or tools from a specific server. Requires mcporter to be installed.",
	Args: []ToolArg{
		{Name: "server", Description: "Optional: MCP server name to list tools from (e.g., 'linear'). If omitted, lists all servers.", Required: false},
		{Name: "schema", Description: "Show schema for tools: 'true' or 'false' (default)", Required: false},
		{Name: "output", Description: "Output format: 'text' (default), 'json'", Required: false},
	},
	Execute: func(args map[string]string) string {
		server := args["server"]
		schema := args["schema"] == "true"
		output := args["output"]

		if output == "" {
			output = "text"
		}

		return listMCP(server, schema, output)
	},
}

var MCPAuth = &ToolDef{
	Name:        "mcp_auth",
	Description: "Authenticate with an MCP server using OAuth or API keys. Requires mcporter to be installed.",
	Args: []ToolArg{
		{Name: "server", Description: "MCP server name or URL to authenticate with", Required: true},
		{Name: "reset", Description: "Reset authentication: 'true' or 'false' (default)", Required: false},
	},
	Execute: func(args map[string]string) string {
		server := args["server"]
		reset := args["reset"] == "true"

		if server == "" {
			return "Error: server is required"
		}

		return authMCP(server, reset)
	},
}

var MCPConfig = &ToolDef{
	Name:        "mcp_config",
	Description: "Manage mcporter configuration (list, add, remove, import servers). Requires mcporter to be installed.",
	Args: []ToolArg{
		{Name: "action", Description: "Action: 'list', 'get', 'add', 'remove', 'import', 'login', 'logout'", Required: true},
		{Name: "server", Description: "Server name (for get/add/remove actions)", Required: false},
		{Name: "config", Description: "Configuration file path (for import action)", Required: false},
		{Name: "url", Description: "Server URL (for add action)", Required: false},
	},
	Execute: func(args map[string]string) string {
		action := args["action"]
		server := args["server"]
		configPath := args["config"]
		url := args["url"]

		if action == "" {
			return "Error: action is required"
		}

		return configMCP(action, server, configPath, url)
	},
}

func callMCP(server string, tool string, toolArgs string, output string) string {
	if !commandExists("mcporter") {
		return "Error: mcporter not found. Install it with: npm install -g @mcporter/cli"
	}

	cmd := exec.Command("mcporter", "call", fmt.Sprintf("%s.%s", server, tool))

	if toolArgs != "" {
		if strings.HasPrefix(strings.TrimSpace(toolArgs), "{") {
			cmd.Args = append(cmd.Args, "--args", toolArgs)
		} else {
			pairs := strings.Fields(toolArgs)
			cmd.Args = append(cmd.Args, pairs...)
		}
	}

	if output == "json" {
		cmd.Args = append(cmd.Args, "--output", "json")
	}

	output_bytes, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Sprintf("Error calling MCP tool: %s\nDetails: %s", err, string(exitErr.Stderr))
		}
		return fmt.Sprintf("Error calling MCP tool: %v", err)
	}

	result := string(output_bytes)

	if output == "json" {
		var jsonData interface{}
		if err := json.Unmarshal([]byte(result), &jsonData); err == nil {
			prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
			return string(prettyJSON)
		}
	}

	return result
}

func listMCP(server string, schema bool, output string) string {
	if !commandExists("mcporter") {
		return "Error: mcporter not found. Install it with: npm install -g @mcporter/cli"
	}

	cmd := exec.Command("mcporter", "list")

	if server != "" {
		cmd.Args = append(cmd.Args, server)
	}

	if schema {
		cmd.Args = append(cmd.Args, "--schema")
	}

	if output == "json" {
		cmd.Args = append(cmd.Args, "--output", "json")
	}

	result, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Sprintf("Error listing MCP servers/tools: %s\nDetails: %s", err, string(exitErr.Stderr))
		}
		return fmt.Sprintf("Error listing MCP servers/tools: %v", err)
	}

	return string(result)
}

func authMCP(server string, reset bool) string {
	if !commandExists("mcporter") {
		return "Error: mcporter not found. Install it with: npm install -g @mcporter/cli"
	}

	cmd := exec.Command("mcporter", "auth", server)

	if reset {
		cmd.Args = append(cmd.Args, "--reset")
	}

	result, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Sprintf("Error authenticating with MCP server: %s\nDetails: %s", err, string(exitErr.Stderr))
		}
		return fmt.Sprintf("Error authenticating with MCP server: %v", err)
	}

	return string(result)
}

func configMCP(action string, server string, configPath string, url string) string {
	if !commandExists("mcporter") {
		return "Error: mcporter not found. Install it with: npm install -g @mcporter/cli"
	}

	cmd := exec.Command("mcporter", "config", action)

	switch action {
	case "get", "remove":
		if server == "" {
			return fmt.Sprintf("Error: server is required for action '%s'", action)
		}
		cmd.Args = append(cmd.Args, server)
	case "add":
		if server == "" || url == "" {
			return "Error: server and url are required for add action"
		}
		cmd.Args = append(cmd.Args, server, url)
	case "import":
		if configPath == "" {
			return "Error: config file path is required for import action"
		}
		cmd.Args = append(cmd.Args, configPath)
	}

	result, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Sprintf("Error managing MCP config: %s\nDetails: %s", err, string(exitErr.Stderr))
		}
		return fmt.Sprintf("Error managing MCP config: %v", err)
	}

	return string(result)
}
