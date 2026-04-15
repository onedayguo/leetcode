package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"leetcode/internal/leetcode"
)

var nonSlugChar = regexp.MustCompile(`[^a-z0-9-]+`)

type ProblemMeta struct {
	QuestionID string `json:"questionId"`
	FrontendID string `json:"frontendId"`
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
	Difficulty string `json:"difficulty"`
	URL        string `json:"url"`
}

func SaveProblem(workspace, site string, detail leetcode.QuestionDetail, force bool) (string, error) {
	folder := ProblemFolder(detail.QuestionFrontendID, detail.TitleSlug)
	dir := filepath.Join(workspace, "problems", folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	meta := ProblemMeta{
		QuestionID: detail.QuestionID,
		FrontendID: detail.QuestionFrontendID,
		Title:      detail.Title,
		TitleSlug:  detail.TitleSlug,
		Difficulty: detail.Difficulty,
		URL:        strings.TrimRight(site, "/") + "/problems/" + detail.TitleSlug + "/",
	}
	if err := writeJSON(filepath.Join(dir, "meta.json"), meta); err != nil {
		return "", err
	}

	readme := buildREADME(meta, detail.ContentHTML)
	if force || !exists(filepath.Join(dir, "README.md")) {
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
			return "", err
		}
	}

	solutionPath := filepath.Join(dir, "solution.go")
	if force || !exists(solutionPath) {
		snippet := strings.TrimSpace(detail.GoSnippet)
		if snippet == "" {
			snippet = "// No golang template returned by LeetCode for this problem.\n"
		}
		// Add package declaration if not present.
		if !strings.HasPrefix(strings.TrimSpace(snippet), "package ") {
			snippet = "package solution\n\n" + snippet
		}
		if err := os.WriteFile(solutionPath, []byte(snippet+"\n"), 0o644); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func FindProblemBySlugOrFolder(workspace, key string) (ProblemMeta, string, error) {
	base := filepath.Join(workspace, "problems")
	entries, err := os.ReadDir(base)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ProblemMeta{}, "", fmt.Errorf("problems directory not found, run pull first")
		}
		return ProblemMeta{}, "", err
	}

	candidates := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			candidates = append(candidates, filepath.Join(base, e.Name()))
		}
	}
	sort.Strings(candidates)

	lowerKey := strings.ToLower(strings.TrimSpace(key))
	for _, dir := range candidates {
		name := strings.ToLower(filepath.Base(dir))
		if name == lowerKey || strings.HasSuffix(name, "-"+lowerKey) {
			m, err := readMeta(dir)
			return m, dir, err
		}
	}

	for _, dir := range candidates {
		m, err := readMeta(dir)
		if err != nil {
			continue
		}
		if strings.EqualFold(m.TitleSlug, lowerKey) {
			return m, dir, nil
		}
	}
	return ProblemMeta{}, "", fmt.Errorf("problem not found: %s", key)
}

func ProblemFolder(frontendID, slug string) string {
	clean := strings.ToLower(strings.TrimSpace(slug))
	clean = strings.ReplaceAll(clean, "_", "-")
	clean = nonSlugChar.ReplaceAllString(clean, "-")
	clean = strings.Trim(clean, "-")
	if clean == "" {
		clean = "unknown"
	}
	if frontendID == "" {
		return clean
	}
	return frontendID + "-" + clean
}

func buildREADME(meta ProblemMeta, contentHTML string) string {
	return fmt.Sprintf("# %s. %s\n\n- Difficulty: %s\n- Link: %s\n\n## Description\n\n%s\n", meta.FrontendID, meta.Title, meta.Difficulty, meta.URL, contentHTML)
}

