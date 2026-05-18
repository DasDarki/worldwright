package http

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"worldwright/backend/internal/store"
)

var allowedMime = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

const maxUploadBytes = 20 * 1024 * 1024

func listAssets(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "50"))
		assets, err := st.ListAssets(c.UserContext(), limit)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"assets": assets})
	}
}

func uploadAsset(st *store.Store, assetsDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "missing file field")
		}
		if fh.Size > maxUploadBytes {
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "file too large")
		}
		mime := fh.Header.Get("Content-Type")
		ext, ok := allowedMime[mime]
		if !ok {
			if guess := mimeFromExt(fh.Filename); guess != "" {
				mime = guess
				ext = allowedMime[mime]
				ok = true
			}
		}
		if !ok {
			return fiber.NewError(fiber.StatusUnsupportedMediaType, "unsupported file type")
		}

		src, err := fh.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		buf, err := io.ReadAll(src)
		if err != nil {
			return err
		}

		var width, height *int
		if cfg, _, derr := image.DecodeConfig(bytes.NewReader(buf)); derr == nil {
			w := cfg.Width
			h := cfg.Height
			width = &w
			height = &h
		}

		name := uuid.NewString() + ext
		if err := os.MkdirAll(assetsDir, 0o755); err != nil {
			return err
		}
		dst := filepath.Join(assetsDir, name)
		if err := os.WriteFile(dst, buf, 0o644); err != nil {
			return err
		}

		asset, err := st.CreateAsset(c.UserContext(), store.NewAsset{
			Filename: fh.Filename,
			Path:     name,
			Mime:     mime,
			Size:     int64(len(buf)),
			Width:    width,
			Height:   height,
		})
		if err != nil {
			_ = os.Remove(dst)
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"asset": asset})
	}
}

func deleteAsset(st *store.Store, assetsDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		a, err := st.AssetByID(c.UserContext(), id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return fiber.NewError(fiber.StatusNotFound, "asset not found")
			}
			return err
		}
		if err := st.DeleteAsset(c.UserContext(), id); err != nil {
			return err
		}
		_ = os.Remove(filepath.Join(assetsDir, a.Path))
		return c.JSON(fiber.Map{"ok": true})
	}
}

func serveAsset(st *store.Store, assetsDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		a, err := st.AssetByID(c.UserContext(), id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return fiber.NewError(fiber.StatusNotFound, "asset not found")
			}
			return err
		}
		c.Set(fiber.HeaderContentType, a.Mime)
		c.Set(fiber.HeaderCacheControl, "public, max-age=604800")
		return c.SendFile(filepath.Join(assetsDir, a.Path), false)
	}
}

func mimeFromExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	}
	return ""
}
