package task

import "fmt"

// ValidateDAG accepts a closed local graph. External dependencies are checked by the application.
func ValidateDAG(nodes []Task) error {
	byID := make(map[string]Task, len(nodes))
	for _, t := range nodes {
		if err := t.Validate(); err != nil {
			return err
		}
		if _, ok := byID[t.ID]; ok {
			return ErrConflict
		}
		byID[t.ID] = t
	}
	color := map[string]uint8{}
	var visit func(string) error
	visit = func(id string) error {
		if color[id] == 1 {
			return fmt.Errorf("%w: cycle at %s", ErrInvalid, id)
		}
		if color[id] == 2 {
			return nil
		}
		color[id] = 1
		seen := map[string]bool{}
		for _, dep := range byID[id].Refs {
			if seen[dep] {
				return fmt.Errorf("%w: duplicate dependency", ErrInvalid)
			}
			seen[dep] = true
			if _, local := byID[dep]; local {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		color[id] = 2
		return nil
	}
	for id := range byID {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
