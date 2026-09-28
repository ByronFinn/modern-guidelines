// Integration tests for the mcpserver package: a real MCP client (the SDK's
// in-memory transport, or the command transport for the full stdio path)
// talks to the server and every tool result is asserted byte-identical to the
// corresponding cli.Run output — the isomorphism requirement made executable.
package mcpserver_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ByronFinn/modern-guidelines/internal/cli"
	"github.com/ByronFinn/modern-guidelines/internal/mcpserver"
	"github.com/ByronFinn/modern-guidelines/internal/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	// Blank-imported so the version detectors register in the test process:
	// file_path list calls then resolve the repository's real manifests,
	// exactly as they do inside a binary that wires the detectors in.
	_ "github.com/ByronFinn/modern-guidelines/internal/detectors"
)

// cliOutput runs the real CLI in-process and returns its output with the
// single trailing newline removed — the expected value for every isomorphism
// assertion.
func cliOutput(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := cli.Run(args, &buf); err != nil {
		t.Fatalf("cli.Run(%q) failed: %v", args, err)
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

// cliError runs the real CLI and returns its user-facing error text.
func cliError(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	err := cli.Run(args, &buf)
	if err == nil {
		t.Fatalf("cli.Run(%q) unexpectedly succeeded:\n%s", args, buf.String())
	}
	return err.Error()
}

// testDeps mirrors the wiring internal/cli installs in its mcp subcommand:
// both route through cli.Run with synthesized flags, so the tool output and
// the CLI output come from the same code path.
func testDeps() mcpserver.Deps {
	return mcpserver.Deps{
		List: func(versionSpec, filePath, language string) (string, error) {
			args := []string{"list"}
			if versionSpec != "" {
				args = append(args, "--version", versionSpec)
			}
			if filePath != "" {
				args = append(args, "--file-path", filePath)
			}
			if language != "" {
				args = append(args, "--lang", language)
			}
			var buf bytes.Buffer
			if err := cli.Run(args, &buf); err != nil {
				return "", err
			}
			return buf.String(), nil
		},
		Explain: func(language string, ids []string) (string, error) {
			args := append([]string{"explain", "--lang", language}, ids...)
			var buf bytes.Buffer
			if err := cli.Run(args, &buf); err != nil {
				return "", err
			}
			return buf.String(), nil
		},
	}
}

// newTestSession connects a real SDK client to the server over the in-memory
// transport. Servers must connect before clients: the client drives the
// initialize handshake during Connect.
func newTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	server := mcpserver.NewServer(testDeps())
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), serverTransport, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "mcpserver-test", Version: "v0"}, nil)
	session, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect (initialize handshake): %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

// callTool calls a tool and fails on protocol-level errors; tool-level errors
// are part of the contract and asserted by the caller through IsError.
func callTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatalf("CallTool %s: protocol error: %v", name, err)
	}
	return result
}

// resultText returns the single text content of a successful tool result.
func resultText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) != 1 {
		t.Fatalf("expected exactly one content block, got %d: %#v", len(result.Content), result.Content)
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %#v", result.Content[0])
	}
	return text.Text
}

// errorText asserts the result is a tool error and returns its text.
func errorText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if !result.IsError {
		t.Fatalf("expected IsError=true, got success:\n%s", resultText(t, result))
	}
	return resultText(t, result)
}

// repoRoot locates the repository root (two levels above this package).
func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("repository root not found at %s: %v", root, err)
	}
	return root
}

// firstRuleID returns the first (newest) rule id of a language dataset.
func firstRuleID(t *testing.T, language string) string {
	t.Helper()
	dataset, ok := registry.DatasetForLanguage(language)
	if !ok {
		t.Fatalf("dataset for %q is not registered", language)
	}
	rules := dataset.Rules()
	if len(rules) == 0 {
		t.Fatalf("dataset for %q has no rules", language)
	}
	return rules[0].ID
}

