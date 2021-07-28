package files

import (
	"io/fs"
	"path/filepath"
	"sort"
)

func Find(root, ext string) (files []string, err error) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	return files, err
}
