package main

import (
	"encoding/json"
	"testing"
)

func TestParseLSPRequest(t *testing.T) {
	tests := []struct {
		name    string
		request Request[json.RawMessage]
		want    LSPParams
		wantErr bool
	}{
		{
			name: "Valid Initialize Request",
			request: Request[json.RawMessage]{
				Jsonrpc: "2.0",
				ID:      1,
				Method:  "initialize",
				Params:  json.RawMessage(`{"processId":1234,"rootUri":"file:///path/to/project"}`),
			},
			want: InitializeParams{
				ProcessID: int64Ptr(1234),
				RootURI:   strPtr("file:///path/to/project"),
			},
			wantErr: false,
		},
		{
			name: "Valid DidOpenTextDocument Request",
			request: Request[json.RawMessage]{
				Jsonrpc: "2.0",
				ID:      2,
				Method:  "textDocument/didOpen",
				Params:  json.RawMessage(`{"textDocument":{"uri":"file:///path/to/file.txt","languageId":"go","version":1,"text":"content"}}`),
			},
			want: DidOpenTextDocumentParams{
				TextDocument: TextDocumentItem{
					URI:        "file:///path/to/file.txt",
					LanguageID: "go",
					Version:    1,
					Text:       "content",
				},
			},
			wantErr: false,
		},
		{
			name: "Unknown Method",
			request: Request[json.RawMessage]{
				Jsonrpc: "2.0",
				ID:      3,
				Method:  "unknownMethod",
				Params:  json.RawMessage(`{}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid JSON in Params",
			request: Request[json.RawMessage]{
				Jsonrpc: "2.0",
				ID:      4,
				Method:  "initialize",
				Params:  json.RawMessage(`{"processId":1234,"rootUri":invalid}`),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLSPRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLSPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !compareLSPParams(got, tt.want) {
				t.Errorf("ParseLSPRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func compareLSPParams(a, b LSPParams) bool {
	switch aVal := a.(type) {
	case InitializeParams:
		bVal, ok := b.(InitializeParams)
		if !ok {
			return false
		}
		return compareInitializeParams(aVal, bVal)
	case DidOpenTextDocumentParams:
		bVal, ok := b.(DidOpenTextDocumentParams)
		if !ok {
			return false
		}
		return compareDidOpenTextDocumentParams(aVal, bVal)
	// Add cases for other concrete types as needed
	default:
		return a == b // Compare other types directly
	}
}

func compareInitializeParams(a, b InitializeParams) bool {
	// Compare the fields of the InitializeParams struct
	return *a.ProcessID == *b.ProcessID &&
		*a.RootURI == *b.RootURI &&
		compareClientCapabilities(a.Capabilities, b.Capabilities)
	// Add comparisons for other fields as needed
}

func compareDidOpenTextDocumentParams(a, b DidOpenTextDocumentParams) bool {
	// Compare the fields of the DidOpenTextDocumentParams struct
	return a.TextDocument.URI == b.TextDocument.URI &&
		a.TextDocument.LanguageID == b.TextDocument.LanguageID &&
		a.TextDocument.Version == b.TextDocument.Version &&
		a.TextDocument.Text == b.TextDocument.Text
}

func compareClientCapabilities(a, b ClientCapabilities) bool {
	// Compare the fields of the ClientCapabilities struct
	// ...

	return true // Return true if a and b are equal
}

