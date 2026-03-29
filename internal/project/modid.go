package project

import (
	"fmt"
	"strings"
)

func ValidateModID(modID string) error {
	modID = strings.TrimSpace(modID)
	if modID == "" {
		return fmt.Errorf("mod ID is required")
	}
	if len(modID) < 2 {
		return fmt.Errorf("mod ID must be at least 2 characters")
	}

	for i, r := range modID {
		isLower := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		isSeparator := r == '_' || r == '-'

		if i == 0 && !isLower {
			return fmt.Errorf("mod ID must start with a lowercase letter")
		}
		if !isLower && !isDigit && !isSeparator {
			return fmt.Errorf("mod ID can only contain lowercase letters, numbers, underscores, and hyphens")
		}
	}

	return nil
}

func NormalizeAutoModID(modName string) string {
	modID := strings.ToLower(modName)
	modID = strings.ReplaceAll(modID, " ", "_")
	modID = strings.ReplaceAll(modID, "-", "_")
	modID = strings.ReplaceAll(modID, ".", "_")

	var result strings.Builder
	for _, r := range modID {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}

	modID = strings.Trim(result.String(), "_")
	if modID == "" {
		return "mymod"
	}
	if modID[0] >= '0' && modID[0] <= '9' {
		modID = "mod_" + modID
	}
	if len(modID) < 2 {
		modID += "_mod"
	}

	return modID
}
