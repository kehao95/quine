package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	skillsDirRel             = ".agents/skills"
	skillFileName            = "SKILL.md"
	maxSkillDescriptionBytes = 1024
	maxSkillFrontmatterBytes = 4096
)

var skillNamePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$`)

// Skill is the prompt-facing index metadata for an external project skill.
type Skill struct {
	Name        string
	Description string
	Source      string
}

type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// LoadSkills discovers external project skills and returns prompt-facing
// metadata only. Skill bodies remain ordinary project files.
func LoadSkills(cfg *Config) ([]Skill, error) {
	root, ok, err := discoverSkillsRoot(cfg)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	promptBase := skillsPromptBaseDir(cfg, root)
	skillsDir := filepath.Join(root, skillsDirRel)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read skills directory %q: %w", skillsDir, err)
	}

	var skills []Skill
	seen := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirName := entry.Name()
		if strings.HasPrefix(dirName, ".") {
			continue
		}
		sourceAbs := filepath.Join(skillsDir, dirName, skillFileName)
		sourceRel, err := filepath.Rel(promptBase, sourceAbs)
		if err != nil {
			return nil, fmt.Errorf("relativize skill path %q: %w", sourceAbs, err)
		}
		sourceRel = filepath.ToSlash(sourceRel)

		skill, err := loadSkillIndex(sourceAbs, sourceRel)
		if err != nil {
			return nil, err
		}
		if previous, exists := seen[skill.Name]; exists {
			return nil, fmt.Errorf("duplicate skill name %q in %s and %s", skill.Name, previous, sourceRel)
		}
		seen[skill.Name] = sourceRel
		skills = append(skills, skill)
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})
	return skills, nil
}

func discoverSkillsRoot(cfg *Config) (string, bool, error) {
	if cfg != nil && cfg.WorkspaceEnabled {
		if strings.TrimSpace(cfg.WorkspaceRoot) == "" {
			return "", false, fmt.Errorf("QUINE_WORKSPACE_ROOT is required when QUINE_AGENTS_SKILLS_ENABLED=1 and workspace physics are enabled")
		}
		return cfg.WorkspaceRoot, true, nil
	}

	start := ""
	if cfg != nil {
		start = strings.TrimSpace(cfg.WorkDir)
	}
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", false, fmt.Errorf("get working directory for skills discovery: %w", err)
		}
	}

	info, err := os.Stat(start)
	if err != nil {
		return "", false, fmt.Errorf("stat skills discovery start %q: %w", start, err)
	}
	if !info.IsDir() {
		start = filepath.Dir(start)
	}

	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false, fmt.Errorf("absolute skills discovery start %q: %w", start, err)
	}
	for {
		candidate := filepath.Join(dir, skillsDirRel)
		if info, err := os.Stat(candidate); err == nil {
			if info.IsDir() {
				return dir, true, nil
			}
			return "", false, fmt.Errorf("%s exists but is not a directory", candidate)
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", false, fmt.Errorf("stat skills directory %q: %w", candidate, err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}

func loadSkillIndex(path string, sourceRel string) (Skill, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Skill{}, fmt.Errorf("skill directory %q is missing %s", filepath.Dir(path), skillFileName)
		}
		return Skill{}, fmt.Errorf("read skill %q: %w", path, err)
	}

	header, err := skillFrontmatterBytes(raw, path)
	if err != nil {
		return Skill{}, err
	}
	var fm skillFrontmatter
	if err := yaml.Unmarshal(header, &fm); err != nil {
		return Skill{}, fmt.Errorf("parse skill frontmatter %q: %w", path, err)
	}

	name := strings.TrimSpace(fm.Name)
	description := strings.TrimSpace(fm.Description)
	if name == "" {
		return Skill{}, fmt.Errorf("skill %s is missing frontmatter field name", sourceRel)
	}
	if !validSkillName(name) {
		return Skill{}, fmt.Errorf("skill %s has invalid name %q", sourceRel, name)
	}
	if description == "" {
		return Skill{}, fmt.Errorf("skill %s is missing frontmatter field description", sourceRel)
	}
	if len([]byte(description)) > maxSkillDescriptionBytes {
		return Skill{}, fmt.Errorf("skill %s description exceeds %d bytes", sourceRel, maxSkillDescriptionBytes)
	}

	return Skill{Name: name, Description: description, Source: sourceRel}, nil
}

func skillsPromptBaseDir(cfg *Config, root string) string {
	base := ""
	if cfg != nil {
		base = strings.TrimSpace(cfg.WorkDir)
	}
	if base == "" {
		base = root
	}
	info, err := os.Stat(base)
	if err == nil && !info.IsDir() {
		base = filepath.Dir(base)
	}
	return base
}

func validSkillName(name string) bool {
	if len(name) > 64 {
		return false
	}
	if strings.Contains(name, "--") {
		return false
	}
	return skillNamePattern.MatchString(name)
}

func skillFrontmatterBytes(raw []byte, path string) ([]byte, error) {
	text := string(raw)
	if !strings.HasPrefix(text, "---\n") {
		return nil, fmt.Errorf("skill %q must start with YAML frontmatter", path)
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return nil, fmt.Errorf("skill %q has unterminated YAML frontmatter", path)
	}
	header := []byte(rest[:end])
	if len(header) > maxSkillFrontmatterBytes {
		return nil, fmt.Errorf("skill %q frontmatter exceeds %d bytes", path, maxSkillFrontmatterBytes)
	}
	return header, nil
}

func envBool01Default(key string, def bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		return def, nil
	}
	switch strings.TrimSpace(v) {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("%s=%q: must be \"0\" or \"1\"", key, v)
	}
}
