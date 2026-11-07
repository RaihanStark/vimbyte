package catalog

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed lessons/*.yaml
var embeddedLessons embed.FS

// Catalog manages all loaded lessons across tracks.
type Catalog struct {
	lessons    map[string]*Lesson
	byTrack    map[string][]*Lesson
	trackNames []string
}

// Load reads all embedded lessons and merges any YAML lessons found in externalDir.
func Load(externalDir string) (*Catalog, error) {
	c := &Catalog{
		lessons: make(map[string]*Lesson),
		byTrack: make(map[string][]*Lesson),
	}

	// 1. Read embedded lessons
	entries, err := embeddedLessons.ReadDir("lessons")
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !isYAML(entry.Name()) {
				continue
			}
			data, err := embeddedLessons.ReadFile("lessons/" + entry.Name())
			if err != nil {
				return nil, fmt.Errorf("reading embedded lesson %s: %w", entry.Name(), err)
			}
			lesson, err := parseLesson(data, entry.Name())
			if err != nil {
				return nil, err
			}
			c.lessons[lesson.ID] = lesson
		}
	}

	// 2. Read external lessons directory if specified
	if externalDir != "" {
		stat, err := os.Stat(externalDir)
		if err == nil && stat.IsDir() {
			err = filepath.WalkDir(externalDir, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if d.IsDir() || !isYAML(d.Name()) {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("reading external lesson %s: %w", path, err)
				}
				lesson, err := parseLesson(data, filepath.Base(path))
				if err != nil {
					return err
				}
				// External lesson overrides or adds
				c.lessons[lesson.ID] = lesson
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("walking external lessons dir: %w", err)
			}
		}
	}

	// 3. Index and sort by track and order
	trackSeen := make(map[string]bool)
	for _, l := range c.lessons {
		c.byTrack[l.Track] = append(c.byTrack[l.Track], l)
		if !trackSeen[l.Track] {
			c.trackNames = append(c.trackNames, l.Track)
			trackSeen[l.Track] = true
		}
	}

	// Canonical track ordering preference: basics, motions, editing, visual, registers, then alphabetical
	canonicalOrder := map[string]int{
		"basics":    1,
		"motions":   2,
		"editing":   3,
		"visual":    4,
		"registers": 5,
	}

	sort.SliceStable(c.trackNames, func(i, j int) bool {
		r1, ok1 := canonicalOrder[c.trackNames[i]]
		r2, ok2 := canonicalOrder[c.trackNames[j]]
		if ok1 && ok2 {
			return r1 < r2
		}
		if ok1 {
			return true
		}
		if ok2 {
			return false
		}
		return c.trackNames[i] < c.trackNames[j]
	})

	for track := range c.byTrack {
		lessons := c.byTrack[track]
		sort.Slice(lessons, func(i, j int) bool {
			if lessons[i].Order != lessons[j].Order {
				return lessons[i].Order < lessons[j].Order
			}
			return lessons[i].ID < lessons[j].ID
		})
		c.byTrack[track] = lessons
	}

	return c, nil
}

func parseLesson(data []byte, filename string) (*Lesson, error) {
	var lesson Lesson
	if err := yaml.Unmarshal(data, &lesson); err != nil {
		return nil, fmt.Errorf("parsing YAML in %s: %w", filename, err)
	}
	if err := lesson.Validate(); err != nil {
		return nil, fmt.Errorf("invalid lesson in %s: %w", filename, err)
	}
	return &lesson, nil
}

func isYAML(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".yaml" || ext == ".yml"
}

// Get returns the lesson with the specified ID or an error if not found.
func (c *Catalog) Get(id string) (*Lesson, error) {
	lesson, ok := c.lessons[id]
	if !ok {
		return nil, fmt.Errorf("lesson with id '%s' not found", id)
	}
	return lesson, nil
}

// List returns all lessons in a given track sorted by order.
func (c *Catalog) List(track string) []*Lesson {
	return c.byTrack[track]
}

// Tracks returns the ordered list of all track names.
func (c *Catalog) Tracks() []string {
	res := make([]string, len(c.trackNames))
	copy(res, c.trackNames)
	return res
}

// FirstInTrack returns the initial lesson for a track or nil if empty.
func (c *Catalog) FirstInTrack(track string) *Lesson {
	lessons := c.byTrack[track]
	if len(lessons) == 0 {
		return nil
	}
	return lessons[0]
}

// NextInTrack finds the lesson immediately following currentID in its track.
func (c *Catalog) NextInTrack(track string, currentID string) *Lesson {
	lessons := c.byTrack[track]
	for i, l := range lessons {
		if l.ID == currentID {
			if i+1 < len(lessons) {
				return lessons[i+1]
			}
			return nil
		}
	}
	return nil
}

// PrevInTrack finds the lesson immediately preceding currentID in its track.
func (c *Catalog) PrevInTrack(track string, currentID string) *Lesson {
	lessons := c.byTrack[track]
	for i, l := range lessons {
		if l.ID == currentID {
			if i-1 >= 0 {
				return lessons[i-1]
			}
			return nil
		}
	}
	return nil
}

// TotalLessons returns the count of all lessons across all tracks.
func (c *Catalog) TotalLessons() int {
	return len(c.lessons)
}
