package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

// resolveDemoPath returns a path to a plain .dem file, transparently decompressing
// rawPath first if it's a .zst archive. The returned cleanup func (if non-nil) removes
// the temporary decompressed file and must be called once the caller is done with it.
func resolveDemoPath(rawPath string) (path string, cleanup func(), err error) {
	if !strings.EqualFold(filepath.Ext(rawPath), ".zst") {
		return rawPath, nil, nil
	}

	in, err := os.Open(rawPath)
	if err != nil {
		return "", nil, err
	}
	defer in.Close()

	dec, err := zstd.NewReader(in)
	if err != nil {
		return "", nil, err
	}
	defer dec.Close()

	tmpDir, err := os.MkdirTemp("", "faceit-voicechat-*")
	if err != nil {
		return "", nil, err
	}

	outPath := filepath.Join(tmpDir, strings.TrimSuffix(filepath.Base(rawPath), filepath.Ext(rawPath)))
	out, err := os.Create(outPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, dec); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	return outPath, func() { os.RemoveAll(tmpDir) }, nil
}

func parseDemo(demoPath string) (ctMask uint32, tMask uint32, err error) {
	f, err := os.Open(demoPath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

	var ctIDs []int
	var tIDs []int

	p.RegisterEventHandler(func(e events.MatchStart) {
		participants := p.GameState().Participants()
		cts := participants.TeamMembers(common.TeamCounterTerrorists)
		ts := participants.TeamMembers(common.TeamTerrorists)

		for _, ct := range cts {
			ctIDs = append(ctIDs, ct.EntityID)
		}
		for _, t := range ts {
			tIDs = append(tIDs, t.EntityID)
		}
	})

	if err := p.ParseToEnd(); err != nil {
		return 0, 0, err
	}

	return idsToBitmask(ctIDs), idsToBitmask(tIDs), nil
}

func idsToBitmask(ids []int) uint32 {
	var mask uint32
	for _, id := range ids {
		if id >= 1 && id <= 32 {
			mask |= 1 << (id - 1)
		}
	}
	return mask
}

func bindCommand(cfg config, ctMask, tMask uint32) string {
	return fmt.Sprintf(
		`bind "%s" "tv_listen_voice_indices %d; tv_listen_voice_indices_h %d"; bind "%s" "tv_listen_voice_indices %d; tv_listen_voice_indices_h %d"; bind "%s" "tv_listen_voice_indices -1; tv_listen_voice_indices_h -1"`,
		cfg.Keys[0], ctMask, ctMask,
		cfg.Keys[1], tMask, tMask,
		cfg.Keys[2],
	)
}

func playdemoName(demoPath string) string {
	base := filepath.Base(demoPath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func playdemoCommand(gameFolder, demoPath string) string {
	name := playdemoName(demoPath)
	if gameFolder == "" || demoPath == "" {
		return "playdemo " + name
	}
	rel, err := filepath.Rel(gameFolder, demoPath)
	if err != nil {
		return "playdemo " + name
	}
	rel = filepath.ToSlash(rel)
	if rel == "." || strings.HasPrefix(rel, "../") {
		return "playdemo " + name
	}
	rel = strings.TrimSuffix(rel, filepath.Ext(rel))
	return "playdemo " + rel
}

func copyDemoToGameFolder(cfg config, demoPath string) (string, error) {
	if cfg.GameFolder == "" {
		return "", fmt.Errorf("game folder is not set")
	}
	dst := filepath.Join(cfg.GameFolder, filepath.Base(demoPath))
	if err := copyFile(demoPath, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func isDemoFilename(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".dem") || strings.HasSuffix(lower, ".dem.zst")
}
