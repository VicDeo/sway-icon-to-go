// Package sway provides a way to process the workspaces and renames them
// according to the name formatter and icon provider basing on the apps running on the workspaces

package sway

import (
	"context"
	"log"
	"strings"

	swayClient "github.com/joshuarubin/go-sway"
)

// Workspaces is a map of workspace number to workspace
type Workspaces map[int64]*Workspace

// WorkspaceNumByName is a map of workspace name to workspace number
type WorkspaceNumByName map[string]int64

// IconProvider is an interface that provides an icon for a given PID and node name
type IconProvider interface {
	GetIcon(pid *uint32, nodeName string) (string, bool)
}

// NameFormatter is an interface that formats a workspace name
type NameFormatter interface {
	Format(workspaceNumber int64, appIcons []string) string
}

// ProcessWorkspaces processes the workspaces and renames them according to the name formatter and icon provider
// basing on the apps running on the workspaces
func ProcessWorkspaces(ctx context.Context, iconProvider IconProvider, nameFormatter NameFormatter) error {
	var commands []string

	sway, err := swayClient.New(ctx)
	if err != nil {
		return err
	}

	// First pre-populate the workspace number by name
	// Sadly Node with type NodeWorkspace does not have
	// this property, so we need to get it from the sway workspaces
	workspaceNumByName := make(WorkspaceNumByName)
	swayWorkspaces, err := sway.GetWorkspaces(ctx)
	if err != nil {
		return err
	}
	for _, workspace := range swayWorkspaces {
		workspaceNumByName[workspace.Name] = workspace.Num
	}

	// Then traverse the tree and populate the workspaces map
	workspaces := make(Workspaces)
	tree, err := sway.GetTree(ctx)
	if err != nil {
		return err
	}
	traverseTree(tree, workspaceNumByName, workspaces, iconProvider)
	//log.Println(workspaceApps)

	// Then iterate over the workspaces and prepare the rename commands
	for _, workspace := range workspaces {
		newName := workspace.GetNewName(nameFormatter)
		if newName != workspace.Name {
			commands = append(commands, getRenameWorkspaceCommand(workspace.Name, newName))
		}
	}

	// Send all commands at once as there could be a mess otherwise
	command := strings.Join(commands, ";")
	log.Println(command)
	if _, error := sway.RunCommand(ctx, command); error != nil {
		log.Println(error)
		return error
	}
	return nil
}

// traverseTree traverses the tree and populates the workspaces map
func traverseTree(node *swayClient.Node, workspaceNumByName WorkspaceNumByName, workspaces Workspaces, iconProvider IconProvider) {
	switch node.Type {
	case swayClient.NodeWorkspace:
		for _, child := range node.Nodes {
			workspace := NewWorkspace(node.Name, workspaceNumByName[node.Name])
			workspaces[workspace.Number] = workspace

			traverseWorkspace(child, workspace.Number, workspaces, iconProvider)
		}
	default:
		for _, child := range node.Nodes {
			traverseTree(child, workspaceNumByName, workspaces, iconProvider)
		}
	}
}

// traverseWorkspace traverses the workspace and populates the workspaces map
func traverseWorkspace(node *swayClient.Node, workspaceNumber int64, workspaces Workspaces, iconProvider IconProvider) {
	if node.Type == swayClient.NodeCon || node.Type == swayClient.NodeFloatingCon {
		icon, found := iconProvider.GetIcon(node.PID, node.Name)
		if !found {
			log.Printf("No app mapping found for %v\n", node)
		}
		workspaces[workspaceNumber].AddAppIcon(icon)
	}
	for _, child := range node.Nodes {
		traverseWorkspace(child, workspaceNumber, workspaces, iconProvider)
	}

	for _, child := range node.FloatingNodes {
		traverseWorkspace(child, workspaceNumber, workspaces, iconProvider)
	}
}
