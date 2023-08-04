package main

type Request[T any] struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  T `json:"params,omitempty"`
}

type LSPParams interface {}
type LSPRequest = Request[LSPParams]

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type InitializeParams struct {
	ProcessID    *int64               `json:"processId,omitempty"`
	RootURI      *string              `json:"rootUri,omitempty"`
	Capabilities ClientCapabilities   `json:"capabilities"`
	// Additional fields as per the LSP specification
}

type ClientCapabilities struct {
	// Define the client capabilities here
}

type InitializeRequest Request[InitializeParams]
