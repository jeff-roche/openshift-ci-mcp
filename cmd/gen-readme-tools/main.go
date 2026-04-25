package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/openshift-eng/openshift-ci-mcp/pkg/tools/domain"
	"github.com/openshift-eng/openshift-ci-mcp/pkg/tools/proxy"
)

func main() {
	root := findProjectRoot()
	readmePath := filepath.Join(root, "README.md")

	domainTable := renderTable(buildDomainServer())
	proxyTable := renderTable(buildProxyServer())

	readme, err := os.ReadFile(readmePath)
	if err != nil {
		fatal("reading README.md: %v", err)
	}

	readme = replaceSection(readme, "DOMAIN TOOLS", domainTable)
	readme = replaceSection(readme, "PROXY TOOLS", proxyTable)

	if err := os.WriteFile(readmePath, readme, 0o644); err != nil {
		fatal("writing README.md: %v", err)
	}
}

func buildDomainServer() *mcpserver.MCPServer {
	s := mcpserver.NewMCPServer("gen", "0.0.0")
	var sippy nopSippy
	var rc nopRC
	var search nopSearch

	domain.RegisterReleaseTools(s, sippy)
	domain.RegisterVariantTools(s, sippy)
	domain.RegisterJobTools(s, sippy)
	domain.RegisterTestTools(s, sippy)
	domain.RegisterComponentTools(s, sippy)
	domain.RegisterPayloadTools(s, sippy, rc)
	domain.RegisterSearchTools(s, search)
	domain.RegisterPullRequestTools(s, sippy)

	return s
}

func buildProxyServer() *mcpserver.MCPServer {
	s := mcpserver.NewMCPServer("gen", "0.0.0")
	var sippy nopSippy
	var rc nopRC
	var search nopSearch

	proxy.RegisterSippyProxy(s, sippy)
	proxy.RegisterReleaseControllerProxy(s, rc)
	proxy.RegisterSearchCIProxy(s, search)

	return s
}

func renderTable(srv *mcpserver.MCPServer) string {
	tools := srv.ListTools()
	names := make([]string, 0, len(tools))
	for name := range tools {
		names = append(names, name)
	}
	sort.Strings(names)

	var buf strings.Builder
	buf.WriteString("| Tool | Description |\n")
	buf.WriteString("| ---- | ----------- |\n")
	for _, name := range names {
		desc := strings.ReplaceAll(tools[name].Tool.Description, "|", "\\|")
		fmt.Fprintf(&buf, "| `%s` | %s |\n", name, desc)
	}
	return buf.String()
}

func replaceSection(content []byte, marker, replacement string) []byte {
	begin := []byte(fmt.Sprintf("<!-- BEGIN %s -->", marker))
	end := []byte(fmt.Sprintf("<!-- END %s -->", marker))

	startIdx := bytes.Index(content, begin)
	if startIdx == -1 {
		fatal("marker %q not found in README.md", string(begin))
	}
	endIdx := bytes.Index(content[startIdx:], end)
	if endIdx == -1 {
		fatal("marker %q not found in README.md", string(end))
	}
	endIdx += startIdx

	var buf bytes.Buffer
	buf.Write(content[:startIdx])
	buf.Write(begin)
	buf.WriteByte('\n')
	buf.WriteString(replacement)
	buf.Write(end)
	buf.Write(content[endIdx+len(end):])

	return buf.Bytes()
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fatal("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fatal("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "gen-readme-tools: "+format+"\n", args...)
	os.Exit(1)
}
