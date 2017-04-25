package utils

import (
    "os"
    "io"
    "path/filepath"
    "errors"
    "io/ioutil"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return
    }

    err = out.Sync()
    if err != nil {
        return
    }

    si, err := os.Stat(src)
    if err != nil {
        return
    }
    err = os.Chmod(dst, si.Mode())
    if err != nil {
        return
    }

    return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
    src = filepath.Clean(src)
    dst = filepath.Clean(dst)

    si, err := os.Stat(src)
    if err != nil {
        return err
    }
    if !si.IsDir() {
        return errors.New("Source is not a directory")
    }

    if _, err = os.Stat(dst); os.IsNotExist(err) {

        err = os.MkdirAll(dst, si.Mode())
        if err != nil {
            return
        }

    }

    entries, err := ioutil.ReadDir(src)
    if err != nil {
        return
    }

    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())

        if entry.IsDir() {
            err = CopyDir(srcPath, dstPath)
            if err != nil {
                return
            }
        } else {
            // Skip symlinks.
            if entry.Mode() & os.ModeSymlink != 0 {
                continue
            }

            err = CopyFile(srcPath, dstPath)
            if err != nil {
                return
            }
        }
    }

    return
}