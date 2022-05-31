package walker

import (
	"errors"
	"io/fs"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ryantate13/walker/internal/mocks"
)

func TestNew(t *testing.T) {
	tests := []struct {
		it     string
		assert func(t *testing.T, w Walker)
	}{
		{
			it: "instantiates a new fs walker",
			assert: func(t *testing.T, w Walker) {
				walkPtr, ok := w.(*walker)
				require.True(t, ok)
				require.NotNil(t, walkPtr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			tt.assert(t, New(&mocks.FakeFS{}))
		})
	}
}

type info struct {
	name  string
	isDir bool
}

type walkResult struct {
	path string
	info *info
	err  error
}

type testCollector struct {
	continueOnStatError bool
	continueOnDirError  bool
	results             []walkResult
}

func (t *testCollector) collect(path string, entry fs.FileInfo, err error) bool {
	result := walkResult{
		path: path,
		info: nil,
		err:  err,
	}
	if entry != nil {
		result.info = &info{
			name:  entry.Name(),
			isDir: entry.IsDir(),
		}
	}
	t.results = append(t.results, result)
	if err == fs.ErrPermission {
		return t.continueOnStatError
	}
	if err != nil {
		return t.continueOnDirError
	}
	return true
}

type scenario struct {
	isDir, statErr, readDirErr bool
}

type mockFS = map[string]scenario

func Test_walker_Walk(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() (string, bool, bool, mockFS)
		assert func(t *testing.T, coll *testCollector)
	}{
		{
			it: "stops iterating if walk func returns false",
			setup: func() (string, bool, bool, mockFS) {
				startDir := "a"
				continueOnStatError := false
				continueOnDirError := true
				f := mockFS{
					"a": {true, true, false},
					"b": {false, false, false},
				}
				return startDir, continueOnStatError, continueOnDirError, f
			},
			assert: func(t *testing.T, coll *testCollector) {
				require.Equal(t, []walkResult{
					{
						path: "a",
						info: nil,
						err:  fs.ErrPermission,
					},
				}, coll.results)
			},
		},
		{
			it: "passes both stat errors and read dir errors to handler func",
			setup: func() (string, bool, bool, mockFS) {
				startDir := "a"
				continueOnStatError := true
				continueOnDirError := false
				f := mockFS{
					"a": {true, false, true},
				}
				return startDir, continueOnStatError, continueOnDirError, f
			},
			assert: func(t *testing.T, coll *testCollector) {
				require.Equal(t, []walkResult{
					{
						path: "a",
						info: &info{
							name:  "a",
							isDir: true,
						},
						err: nil,
					},
					{
						path: "a",
						info: &info{
							name:  "a",
							isDir: true,
						},
						err: &fs.PathError{},
					},
				}, coll.results)
			},
		},
		{
			it: "recursively walks a directory and continues past errors when configured to",
			setup: func() (string, bool, bool, mockFS) {
				startDir := "a"
				continueOnStatError := true
				continueOnDirError := true
				f := mockFS{
					"a":     {true, false, false},
					"a/b":   {true, false, true},
					"a/b/c": {false, false, false},
					"a/d":   {false, false, false},
				}
				return startDir, continueOnStatError, continueOnDirError, f
			},
			assert: func(t *testing.T, coll *testCollector) {
				require.Equal(t, []walkResult{
					{
						path: "a",
						info: &info{
							name:  "a",
							isDir: true,
						},
						err: nil,
					},
					{
						path: "a/b",
						info: &info{
							name:  "b",
							isDir: true,
						},
						err: nil,
					},
					{
						path: "a/b",
						info: &info{
							name:  "b",
							isDir: true,
						},
						err: &fs.PathError{},
					},
					{
						path: "a/d",
						info: &info{
							name:  "d",
							isDir: false,
						},
						err: nil,
					},
				}, coll.results)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			dir, continueOnStatError, continueOnDirError, fsLayer := tt.setup()
			f := &mocks.FakeFS{
				OpenStub: func(s string) (fs.File, error) {
					entry, ok := fsLayer[s]
					if !ok {
						return nil, fs.ErrNotExist
					}
					return &mocks.FakeReadDirFile{
						ReadDirStub: func(i int) ([]fs.DirEntry, error) {
							if entry.readDirErr {
								return nil, &fs.PathError{}
							}
							entries := make([]fs.DirEntry, 0)
							for k, v := range fsLayer {
								sDepth := strings.Count(s, "/")
								kDepth := strings.Count(k, "/")
								if strings.HasPrefix(k, s) && k != s && kDepth == sDepth+1 {
									entries = append(entries, func(k string, v scenario) fs.DirEntry {
										return &mocks.FakeDirEntry{
											IsDirStub: func() bool {
												return v.isDir
											},
											NameStub: func() string {
												return path.Base(k)
											},
										}
									}(k, v))
								}
							}
							sort.Slice(entries, func(i, j int) bool {
								return entries[i].Name() < entries[j].Name()
							})
							return entries, nil
						},
						StatStub: func() (fs.FileInfo, error) {
							if entry.statErr {
								return nil, fs.ErrPermission
							}
							return &mocks.FakeFileInfo{
								IsDirStub: func() bool {
									return entry.isDir
								},
								NameStub: func() string {
									return path.Base(s)
								},
							}, nil
						},
					}, nil
				},
			}
			testWalker := New(f)
			collector := &testCollector{
				continueOnStatError: continueOnStatError,
				continueOnDirError:  continueOnDirError,
			}
			testWalker.Walk(dir, collector.collect)
			tt.assert(t, collector)
		})
	}
}

func Test_walker_MustWalk(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() (string, fs.FS)
		assert func(t *testing.T, err interface{})
	}{
		{
			it: "panics if an error is encountered",
			setup: func() (string, fs.FS) {
				return "test", &mocks.FakeFS{
					OpenStub: func(s string) (fs.File, error) {
						return nil, errors.New("test")
					},
				}
			},
			assert: func(t *testing.T, err interface{}) {
				actualErr, ok := err.(error)
				require.True(t, ok)
				require.Equal(t, "test", actualErr.Error())
			},
		},
		{
			it: "behaves identically to walk in the absence of errors",
			setup: func() (string, fs.FS) {
				return "test", &mocks.FakeFS{
					OpenStub: func(s string) (fs.File, error) {
						return &mocks.FakeReadDirFile{
							StatStub: func() (fs.FileInfo, error) {
								return &mocks.FakeFileInfo{
									IsDirStub: func() bool {
										return false
									},
									NameStub: func() string {
										return s
									},
								}, nil
							},
						}, nil
					},
				}
			},
			assert: func(t *testing.T, err interface{}) {
				require.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			defer func() {
				err := recover()
				tt.assert(t, err)
			}()
			dir, testFS := tt.setup()
			w := New(testFS)
			w.MustWalk(dir, func(_ string, _ fs.FileInfo) bool {
				return true
			})
		})
	}
}
