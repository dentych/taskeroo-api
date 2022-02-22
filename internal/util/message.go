package util

import (
	"bytes"
)

func CommonTaskMessage(taskTitles []string) string {
	if taskTitles == nil {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("Fællesopgaver der skal udføres i dag:\n")
	for _, title := range taskTitles {
		buf.WriteString("• ")
		buf.WriteString(title)
		buf.WriteString("\n")
	}

	buf.WriteRune('\n')
	buf.WriteString("https://taskeroo.tychsen.me")

	return buf.String()
}

func AssignedTasksMessage(taskTitles []string) string {
	if taskTitles == nil {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("Du har følgende tildelte opgaver, som skal udføres i dag:\n")
	for _, title := range taskTitles {
		buf.WriteString("• ")
		buf.WriteString(title)
		buf.WriteString("\n")
	}

	buf.WriteRune('\n')
	buf.WriteString("https://taskeroo.tychsen.me")

	return buf.String()
}
