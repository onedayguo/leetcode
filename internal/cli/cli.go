package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"leetcode/internal/config"
	"leetcode/internal/leetcode"
	"leetcode/internal/project"
)

func Run(args []string, workspace string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	cfg, err := config.Load(workspace)
	if err != nil {
		return err
	}
	client := leetcode.NewClient(cfg.Site, cfg.Session, cfg.CSRFToken)

	switch args[0] {
	case "pull":
		return runPull(client, cfg.Site, workspace, args[1:])
	case "submit":
		return runSubmit(client, cfg, workspace, args[1:])
	case "status":
		return runStatus(client, cfg, args[1:])
	case "fixall":
		return runFixAll(workspace, args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runPull(client *leetcode.Client, site, workspace string, args []string) error {
	if len(args) == 0 || args[0] != "top100" {
		return errors.New("usage: pull top100 [--limit N] [--force]")
	}

	fs := flag.NewFlagSet("pull top100", flag.ContinueOnError)
	limit := fs.Int("limit", 100, "number of problems to pull")
	force := fs.Bool("force", false, "overwrite existing README/solution")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	items, err := client.FetchTop100(*limit)
	if err != nil {
		return err
	}

	fmt.Printf("pulling %d problems from %s\n", len(items), site)
	for i, q := range items {
		detail, err := client.FetchQuestionDetail(q.TitleSlug)
		if err != nil {
			fmt.Printf("[%d/%d] skip %s: %v\n", i+1, len(items), q.TitleSlug, err)
			continue
		}
		dir, err := project.SaveProblem(workspace, site, detail, *force)
		if err != nil {
			fmt.Printf("[%d/%d] save failed %s: %v\n", i+1, len(items), q.TitleSlug, err)
			continue
		}
		fmt.Printf("[%d/%d] %s -> %s\n", i+1, len(items), q.TitleSlug, dir)
	}

	return nil
}

func runSubmit(client *leetcode.Client, cfg config.Config, workspace string, args []string) error {
	if err := cfg.ValidateForSubmit(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("usage: submit <slug-or-folder> [--file path] [--lang golang] [--wait 120]")
	}

	key := args[0]
	fs := flag.NewFlagSet("submit", flag.ContinueOnError)
	file := fs.String("file", "", "path to solution file")
	lang := fs.String("lang", "golang", "language slug")
	waitSec := fs.Int("wait", 120, "seconds to wait for judging result")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	meta, dir, err := project.FindProblemBySlugOrFolder(workspace, key)
	if err != nil {
		return err
	}

	sourcePath := *file
	if sourcePath == "" {
		sourcePath = filepath.Join(dir, "solution.go")
	} else if !filepath.IsAbs(sourcePath) {
		sourcePath = filepath.Join(workspace, sourcePath)
	}

	codeBytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Extract submission code (remove package, imports, main() for LeetCode)
	submitCode := project.ExtractSubmissionCode(string(codeBytes))

	sid, err := client.Submit(meta.TitleSlug, meta.QuestionID, *lang, submitCode)
	if err != nil {
		return err
	}
	fmt.Printf("submitted: id=%s problem=%s\n", sid, meta.TitleSlug)

	res, err := client.WaitSubmission(sid, time.Duration(*waitSec)*time.Second)
	if err != nil {
		fmt.Printf("polling interrupted: %v\n", err)
		return nil
	}
	if rp, mp, err := client.FetchSubmissionPercentiles(sid); err == nil {
		res.RuntimePercentile = rp
		res.MemoryPercentile = mp
	}
	printResult(res)

	// Save submission record
	if err := project.SaveSubmissionRecord(dir, sid, res.State, res.Runtime, res.Memory, res.RuntimePercentile, res.MemoryPercentile, res.TotalCorrect, res.TotalTestcases, res.StatusCode, *lang, ""); err != nil {
		fmt.Printf("warning: failed to save submission record: %v\n", err)
	}
	return nil
}

func runStatus(client *leetcode.Client, cfg config.Config, args []string) error {
	if err := cfg.ValidateForSubmit(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("usage: status <submission-id>")
	}

	sid := strings.TrimSpace(args[0])
	if _, err := strconv.Atoi(sid); err != nil {
		return fmt.Errorf("invalid submission id: %s", sid)
	}

	res, err := client.CheckSubmission(sid)
	if err != nil {
		return err
	}
	printResult(res)
	return nil
}

func runFixAll(workspace string, args []string) error {
	fmt.Println("Fixing all solution files to enforce package solution...")
	return project.EnsureProblemPackageForExistingProblems(workspace)
}

func printResult(res leetcode.SubmissionResult) {
	fmt.Printf("state=%s status=%s\n", res.State, res.StatusMsg)
	runtime := metricString(res.Runtime)
	memory := metricString(res.Memory)

	// Append percentile info if available
	if res.RuntimePercentile > 0 {
		runtime = fmt.Sprintf("%s (Beats %.2f%% of submissions)", runtime, res.RuntimePercentile)
	}
	if res.MemoryPercentile > 0 {
		memory = fmt.Sprintf("%s (Beats %.2f%% of submissions)", memory, res.MemoryPercentile)
	}

	if runtime != "" || memory != "" {
		fmt.Printf("runtime=%s memory=%s\n", runtime, memory)
	}
	if res.TotalTestcases > 0 {
		fmt.Printf("tests=%d/%d\n", res.TotalCorrect, res.TotalTestcases)
	}
	if res.CompileError != "" {
		fmt.Printf("compile_error:\n%s\n", res.CompileError)
	}
	if res.RuntimeError != "" {
		fmt.Printf("runtime_error:\n%s\n", res.RuntimeError)
	}
	if res.LastTestcase != "" {
		fmt.Printf("last_testcase: %s\n", res.LastTestcase)
	}
}

func metricString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	default:
		return fmt.Sprint(x)
	}
}

func printUsage() {
	fmt.Println(`LeetCode local helper

Usage:
  go run . pull top100 [--limit 100] [--force]
  go run . submit <slug-or-folder> [--file path] [--lang golang] [--wait 120]
  go run . status <submission-id>
  go run . fixall

Auth for submit:
  export LEETCODE_SESSION=... 
  export LEETCODE_CSRF_TOKEN=...
Optional:
  export LC_SITE=https://leetcode.com
  # or create .leetcode.json in workspace root`)
}