func readMeta(dir string) (ProblemMeta, error) {
	b, err := os.ReadFile(filepath.Join(dir, "meta.json"))
	if err != nil {
		return ProblemMeta{}, err
	}
	var m ProblemMeta
	if err := json.Unmarshal(b, &m); err != nil {
		return ProblemMeta{}, err
	}
	return m, nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type SubmissionRecord struct {
	SubmissionID      string      `json:"submissionId"`
	Status            string      `json:"status"`
	Runtime           interface{} `json:"runtime"`
	RuntimeMs         float64     `json:"runtimeMs,omitempty"`
	RuntimePercentile float64     `json:"runtimePercentile,omitempty"`
	Memory            interface{} `json:"memory"`
	MemoryMB          float64     `json:"memoryMB,omitempty"`
	MemoryPercentile  float64     `json:"memoryPercentile,omitempty"`
	TotalCorrect      int         `json:"totalCorrect"`
	TotalTestcases    int         `json:"totalTestcases"`
	StatusCode        int         `json:"statusCode"`
	Timestamp         int64       `json:"timestamp"`
	Language          string      `json:"language"`
	UpdatedAt         string      `json:"updatedAt,omitempty"`
	UpdatedBy         string      `json:"updatedBy,omitempty"`
	UpdateNote        string      `json:"updateNote,omitempty"`
	Note              string      `json:"note,omitempty"`
}

func SaveSubmissionRecord(problemDir, sid, status string, runtime, memory interface{}, runtimePercentile, memoryPercentile float64, totalCorrect, totalTestcases int, statusCode int, lang, note string) error {
	submissionsDir := filepath.Join(problemDir, "submissions")
	if err := os.MkdirAll(submissionsDir, 0o755); err != nil {
		return err
	}

	record := SubmissionRecord{
		SubmissionID:      sid,
		Status:            status,
		Runtime:           runtime,
		RuntimePercentile: runtimePercentile,
		Memory:            memory,
		MemoryPercentile:  memoryPercentile,
		TotalCorrect:      totalCorrect,
		TotalTestcases:    totalTestcases,
		StatusCode:        statusCode,
		Timestamp:         time.Now().Unix(),
		Language:          lang,
		Note:              note,
	}

	if runtimeMs, ok := parseNumeric(runtime); ok {
		record.RuntimeMs = runtimeMs
	}
	if memoryMB, ok := parseNumeric(memory); ok {
		record.MemoryMB = memoryMB
	}

	filename := filepath.Join(submissionsDir, sid+".json")
	if err := writeJSON(filename, record); err != nil {
		return err
	}

	// Also write to a summary file
	summaryPath := filepath.Join(problemDir, "submission_summary.txt")
	runtimeLine := fmt.Sprintf("Runtime: %v", runtime)
	if runtimePercentile > 0 {
		runtimeLine = fmt.Sprintf("Runtime: %v (Beats %.2f%% of submissions)", runtime, runtimePercentile)
	}
	memoryLine := fmt.Sprintf("Memory: %v", memory)
	if memoryPercentile > 0 {
		memoryLine = fmt.Sprintf("Memory: %v (Beats %.2f%% of submissions)", memory, memoryPercentile)
	}

	summary := fmt.Sprintf(`Latest Submission Summary
=========================
ID: %s
Status: %s
%s
%s
Tests: %d/%d
Language: %s
Time: %s
Note: %s
`, sid, status, runtimeLine, memoryLine, totalCorrect, totalTestcases, lang, time.Unix(record.Timestamp, 0).Format("2006-01-02 15:04:05"), note)

	return os.WriteFile(summaryPath, []byte(summary), 0o644)
}

func parseNumeric(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case string:
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// ExtractSubmissionCode removes package, imports, and main() from Go code,
// leaving only the problem-solving functions for LeetCode submission.
func ExtractSubmissionCode(fullCode string) string {
	lines := strings.Split(fullCode, "\n")
	var result []string
	inImportBlock := false
	inMain := false
	mainBraceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if inImportBlock {
			if strings.Contains(trimmed, ")") {
				inImportBlock = false
			}
			continue
		}

		if strings.HasPrefix(trimmed, "package ") {
			continue
		}
		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
			continue
		}
		if strings.HasPrefix(trimmed, "import ") {
			continue
		}

		if !inMain && strings.HasPrefix(trimmed, "func main(") {
			inMain = true
			mainBraceDepth = strings.Count(line, "{") - strings.Count(line, "}")
			if mainBraceDepth <= 0 {
				inMain = false
			}
			continue
		}

		if inMain {
			mainBraceDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if mainBraceDepth <= 0 {
				inMain = false
			}
			continue
		}

		result = append(result, line)
	}

	// Remove leading empty lines
	for len(result) > 0 && strings.TrimSpace(result[0]) == "" {
		result = result[1:]
	}

	// Remove trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n") + "\n"
}

// EnsureProblemPackageForExistingProblems enforces package solution in all solution.go files.
// This is a one-time utility function to update existing problem files.
func EnsureProblemPackageForExistingProblems(workspace string) error {
	problemsDir := filepath.Join(workspace, "problems")
	entries, err := os.ReadDir(problemsDir)
	if err != nil {
		return err
	}

	updated := 0
	skipped := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		solFile := filepath.Join(problemsDir, entry.Name(), "solution.go")
		content, err := os.ReadFile(solFile)
		if err != nil {
			continue // Skip if no solution.go
		}

		contentStr := string(content)
		trimmed := strings.TrimSpace(contentStr)

		// Already in desired package.
		if strings.HasPrefix(trimmed, "package solution") {
			skipped++
			continue
		}

		// Replace existing package declaration or prepend desired package.
		newContent := contentStr
		if strings.HasPrefix(trimmed, "package ") {
			lines := strings.Split(contentStr, "\n")
			replaced := false
			for i, line := range lines {
				lineTrimmed := strings.TrimSpace(line)
				if lineTrimmed == "" {
					continue
				}
				if strings.HasPrefix(lineTrimmed, "package ") {
					lines[i] = "package solution"
					replaced = true
				}
				break
			}
			if replaced {
				newContent = strings.Join(lines, "\n")
			}
		} else {
			newContent = "package solution\n\n" + contentStr
		}
		if err := os.WriteFile(solFile, []byte(newContent), 0o644); err != nil {
			continue
		}

		updated++
	}

	fmt.Printf("Enforced package solution on %d files, skipped %d files\n", updated, skipped)
	return nil
}
