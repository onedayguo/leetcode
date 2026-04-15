package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	session    string
	csrfToken  string
}

type TopQuestion struct {
	QuestionID string
	FrontendID string
	Title      string
	TitleSlug  string
	Difficulty string
}

type QuestionDetail struct {
	QuestionID         string
	QuestionFrontendID string
	Title              string
	TitleSlug          string
	ContentHTML        string
	Difficulty         string
	GoSnippet          string
}

type SubmissionResult struct {
	State             string  `json:"state"`
	StatusMsg         string  `json:"status_msg"`
	StatusCode        int     `json:"status_code"`
	RunSuccess        bool    `json:"run_success"`
	TotalCorrect      int     `json:"total_correct"`
	TotalTestcases    int     `json:"total_testcases"`
	Runtime           any     `json:"status_runtime"`
	Memory            any     `json:"memory"`
	RuntimePercentile float64 `json:"runtimePercentile,omitempty"`
	MemoryPercentile  float64 `json:"memoryPercentile,omitempty"`
	CompileError      string  `json:"compile_error"`
	RuntimeError      string  `json:"runtime_error"`
	LastTestcase      string  `json:"last_testcase"`
}

func NewClient(baseURL, session, csrfToken string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		session:   session,
		csrfToken: csrfToken,
	}
}

func (c *Client) FetchTop100(limit int) ([]TopQuestion, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	isCN := strings.Contains(c.baseURL, "leetcode.cn")
	if isCN {
		return c.fetchHot100CN(limit)
	}
	return c.fetchTop100COM(limit)
}

