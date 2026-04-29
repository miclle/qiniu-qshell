package operations

import (
	"context"
	"os"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// DeleteInfo holds parameters for deleting injection rules.
type DeleteInfo struct {
	RuleIDs []string
	Yes     bool
}

// Delete deletes one or more injection rules.
func Delete(info DeleteInfo) {
	if len(info.RuleIDs) == 0 {
		sbClient.PrintError("at least one rule ID is required")
		return
	}

	client, err := sbClient.NewInjectionRuleClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()

	if !info.Yes && !sbClient.Confirm("Are you sure you want to delete %d injection rule(s)?", len(info.RuleIDs)) {
		return
	}

	hasError := false
	for _, id := range info.RuleIDs {
		if dErr := client.DeleteInjectionRule(ctx, id); dErr != nil {
			sbClient.PrintError("delete injection rule %s failed: %v", id, dErr)
			hasError = true
			continue
		}
		sbClient.PrintSuccess("Injection rule %s deleted", id)
	}
	if hasError {
		os.Exit(1)
	}
}
