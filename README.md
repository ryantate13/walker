![📁: Walker](./banner.svg)

# Walker

Walker is a go (1.18+) package that provides recursive directory iteration in which directories are processed prior to files, in a depth first vs. breadth first fashion.

For example, given the following directory layout:

```console
$ find /a -type f -name '*.txt'
/a/baz.txt
/a/b/bar.txt
/a/b/c/foo.txt
```

The order in which files are processed will be `foo.txt`, `bar.txt` and finally `baz.txt`.

## Usage:

The following example will recursively print out all paths of files and directories on the supplied filesystem. The walker is instantiated with an instance of the [`fs.FS`](https://pkg.go.dev/io/fs#FS) interface and given a directory to walk and a callback function. The callback is invoked with the full path to each file and directory, the corresponding [`fs.FileInfo`](https://pkg.go.dev/io/fs#FileInfo) related to it, and optionally an error. For directories, there are two types of errors that can be encountered, errors calling [`Stat`](https://pkg.go.dev/io/fs#Stat) on the directory and errors listing the directory's contents. In the latter case, the callback will be invoked twice, once with a `nil` error and again when the call to [`ReadDir`](https://pkg.go.dev/io/fs#ReadDir) fails. Returning `false` from the supplied callback will halt any further iteration.

```go
package main

import (
    "fmt"
    "io/fs"
    "os"
    
    "github.com/ryantate13/walker"
)

func main() {
    w := walker.New(os.DirFS("/some/dir"))
    w.Walk("path/to/subdirectory", func(path string, entry fs.FileInfo, err error) bool {
        if err != nil {
            // handle error
        }
        fmt.Println(path)
        return true // return false from handler to stop iteration
    })
}
```

As an alternative to calling `Walk`, in cases where errors encountered should generate a panic, there is the alternative `MustWalk` which is identical to `Walk` except that any error will automatically induce a panic.

```go
w.MustWalk("path/to/subdirectory", func(path string, entry fs.FileInfo) bool {
    // no error handling is necessary but any error encountered will trigger a panic
    fmt.Println(path)
    return true
})
```

For more examples of usage, refer to the [`examples`](./examples) directory. For an example of how walker can be used to traverse filesystem-like hierarchical objects other than a directory on the actual filesystem, [refer to the test suite](./walker_test.go#L148).