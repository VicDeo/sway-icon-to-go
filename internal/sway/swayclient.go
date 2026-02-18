package sway

import (
	"context"
	"log/slog"

	sc "github.com/joshuarubin/go-sway"
)

const (
	// Scratchpad workspace name is "__i3_scratch".
	// See https://pkg.go.dev/github.com/joshuarubin/go-sway@v1.2.0#Node for more details.
	ScratchpadWorkspaceName = "__i3_scratch"
)

// SwayClient is an interface that provides a way to interact with the Sway window manager.
type SwayClient interface {
	CollectWorkspaces() (Workspaces, error)
	RenameWorkspaces(workspaces Workspaces, nameFormatter NameFormatter) error
}

// WorkspaceNumByName is a map of workspace name to workspace number.
type WorkspaceNumByName map[string]int64

// NewSwayClient creates a new SwayClient instance.
func NewSwayClient(ctx context.Context) (SwayClient, error) {
	client, err := sc.New(ctx)
	if err != nil {
		return nil, err
	}

	s := &swayClient{
		ctx:    ctx,
		client: client,
	}

	// First pre-populate the workspace number by name
	s.workspaceNumByName = make(WorkspaceNumByName)
	swayWorkspaces, err := s.client.GetWorkspaces(s.ctx)
	if err != nil {
		return nil, err
	}
	for _, workspace := range swayWorkspaces {
		s.workspaceNumByName[workspace.Name] = workspace.Num
	}

	return s, nil
}

// swayClient is a struct that implements the SwayClient interface.
type swayClient struct {
	ctx    context.Context
	client sc.Client
	// workspaceNumByName is a map of workspace name to workspace number
	// Sadly Node with type NodeWorkspace does not have
	// this property, so we need to get it from the sway workspaces
	workspaceNumByName WorkspaceNumByName
}

// CollectWorkspaces collects the workspaces from the Sway window manager.
func (s *swayClient) CollectWorkspaces() (Workspaces, error) {
	workspaces := make(Workspaces, 0)
	tree, err := s.client.GetTree(s.ctx)
	if err != nil {
		return nil, err
	}

	s.traverseTree(tree, workspaces)
	return workspaces, nil
}

// RenameWorkspaces renames the workspaces.
func (s *swayClient) RenameWorkspaces(workspaces Workspaces, nameFormatter NameFormatter) error {
	renameCommand := workspaces.ToRenameCommand(nameFormatter)
	if renameCommand == "" {
		// No changes to the workspaces, so we can return early
		return nil
	}
	// Send all commands at once as there could be a mess otherwise
	if _, err := s.client.RunCommand(s.ctx, renameCommand); err != nil {
		slog.Error("Error while renaming workspaces", "error", err)
		return err
	}
	return nil
}

// traverseTree traverses the tree and populates the initial workspaces map.
func (s *swayClient) traverseTree(node *sc.Node, workspaces Workspaces) {
	switch node.Type {
	case sc.NodeWorkspace:
		if node.Name == ScratchpadWorkspaceName {
			slog.Debug("Ignoring scratchpad workspace", "name", node.Name)
			return
		}

		workspaceNum, ok := s.workspaceNumByName[node.Name]
		if !ok {
			// Workspace not found in workspaceNumByName, so we skip it
			slog.Warn("Workspace not found in workspaceNumByName", "name", node.Name)
			return
		}

		workspace := NewWorkspace(node.Name, workspaceNum)
		workspaces[workspace.Number] = workspace
		for _, child := range node.Nodes {
			s.traverseWorkspace(child, workspace.Number, workspaces)
		}
	default:
		for _, child := range node.Nodes {
			s.traverseTree(child, workspaces)
		}
	}
}

// traverseWorkspace traverses the workspace and populates the workspaces map.
func (s *swayClient) traverseWorkspace(node *sc.Node, workspaceNumber int64, workspaces Workspaces) {
	if node.Type == sc.NodeCon || node.Type == sc.NodeFloatingCon {
		// Ignore ghost nodes that we can't resolve anyway
		if !(node.PID == nil && node.Name == "") {
			windowInfo := WindowInfo{
				PID:   node.PID,
				Title: node.Name,
			}
			workspaces[workspaceNumber].AddWindow(windowInfo)
		}
	}
	for _, child := range node.Nodes {
		s.traverseWorkspace(child, workspaceNumber, workspaces)
	}

	for _, child := range node.FloatingNodes {
		s.traverseWorkspace(child, workspaceNumber, workspaces)
	}
}
