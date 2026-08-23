package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type PathSegment struct {
	Name     string
	Wildcard bool
	Index    *int
}

func ParsePath(path string) ([]PathSegment, error) {
	if path == "" {
		return nil, fmt.Errorf("%w: empty path", ErrInvalid)
	}
	out := []PathSegment{}
	for _, part := range strings.Split(path, ".") {
		seg := PathSegment{Name: part}
		if i := strings.Index(part, "["); i >= 0 {
			if !strings.HasSuffix(part, "]") {
				return nil, fmt.Errorf("%w: malformed path", ErrInvalid)
			}
			seg.Name = part[:i]
			idx := part[i+1 : len(part)-1]
			if idx == "*" {
				seg.Wildcard = true
			} else {
				n, e := strconv.Atoi(idx)
				if e != nil || n < 0 {
					return nil, fmt.Errorf("%w: array index", ErrInvalid)
				}
				seg.Index = &n
			}
		}
		if seg.Name == "" {
			return nil, fmt.Errorf("%w: empty segment", ErrInvalid)
		}
		out = append(out, seg)
	}
	return out, nil
}
func DeepCopyJSON(in map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range in {
		out[k] = cloneJSON(v)
	}
	return out
}
func cloneJSON(v any) any {
	switch x := v.(type) {
	case map[string]any:
		return DeepCopyJSON(x)
	case []any:
		o := make([]any, len(x))
		for i, v := range x {
			o[i] = cloneJSON(v)
		}
		return o
	default:
		return x
	}
}
func TransformPath(root map[string]any, path string, fn func(any) (any, bool)) (int, error) {
	p, e := ParsePath(path)
	if e != nil {
		return 0, e
	}
	return transformMap(root, p, fn), nil
}
func transformMap(m map[string]any, p []PathSegment, fn func(any) (any, bool)) int {
	if len(p) == 0 {
		return 0
	}
	s := p[0]
	v, ok := m[s.Name]
	if !ok {
		return 0
	}
	if len(p) == 1 {
		if s.Index != nil {
			a, ok := v.([]any)
			if !ok || *s.Index >= len(a) {
				return 0
			}
			n, k := fn(a[*s.Index])
			if k {
				a[*s.Index] = n
			} else {
				a = append(a[:*s.Index], a[*s.Index+1:]...)
			}
			m[s.Name] = a
			return 1
		}
		if s.Wildcard {
			a, ok := v.([]any)
			if !ok {
				return 0
			}
			c := 0
			for i, x := range a {
				n, k := fn(x)
				if k {
					a[i] = n
				} else {
					a = append(a[:i], a[i+1:]...)
					i--
				}
				c++
			}
			m[s.Name] = a
			return c
		}
		n, k := fn(v)
		if k {
			m[s.Name] = n
		} else {
			delete(m, s.Name)
		}
		return 1
	}
	rest := p[1]
	if s.Wildcard {
		a, ok := v.([]any)
		if !ok {
			return 0
		}
		c := 0
		for _, x := range a {
			if child, ok := x.(map[string]any); ok {
				c += transformMap(child, p[1:], fn)
			}
		}
		return c
	}
	if s.Index != nil {
		a, ok := v.([]any)
		if !ok || *s.Index >= len(a) {
			return 0
		}
		child, ok := a[*s.Index].(map[string]any)
		if !ok {
			return 0
		}
		return transformMap(child, p[1:], fn)
	}
	child, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	_ = rest
	return transformMap(child, p[1:], fn)
}
