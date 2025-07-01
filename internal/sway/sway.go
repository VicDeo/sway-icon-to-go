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

type WorkspaceApps map[string][]swayClient.Node

func ProcessWorkspaces(ctx context.Context) error {
	var commands []string

	log.Println("Processing workspaces")
	sway, err := swayClient.New(ctx)
	if err != nil {
		return err
	}

	workspaceApps := make(WorkspaceApps)
	tree, err := sway.GetTree(ctx)
	if err != nil {
		return err
	}
	traverseTree(tree, workspaceApps)
	//log.Println(workspaceApps)

	for workspaceName, workspaceNodes := range workspaceApps {
		workspaceNumber := getWorkspaceNumber(ctx, sway, workspaceName)
		newName := getNewWorkspaceName(workspaceNumber, workspaceNodes)
		if newName != workspaceName {
			commands = append(commands, getRenameWorkspaceCommand(workspaceName, newName))
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

func traverseTree(node *swayClient.Node, workspaceApps WorkspaceApps) {
	switch node.Type {
	case swayClient.NodeWorkspace:
		for _, child := range node.Nodes {
			traverseWorkspace(child, node.Name, workspaceApps)
		}
	default:
		for _, child := range node.Nodes {
			traverseTree(child, workspaceApps)
		}
	}
}

func traverseWorkspace(node *swayClient.Node, workspaceName string, workspaceApps WorkspaceApps) {
	if node.Type == swayClient.NodeCon || node.Type == swayClient.NodeFloatingCon {
		workspaceApps[workspaceName] = append(workspaceApps[workspaceName], *node)
	}
	for _, child := range node.Nodes {
		traverseWorkspace(child, workspaceName, workspaceApps)
	}

	for _, child := range node.FloatingNodes {
		traverseWorkspace(child, workspaceName, workspaceApps)
	}
}

func getWorkspaceNumber(ctx context.Context, sway swayClient.Client, workspaceName string) int64 {
	workspaces, err := sway.GetWorkspaces(ctx)
	if err != nil {
		return 0
	}
	for _, workspace := range workspaces {
		if workspace.Name == workspaceName {
			return workspace.Num
		}
	}
	return 0
}

func getNewWorkspaceName(workspaceNumber int64, workspaceNodes []swayClient.Node) string {
	apps := []string{}
	for _, app := range workspaceNodes {
		// Ignore ghost nodes
		if app.Name == "" {
			continue
		}
		apps = append(apps, getAppTitle(app))
	}
	return config.BuildName(workspaceNumber, apps)
}

func getAppTitle(app swayClient.Node) string {
	name, err := getExecutableName(app.PID)
	log.Println("executable name is", name)
	if err != nil || name == "" {
		fmt.Println(err)
		name = app.Name
	}
	icon := config.GetAppIcon(name)
	if !config.IsNoMatchIcon(icon) {
		return icon
	}
	// No app mapping found, use the app name
	fmt.Println("No app mapping found for ", app)
	return config.TrimAppName(name)
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
