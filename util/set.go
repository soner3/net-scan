/*
Copyright © 2025 Soner Astan <sonerastan@icloud.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package util

import "slices"

type Set[T comparable] map[T]struct{}

func (s *Set[T]) Add(item T) {
	(*s)[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	delete(*s, item)
}

func (s *Set[T]) ToSortedSlice(cmp func(a, b T) int) *[]T {
	slice := make([]T, 0, len(*s))
	for item := range *s {
		slice = append(slice, item)
	}

	slices.SortFunc(slice, cmp)
	return &slice

}
