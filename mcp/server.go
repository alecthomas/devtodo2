package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/devtodo2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TaskInfo represents a task in JSON-serializable form.
type TaskInfo struct {
	Index       string     `json:"index"`
	Text        string     `json:"text"`
	Priority    string     `json:"priority"`
	Created     time.Time  `json:"created"`
	Completed   *time.Time `json:"completed,omitempty"`
	HasChildren bool       `json:"hasChildren,omitempty"`
	Children    []TaskInfo `json:"children,omitempty"`
}

// ListResult is returned by the list tool.
type ListResult struct {
	Title string     `json:"title,omitempty"`
	File  string     `json:"file"`
	Tasks []TaskInfo `json:"tasks"`
}

// AddResult is returned by the add tool.
type AddResult struct {
	Index string   `json:"index"`
	Task  TaskInfo `json:"task"`
}

// Server wraps an MCP server with devtodo2 tools.
type Server struct {
	server     *mcp.Server
	file       string
	legacyFile string
}

// NewServer creates a new MCP server for devtodo2.
func NewServer(file, legacyFile string) *Server {
	s := &Server{
		file:       file,
		legacyFile: legacyFile,
	}
	s.server = mcp.NewServer(&mcp.Implementation{
		Name:    "devtodo2",
		Version: "2.2.0",
	}, nil)
	s.registerTools()
	return s
}

// Run starts the server with stdio transport.
func (s *Server) Run(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

func (s *Server) load() (devtodo2.TaskList, error) {
	if f, err := os.Open(s.file); err == nil {
		defer f.Close()
		return devtodo2.NewJSONIO().Deserialize(f)
	}
	if f, err := os.Open(s.legacyFile); err == nil {
		defer f.Close()
		return devtodo2.NewLegacyIO().Deserialize(f)
	}
	return nil, nil
}

func (s *Server) save(tasks devtodo2.TaskList) error {
	path := s.file
	previous := path + "~"
	temp := path + "~~"

	f, err := os.Create(temp)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	if err := devtodo2.NewJSONIO().Serialize(f, tasks); err != nil {
		f.Close()
		os.Remove(temp)
		return fmt.Errorf("serialize: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(temp)
		return fmt.Errorf("close temp file: %w", err)
	}

	if _, err := os.Stat(path); err == nil {
		if err := os.Rename(path, previous); err != nil {
			return fmt.Errorf("backup: %w", err)
		}
	}

	return os.Rename(temp, path)
}

func (s *Server) registerTools() {
	type ListArgs struct {
		All bool `json:"all,omitempty" jsonschema:"include completed tasks"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_list",
		Description: "List tasks in the current directory's todo list",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ListArgs) (*mcp.CallToolResult, any, error) {
		return s.handleList(args.All)
	})

	type AddArgs struct {
		Text     string `json:"text" jsonschema:"task text"`
		Priority string `json:"priority,omitempty" jsonschema:"task priority (veryhigh/high/medium/low/verylow)"`
		Parent   string `json:"parent,omitempty" jsonschema:"parent task index to nest under (e.g. 1 or 1.2)"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_add",
		Description: "Add a new task to the todo list",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddArgs) (*mcp.CallToolResult, any, error) {
		return s.handleAdd(args.Text, args.Priority, args.Parent)
	})

	type DoneArgs struct {
		Indices []string `json:"indices" jsonschema:"task indices to mark as done (e.g. 1 or 2.1)"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_done",
		Description: "Mark tasks as completed",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args DoneArgs) (*mcp.CallToolResult, any, error) {
		return s.handleDone(args.Indices)
	})

	type UndoneArgs struct {
		Indices []string `json:"indices" jsonschema:"task indices to mark as not done"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_undone",
		Description: "Mark tasks as not completed",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args UndoneArgs) (*mcp.CallToolResult, any, error) {
		return s.handleUndone(args.Indices)
	})

	type EditArgs struct {
		Index    string `json:"index" jsonschema:"task index to edit"`
		Text     string `json:"text,omitempty" jsonschema:"new task text"`
		Priority string `json:"priority,omitempty" jsonschema:"new priority"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_edit",
		Description: "Edit a task's text and/or priority",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args EditArgs) (*mcp.CallToolResult, any, error) {
		return s.handleEdit(args.Index, args.Text, args.Priority)
	})

	type RemoveArgs struct {
		Indices []string `json:"indices" jsonschema:"task indices to remove"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_remove",
		Description: "Remove tasks from the list",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args RemoveArgs) (*mcp.CallToolResult, any, error) {
		return s.handleRemove(args.Indices)
	})

	type ReparentArgs struct {
		Index  string `json:"index" jsonschema:"task index to move"`
		Parent string `json:"parent,omitempty" jsonschema:"new parent index (empty for root)"`
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "todo_reparent",
		Description: "Move a task under a different parent",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ReparentArgs) (*mcp.CallToolResult, any, error) {
		return s.handleReparent(args.Index, args.Parent)
	})
}

