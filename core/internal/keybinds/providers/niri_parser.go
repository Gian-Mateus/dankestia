package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sblinch/kdl-go"
	"github.com/sblinch/kdl-go/document"
)

type NiriKeyBinding struct {
	Mods            []string
	Key             string
	Action          string
	Args            []string
	Description     string
	HideOnOverlay   bool
	CooldownMs      int
	AllowWhenLocked bool
	AllowInhibiting *bool
	Repeat          *bool
	Source          string
}

type NiriSection struct {
	Name     string
	Keybinds []NiriKeyBinding
	Children []NiriSection
}

type NiriParser struct {
	configDir          string
	processedFiles     map[string]bool
	bindMap            map[string]*NiriKeyBinding
	bindOrder          []string
	currentSource      string
	dankestiaBindsIncluded   bool
	dankestiaBindsExists     bool
	includeCount       int
	dankestiaIncludePos      int
	bindsBeforeDANKESTIA     int
	bindsAfterDANKESTIA      int
	dankestiaBindKeys        map[string]bool
	configBindKeys     map[string]bool
	dankestiaProcessed       bool
	dankestiaBindMap         map[string]*NiriKeyBinding
	conflictingConfigs map[string]*NiriKeyBinding
}

func parseKDL(data []byte) (*document.Document, error) {
	return kdl.Parse(strings.NewReader(normalizeKDLBraces(string(data))))
}

func normalizeKDLBraces(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	var prev byte
	n := len(input)
	for i := 0; i < n; {
		c := input[i]

		switch {
		case c == '"':
			end := findStringEnd(input, i)
			sb.WriteString(input[i:end])
			prev = '"'
			i = end
		case c == '/' && i+1 < n && input[i+1] == '/':
			end := findLineCommentEnd(input, i)
			sb.WriteString(input[i:end])
			prev = '\n'
			i = end
		case c == '/' && i+1 < n && input[i+1] == '*':
			end := findBlockCommentEnd(input, i)
			sb.WriteString(input[i:end])
			prev = '/'
			i = end
		case c == '{' && prev != 0 && !isBraceAdjacentSpace(prev):
			sb.WriteByte(' ')
			sb.WriteByte(c)
			prev = c
			i++
		default:
			sb.WriteByte(c)
			prev = c
			i++
		}
	}

	return sb.String()
}

func findStringEnd(s string, start int) int {
	n := len(s)
	for i := start + 1; i < n; {
		switch s[i] {
		case '\\':
			i += 2
		case '"':
			return i + 1
		default:
			i++
		}
	}
	return n
}

func findLineCommentEnd(s string, start int) int {
	for i := start + 2; i < len(s); i++ {
		if s[i] == '\n' {
			return i
		}
	}
	return len(s)
}

func findBlockCommentEnd(s string, start int) int {
	n := len(s)
	depth := 1
	for i := start + 2; i < n && depth > 0; {
		switch {
		case i+1 < n && s[i] == '/' && s[i+1] == '*':
			depth++
			i += 2
		case i+1 < n && s[i] == '*' && s[i+1] == '/':
			depth--
			i += 2
			if depth == 0 {
				return i
			}
		default:
			i++
		}
	}
	return n
}

func isBraceAdjacentSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '{':
		return true
	}
	return false
}

func NewNiriParser(configDir string) *NiriParser {
	return &NiriParser{
		configDir:          configDir,
		processedFiles:     make(map[string]bool),
		bindMap:            make(map[string]*NiriKeyBinding),
		bindOrder:          []string{},
		currentSource:      "",
		dankestiaIncludePos:      -1,
		dankestiaBindKeys:        make(map[string]bool),
		configBindKeys:     make(map[string]bool),
		dankestiaBindMap:         make(map[string]*NiriKeyBinding),
		conflictingConfigs: make(map[string]*NiriKeyBinding),
	}
}

func normalizeNiriBindKey(key string) string {
	parts := strings.Split(key, "+")
	for i := range parts {
		parts[i] = strings.ToLower(strings.TrimSpace(parts[i]))
	}
	return strings.Join(parts, "+")
}

func (p *NiriParser) Parse() (*NiriSection, error) {
	dankestiaBindsPath := filepath.Join(p.configDir, "dankestia", "binds.kdl")
	if _, err := os.Stat(dankestiaBindsPath); err == nil {
		p.dankestiaBindsExists = true
	}

	configPath := filepath.Join(p.configDir, "config.kdl")
	section, err := p.parseFile(configPath, "")
	if err != nil {
		return nil, err
	}

	if p.dankestiaBindsExists && !p.dankestiaProcessed {
		p.parseDANKESTIABindsDirectly(dankestiaBindsPath, section)
	}

	section.Keybinds = p.finalizeBinds()
	return section, nil
}

