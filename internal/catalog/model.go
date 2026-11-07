package catalog

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// KeyBinding represents a single keyboard shortcut and its explanation.
type KeyBinding struct {
	Key  string `yaml:"key" json:"key"`
	Desc string `yaml:"desc" json:"desc"`
}

// UnmarshalYAML supports both structured objects {key, desc} and concise strings "key: desc".
func (kb *KeyBinding) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		raw := strings.TrimSpace(node.Value)
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) == 2 {
			kb.Key = strings.TrimSpace(parts[0])
			kb.Desc = strings.TrimSpace(parts[1])
			return nil
		}
		kb.Key = raw
		kb.Desc = ""
		return nil
	}

	type rawKeyBinding KeyBinding
	var raw rawKeyBinding
	if err := node.Decode(&raw); err != nil {
		return err
	}
	kb.Key = strings.TrimSpace(raw.Key)
	kb.Desc = strings.TrimSpace(raw.Desc)
	return nil
}

// Lesson encapsulates a single bite-sized Vim learning module.
type Lesson struct {
	ID           string       `yaml:"id" json:"id"`
	Track        string       `yaml:"track" json:"track"`
	Order        int          `yaml:"order" json:"order"`
	EstSeconds   int          `yaml:"est_seconds" json:"est_seconds"`
	Title        string       `yaml:"title" json:"title"`
	Summary      string       `yaml:"summary" json:"summary"`
	Keys         []KeyBinding `yaml:"keys" json:"keys"`
	Body         string       `yaml:"body" json:"body"`
	Drill        string       `yaml:"drill" json:"drill"`
	Instructions string       `yaml:"instructions" json:"instructions"`
	Solution     string       `yaml:"solution" json:"solution"`
}

// Validate checks that required lesson fields are present and valid.
func (l *Lesson) Validate() error {
	if strings.TrimSpace(l.ID) == "" {
		return errors.New("lesson id cannot be empty")
	}
	if strings.TrimSpace(l.Track) == "" {
		return fmt.Errorf("lesson '%s': track cannot be empty", l.ID)
	}
	if l.Order <= 0 {
		return fmt.Errorf("lesson '%s': order must be positive", l.ID)
	}
	if strings.TrimSpace(l.Title) == "" {
		return fmt.Errorf("lesson '%s': title cannot be empty", l.ID)
	}
	if strings.TrimSpace(l.Body) == "" {
		return fmt.Errorf("lesson '%s': body cannot be empty", l.ID)
	}
	if strings.TrimSpace(l.Instructions) == "" {
		return fmt.Errorf("lesson '%s': instructions cannot be empty", l.ID)
	}
	return nil
}
