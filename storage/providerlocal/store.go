package providerlocal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/3086953492/gokit/storage"
)

// ProviderName 本地文件系统 Provider 名称。
const ProviderName = "local"

const (
	defaultDirPerm  = 0o755
	defaultFilePerm = 0o644
)

// Store 本地文件系统存储实现。
type Store struct {
	root     string
	baseURL  string
	baseHost string
	basePath string
	dirPerm  os.FileMode
	filePerm os.FileMode
}

type urlStore struct {
	*Store
}

type listEntry struct {
	token string
	meta  *storage.ObjectMeta
}

type sectionReadCloser struct {
	reader io.Reader
	file   *os.File
}

func (s *sectionReadCloser) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

func (s *sectionReadCloser) Close() error {
	return s.file.Close()
}

// New 创建本地文件系统存储实现。
func New(cfg Config) (storage.Store, error) {
	if cfg.Root == "" {
		return nil, fmt.Errorf("%w: Root is required", storage.ErrInvalidConfig)
	}

	absRoot, err := filepath.Abs(cfg.Root)
	if err != nil {
		return nil, fmt.Errorf("%w: resolve Root: %v", storage.ErrInvalidConfig, err)
	}

	store := &Store{
		root:     filepath.Clean(absRoot),
		dirPerm:  defaultOrMode(cfg.DirPerm, defaultDirPerm),
		filePerm: defaultOrMode(cfg.FilePerm, defaultFilePerm),
	}

	if cfg.BaseURL == "" {
		return store, nil
	}

	baseURL, baseHost, basePath, err := normalizeBaseURL(cfg.BaseURL)
	if err != nil {
		return nil, err
	}
	store.baseURL = baseURL
	store.baseHost = baseHost
	store.basePath = basePath

	return &urlStore{Store: store}, nil
}

// Upload 上传对象到本地文件系统。
func (s *Store) Upload(ctx context.Context, key string, r io.Reader, opts *storage.WriteOptions) (*storage.ObjectMeta, error) {
	filePath, err := s.resolveKeyPath(key)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), s.dirPerm); err != nil {
		return nil, wrapPathError("mkdir", err)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(filePath), ".storage-*")
	if err != nil {
		return nil, wrapPathError("create temp file", err)
	}
	tmpName := tmpFile.Name()

	success := false
	defer func() {
		if success {
			return
		}
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
	}()

	written, err := io.Copy(tmpFile, &contextReader{ctx: ctx, reader: r})
	if err != nil {
		return nil, wrapPathError("write file", err)
	}

	if err := tmpFile.Chmod(s.filePerm); err != nil {
		return nil, wrapPathError("chmod file", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, wrapPathError("close file", err)
	}
	if err := replaceFile(tmpName, filePath); err != nil {
		return nil, wrapPathError("rename file", err)
	}

	success = true

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, wrapPathError("stat file", err)
	}

	return s.newObjectMeta(key, info, written, contentTypeFromUpload(key, opts)), nil
}

// Download 从本地文件系统下载对象。
func (s *Store) Download(ctx context.Context, key string, opts *storage.ReadOptions) (io.ReadCloser, *storage.ObjectMeta, error) {
	filePath, err := s.resolveKeyPath(key)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, wrapPathError("open file", err)
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, nil, wrapPathError("stat file", err)
	}

	if err := ctx.Err(); err != nil {
		_ = file.Close()
		return nil, nil, err
	}

	meta := s.newObjectMeta(key, info, info.Size(), detectContentType(key))
	if opts == nil || opts.Range == "" {
		return file, meta, nil
	}

	offset, length, err := parseRange(opts.Range, info.Size())
	if err != nil {
		_ = file.Close()
		return nil, nil, err
	}

	reader := io.NewSectionReader(file, offset, length)
	return &sectionReadCloser{reader: reader, file: file}, meta, nil
}

