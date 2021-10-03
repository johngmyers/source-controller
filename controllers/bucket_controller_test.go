/*
Copyright 2021 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBucketReconciler_checksum(t *testing.T) {
	tests := []struct {
		name       string
		beforeFunc func(root string)
		want       string
		wantErr    bool
	}{
		{
			name: "empty root",
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "with file",
			beforeFunc: func(root string) {
				mockFile(root, "a/b/c.txt", "a dummy string")
			},
			want: "65a8ff8c0feac7c8710ece851923f173e45cb89dc264378ed291074e5759aa61",
		},
		{
			name: "with file in different path",
			beforeFunc: func(root string) {
				mockFile(root, "a/b.txt", "a dummy string")
			},
			want: "361fe57151a3893cba2186027f303a4f7c7c85623d05c99b003c88522bf88db4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := os.MkdirTemp("", "bucket-checksum-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(root)
			if tt.beforeFunc != nil {
				tt.beforeFunc(root)
			}
			got, err := (&BucketReconciler{}).checksum(root)
			if (err != nil) != tt.wantErr {
				t.Errorf("checksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockFile(root, path, content string) error {
	filePath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		panic(err)
	}
	return nil
}
