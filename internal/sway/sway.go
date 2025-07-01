package sway

import (
	"context"
	"fmt"
	"log"
	"strings"

	swayClient "github.com/joshuarubin/go-sway"
)

type UnifiedWorkspaces []UnifiedWorkspace

func ProcessWorkspaces(ctx context.Context) error {
	var commands []string

	log.Println("Processing workspaces")
	sway, err := swayClient.New(ctx)
	if err != nil {
		return err
	}

	unifiedWorkspaces, err := getWorkspaces(sway, ctx)
	if err != nil {
		return err
	}
	// placement rules create a workspace with a name that has just digits only
	// Thus a newly created window will be placed to the workspace '2'
	// and we end up with the duplicated workspaces '2' and '2:icon'
	// Merge such workspaces on the first pass
	for _, w := range unifiedWorkspaces {
		fmt.Println(w.Workspace.Name, w.MergeWith)
		if w.Workspace.Num < 0 || w.MergeWith == nil {
			continue
		}
		targetName := w.MergeWith.Name
		for _, l := range getLeafNodes(w.Node) {
			commands = append(commands, getMoveAppToWorkspaceCommand(l.ID, targetName))
			fmt.Println(l.ID, targetName)

		}
		if w.Workspace.Focused {
			commands = append(commands, getFocusCommand(targetName))
		}
	}
	for _, w := range unifiedWorkspaces {
		if w.Workspace.Num < 0 || w.MergeWith != nil {
			continue
		}
		newName := w.GetNewName()
		if w.Workspace.Name != newName {
			commands = append(commands, getRenameWorkspaceCommand(w.Workspace.Name, newName))
		}
	}

	// Send all commands at once as there could be a mess otherwise
	command := strings.Join(commands, ";")
	fmt.Println(command)
	if _, error := sway.RunCommand(ctx, command); error != nil {
		fmt.Println(error)
		return error
	}
	return nil
}

func getWorkspaces(sway swayClient.Client, ctx context.Context) (UnifiedWorkspaces, error) {
	unifiedWorkspaces := []UnifiedWorkspace{}

	// Get the tree structure
	tree, err := sway.GetTree(ctx)
	if err != nil {
		return nil, err
	}

	// Get workspace metadata
	workspaces, err := sway.GetWorkspaces(ctx)
	if err != nil {
		return nil, err
	}

	// Find workspace nodes in the tree
	workspaceNodes := findWorkspaceNodes(tree)

	for _, w := range workspaces {
		unifiedWorkspaces = append(unifiedWorkspaces, UnifiedWorkspace{
			Workspace: w,
			Node:      getNodeByName(workspaceNodes, w.Name),
			MergeWith: nil,
		})
	}
	// Searching for duplicates e.g. '2' and '2:'
	for i := 0; i < len(unifiedWorkspaces); i++ {
		outer := unifiedWorkspaces[i]
		for j := i + 1; j < len(unifiedWorkspaces); j++ {
			inner := unifiedWorkspaces[j]
			if outer.IsDuplicateOf(inner) {
				if !outer.IsCustom() {
					unifiedWorkspaces[i].MergeWith = &inner.Workspace
				} else {
					unifiedWorkspaces[j].MergeWith = &outer.Workspace
				}
			}
		}
	}
	return unifiedWorkspaces, nil
}

// findWorkspaceNodes traverses the tree and returns all workspace nodes
func findWorkspaceNodes(root *swayClient.Node) []*swayClient.Node {
	var workspaceNodes []*swayClient.Node

	// Use the TraverseNodes method to find all workspace nodes
	root.TraverseNodes(func(n *swayClient.Node) bool {
		if n.Type == swayClient.NodeWorkspace {
			workspaceNodes = append(workspaceNodes, n)
		}
		return false // Continue traversing
	})

	return workspaceNodes
}

// getNodeByName finds a node by name in a slice of nodes
func getNodeByName(nodes []*swayClient.Node, name string) *swayClient.Node {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

// CollectUniqueAppNames traverses all workspaces and collects unique app names per workspace
func CollectUniqueAppNames(ctx context.Context) (map[string][]string, error) {
	sway, err := swayClient.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sway: %v", err)
	}

	// Get the tree structure
	tree, err := sway.GetTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %v", err)
	}

	// Get workspace metadata
	workspaces, err := sway.GetWorkspaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces: %v", err)
	}

	// Map to store workspace name -> unique app names
	workspaceApps := make(map[string][]string)

	// Find all workspace nodes in the tree
	workspaceNodes := findWorkspaceNodes(tree)

	// For each workspace, collect unique app names
	for _, workspace := range workspaces {
		workspaceNode := getNodeByName(workspaceNodes, workspace.Name)
		if workspaceNode == nil {
			continue
		}

		// Collect all leaf nodes (windows) in this workspace
		appNames := collectAppNamesFromNode(workspaceNode)

		// Store unique app names for this workspace
		workspaceApps[workspace.Name] = appNames
	}

	return workspaceApps, nil
}

// collectAppNamesFromNode recursively collects app names from a node and its children
func collectAppNamesFromNode(node *swayClient.Node) []string {
	var appNames []string
	appNameSet := make(map[string]bool)

	// Traverse all nodes to find leaf nodes (windows)
	node.TraverseNodes(func(n *swayClient.Node) bool {
		// Check if this is a leaf node (window)
		if len(n.Nodes) == 0 && len(n.FloatingNodes) == 0 {
			appName := getAppNameFromNode(n)
			if appName != "" && !appNameSet[appName] {
				appNameSet[appName] = true
				appNames = append(appNames, appName)
			}
		}
		return false // Continue traversing
	})

	return appNames
}

// getAppNameFromNode extracts the app name from a window node
func getAppNameFromNode(node *swayClient.Node) string {
	// Try to get app name from various properties in order of preference
	if node.AppID != nil && *node.AppID != "" {
		return *node.AppID
	}

	if node.WindowProperties != nil {
		// Try class first, then instance, then title
		if node.WindowProperties.Class != "" {
			return node.WindowProperties.Class
		}
		if node.WindowProperties.Instance != "" {
			return node.WindowProperties.Instance
		}
		if node.WindowProperties.Title != "" {
			return node.WindowProperties.Title
		}
	}

	// Fallback to node name
	if node.Name != "" {
		return node.Name
	}

	return ""
}

// getLeafNodes returns all leaf nodes from a sway Node
func getLeafNodes(node *swayClient.Node) []*swayClient.Node {
	var leafNodes []*swayClient.Node

	// Use the TraverseNodes method to find all leaf nodes
	node.TraverseNodes(func(n *swayClient.Node) bool {
		if len(n.Nodes) == 0 && len(n.FloatingNodes) == 0 {
			leafNodes = append(leafNodes, n)
		}
		return false // Continue traversing
	})

	return leafNodes
}