func (s *Server) handleList(all bool) (*mcp.CallToolResult, any, error) {
	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		tasks = devtodo2.NewTaskList()
	}

	absPath, _ := filepath.Abs(s.file)
	result := ListResult{
		Title: tasks.Title(),
		File:  absPath,
		Tasks: collectTasks(tasks, "", all),
	}

	return jsonResult(result)
}

func (s *Server) handleAdd(text, priority, parent string) (*mcp.CallToolResult, any, error) {
	if text == "" {
		return nil, nil, fmt.Errorf("text is required")
	}

	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		tasks = devtodo2.NewTaskList()
	}

	var graft devtodo2.TaskNode = tasks
	if parent != "" {
		graft = tasks.Find(parent)
		if graft == nil {
			return nil, nil, fmt.Errorf("parent task %q not found", parent)
		}
	}

	p := devtodo2.MEDIUM
	if priority != "" {
		p = devtodo2.PriorityFromString(priority)
	}

	task := graft.Create(text, p)
	if err := s.save(tasks); err != nil {
		return nil, nil, fmt.Errorf("save: %w", err)
	}

	result := AddResult{
		Index: taskIndex(task),
		Task:  taskToInfo(task),
	}
	return jsonResult(result)
}

func (s *Server) handleDone(indices []string) (*mcp.CallToolResult, any, error) {
	if len(indices) == 0 {
		return nil, nil, fmt.Errorf("at least one index is required")
	}

	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		return nil, nil, fmt.Errorf("no task list found")
	}

	var completed []TaskInfo
	for _, idx := range indices {
		task := tasks.Find(idx)
		if task == nil {
			return nil, nil, fmt.Errorf("task %q not found", idx)
		}
		task.SetCompleted()
		completed = append(completed, taskToInfo(task))
	}

	if err := s.save(tasks); err != nil {
		return nil, nil, fmt.Errorf("save: %w", err)
	}

	return jsonResult(map[string]any{"completed": completed})
}

func (s *Server) handleUndone(indices []string) (*mcp.CallToolResult, any, error) {
	if len(indices) == 0 {
		return nil, nil, fmt.Errorf("at least one index is required")
	}

	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		return nil, nil, fmt.Errorf("no task list found")
	}

	var reopened []TaskInfo
	for _, idx := range indices {
		task := tasks.Find(idx)
		if task == nil {
			return nil, nil, fmt.Errorf("task %q not found", idx)
		}
		task.SetCompletionTime(time.Time{})
		reopened = append(reopened, taskToInfo(task))
	}

	if err := s.save(tasks); err != nil {
		return nil, nil, fmt.Errorf("save: %w", err)
	}

	return jsonResult(map[string]any{"reopened": reopened})
}