func TestInitializeAndListTools(t *testing.T) {
	session := newTestSession(t)

	result, err := session.ListTools(context.Background(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(result.Tools) != 2 {
		t.Fatalf("expected exactly 2 tools, got %d: %#v", len(result.Tools), result.Tools)
	}
	tools := make(map[string]*mcp.Tool, len(result.Tools))
	for _, tool := range result.Tools {
		tools[tool.Name] = tool
	}

	list, ok := tools[mcpserver.ListToolName]
	if !ok {
		t.Fatalf("tool %q is missing; got %v", mcpserver.ListToolName, tools)
	}
	for _, want := range []string{
		"Call before writing or editing Python, Go, or TypeScript code",
		"file_path",
		"axis=version",
		"provenance",
	} {
		if !strings.Contains(list.Description, want) {
			t.Errorf("list_guidelines description misses %q:\n%s", want, list.Description)
		}
	}
	assertReadOnly(t, list)

	explain, ok := tools[mcpserver.ExplainToolName]
	if !ok {
		t.Fatalf("tool %q is missing; got %v", mcpserver.ExplainToolName, tools)
	}
	if !strings.Contains(explain.Description, "list_guidelines") {
		t.Errorf("explain_guideline description should point at list_guidelines:\n%s", explain.Description)
	}
	assertReadOnly(t, explain)

	// The explain schema must mark both parameters required so models always
	// pass them.
	schema, ok := explain.InputSchema.(map[string]any)
	if !ok {
		t.Fatalf("explain_guideline input schema is not a JSON object: %#v", explain.InputSchema)
	}
	required, _ := schema["required"].([]any)
	for _, want := range []string{"language", "ids"} {
		found := false
		for _, name := range required {
			if name == want {
				found = true
			}
		}
		if !found {
			t.Errorf("explain_guideline schema must require %q; required = %v", want, required)
		}
	}
}

func assertReadOnly(t *testing.T, tool *mcp.Tool) {
	t.Helper()
	if tool.Annotations == nil {
		t.Fatalf("tool %q has no annotations", tool.Name)
	}
	if !tool.Annotations.ReadOnlyHint {
		t.Errorf("tool %q must advertise ReadOnlyHint", tool.Name)
	}
	for _, hint := range []struct {
		name  string
		value *bool
	}{
		{"DestructiveHint", tool.Annotations.DestructiveHint},
		{"OpenWorldHint", tool.Annotations.OpenWorldHint},
	} {
		if hint.value == nil || *hint.value {
			t.Errorf("tool %q must advertise %s=false", tool.Name, hint.name)
		}
	}
}

func TestListGuidelinesByFilePath(t *testing.T) {
	session := newTestSession(t)
	mainGo := filepath.Join(repoRoot(t), "main.go")

	result := callTool(t, session, mcpserver.ListToolName, map[string]any{"file_path": mainGo})
	if result.IsError {
		t.Fatalf("list_guidelines failed: %s", resultText(t, result))
	}
	got := resultText(t, result)
	if expected := cliOutput(t, "list", "--file-path", mainGo); got != expected {
		t.Errorf("list_guidelines output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
	}
	// The provenance lines (PRD requirement 5) must ride along as text.
	for _, want := range []string{"# go 1.", "source: go.mod go directive", "fidelity: high)"} {
		if !strings.Contains(got, want) {
			t.Errorf("list_guidelines output misses provenance %q:\n%s", want, got)
		}
	}
}

func TestListGuidelinesByLanguage(t *testing.T) {
	session := newTestSession(t)
	t.Run("language_overrides_unknown_extension", func(t *testing.T) {
		// go.sum has no registered extension, so the language must come from
		// the language parameter, mirroring `list --file-path go.sum --lang go`.
		goSum := filepath.Join(repoRoot(t), "go.sum")
		result := callTool(t, session, mcpserver.ListToolName, map[string]any{"file_path": goSum, "language": "go"})
		if result.IsError {
			t.Fatalf("list_guidelines failed: %s", resultText(t, result))
		}
		got := resultText(t, result)
		if expected := cliOutput(t, "list", "--file-path", goSum, "--lang", "go"); got != expected {
			t.Errorf("list_guidelines output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
		if !strings.Contains(got, "# go 1.") {
			t.Errorf("list_guidelines output misses the go provenance line:\n%s", got)
		}
	})
	t.Run("language_alone_resolves_no_versions", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ListToolName, map[string]any{"language": "go"}))
		if !strings.Contains(got, "file_path") {
			t.Errorf("error should point at file_path/version: %q", got)
		}
	})
}

func TestListGuidelinesByVersion(t *testing.T) {
	session := newTestSession(t)
	t.Run("single_axis", func(t *testing.T) {
		result := callTool(t, session, mcpserver.ListToolName, map[string]any{"version": "go=1.24"})
		if result.IsError {
			t.Fatalf("list_guidelines failed: %s", resultText(t, result))
		}
		got := resultText(t, result)
		if expected := cliOutput(t, "list", "--version", "go=1.24"); got != expected {
			t.Errorf("list_guidelines output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
		if !strings.Contains(got, "# go 1.24 (source: --version flag, fidelity: high)") {
			t.Errorf("list_guidelines output misses the --version provenance line:\n%s", got)
		}
	})
	t.Run("multi_axis", func(t *testing.T) {
		version := "typescript=5.4,node=18"
		result := callTool(t, session, mcpserver.ListToolName, map[string]any{"version": version})
		if result.IsError {
			t.Fatalf("list_guidelines failed: %s", resultText(t, result))
		}
		got := resultText(t, result)
		if expected := cliOutput(t, "list", "--version", version); got != expected {
			t.Errorf("list_guidelines output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
		for _, want := range []string{"# typescript 5.4 (source: --version flag, fidelity: high)", "# node 18 (source: --version flag, fidelity: high)", "# axis: typescript", "# axis: node"} {
			if !strings.Contains(got, want) {
				t.Errorf("list_guidelines output misses %q:\n%s", want, got)
			}
		}
	})
}

func TestListGuidelinesArgumentErrors(t *testing.T) {
	session := newTestSession(t)
	t.Run("no_version_source", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ListToolName, map[string]any{}))
		if !strings.Contains(got, "file_path") {
			t.Errorf("error should point at file_path/language/version: %q", got)
		}
	})
	t.Run("version_conflicts_with_file_path", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ListToolName, map[string]any{"version": "go=1.24", "file_path": "main.go"}))
		if !strings.Contains(got, "only one version source") {
			t.Errorf("unexpected conflict error: %q", got)
		}
	})
	t.Run("version_conflicts_with_language", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ListToolName, map[string]any{"version": "go=1.24", "language": "go"}))
		if !strings.Contains(got, "language") {
			t.Errorf("unexpected conflict error: %q", got)
		}
	})
	t.Run("unknown_axis_matches_cli_error", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ListToolName, map[string]any{"version": "foo=1.2"}))
		if expected := cliError(t, "list", "--version", "foo=1.2"); got != expected {
			t.Errorf("error differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
	})
}