func (p *NiriParser) parseDANKESTIABindsDirectly(dankestiaBindsPath string, section *NiriSection) {
	data, err := os.ReadFile(dankestiaBindsPath)
	if err != nil {
		return
	}

	doc, err := parseKDL(data)
	if err != nil {
		return
	}

	prevSource := p.currentSource
	p.currentSource = dankestiaBindsPath
	baseDir := filepath.Dir(dankestiaBindsPath)
	p.processNodes(doc.Nodes, section, baseDir)
	p.currentSource = prevSource
	p.dankestiaProcessed = true
}

func (p *NiriParser) finalizeBinds() []NiriKeyBinding {
	binds := make([]NiriKeyBinding, 0, len(p.bindOrder))
	for _, key := range p.bindOrder {
		if kb, ok := p.bindMap[key]; ok {
			binds = append(binds, *kb)
		}
	}
	return binds
}

func (p *NiriParser) addBind(kb *NiriKeyBinding) {
	key := p.formatBindKey(kb)
	normalizedKey := normalizeNiriBindKey(key)
	isDANKESTIABind := strings.Contains(kb.Source, "dankestia/binds.kdl")

	if isDANKESTIABind {
		p.dankestiaBindKeys[normalizedKey] = true
		p.dankestiaBindMap[normalizedKey] = kb
	} else if p.dankestiaBindKeys[normalizedKey] {
		p.bindsAfterDANKESTIA++
		p.conflictingConfigs[normalizedKey] = kb
		p.configBindKeys[normalizedKey] = true
		return
	} else {
		p.configBindKeys[normalizedKey] = true
	}

	if _, exists := p.bindMap[normalizedKey]; !exists {
		p.bindOrder = append(p.bindOrder, normalizedKey)
	}
	p.bindMap[normalizedKey] = kb
}

func (p *NiriParser) formatBindKey(kb *NiriKeyBinding) string {
	parts := make([]string, 0, len(kb.Mods)+1)
	parts = append(parts, kb.Mods...)
	parts = append(parts, kb.Key)
	return strings.Join(parts, "+")
}

func (p *NiriParser) parseFile(filePath, sectionName string) (*NiriSection, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path %s: %w", filePath, err)
	}

	if p.processedFiles[absPath] {
		return &NiriSection{Name: sectionName}, nil
	}
	p.processedFiles[absPath] = true

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", absPath, err)
	}

	doc, err := parseKDL(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse KDL in %s: %w", absPath, err)
	}

	section := &NiriSection{
		Name: sectionName,
	}

	prevSource := p.currentSource
	p.currentSource = absPath
	baseDir := filepath.Dir(absPath)
	p.processNodes(doc.Nodes, section, baseDir)
	p.currentSource = prevSource

	return section, nil
}

func (p *NiriParser) processNodes(nodes []*document.Node, section *NiriSection, baseDir string) {
	for _, node := range nodes {
		name := node.Name.String()

		switch name {
		case "include":
			p.handleInclude(node, section, baseDir)
		case "binds":
			p.extractBinds(node, section, "")
		case "recent-windows":
			p.handleRecentWindows(node, section)
		}
	}
}

func (p *NiriParser) handleInclude(node *document.Node, section *NiriSection, baseDir string) {
	if len(node.Arguments) == 0 {
		return
	}

	includePath := strings.Trim(node.Arguments[0].String(), "\"")
	isDANKESTIAInclude := includePath == "dankestia/binds.kdl" || strings.HasSuffix(includePath, "/dankestia/binds.kdl")

	p.includeCount++
	if isDANKESTIAInclude {
		p.dankestiaBindsIncluded = true
		p.dankestiaIncludePos = p.includeCount
		p.bindsBeforeDANKESTIA = len(p.bindMap)
	}

	fullPath := filepath.Join(baseDir, includePath)
	if filepath.IsAbs(includePath) {
		fullPath = includePath
	}

	if isDANKESTIAInclude {
		p.dankestiaProcessed = true
	}

	includedSection, err := p.parseFile(fullPath, "")
	if err != nil {
		return
	}

	section.Children = append(section.Children, includedSection.Children...)
}

func (p *NiriParser) HasDANKESTIABindsIncluded() bool {
	return p.dankestiaBindsIncluded
}

func (p *NiriParser) handleRecentWindows(node *document.Node, section *NiriSection) {
	if node.Children == nil {
		return
	}

	for _, child := range node.Children {
		if child.Name.String() != "binds" {
			continue
		}
		p.extractBinds(child, section, "Alt-Tab")
	}
}

