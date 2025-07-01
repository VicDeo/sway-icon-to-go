package sway

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sway-icon-to-go/internal/config"

	swayClient "github.com/joshuarubin/go-sway"
)

type Workspaces map[int64]*Workspace
type WorkspaceNumByName map[string]int64

func ProcessWorkspaces(ctx context.Context) error {
	var commands []string

	log.Println("Processing workspaces")
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
	traverseTree(tree, workspaceNumByName, workspaces)
	//log.Println(workspaceApps)

	// Then iterate over the workspaces and prepare the rename commands
	for _, workspace := range workspaces {
		newName := workspace.GetNewName()
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

func traverseTree(node *swayClient.Node, workspaceNumByName WorkspaceNumByName, workspaces Workspaces) {
	switch node.Type {
	case swayClient.NodeWorkspace:
		for _, child := range node.Nodes {
			workspace := NewWorkspace(node.Name, workspaceNumByName[node.Name])
			workspaces[workspace.Number] = workspace

			traverseWorkspace(child, workspace.Number, workspaces)
		}
	default:
		for _, child := range node.Nodes {
			traverseTree(child, workspaceNumByName, workspaces)
		}
	}
}

func traverseWorkspace(node *swayClient.Node, workspaceNumber int64, workspaces Workspaces) {
	if node.Type == swayClient.NodeCon || node.Type == swayClient.NodeFloatingCon {
		// Ignore ghost nodes
		if node.AppID != nil || node.Name != nil {
			workspaces[workspaceNumber].AddAppIcon(getAppIcon(*node))
		}
	}
	for _, child := range node.Nodes {
		traverseWorkspace(child, workspaceNumber, workspaces)
	}

	for _, child := range node.FloatingNodes {
		traverseWorkspace(child, workspaceNumber, workspaces)
	}
}

func getAppIcon(app swayClient.Node) string {
	name, err := getExecutableName(app.PID)
	if err != nil || name == "" {
		fmt.Println(err)
		name = app.Name
	}
	icon, found := config.GetAppIcon(name)
	if !found {
		fmt.Println("No app mapping found for ", app)
	}
	return icon
}

func getExecutableName(pid *uint32) (string, error) {
	pidInt := int(*pid)
	exePath := filepath.Join("/proc", strconv.Itoa(pidInt), "exe")
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Base(realPath), nil
}