// Delete 删除本地文件系统中的对象。
func (s *Store) Delete(ctx context.Context, key string, opts *storage.DeleteOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	filePath, err := s.resolveKeyPath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return wrapPathError("delete file", err)
	}
	return nil
}

// List 列举本地文件系统中的对象。
func (s *Store) List(ctx context.Context, prefix string, opts *storage.ListOptions) (*storage.ListResult, error) {
	normalizedPrefix, err := normalizePrefix(prefix)
	if err != nil {
		return nil, err
	}

	listOpts := applyListDefaults(opts)
	keys, err := s.collectKeys(ctx)
	if err != nil {
		return nil, err
	}

	entries, err := s.buildListEntries(ctx, keys, normalizedPrefix, listOpts)
	if err != nil {
		return nil, err
	}

	result := &storage.ListResult{
		Objects:        make([]*storage.ObjectMeta, 0),
		CommonPrefixes: make([]string, 0),
	}

	if len(entries) > listOpts.MaxKeys {
		result.IsTruncated = true
		result.NextMarker = entries[listOpts.MaxKeys-1].token
		entries = entries[:listOpts.MaxKeys]
	}

	for _, entry := range entries {
		if entry.meta != nil {
			result.Objects = append(result.Objects, entry.meta)
			continue
		}
		result.CommonPrefixes = append(result.CommonPrefixes, entry.token)
	}

	return result, nil
}

// Exists 检查本地文件系统中的对象是否存在。
func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	filePath, err := s.resolveKeyPath(key)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(filePath)
	if err == nil {
		return info.Mode().IsRegular(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, wrapPathError("stat file", err)
}

// Head 获取本地文件系统中的对象元信息。
func (s *Store) Head(ctx context.Context, key string) (*storage.ObjectMeta, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	filePath, err := s.resolveKeyPath(key)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, wrapPathError("stat file", err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: %s", storage.ErrNotFound, key)
	}

	return s.newObjectMeta(key, info, info.Size(), detectContentType(key)), nil
}

// AllowedHosts 返回当前 Store 允许的域名列表。
func (s *urlStore) AllowedHosts() []string {
	return []string{s.baseHost}
}

// KeyFromURL 从已解析的 URL 提取对象 key。
func (s *urlStore) KeyFromURL(u *url.URL) (string, error) {
	pathValue := u.Path
	if pathValue == "" || pathValue == "/" {
		return "", fmt.Errorf("%w: empty path", storage.ErrInvalidURL)
	}

	if s.basePath != "" {
		if pathValue == s.basePath {
			return "", fmt.Errorf("%w: empty key path", storage.ErrInvalidURL)
		}

		prefix := s.basePath + "/"
		if !strings.HasPrefix(pathValue, prefix) {
			return "", fmt.Errorf("%w: path %q", storage.ErrInvalidURL, pathValue)
		}
		pathValue = strings.TrimPrefix(pathValue, prefix)
	} else {
		pathValue = strings.TrimPrefix(pathValue, "/")
	}

	key, err := storage.UnescapeKey(pathValue)
	if err != nil {
		return "", fmt.Errorf("%w: %v", storage.ErrInvalidURL, err)
	}
	if key == "" {
		return "", fmt.Errorf("%w: empty key after decode", storage.ErrInvalidURL)
	}
	if _, err := s.resolveKeyPath(key); err != nil {
		return "", err
	}

	return key, nil
}

func (s *Store) newObjectMeta(key string, info os.FileInfo, size int64, contentType string) *storage.ObjectMeta {
	meta := &storage.ObjectMeta{
		Key:          key,
		Size:         size,
		ContentType:  contentType,
		LastModified: info.ModTime(),
	}
	if s.baseURL != "" {
		meta.URL = s.objectURL(key)
	}
	return meta
}

func (s *Store) objectURL(key string) string {
	return s.baseURL + "/" + storage.EscapeKey(key)
}

func (s *Store) resolveKeyPath(key string) (string, error) {
	logicalKey, err := normalizeKey(key)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(s.root, filepath.FromSlash(logicalKey))
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("%w: resolve key path: %v", storage.ErrInvalidKey, err)
	}

	rel, err := filepath.Rel(s.root, absPath)
	if err != nil {
		return "", fmt.Errorf("%w: resolve relative path: %v", storage.ErrInvalidKey, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: key escapes storage root", storage.ErrInvalidKey)
	}

	return absPath, nil
}

func (s *Store) collectKeys(ctx context.Context) ([]string, error) {
	keys := make([]string, 0)

	if _, err := os.Stat(s.root); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return keys, nil
		}
		return nil, wrapPathError("stat storage root", err)
	}

	err := filepath.WalkDir(s.root, func(currentPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return wrapPathError("walk storage root", walkErr)
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(s.root, currentPath)
		if err != nil {
			return fmt.Errorf("%w: build relative path: %v", storage.ErrBackendUnavailable, err)
		}
		keys = append(keys, filepath.ToSlash(relPath))
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(keys)
	return keys, nil
}

func (s *Store) buildListEntries(ctx context.Context, keys []string, prefix string, opts *storage.ListOptions) ([]listEntry, error) {
	entries := make([]listEntry, 0)
	commonPrefixSeen := make(map[string]struct{})

	for _, key := range keys {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}

		if opts.Delimiter != "" {
			rest := strings.TrimPrefix(key, prefix)
			if idx := strings.Index(rest, opts.Delimiter); idx >= 0 {
				commonPrefix := prefix + rest[:idx+len(opts.Delimiter)]
				if commonPrefix <= opts.Marker {
					continue
				}
				if _, ok := commonPrefixSeen[commonPrefix]; ok {
					continue
				}
				commonPrefixSeen[commonPrefix] = struct{}{}
				entries = append(entries, listEntry{token: commonPrefix})
				continue
			}
		}

		if key <= opts.Marker {
			continue
		}

		info, err := os.Stat(filepath.Join(s.root, filepath.FromSlash(key)))
		if err != nil {
			return nil, wrapPathError("stat file", err)
		}
		if !info.Mode().IsRegular() {
			continue
		}

		entries = append(entries, listEntry{
			token: key,
			meta:  s.newObjectMeta(key, info, info.Size(), detectContentType(key)),
		})
	}

	return entries, nil
}

func normalizeKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("%w: key cannot be empty", storage.ErrInvalidKey)
	}
	if strings.Contains(key, "\\") {
		return "", fmt.Errorf("%w: key must use '/' separator", storage.ErrInvalidKey)
	}

	cleaned := path.Clean(key)
	if cleaned == "." || cleaned == "" {
		return "", fmt.Errorf("%w: key cannot be empty", storage.ErrInvalidKey)
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return "", fmt.Errorf("%w: invalid key path", storage.ErrInvalidKey)
	}
	if filepath.IsAbs(filepath.FromSlash(cleaned)) {
		return "", fmt.Errorf("%w: absolute key path is not allowed", storage.ErrInvalidKey)
	}

	return cleaned, nil
}

func normalizePrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", nil
	}
	if strings.Contains(prefix, "\\") {
		return "", fmt.Errorf("%w: prefix must use '/' separator", storage.ErrInvalidKey)
	}

	hasTrailingSlash := strings.HasSuffix(prefix, "/")
	cleaned := path.Clean(prefix)
	if cleaned == "." {
		return "", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return "", fmt.Errorf("%w: invalid prefix path", storage.ErrInvalidKey)
	}
	if filepath.IsAbs(filepath.FromSlash(cleaned)) {
		return "", fmt.Errorf("%w: absolute prefix path is not allowed", storage.ErrInvalidKey)
	}
	if hasTrailingSlash && !strings.HasSuffix(cleaned, "/") {
		cleaned += "/"
	}

	return cleaned, nil
}