func (p *NiriParser) extractBinds(node *document.Node, section *NiriSection, subcategory string) {
	if node.Children == nil {
		return
	}

	for _, child := range node.Children {
		kb := p.parseKeybindNode(child, subcategory)
		if kb == nil {
			continue
		}
		p.addBind(kb)
	}
}

func (p *NiriParser) parseKeybindNode(node *document.Node, _ string) *NiriKeyBinding {
	keyCombo := node.Name.String()
	if keyCombo == "" {
		return nil
	}

	mods, key := p.parseKeyCombo(keyCombo)

	var action string
	var args []string
	if len(node.Children) > 0 {
		actionNode := node.Children[0]
		action = actionNode.Name.String()
		for _, arg := range actionNode.Arguments {
			args = append(args, arg.ValueString())
		}
		if actionNode.Properties != nil {
			for _, propName := range []string{"focus", "show-pointer", "write-to-disk", "skip-confirmation", "delay-ms"} {
				if val, ok := actionNode.Properties.Get(propName); ok {
					args = append(args, propName+"="+val.String())
				}
			}
		}
	}

	var description string
	var hideOnOverlay bool
	var cooldownMs int
	var allowWhenLocked bool
	var allowInhibiting *bool
	var repeat *bool
	if node.Properties != nil {
		if val, ok := node.Properties.Get("hotkey-overlay-title"); ok {
			switch val.ValueString() {
			case "null", "":
				hideOnOverlay = true
			default:
				description = val.ValueString()
			}
		}
		if val, ok := node.Properties.Get("cooldown-ms"); ok {
			cooldownMs, _ = strconv.Atoi(val.String())
		}
		if val, ok := node.Properties.Get("allow-when-locked"); ok {
			allowWhenLocked = val.String() == "true"
		}
		if val, ok := node.Properties.Get("allow-inhibiting"); ok {
			v := val.String() == "true"
			allowInhibiting = &v
		}
		if val, ok := node.Properties.Get("repeat"); ok {
			v := val.String() == "true"
			repeat = &v
		}
	}

	return &NiriKeyBinding{
		Mods:            mods,
		Key:             key,
		Action:          action,
		Args:            args,
		Description:     description,
		HideOnOverlay:   hideOnOverlay,
		CooldownMs:      cooldownMs,
		AllowWhenLocked: allowWhenLocked,
		AllowInhibiting: allowInhibiting,
		Repeat:          repeat,
		Source:          p.currentSource,
	}
}

func (p *NiriParser) parseKeyCombo(combo string) ([]string, string) {
	parts := strings.Split(combo, "+")

	switch len(parts) {
	case 0:
		return nil, combo
	case 1:
		return nil, parts[0]
	default:
		return parts[:len(parts)-1], parts[len(parts)-1]
	}
}

type NiriParseResult struct {
	Section            *NiriSection
	DANKESTIABindsIncluded   bool
	DANKESTIAStatus          *DANKESTIABindsStatusInfo
	ConflictingConfigs map[string]*NiriKeyBinding
}

type DANKESTIABindsStatusInfo struct {
	Exists          bool
	Included        bool
	IncludePosition int
	TotalIncludes   int
	BindsAfterDANKESTIA   int
	Effective       bool
	OverriddenBy    int
	StatusMessage   string
}

func (p *NiriParser) buildDANKESTIAStatus() *DANKESTIABindsStatusInfo {
	status := &DANKESTIABindsStatusInfo{
		Exists:          p.dankestiaBindsExists,
		Included:        p.dankestiaBindsIncluded,
		IncludePosition: p.dankestiaIncludePos,
		TotalIncludes:   p.includeCount,
		BindsAfterDANKESTIA:   p.bindsAfterDANKESTIA,
	}

	switch {
	case !p.dankestiaBindsExists:
		status.Effective = false
		status.StatusMessage = "dankestia/binds.kdl does not exist"
	case !p.dankestiaBindsIncluded:
		status.Effective = false
		status.StatusMessage = "dankestia/binds.kdl is not included in config.kdl"
	case p.bindsAfterDANKESTIA > 0:
		status.Effective = true
		status.OverriddenBy = p.bindsAfterDANKESTIA
		status.StatusMessage = "Some DANKESTIA binds may be overridden by config binds"
	default:
		status.Effective = true
		status.StatusMessage = "DANKESTIA binds are active"
	}

	return status
}

func ParseNiriKeys(configDir string) (*NiriParseResult, error) {
	parser := NewNiriParser(configDir)
	section, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return &NiriParseResult{
		Section:            section,
		DANKESTIABindsIncluded:   parser.HasDANKESTIABindsIncluded(),
		DANKESTIAStatus:          parser.buildDANKESTIAStatus(),
		ConflictingConfigs: parser.conflictingConfigs,
	}, nil
}
