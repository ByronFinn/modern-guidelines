// Package mcpserver implements the stdio MCP server of the modern-guidelines
// hub (PRD-0000 requirement 9): a read-only server that exposes the
// version-gated guideline data through two tools, list_guidelines and
// explain_guideline, so coding agents can pull up-to-date guidance despite
// their knowledge cutoff.
//
// The server makes no outbound network requests: it speaks newline-delimited
// JSON-RPC over the injected streams and answers from the in-process
// guideline data.
//
// Rendering stays single-source. The real list and explain rendering lives in
// internal/cli (which delegates to internal/guidelines); this package cannot
// import internal/cli, because the cli mcp subcommand imports this package —
// a direct dependency would be an import cycle. Instead the CLI injects its
// own implementations through [Deps] when it starts the server, so both entry
// points share one code path and their output stays byte-identical. This
// package owns only the MCP protocol layer: tool schemas, argument mapping,
// and the read-only annotations.
package mcpserver

import (
	"context"
	"errors"
	"io"
	"runtime/debug"
	"strings"

	"github.com/ByronFinn/modern-guidelines/internal/guidelines"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Deps carries the render entry points the server calls. The CLI injects its
// list and explain implementations; keeping them as function values avoids
// the import cycle described in the package comment while preventing any
// duplication of the rendering logic.
type Deps struct {
	// List renders the guidelines of one language for a resolved version
	// source, mirroring the CLI list command. Its parameters map one to one
	// onto the CLI flags: versionSpec → --version axis=ver[,...], filePath →
	// --file-path, language → --lang. At most one of versionSpec and filePath
	// is non-empty (the caller enforces the tool's conflict rules before
	// dispatching; the injected implementation re-validates through the CLI
	// flag parsing). The returned text ends with a single newline, exactly as
	// the CLI writes it.
	List func(versionSpec, filePath, language string) (string, error)
	// Explain renders the detail blocks for guideline ids of one language,
	// mirroring the CLI explain command: language → --lang, ids → the
	// positional references (each may be a bare id or a lang:id override).
	// language is never empty. The returned text ends with a single newline.
	Explain func(language string, ids []string) (string, error)
}

// Tool names advertised on tools/list.
const (
	ListToolName    = "list_guidelines"
	ExplainToolName = "explain_guideline"
)

// Tool descriptions state when an agent should call the tool (PRD
// requirement 9) and that the output matches the CLI.
const listToolDescription = "Call before writing or editing Python, Go, or TypeScript code to get the version-gated modern guidelines of the project. " +
	"Prefer passing only file_path: the language is inferred from the file extension and every axis version is resolved from the project's own manifests (go.mod, pyproject.toml, package.json, ...). " +
	"Alternatively pass version as comma-separated axis=version entries covering every axis of one language (python=3.12 or typescript=5.4,node=18). " +
	"language only overrides the extension-based inference when file_path has an unknown extension; it resolves no versions on its own. " +
	"The output starts with one provenance line per axis (version, source, fidelity) followed by one \"id: guideline\" line per applicable rule, newest first — identical to the mg list CLI."

const explainToolDescription = "Call to get the details behind guideline ids returned by list_guidelines: summary, details, before/after examples, references, and the autofix mapping. " +
	"Pass the language of the ids and the ids themselves; a lang:id entry overrides the language for that id. " +
	"The output is identical to the mg explain CLI."

// listArgs is the input schema of list_guidelines: every parameter optional,
// exactly one version source expected.
type listArgs struct {
	FilePath string `json:"file_path,omitempty" jsonschema:"Path of a project file; the recommended single argument: infers the language from the file extension and resolves every axis version from the project's manifests"`
	Language string `json:"language,omitempty" jsonschema:"Guideline dataset language (go, python, typescript) overriding the extension-based inference of file_path; resolves no versions on its own"`
	Version  string `json:"version,omitempty" jsonschema:"Explicit axis versions as comma-separated axis=version entries covering every axis of one language, e.g. python=3.12 or typescript=5.4,node=18"`
}

// explainArgs is the input schema of explain_guideline.
type explainArgs struct {
	Language string   `json:"language" jsonschema:"Language of the guideline ids (go, python, typescript); an explicit lang:id entry overrides it for that id"`
	IDs      []string `json:"ids" jsonschema:"Guideline ids returned by list_guidelines; each entry is a bare id or a lang:id reference"`
}

// NewServer builds the MCP server with both read-only tools registered. It
// panics when deps is incomplete: an unusable server must fail at startup,
// not on the first tool call.
func NewServer(deps Deps) *mcp.Server {
	if deps.List == nil || deps.Explain == nil {
		panic("mcpserver: Deps.List and Deps.Explain must both be non-nil")
	}
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "modern-guidelines",
		Version: serverVersion(),
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        ListToolName,
		Description: listToolDescription,
		Annotations: readOnlyAnnotations(),
	}, func(_ context.Context, _ *mcp.CallToolRequest, args listArgs) (*mcp.CallToolResult, any, error) {
		text, err := listText(deps, args)
		if err != nil {
			return nil, nil, err
		}
		return textResult(text), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        ExplainToolName,
		Description: explainToolDescription,
		Annotations: readOnlyAnnotations(),
	}, func(_ context.Context, _ *mcp.CallToolRequest, args explainArgs) (*mcp.CallToolResult, any, error) {
		text, err := explainText(deps, args)
		if err != nil {
			return nil, nil, err
		}
		return textResult(text), nil, nil
	})
	return server
}

