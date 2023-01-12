// Copyright Ivan Sukharev

// The MIT License (MIT)

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package fitree provides function to read folders tree with files

package tree

import (
	"os"
)

const StartDepth = 1

// Список Узлов Дерева
type TreeNodeListStruct struct {
	List    []TreeNodeStruct
	Reverse map[string]int
}

// Добавляем Узел
func (list *TreeNodeListStruct) Add(node TreeNodeStruct) {
	list.List = append(list.List, node)

	if list.Reverse == nil {
		list.Reverse = make(map[string]int)
	}

	list.Reverse[node.Path] = len(list.List) - 1
}

// Сканируем Дерево Каталогов
func (list *TreeNodeListStruct) Scan(path string, depth int, ignoreErr bool) error {
	path += string(os.PathSeparator)

	node := TreeNodeStruct{}

	fillErr := node.Fill(path, depth)
	if fillErr != nil && !ignoreErr {
		return fillErr
	}

	list.Add(node)

	depth++

	for _, f := range node.SubDirs {
		scanErr := list.Scan(path+f.Name(), depth, ignoreErr)
		if scanErr != nil && !ignoreErr {
			return scanErr
		}
	}

	return nil
}

// Получить ноду по индексу
func (list *TreeNodeListStruct) GetNode(index int) (*TreeNodeStruct, error) {
	if len(list.List) < index+1 {
		return nil, &NodeIndexDontExistsError{Index: index}
	}

	return &list.List[index], nil
}
