package utilities

import "regexp"

func parseTableName(command string) string {
  insertRegex := regexp.MustCompile(`(?i)insert into\s+([a-zA-Z0-9_]+)`)
  updateRegex := regexp.MustCompile(`(?i)update\s+([a-zA-Z0-9_]+)`)

  // parse and return table name being modified
  if matches := insertRegex.FindStringSubmatch(command); len(matches) > 1 {
    return matches[1]
  }
  if matches := updateRegex.FindStringSubmatch(command); len(matches) > 1 {
    return matches[1]
  }
  return ""
}

func extractUserName(command string) string {
  userRegex := regexp.MustCompile(`(?i)user\s+([a-zA-Z0-9_]+)`)
  if matches := userRegex.FindStringSubmatch(command); len(matches) > 1 {
    return matches[1]
  }
  return ""
}
 q