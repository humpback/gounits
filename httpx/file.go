package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func GetFile(ctx context.Context, save string, path string, query url.Values, headers map[string][]string) error {

	return DefaultClient.GetFile(ctx, save, path, query, headers)
}

func GetFdWith(ctx context.Context, fd *os.File, path string, query url.Values, headers map[string][]string) error {

	return DefaultClient.GetFdWith(ctx, fd, path, query, headers)
}

func (client *HttpClient) GetFile(ctx context.Context, save string, path string, query url.Values, headers map[string][]string) error {

	fpath, err := filepath.Abs(save)
	if err != nil {
		return fmt.Errorf("client pull file save path invalid, %s", save)
	}

	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("client pull file open error, %s", err.Error())
	}

	if err := client.GetFdWith(ctx, fd, path, query, headers); err != nil {
		fd.Close()
		os.Remove(path)
		return err
	}
	fd.Close()
	return nil
}

func (client *HttpClient) GetFdWith(ctx context.Context, fd *os.File, path string, query url.Values, headers map[string][]string) error {

	if fd == nil {
		return fmt.Errorf("client get file fd invalid, %s", path)
	}

	resp, err := client.Get(ctx, path, query, headers)
	if err != nil {
		return err
	}

	defer resp.Close()
	statusCode := resp.StatusCode()
	if statusCode == http.StatusOK {
		_, err := io.Copy(fd, resp.body)
		return err
	}
	return fmt.Errorf("client get file fail %d, %s", statusCode, resp.Status())
}
