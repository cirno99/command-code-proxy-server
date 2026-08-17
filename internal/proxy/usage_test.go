package proxy

	"testing"
)

// 构造模拟 CommandCode 上游的流式响应，包含 totalUsage 与 inputTokenDetails
func mockCCStreamBody() io.ReadCloser {
	var sb strings.Builder
	sb.WriteString(`{"type":"text-delta","text":"你好，"}` + "\n")
	sb.WriteString(`{"type":"text-delta","text":"世界！"}` + "\n")
	sb.WriteString(`{"type":"finish","finishReason":"stop","totalUsage":{` +
		`"inputTokens":150,"outputTokens":30,` +
		`"inputTokenDetails":{"cacheReadTokens":100,"cacheWriteTokens":20,"noCacheTokens":30}}}` + "\n")
	return io.NopCloser(strings.NewReader(sb.String()))
}

// 流式路径：finish 事件必须透传带 usage 的 chunk（空 choices + usage 字段）
func TestStreamResponseForwardsUsage(t *testing.T) {
	p := NewProxy("test-key")
	ccResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       mockCCStreamBody(),
	}

	rec := httptest.NewRecorder()
	p.StreamResponse(rec, httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil),
		ccResp, "chatcmpl-test", "deepseek/deepseek-v4-pro", 1700000000)

	body := rec.Body.String()

	// 1. usage chunk 必须存在，且 choices 为空数组
	if !strings.Contains(body, `"usage":{"prompt_tokens":150,"completion_tokens":30,"total_tokens":180`) {
		t.Fatalf("流式响应缺少 usage chunk:\n%s", body)
	}
	// 2. 缓存命中明细必须透传
	if !strings.Contains(body, `"prompt_tokens_details":{"cached_tokens":100,"cache_write_tokens":20}`) {
		t.Fatalf("流式响应缺少缓存命中明细:\n%s", body)
	}
	// 3. finish chunk 与 [DONE] 仍然存在
	if !strings.Contains(body, `"finish_reason":"stop"`) {
		t.Fatalf("流式响应缺少 finish_reason:\n%s", body)
	}
	if !strings.Contains(body, "data: [DONE]") {
		t.Fatalf("流式响应缺少 [DONE] 结束标记:\n%s", body)
	}
	// 4. 文本内容必须正常透传
	if !strings.Contains(body, "你好，世界！") {
		t.Fatalf("流式响应缺少文本内容:\n%s", body)
	}
}

// 非流式路径：usage 与缓存明细必须出现在最终 JSON 响应中
func TestNonStreamResponseForwardsUsage(t *testing.T) {
	p := NewProxy("test-key")
	ccResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       mockCCStreamBody(),
	}

	rec := httptest.NewRecorder()
	p.NonStreamResponse(rec, ccResp, "chatcmpl-test", "deepseek/deepseek-v4-pro", 1700000000)

	body := rec.Body.String()

	if !strings.Contains(body, `"prompt_tokens":150`) {
		t.Fatalf("非流式响应缺少 prompt_tokens:\n%s", body)
	}
	if !strings.Contains(body, `"completion_tokens":30`) {
		t.Fatalf("非流式响应缺少 completion_tokens:\n%s", body)
	}
	if !strings.Contains(body, `"total_tokens":180`) {
		t.Fatalf("非流式响应缺少 total_tokens:\n%s", body)
	}
	if !strings.Contains(body, `"cached_tokens":100`) {
		t.Fatalf("非流式响应缺少 cached_tokens:\n%s", body)
	}
	if !strings.Contains(body, `"cache_write_tokens":20`) {
		t.Fatalf("非流式响应缺少 cache_write_tokens:\n%s", body)
	}
}

// 上游不返回 usage 时：不 panic，且响应中不包含 usage 相关字段
func TestNoUsageGraceful(t *testing.T) {
	p := NewProxy("test-key")
	body := io.NopCloser(strings.NewReader(
		`{"type":"text-delta","text":"hi"}` + "\n" +
			`{"type":"finish","finishReason":"stop"}` + "\n"))
	ccResp := &http.Response{StatusCode: http.StatusOK, Body: body}

	rec := httptest.NewRecorder()
	p.NonStreamResponse(rec, ccResp, "chatcmpl-test", "m", 1700000000)

	// 不 panic 即通过；usage 为 0 也应输出
	if !strings.Contains(rec.Body.String(), `"prompt_tokens":0`) {
		t.Fatalf("无 usage 时响应异常:\n%s", rec.Body.String())
	}
}
