package main

import (
	"fmt"

	"github.com/adam-eques/mcpc/mcp"
)

// printResult writes a tool result's content blocks to stdout.
func printResult(res *mcp.CallToolResult) {
	for _, content := range res.Content {
		switch v := content.(type) {
		case mcp.TextContent:
			fmt.Println(v.Text)
		case mcp.ImageContent:
			fmt.Printf("[image %s, %d bytes base64]\n", v.MimeType, len(v.Data))
		default:
			fmt.Println("[unsupported content]")
		}
	}
}
