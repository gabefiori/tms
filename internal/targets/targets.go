package targets

import (
	"bytes"
	"errors"
	"io"
	"maps"
	"slices"
	"sort"
	"strings"
	"sync"
)

type InputTarget struct {
	Path  string `json:"path"`
	Depth uint8  `json:"depth"`
}

type Targets struct {
	Buf *bytes.Buffer
}

func New() *Targets {
	return &Targets{}
}

func Collect(input []InputTarget) (*Targets, error) {
	tg := New()
	err := tg.Collect(input)

	if err != nil {
		return nil, err
	}

	return tg, nil
}

func (t *Targets) Print(w io.Writer) error {
	if t.Buf.Len() == 0 {
		return nil
	}

	_, err := io.Copy(w, t.Buf)

	return err
}

func (t *Targets) Filter(filter string) {
	buf := new(bytes.Buffer)

	for {
		line, err := t.Buf.ReadString('\n')

		if strings.Contains(line, filter) {
			buf.WriteString(line)
		}

		if err != nil {
			break
		}
	}

	t.Buf = buf
}

func (t *Targets) Collect(input []InputTarget) error {
	var (
		wg sync.WaitGroup
		mu sync.Mutex

		targetsMap = make(map[string]struct{})
		errs       []error
	)

	for _, inputTarget := range input {
		wg.Add(1)

		go func(it InputTarget) {
			defer wg.Done()

			found, err := Find(it.Path, it.Depth)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, err)
				return
			}

			for _, f := range found {
				targetsMap[f] = struct{}{}
			}
		}(inputTarget)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	sortedTargets := slices.Collect(maps.Keys(targetsMap))

	sort.Slice(sortedTargets, func(i, j int) bool {
		return sortedTargets[i] > sortedTargets[j]
	})

	t.Buf = bytes.NewBufferString(strings.Join(sortedTargets, "\n"))

	return nil
}
