package validation

import (
	"regexp"
	"strings"
	"time"
	"todo/internal/domain/model"
)

type MacroResult struct {
	Title    string
	Priority *model.TaskPriority
	Deadline *time.Time
}

func ParseTaskMacros(title string) MacroResult {
	result := MacroResult{Title: title}

	priorityMap := map[string]model.TaskPriority{
		"!1": model.PriorityCritical,
		"!2": model.PriorityHigh,
		"!3": model.PriorityMedium,
		"!4": model.PriorityLow,
	}
	for macro, priority := range priorityMap {
		if strings.Contains(result.Title, macro) {
			result.Priority = &priority
			result.Title = strings.ReplaceAll(result.Title, macro, "")
			break
		}
	}

	re := regexp.MustCompile(`!before\s+(\d{2}[.-]\d{2}[.-]\d{4})`)
	matches := re.FindStringSubmatch(result.Title)
	if len(matches) == 2 {
		dateStr := matches[1]
		layouts := []string{"02.01.2006", "02-01-2006"}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, dateStr); err == nil {
				result.Deadline = &t
				break
			}
		}
		result.Title = strings.Replace(result.Title, matches[0], "", 1)
	}

	result.Title = strings.TrimSpace(result.Title)
	return result
}
