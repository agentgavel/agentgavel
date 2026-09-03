package oracle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ParseDirective extracts the probe directive from headers or chat messages.
func ParseDirective(h http.Header, messages []ChatMessage) (*Directive, error) {
	if raw := strings.TrimSpace(h.Get(HeaderProbeDirective)); raw != "" {
		var d Directive
		if err := json.Unmarshal([]byte(raw), &d); err != nil {
			return nil, fmt.Errorf("invalid %s: %w", HeaderProbeDirective, err)
		}
		if err := d.validate(); err != nil {
			return nil, err
		}
		return &d, nil
	}
	for _, m := range messages {
		if !strings.EqualFold(m.Role, "system") {
			continue
		}
		body := strings.TrimSpace(m.Content)
		if !strings.HasPrefix(body, SystemPrefixProbeDirective) {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(body, SystemPrefixProbeDirective))
		var d Directive
		if err := json.Unmarshal([]byte(raw), &d); err != nil {
			return nil, fmt.Errorf("invalid system probe directive: %w", err)
		}
		if err := d.validate(); err != nil {
			return nil, err
		}
		return &d, nil
	}
	return nil, errMissingDirective
}

func (d *Directive) validate() error {
	if d == nil || strings.TrimSpace(d.ToolName) == "" {
		return fmt.Errorf("tool_name is required")
	}
	if d.Arguments == nil {
		d.Arguments = map[string]any{}
	}
	return nil
}

var errMissingDirective = fmt.Errorf("missing probe directive")

// ChatMessage is a minimal OpenAI/Anthropic-shaped message.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Handler serves Oracle HTTP endpoints.
type Handler struct {
	Mux *http.ServeMux
}

// NewHandler registers OpenAI and Anthropic-shaped routes plus healthz.
func NewHandler() *Handler {
	h := &Handler{Mux: http.NewServeMux()}
	h.Mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	h.Mux.HandleFunc("/v1/chat/completions", h.openAIChatCompletions)
	h.Mux.HandleFunc("/v1/messages", h.anthropicMessages)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Mux.ServeHTTP(w, r)
}

type openAIRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

func (h *Handler) openAIChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}
	var req openAIRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	dir, err := ParseDirective(r.Header, req.Messages)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"message": err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}
	argsJSON, _ := json.Marshal(dir.Arguments)
	writeJSON(w, http.StatusOK, map[string]any{
		"id":     "chatcmpl-agentgavel-oracle",
		"object": "chat.completion",
		"model":  req.Model,
		"choices": []any{
			map[string]any{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": nil,
					"tool_calls": []any{
						map[string]any{
							"id":   "call_oracle_1",
							"type": "function",
							"function": map[string]any{
								"name":      dir.ToolName,
								"arguments": string(argsJSON),
							},
						},
					},
				},
				"finish_reason": "tool_calls",
			},
		},
	})
}

type anthropicRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	System   any           `json:"system"`
}

func (h *Handler) anthropicMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}
	var req anthropicRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	msgs := append([]ChatMessage{}, req.Messages...)
	switch s := req.System.(type) {
	case string:
		if s != "" {
			msgs = append([]ChatMessage{{Role: "system", Content: s}}, msgs...)
		}
	case []any:
		// ignore structured system blocks for directive extraction beyond string form
	}
	dir, err := ParseDirective(r.Header, msgs)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"type": "error",
			"error": map[string]any{
				"type":    "invalid_request_error",
				"message": err.Error(),
			},
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":    "msg_agentgavel_oracle",
		"type":  "message",
		"role":  "assistant",
		"model": req.Model,
		"content": []any{
			map[string]any{
				"type":  "tool_use",
				"id":    "toolu_oracle_1",
				"name":  dir.ToolName,
				"input": dir.Arguments,
			},
		},
		"stop_reason": "tool_use",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
