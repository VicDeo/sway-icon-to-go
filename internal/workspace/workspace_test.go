package workspace

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nameFormatter struct {
}

func (nf *nameFormatter) Format(workspaceNumber int64, appIcons []string) string {
	return fmt.Sprintf("%d: %s", workspaceNumber, strings.Join(appIcons, "|"))
}

func TestWorkspace_ToRenameCommand(t *testing.T) {
	nameFormatter := &nameFormatter{}
	workspace := NewWorkspace("1: app1|app2|app3", 1)
	workspace.AddAppIcon("New app1")
	workspace.AddAppIcon("New app2")
	workspace.AddAppIcon("New app3")
	command := workspace.ToRenameCommand(nameFormatter)
	assert.Equal(t, "rename workspace \"1: app1|app2|app3\" to \"1: New app1|New app2|New app3\"", command)
}

func TestWorkspaces_ToRenameCommand(t *testing.T) {
	nameFormatter := &nameFormatter{}
	workspaces := Workspaces{}
	workspaces[1] = NewWorkspace("1: app1|app2|app3", 1)
	workspaces[1].AddAppIcon("New app1")
	workspaces[1].AddAppIcon("New app2")
	workspaces[1].AddAppIcon("New app3")
	workspaces[2] = NewWorkspace("2: app4|app5|app6", 2)
	workspaces[2].AddAppIcon("New app4")
	workspaces[2].AddAppIcon("New app5")
	workspaces[2].AddAppIcon("New app6")
	command := workspaces.ToRenameCommand(nameFormatter)
	assert.Equal(t, "rename workspace \"1: app1|app2|app3\" to \"1: New app1|New app2|New app3\";rename workspace \"2: app4|app5|app6\" to \"2: New app4|New app5|New app6\"", command)
}