func (s *Server) handleEdit(index, text, priority string) (*mcp.CallToolResult, any, error) {
	if index == "" {
		return nil, nil, fmt.Errorf("index is required")
	}
	if text == "" && priority == "" {
		return nil, nil, fmt.Errorf("at least one of text or priority is required")
	}

	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		return nil, nil, fmt.Errorf("no task list found")
	}

	task := tasks.Find(index)
	if task == nil {
		return nil, nil, fmt.Errorf("task %q not found", index)
	}

	if text != "" {
		task.SetText(text)
	}
	if priority != "" {
		task.SetPriority(devtodo2.PriorityFromString(priority))
	}

	if err := s.save(tasks); err != nil {
		return nil, nil, fmt.Errorf("save: %w", err)
	}

	return jsonResult(map[string]any{"task": taskToInfo(task)})
}

func (s *Server) handleRemove(indices []string) (*mcp.CallToolResult, any, error) {
	if len(indices) == 0 {
		return nil, nil, fmt.Errorf("at least one index is required")
	}

	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		return nil, nil, fmt.Errorf("no task list found")
	}

	var removed []string
	for _, idx := range indices {
		task := tasks.Find(idx)
		if task == nil {
			return nil, nil, fmt.Errorf("task %q not found", idx)
		}
		task.Delete()
		removed = append(removed, idx)
	}

	if err := s.save(tasks); err != nil {
		return nil, nil, fmt.Errorf("save: %w", err)
	}

	return jsonResult(map[string]any{"removed": removed})
}

func (s *Server) handleReparent(index, parent string) (*mcp.CallToolResult, any, error) {
	if index == "" {
		return nil, nil, fmt.Errorf("index is required")
	}

	tasks, err := s.load()
	if err != nil {
		return nil, nil, fmt.Errorf("load tasks: %w", err)
	}
	if tasks == nil {
		return nil, nil, fmt.Errorf("no task list found")
	}

	task := tasks.Find(index)
	if task == nil {
		return nil, nil, fmt.Errorf("task %q not found", index)
	}

	var below devtodo2.TaskNode = tasks
	if parent != "" {
		below = tasks.Find(parent)
		if below == nil {
			return nil, nil, fmt.Errorf("parent %q not found", parent)
		}
	}

	devtodo2.ReparentTask(task, below)

	if err := s.save(tasks); err != nil {
		return nil, nil, fmt.Errorf("save: %w", err)
	}

	return jsonResult(map[string]any{"task": taskToInfo(task)})
}

func taskIndex(task devtodo2.Task) string {
	var parts []string
	for n := devtodo2.TaskNode(task); n != nil && n.ID() >= 0; n = n.Parent() {
		if p := n.Parent(); p != nil {
			for i := 0; i < p.Len(); i++ {
				if p.At(i).Equal(n) {
					parts = append([]string{fmt.Sprintf("%d", i+1)}, parts...)
					break
				}
			}
		}
	}
	return strings.Join(parts, ".")
}

func taskToInfo(task devtodo2.Task) TaskInfo {
	info := TaskInfo{
		Index:       taskIndex(task),
		Text:        task.Text(),
		Priority:    task.Priority().String(),
		Created:     task.CreationTime(),
		HasChildren: task.Len() > 0,
	}
	if !task.CompletionTime().IsZero() {
		t := task.CompletionTime()
		info.Completed = &t
	}
	return info
}

func collectTasks(node devtodo2.TaskNode, prefix string, all bool) []TaskInfo {
	var tasks []TaskInfo
	for i := 0; i < node.Len(); i++ {
		task := node.At(i)
		if !all && !task.CompletionTime().IsZero() {
			continue
		}
		var idx string
		if prefix != "" {
			idx = fmt.Sprintf("%s.%d", prefix, i+1)
		} else {
			idx = fmt.Sprintf("%d", i+1)
		}
		info := TaskInfo{
			Index:       idx,
			Text:        task.Text(),
			Priority:    task.Priority().String(),
			Created:     task.CreationTime(),
			HasChildren: task.Len() > 0,
		}
		if !task.CompletionTime().IsZero() {
			t := task.CompletionTime()
			info.Completed = &t
		}
		info.Children = collectTasks(task, idx, all)
		tasks = append(tasks, info)
	}
	return tasks
}

func jsonResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}
