package ai

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Summarize 生成摘要和关键词：优先尝试真实LLM，失败时回退到本地实现
func Summarize(text string) (string, []string) {
	clean := strings.TrimSpace(text)
	if clean == "" {
		return "", nil
	}

	// 先尝试LLM（如果配置了密钥）
	if summary, keywords, err := SummarizeWithLLM(clean); err == nil && (summary != "" || len(keywords) > 0) {
		return summary, keywords
	}

	// 本地简易摘要与关键词提取（降级）
	// 句子分割（中文和英文常见标点）
	sentences := splitSentences(clean)
	// 摘要：取前3句，最多300字符
	summary := strings.Join(sentences[:min(3, len(sentences))], "")
	if len(summary) > 300 {
		summary = summary[:300]
	}

	// 关键词：简单频次统计，过滤短词和纯数字，取前8
	tokens := tokenize(clean)
	freq := make(map[string]int)
	for _, t := range tokens {
		if len([]rune(t)) < 2 {
			continue
		}
		if isNumeric(t) {
			continue
		}
		freq[strings.ToLower(t)]++
	}
	type kv struct {
		k string
		v int
	}
	var arr []kv
	for k, v := range freq {
		arr = append(arr, kv{k, v})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].v > arr[j].v })
	keywords := make([]string, 0, min(8, len(arr)))
	for i := 0; i < min(8, len(arr)); i++ {
		keywords = append(keywords, arr[i].k)
	}

	return summary, keywords
}

func splitSentences(text string) []string {
	seps := []rune{'。', '！', '？', '.', '!', '?', '；', ';'}
	var res []string
	start := 0
	runes := []rune(text)
	for i, r := range runes {
		for _, s := range seps {
			if r == s {
				res = append(res, string(runes[start : i+1]))
				start = i + 1
				break
			}
		}
	}
	if start < len(runes) {
		res = append(res, string(runes[start:]))
	}
	return compactStrings(res)
}

func tokenize(text string) []string {
	// 以空白和常见标点分割
	re := regexp.MustCompile(`[\s,.;:，。！？；：()（）\[\]{}"“”'` + "`" + `-]+`)
	parts := re.Split(text, -1)
	return compactStrings(parts)
}

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func compactStrings(arr []string) []string {
	res := make([]string, 0, len(arr))
	for _, s := range arr {
		s = strings.TrimSpace(s)
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

func min(a, b int) int {
	if a < b { return a }
	return b
}