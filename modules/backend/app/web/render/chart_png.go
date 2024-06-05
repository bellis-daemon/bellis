package render

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/minoic/glgf"
)

const (
	HTML               = "html"
	FileProtocol       = "file://"
	EchartsInstanceDom = "div[_echarts_instance_]"
	CanvasJs           = "echarts.getInstanceByDom(document.querySelector('div[_echarts_instance_]'))" +
		".getDataURL({type: '%s', pixelRatio: %d, excludeComponents: ['toolbox']})"
)

type SnapshotConfig struct {
	// RenderContent the content bytes of charts after rendered
	RenderContent []byte
	// Quality the generated image quality, aka pixelRatio
	Quality int
	// KeepHtml whether keep the generated html also, default false
	KeepHtml bool
	// Timeout  the timeout config
	Timeout time.Duration
	ExtName string
}

type SnapshotConfigOption func(config *SnapshotConfig)

func NewSnapshotConfig(content []byte, opts ...SnapshotConfigOption) *SnapshotConfig {
	config := &SnapshotConfig{
		RenderContent: content,
		Quality:       1,
		KeepHtml:      false,
		Timeout:       0,
		ExtName:       "png",
	}

	for _, o := range opts {
		o(config)
	}
	return config
}

func MakeChartSnapshotPng(content []byte) ([]byte, error) {
	return makeSnapshot(NewSnapshotConfig(content))
}

func MakeChartSnapshotJpg(content []byte) ([]byte, error) {
	conf := NewSnapshotConfig(content)
	conf.ExtName = "jpg"
	return makeSnapshot(conf)
}

func makeSnapshot(config *SnapshotConfig) ([]byte, error) {
	content := config.RenderContent
	quality := config.Quality
	keepHtml := config.KeepHtml
	timeout := config.Timeout

	allocatorContext, _ := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:9222")

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	htmlFullPath := filepath.Join(wd, fmt.Sprintf("%d.html", time.Now().UnixNano()))

	if !keepHtml {
		defer func() {
			err := os.Remove(htmlFullPath)
			if err != nil {
				glgf.Errorf("Failed to delete the file(%s), err: %s\n", htmlFullPath, err)
			}
		}()
	}

	err = os.WriteFile(htmlFullPath, content, 0o644)
	if err != nil {
		return nil, err
	}

	if quality < 1 {
		quality = 1
	}

	var base64Data string
	executeJS := fmt.Sprintf(CanvasJs, config.ExtName, quality)

	err = chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("%s%s", FileProtocol, htmlFullPath)),
		chromedp.WaitVisible(EchartsInstanceDom, chromedp.ByQuery),
		chromedp.Evaluate(executeJS, &base64Data),
	)
	if err != nil {
		return nil, err
	}

	imgContent, err := base64.StdEncoding.DecodeString(strings.Split(base64Data, ",")[1])
	if err != nil {
		return nil, err
	}

	return imgContent, nil
}