func (c *Client) fetchTop100COM(limit int) ([]TopQuestion, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/problems/top-100-liked/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if c.session != "" || c.csrfToken != "" {
		req.Header.Set("Cookie", "LEETCODE_SESSION="+c.session+"; csrftoken="+c.csrfToken)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch top100 failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	var payload struct {
		StatStatusPairs []struct {
			Stat struct {
				QuestionID         int    `json:"question_id"`
				QuestionTitle      string `json:"question__title"`
				QuestionTitleSlug  string `json:"question__title_slug"`
				FrontendQuestionID int    `json:"frontend_question_id"`
			} `json:"stat"`
			Difficulty struct {
				Level int `json:"level"`
			} `json:"difficulty"`
		} `json:"stat_status_pairs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make([]TopQuestion, 0, limit)
	for _, p := range payload.StatStatusPairs {
		if len(out) >= limit {
			break
		}
		if p.Stat.QuestionTitleSlug == "" || p.Stat.FrontendQuestionID == 0 {
			continue
		}
		out = append(out, TopQuestion{
			QuestionID: fmt.Sprintf("%d", p.Stat.QuestionID),
			FrontendID: fmt.Sprintf("%d", p.Stat.FrontendQuestionID),
			Title:      p.Stat.QuestionTitle,
			TitleSlug:  p.Stat.QuestionTitleSlug,
			Difficulty: diffFromLevel(p.Difficulty.Level),
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("top100 list is empty from %s; try setting LEETCODE_SESSION/LEETCODE_CSRF_TOKEN", c.baseURL)
	}
	return out, nil
}

func (c *Client) fetchHot100CN(limit int) ([]TopQuestion, error) {
	query := `query problemsetQuestionList($categorySlug: String, $skip: Int, $limit: Int, $filters: QuestionListFilterInput) {
  problemsetQuestionList(categorySlug: $categorySlug, skip: $skip, limit: $limit, filters: $filters) {
    total
    questions {
      frontendQuestionId
      title
      titleSlug
      difficulty
    }
  }
}`
	vars := map[string]any{
		"categorySlug": "all-code-essentials",
		"skip":         0,
		"limit":        limit,
		"filters": map[string]any{
			"listId": "2cktkvj",
		},
	}

	var resp struct {
		Data struct {
			ProblemsetQuestionList struct {
				Questions []struct {
					FrontendID string `json:"frontendQuestionId"`
					Title      string `json:"title"`
					TitleSlug  string `json:"titleSlug"`
					Difficulty string `json:"difficulty"`
				} `json:"questions"`
			} `json:"problemsetQuestionList"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := c.postGraphQL(query, vars, &resp); err != nil {
		return nil, err
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", resp.Errors[0].Message)
	}

	out := make([]TopQuestion, 0, len(resp.Data.ProblemsetQuestionList.Questions))
	for _, q := range resp.Data.ProblemsetQuestionList.Questions {
		out = append(out, TopQuestion{
			QuestionID: q.TitleSlug,
			FrontendID: q.FrontendID,
			Title:      q.Title,
			TitleSlug:  q.TitleSlug,
			Difficulty: q.Difficulty,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("hot100 list is empty from %s; try setting LEETCODE_SESSION/LEETCODE_CSRF_TOKEN", c.baseURL)
	}
	return out, nil
}

func (c *Client) FetchQuestionDetail(slug string) (QuestionDetail, error) {
	query := `query questionData($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    questionId
    questionFrontendId
    title
    titleSlug
    content
    difficulty
    codeSnippets {
      lang
      langSlug
      code
    }
  }
}`
	var resp struct {
		Data struct {
			Question struct {
				QuestionID         string `json:"questionId"`
				QuestionFrontendID string `json:"questionFrontendId"`
				Title              string `json:"title"`
				TitleSlug          string `json:"titleSlug"`
				Content            string `json:"content"`
				Difficulty         string `json:"difficulty"`
				CodeSnippets       []struct {
					LangSlug string `json:"langSlug"`
					Code     string `json:"code"`
				} `json:"codeSnippets"`
			} `json:"question"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := c.postGraphQL(query, map[string]string{"titleSlug": slug}, &resp); err != nil {
		return QuestionDetail{}, err
	}
	if len(resp.Errors) > 0 {
		return QuestionDetail{}, fmt.Errorf("graphql error: %s", resp.Errors[0].Message)
	}
	if resp.Data.Question.QuestionID == "" {
		return QuestionDetail{}, fmt.Errorf("question not found: %s", slug)
	}

	goSnippet := ""
	for _, s := range resp.Data.Question.CodeSnippets {
		if s.LangSlug == "golang" {
			goSnippet = s.Code
			break
		}
	}

	return QuestionDetail{
		QuestionID:         resp.Data.Question.QuestionID,
		QuestionFrontendID: resp.Data.Question.QuestionFrontendID,
		Title:              resp.Data.Question.Title,
		TitleSlug:          resp.Data.Question.TitleSlug,
		ContentHTML:        resp.Data.Question.Content,
		Difficulty:         resp.Data.Question.Difficulty,
		GoSnippet:          goSnippet,
	}, nil
}

func (c *Client) Submit(slug, questionID, lang, code string) (string, error) {
	payload := map[string]string{
		"lang":        lang,
		"question_id": questionID,
		"typed_code":  code,
	}
	b, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/problems/%s/submit/", c.baseURL, slug)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	c.attachAuthHeaders(req, slug)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rb, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("submit failed: status=%d body=%s", resp.StatusCode, string(rb))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var body struct {
		SubmissionID json.Number `json:"submission_id"`
	}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return "", fmt.Errorf("submit decode failed: %w, raw=%s", err, string(bodyBytes))
	}
	if body.SubmissionID == "" {
		return "", fmt.Errorf("submit failed: empty submission_id")
	}
	return body.SubmissionID.String(), nil
}

func (c *Client) CheckSubmission(submissionID string) (SubmissionResult, error) {
	url := fmt.Sprintf("%s/submissions/detail/%s/check/", c.baseURL, submissionID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return SubmissionResult{}, err
	}
	c.attachAuthHeaders(req, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SubmissionResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return SubmissionResult{}, fmt.Errorf("check failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	var result SubmissionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SubmissionResult{}, err
	}
	return result, nil
}

func (c *Client) WaitSubmission(submissionID string, timeout time.Duration) (SubmissionResult, error) {
	deadline := time.Now().Add(timeout)
	for {
		res, err := c.CheckSubmission(submissionID)
		if err != nil {
			return SubmissionResult{}, err
		}
		if res.State == "SUCCESS" {
			return res, nil
		}
		if time.Now().After(deadline) {
			return res, fmt.Errorf("wait timeout, latest state=%s", res.State)
		}
		time.Sleep(2 * time.Second)
	}
}

func (c *Client) FetchSubmissionPercentiles(submissionID string) (float64, float64, error) {
	id, err := strconv.Atoi(strings.TrimSpace(submissionID))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid submission id: %s", submissionID)
	}

	query := `query {
  submissionDetail(submissionId: ` + strconv.Itoa(id) + `) {
    runtimePercentile
    memoryPercentile
  }
}`

	var resp struct {
		Data struct {
			SubmissionDetail struct {
				RuntimePercentile float64 `json:"runtimePercentile"`
				MemoryPercentile  float64 `json:"memoryPercentile"`
			} `json:"submissionDetail"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.postGraphQL(query, map[string]any{}, &resp); err != nil {
		return 0, 0, err
	}
	if len(resp.Errors) > 0 {
		return 0, 0, fmt.Errorf("graphql error: %s", resp.Errors[0].Message)
	}

	return resp.Data.SubmissionDetail.RuntimePercentile, resp.Data.SubmissionDetail.MemoryPercentile, nil
}

func (c *Client) postGraphQL(query string, variables any, out any) error {
	payload := map[string]any{
		"query":     query,
		"variables": variables,
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/graphql/", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if c.csrfToken != "" {
		req.Header.Set("X-CSRFToken", c.csrfToken)
	}
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/")
	if c.session != "" || c.csrfToken != "" {
		req.Header.Set("Cookie", "LEETCODE_SESSION="+c.session+"; csrftoken="+c.csrfToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("graphql failed: status=%d body=%s", resp.StatusCode, string(rb))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) attachAuthHeaders(req *http.Request, slug string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRFToken", c.csrfToken)
	if slug != "" {
		req.Header.Set("Referer", c.baseURL+"/problems/"+slug+"/")
	} else {
		req.Header.Set("Referer", c.baseURL+"/")
	}
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if c.session != "" || c.csrfToken != "" {
		req.Header.Set("Cookie", "LEETCODE_SESSION="+c.session+"; csrftoken="+c.csrfToken)
	}
}

func diffFromLevel(level int) string {
	switch level {
	case 1:
		return "Easy"
	case 2:
		return "Medium"
	case 3:
		return "Hard"
	default:
		return "Unknown"
	}
}