func TestExplainGuideline(t *testing.T) {
	session := newTestSession(t)
	t.Run("bare_ids", func(t *testing.T) {
		id := firstRuleID(t, "go")
		result := callTool(t, session, mcpserver.ExplainToolName, map[string]any{"language": "go", "ids": []string{id}})
		if result.IsError {
			t.Fatalf("explain_guideline failed: %s", resultText(t, result))
		}
		got := resultText(t, result)
		if expected := cliOutput(t, "explain", "--lang", "go", id); got != expected {
			t.Errorf("explain_guideline output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
		for _, want := range []string{"Since: go", "Summary:", "Details:", "Autofix:"} {
			if !strings.Contains(got, want) {
				t.Errorf("explain_guideline output misses %q:\n%s", want, got)
			}
		}
	})
	t.Run("lang_id_overrides_language", func(t *testing.T) {
		reference := "python:" + firstRuleID(t, "python")
		result := callTool(t, session, mcpserver.ExplainToolName, map[string]any{"language": "go", "ids": []string{reference}})
		if result.IsError {
			t.Fatalf("explain_guideline failed: %s", resultText(t, result))
		}
		got := resultText(t, result)
		if expected := cliOutput(t, "explain", "--lang", "go", reference); got != expected {
			t.Errorf("explain_guideline output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
		if !strings.Contains(got, "Since: python") {
			t.Errorf("explain_guideline output misses the python axis line:\n%s", got)
		}
	})
	t.Run("unknown_ids_matches_cli_error", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ExplainToolName, map[string]any{"language": "go", "ids": []string{"no-such-id"}}))
		if expected := cliError(t, "explain", "--lang", "go", "no-such-id"); got != expected {
			t.Errorf("error differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
		}
	})
	t.Run("missing_language_rejected_by_schema", func(t *testing.T) {
		got := errorText(t, callTool(t, session, mcpserver.ExplainToolName, map[string]any{"ids": []string{"any"}}))
		if !strings.Contains(got, "language") {
			t.Errorf("schema validation error should mention language: %q", got)
		}
	})
}

// TestServerOverCommandTransport exercises the full shipped path end to end:
// the mcp subcommand of the real binary, served over os.Stdin/os.Stdout
// through mcpserver.Run, driven by the SDK client over the command transport.
// The version source is used because it is deterministic and, unlike
// file_path detection, independent of which packages the binary links in.
func TestServerOverCommandTransport(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("go toolchain not available: %v", err)
	}
	root := repoRoot(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.Command("go", "run", ".", "mcp")
	cmd.Dir = root
	client := mcp.NewClient(&mcp.Implementation{Name: "mcpserver-e2e", Version: "v0"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect to `go run . mcp`: %v", err)
	}
	defer session.Close()

	result := callTool(t, session, mcpserver.ListToolName, map[string]any{"version": "go=1.24"})
	if result.IsError {
		t.Fatalf("list_guidelines over stdio failed: %s", resultText(t, result))
	}
	got := resultText(t, result)
	if expected := cliOutput(t, "list", "--version", "go=1.24"); got != expected {
		t.Errorf("stdio list_guidelines output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
	}

	id := firstRuleID(t, "go")
	result = callTool(t, session, mcpserver.ExplainToolName, map[string]any{"language": "go", "ids": []string{id}})
	if result.IsError {
		t.Fatalf("explain_guideline over stdio failed: %s", resultText(t, result))
	}
	got = resultText(t, result)
	if expected := cliOutput(t, "explain", "--lang", "go", id); got != expected {
		t.Errorf("stdio explain_guideline output differs from the CLI:\n--- tool ---\n%s\n--- cli ---\n%s", got, expected)
	}
}
