// Copyright (c) 2018, The GoKi Authors. All rights reserved.

// BSD 3-Clause License

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// * Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.

// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package dirs provides various utility functions in dealing with directories
// such as a list of all the files with a given (set of) extensions and
// finding paths within the Go source directory (GOPATH, etc)
package dirfile

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Filter func(File) bool

// ExtFiles returns all the DirEntry's for files with given extension(s) in directory
// in sorted order (if exts is empty then all files are returned).
// In case of error, returns nil.
func ExtFiles(path string, exts ...string) []os.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	return FilterDirEntriesByExt(files, exts...)
}

func FilterDirEntriesByExt(files []os.DirEntry, exts ...string) []os.DirEntry {
	if len(exts) == 0 {
		return files
	}
	sz := len(files)
	if sz == 0 {
		return nil
	}
	for i := sz - 1; i >= 0; i-- {
		fn := files[i]
		ext := filepath.Ext(fn.Name())
		keep := false
		for _, ex := range exts {
			if strings.EqualFold(ext, ex) {
				keep = true
				break
			}
		}
		if !keep {
			files = append(files[:i], files[i+1:]...)
		}
	}
	return files
}

func FilterFilesByExt(files []File, exts ...string) []File {
	if len(exts) == 0 {
		return files
	}
	sz := len(files)
	if sz == 0 {
		return nil
	}
	for i := sz - 1; i >= 0; i-- {
		fn := files[i]
		keep := false
		for _, ex := range exts {
			if strings.EqualFold(fn.Extension, ex) {
				keep = true
				break
			}
		}
		if !keep {
			files = append(files[:i], files[i+1:]...)
		}
	}
	return files
}

// ExtFileNames returns all the file names with given extension(s) in directory
// in sorted order (if exts is empty then all files are returned)
func ExtFileNames(path string, exts ...string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	files, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil
	}
	return FilterFileNamesByExt(files, exts...)
}

func FilterFileNamesByExt(files []string, exts ...string) []string {
	if len(exts) == 0 {
		sort.StringSlice(files).Sort()
		return files
	}
	sz := len(files)
	if sz == 0 {
		return nil
	}
	for i := sz - 1; i >= 0; i-- {
		fn := files[i]
		ext := filepath.Ext(fn)
		keep := false
		for _, ex := range exts {
			if strings.EqualFold(ext, ex) {
				keep = true
				break
			}
		}
		if !keep {
			files = append(files[:i], files[i+1:]...)
		}
	}
	sort.StringSlice(files).Sort()
	return files
}

// Dirs returns a slice of all the directories within a given directory
func Dirs(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	var fnms []string
	for _, fi := range files {
		if fi.IsDir() {
			fnms = append(fnms, fi.Name())
		}
	}
	return fnms
}

// AllFiles returns a slice of all the files, recursively, within a given directory
// Due to the nature of the filepath.Walk function, the first entry will be the
// directory itself, for reference -- just skip past that if you don't need it.
func AllFiles(path string) ([]string, error) {
	var fnms []string
	er := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fnms = append(fnms, path)
		return nil
	})
	return fnms, er
}

// HasFile returns true if given directory has given file (exact match)
func HasFile(path, file string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, fn := range files {
		if fn.Name() == file {
			return true
		}
	}
	return false
}

// note: rejected from std lib, but often need: https://github.com/golang/go/issues/25012
// https://github.com/golang/go/issues/5366

// SplitExt returns the base of the file name without extension, and the extension
func SplitExt(fname string) (fbase, ext string) {
	ext = filepath.Ext(fname)
	fbase = strings.TrimSuffix(fname, ext)
	return
}

func FilterFiles(files []File, filter Filter) []File {
	var keep []File
	for _, fn := range files {
		if filter(fn) {
			keep = append(keep, fn)
		}
	}
	return keep
}

func ExtFilter(exts ...string) Filter {
	filter := func(file File) bool {
		for _, ex := range exts {
			if strings.EqualFold(file.Extension, ex) {
				return true
			}
		}
		return false
	}
	return filter
}

func MimeFilter(mimes ...string) Filter {
	filter := func(file File) bool {
		for _, mt := range mimes {
			if strings.Contains(file.Mime, mt) {
				return true
			}
		}
		return false
	}
	return filter
}
