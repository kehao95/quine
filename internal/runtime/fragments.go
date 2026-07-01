package runtime

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type promptFragment struct {
	Name    string
	Content string
}

func (r *Runtime) currentSystemPrompt() (string, error) {
	if err := r.syncPromptFragments(); err != nil {
		return "", err
	}
	baseData, err := os.ReadFile(r.currentIncarnationPromptFile("00-runtime.md"))
	if err != nil {
		return "", fmt.Errorf("read runtime context file: %w", err)
	}
	fragments, err := loadPromptFragments(r.currentIncarnationPromptRoot())
	if err != nil {
		return "", err
	}
	return assembleSystemPrompt(strings.TrimSpace(string(baseData)), fragments), nil
}

func (r *Runtime) syncPromptFragments() error {
	if err := r.ensureIncarnation(); err != nil {
		return err
	}
	contextRoot := r.currentIncarnationContextRoot()
	root := r.currentIncarnationPromptRoot()
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("mkdir prompt context root %s: %w", root, err)
	}
	for _, legacy := range []string{"runtime.md", "mission.md", "memory.md", "AGENTS.md", "SKILLS.md", "fragments.d"} {
		if err := os.RemoveAll(filepath.Join(contextRoot, legacy)); err != nil {
			return fmt.Errorf("remove legacy context fragment %s: %w", legacy, err)
		}
	}

	runtimePath := filepath.Join(root, "00-runtime.md")
	if err := writeTextFile(runtimePath, BuildSystemPrompt(r.cfg, r.originalInput, r.hasMaterial)); err != nil {
		return err
	}

	missionPath := filepath.Join(root, "40-mission.md")
	if strings.TrimSpace(r.originalInput) == "" {
		if err := os.RemoveAll(missionPath); err != nil {
			return fmt.Errorf("remove mission prompt fragment %s: %w", missionPath, err)
		}
	} else {
		if err := writeTextFile(missionPath, r.originalInput); err != nil {
			return err
		}
	}

	memoryPath := filepath.Join(root, "30-memory.md")
	if _, err := os.Stat(memoryPath); errors.Is(err, os.ErrNotExist) {
		if err := writeTextFile(memoryPath, ""); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("stat %s: %w", memoryPath, err)
	}

	agentsPath := filepath.Join(root, "10-agents.md")
	if !r.cfg.AgentsMDEnabled || strings.TrimSpace(r.cfg.AgentsMDPath) == "" {
		if err := os.RemoveAll(agentsPath); err != nil {
			return fmt.Errorf("remove %s: %w", agentsPath, err)
		}
	} else if err := replaceSymlink(agentsPath, r.cfg.AgentsMDPath); err != nil {
		return err
	}

	skillsPath := filepath.Join(root, "20-skills.md")
	if !r.cfg.AgentsSkillsEnabled || len(r.cfg.Skills) == 0 {
		if err := os.RemoveAll(skillsPath); err != nil {
			return fmt.Errorf("remove %s: %w", skillsPath, err)
		}
	} else if err := writeTextFile(skillsPath, renderSkillsFragment(r.cfg.Skills)); err != nil {
		return err
	}

	return nil
}

func loadPromptFragments(root string) ([]promptFragment, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	fragments := make([]promptFragment, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".md") || name == "00-runtime.md" {
			continue
		}
		path := filepath.Join(root, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read prompt fragment %s: %w", path, err)
		}
		fragments = append(fragments, promptFragment{
			Name:    name,
			Content: strings.TrimSpace(string(data)),
		})
	}
	return fragments, nil
}

func assembleSystemPrompt(base string, fragments []promptFragment) string {
	if len(fragments) == 0 {
		return base
	}

	var sb strings.Builder
	sb.WriteString(strings.TrimRight(base, "\n"))
	sb.WriteString("\n\n")
	// Separate blocks by emission, not loop index: an earlier fragment that
	// renders empty (the default empty memory fragment is the routine case)
	// must not leave a stray double blank line before the first real block.
	wrote := false
	for _, fragment := range fragments {
		block := renderPromptFragmentBlock(fragment)
		if strings.TrimSpace(block) == "" {
			continue
		}
		if wrote {
			sb.WriteString("\n\n")
		}
		wrote = true
		sb.WriteString(block)
	}
	sb.WriteString("\n")
	return sb.String()
}

func renderPromptFragmentBlock(fragment promptFragment) string {
	content := strings.TrimSpace(fragment.Content)
	if content == "" {
		return ""
	}
	switch fragment.Name {
	case "10-agents.md":
		return "### AGENTS.md\n" + content
	case "20-skills.md":
		return "### SKILLS.md\n" + content
	case "30-memory.md":
		return "### Memory\n" + content
	case "40-mission.md":
		return "### Your Mission\n" + content
	case "45-fork-assignment.md":
		return "### Fork Assignment\n" + content
	default:
		return "### " + fragment.Name + "\n" + content
	}
}
