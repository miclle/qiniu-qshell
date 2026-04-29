package operations

import (
	"context"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// DeleteInfo holds parameters for deleting templates.
type DeleteInfo struct {
	TemplateIDs []string // One or more template IDs to delete
	Yes         bool     // Skip confirmation
}

// Delete deletes one or more templates.
func Delete(info DeleteInfo) {
	if len(info.TemplateIDs) == 0 {
		id, ok := templateIDFromCwdConfig()
		if !ok {
			return
		}
		if id != "" {
			info.TemplateIDs = []string{id}
		}
	}
	if len(info.TemplateIDs) == 0 {
		sbClient.PrintError("at least one template ID is required (positional args or qshell.sandbox.toml)")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()

	if !info.Yes && !sbClient.Confirm("Are you sure you want to delete %d template(s)?", len(info.TemplateIDs)) {
		return
	}

	for _, id := range info.TemplateIDs {
		if dErr := client.DeleteTemplate(ctx, id); dErr != nil {
			sbClient.PrintError("delete template %s failed: %v", id, dErr)
			continue
		}
		sbClient.PrintSuccess("Template %s deleted", id)
	}
}
