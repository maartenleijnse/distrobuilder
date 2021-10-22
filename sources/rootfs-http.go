package sources

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/lxc/distrobuilder/shared"
	lxd "github.com/lxc/lxd/shared"
)

type rootfs struct {
	common
}

// Run downloads a tarball.
func (s *rootfs) Run() error {
	URL, err := url.Parse(s.definition.Source.URL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %w", err)
	}

	var fpath string
	var filename string

	if URL.Scheme == "file" {
		fpath = filepath.Dir(URL.Path)
		filename = filepath.Base(URL.Path)
	} else {
		fpath, err = shared.DownloadHash(s.definition.Image, s.definition.Source.URL, "", nil)
		if err != nil {
			return fmt.Errorf("Failed to download %q: %w", s.
				definition.Source.URL, err)
		}

		filename = path.Base(s.definition.Source.URL)
	}

	s.logger.Infow("Unpacking image", "file", filepath.Join(fpath, filename))

	// Unpack
	err = lxd.Unpack(filepath.Join(fpath, filename), s.rootfsDir, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to unpack %q: %w", filepath.Join(fpath, filename), err)
	}

	return nil
}
