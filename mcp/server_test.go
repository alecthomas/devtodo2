package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPServer(t *testing.T) {
	// Create temp directory
	dir := t.TempDir()
	file := filepath.Join(dir, ".todo2")
	legacyFile := filepath.Join(dir, ".todo")

	// Create server
	server := NewServer(file, legacyFile)

	// Connect using in-memory transport
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0"}, nil)

	t1, t2 := mcp.NewInMemoryTransports()
	ctx := context.Background()

	serverSession, err := server.server.Connect(ctx, t1, nil)
	assert.NoError(t, err)
	defer serverSession.Close()

	clientSession, err := client.Connect(ctx, t2, nil)
	assert.NoError(t, err)
	defer clientSession.Close()

	// Test list tools
	tools, err := clientSession.ListTools(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(tools.Tools))

	toolNames := make(map[string]bool)
	for _, tool := range tools.Tools {
		toolNames[tool.Name] = true
	}
	assert.True(t, toolNames["todo_list"])
	assert.True(t, toolNames["todo_add"])
	assert.True(t, toolNames["todo_done"])

	// Test add task
	addResult, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "todo_add",
		Arguments: json.RawMessage(`{"text": "Test task 1"}`),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(addResult.Content))

	// Verify file was created
	_, err = os.Stat(file)
	assert.NoError(t, err)

	// Test list tasks
	listResult, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "todo_list",
		Arguments: json.RawMessage(`{}`),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(listResult.Content))

	// Parse list result
	textContent, ok := listResult.Content[0].(*mcp.TextContent)
	assert.True(t, ok)
	var lr ListResult
	err = json.Unmarshal([]byte(textContent.Text), &lr)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(lr.Tasks))
	assert.Equal(t, "Test task 1", lr.Tasks[0].Text)
	assert.Equal(t, "1", lr.Tasks[0].Index)

	// Test mark done
	doneResult, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "todo_done",
		Arguments: json.RawMessage(`{"indices": ["1"]}`),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(doneResult.Content))

	// Verify task is done (not in default list)
	listResult, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "todo_list",
		Arguments: json.RawMessage(`{}`),
	})
	assert.NoError(t, err)
	textContent, ok = listResult.Content[0].(*mcp.TextContent)
	assert.True(t, ok)
	err = json.Unmarshal([]byte(textContent.Text), &lr)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(lr.Tasks))

	// Verify task shows with all=true
	listResult, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "todo_list",
		Arguments: json.RawMessage(`{"all": true}`),
	})
	assert.NoError(t, err)
	textContent, ok = listResult.Content[0].(*mcp.TextContent)
	assert.True(t, ok)
	err = json.Unmarshal([]byte(textContent.Text), &lr)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(lr.Tasks))
}
