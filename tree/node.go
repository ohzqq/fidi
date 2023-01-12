// Copyright Ivan Sukharev

// The MIT License (MIT)

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package tree

import (
	"io/fs"
	"os"
)

// Один Узел Дерева
type TreeNodeStruct struct {
	Path         string
	Name         string
	Depth        int
	Files        []os.FileInfo
	FilesCount   int
	SubDirs      []os.FileInfo
	SubDirsCount int
}

// Заполняем структуру при сканировании дерева
func (node *TreeNodeStruct) Fill(path string, depth int) error {
	DirInfo, DirInfoErr := os.Lstat(path)
	if DirInfoErr != nil {
		return DirInfoErr
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	files := make([]fs.FileInfo, 0, len(entries))

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		files = append(files, info)
	}

	node.Path = path
	node.Name = DirInfo.Name()
	node.Depth = depth

	for _, f := range files {
		if f.IsDir() {
			node.SubDirs = append(node.SubDirs, f)
		} else {
			node.Files = append(node.Files, f)
		}
	}

	node.FilesCount = len(node.Files)
	node.SubDirsCount = len(node.SubDirs)

	return nil
}