func normalizeBaseURL(raw string) (baseURL string, host string, basePath string, err error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", "", fmt.Errorf("%w: parse BaseURL: %v", storage.ErrInvalidConfig, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", "", fmt.Errorf("%w: BaseURL must use http or https", storage.ErrInvalidConfig)
	}
	if parsed.Host == "" {
		return "", "", "", fmt.Errorf("%w: BaseURL host is required", storage.ErrInvalidConfig)
	}

	basePath = path.Clean(strings.TrimSuffix(parsed.Path, "/"))
	if basePath == "." {
		basePath = ""
	}
	if basePath == "/" {
		basePath = ""
	}

	parsed.Path = basePath
	parsed.RawPath = ""

	return strings.TrimSuffix(parsed.String(), "/"), parsed.Host, basePath, nil
}

func defaultOrMode(mode os.FileMode, fallback os.FileMode) os.FileMode {
	if mode == 0 {
		return fallback
	}
	return mode
}

func replaceFile(src string, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	if _, err := os.Stat(dst); err != nil {
		return err
	}
	if err := os.Remove(dst); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

func contentTypeFromUpload(key string, opts *storage.WriteOptions) string {
	if opts != nil && opts.ContentType != "" {
		return opts.ContentType
	}
	return detectContentType(key)
}

func detectContentType(key string) string {
	return mime.TypeByExtension(filepath.Ext(key))
}

func wrapPathError(action string, err error) error {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("%s: %w", action, storage.ErrNotFound)
	case errors.Is(err, os.ErrPermission):
		return fmt.Errorf("%s: %w", action, storage.ErrPermissionDenied)
	default:
		return fmt.Errorf("%s: %w", action, err)
	}
}

func applyListDefaults(opts *storage.ListOptions) *storage.ListOptions {
	if opts == nil {
		return &storage.ListOptions{MaxKeys: storage.DefaultMaxKeys}
	}

	cloned := *opts
	if cloned.MaxKeys <= 0 {
		cloned.MaxKeys = storage.DefaultMaxKeys
	}
	return &cloned
}

func parseRange(raw string, size int64) (offset int64, length int64, err error) {
	rangeErr := func() error { return fmt.Errorf("%w: %q", storage.ErrInvalidRange, raw) }

	const prefix = "bytes="
	if !strings.HasPrefix(raw, prefix) {
		return 0, 0, rangeErr()
	}

	spec := strings.TrimPrefix(raw, prefix)
	if strings.Count(spec, "-") != 1 {
		return 0, 0, rangeErr()
	}

	parts := strings.SplitN(spec, "-", 2)
	startPart := strings.TrimSpace(parts[0])
	endPart := strings.TrimSpace(parts[1])

	switch {
	case startPart == "" && endPart == "":
		return 0, 0, rangeErr()
	case startPart == "":
		lastBytes, err := strconv.ParseInt(endPart, 10, 64)
		if err != nil || lastBytes <= 0 {
			return 0, 0, rangeErr()
		}
		if lastBytes > size {
			lastBytes = size
		}
		return size - lastBytes, lastBytes, nil
	case endPart == "":
		start, err := strconv.ParseInt(startPart, 10, 64)
		if err != nil || start < 0 || start > size {
			return 0, 0, rangeErr()
		}
		return start, size - start, nil
	default:
		start, err := strconv.ParseInt(startPart, 10, 64)
		if err != nil || start < 0 {
			return 0, 0, rangeErr()
		}
		end, err := strconv.ParseInt(endPart, 10, 64)
		if err != nil || end < start {
			return 0, 0, rangeErr()
		}
		if start >= size {
			return 0, 0, rangeErr()
		}
		if end >= size {
			end = size - 1
		}
		return start, end - start + 1, nil
	}
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

var _ storage.Store = (*Store)(nil)
var _ storage.Store = (*urlStore)(nil)
var _ storage.URLKeyResolver = (*urlStore)(nil)