// Run serves the MCP protocol over the given streams until the context is
// cancelled or the peer closes the connection. The streams are normally
// os.Stdin and os.Stdout: the server is read-only and initiates no outbound
// connections. Diagnostics go to the process's stderr, never the protocol
// stream.
func Run(ctx context.Context, stdin io.Reader, stdout io.Writer, deps Deps) error {
	server := NewServer(deps)
	// IOTransport is StdioTransport over injected streams; it keeps the same
	// newline-delimited JSON wire format while staying testable.
	return server.Run(ctx, &mcp.IOTransport{
		Reader: io.NopCloser(stdin),
		Writer: nopWriteCloser{stdout},
	})
}

// listText validates the tool arguments against the one-version-source rules
// the CLI enforces at flag level, then delegates to the injected list
// implementation. Returned errors surface as MCP tool errors (IsError=true),
// letting the calling agent self-correct.
func listText(deps Deps, args listArgs) (string, error) {
	version := strings.TrimSpace(args.Version)
	filePath := strings.TrimSpace(args.FilePath)
	language := strings.TrimSpace(args.Language)
	if version != "" && filePath != "" {
		return "", errors.New("pass only one version source: file_path or version")
	}
	if version != "" && language != "" {
		return "", errors.New("the axis names of version already select the language; drop language, or pass file_path/language without version")
	}
	if version == "" && filePath == "" {
		return "", errors.New("pass one of file_path (recommended, optionally with language) or version; language alone resolves no versions")
	}
	text, err := deps.List(version, filePath, language)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(text, "\n"), nil
}

// explainText requires both tool parameters — the tool contract narrows the
// CLI, where --lang may be omitted when every reference is a lang:id — then
// normalizes ids exactly like the CLI before delegating.
func explainText(deps Deps, args explainArgs) (string, error) {
	language := strings.TrimSpace(args.Language)
	if language == "" {
		return "", errors.New("language is required, for example go, python, or typescript")
	}
	ids := guidelines.NormalizeIDs(args.IDs)
	if len(ids) == 0 {
		return "", errors.New("pass at least one guideline id from list_guidelines")
	}
	text, err := deps.Explain(language, ids)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(text, "\n"), nil
}

// textResult wraps rendered text as the single text content of a tool result.
func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}
}

// readOnlyAnnotations marks both tools as read-only and closed-world per the
// PRD: the server only reads local files to resolve versions and never
// mutates anything.
func readOnlyAnnotations() *mcp.ToolAnnotations {
	no := false
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		IdempotentHint:  true,
		DestructiveHint: &no,
		OpenWorldHint:   &no,
	}
}

// serverVersion mirrors the CLI's build-info version detection for the
// initialize handshake. These few lines are metadata, not rendering: keeping
// them local avoids the import cycle with internal/cli described in the
// package comment.
func serverVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "dev"
	}
	return info.Main.Version
}

// nopWriteCloser adapts an io.Writer to the io.WriteCloser IOTransport
// requires without ever closing the protocol peer's stream.
type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
