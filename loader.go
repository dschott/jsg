package main

import (
	gourl "net/url"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Loader struct {
	fileLoader    jsonschema.FileLoader
	retrievalURIs map[string]string
}

func (l *Loader) AddRetrievalPath(baselURI string, path string) error {
	if l.retrievalURIs == nil {
		l.retrievalURIs = make(map[string]string)
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	l.retrievalURIs[strings.TrimSuffix(baselURI, "/")] = "file://" + path
	return nil
}

func (l *Loader) Load(url string) (any, error) {
	doc, lerr := l.fileLoader.Load(url)
	if lerr == nil || l.retrievalURIs == nil || strings.HasPrefix(url, "file://") {
		return doc, lerr
	}

	for baseURI, retrievalURI := range l.retrievalURIs {
		if strings.HasPrefix(url, baseURI) {
			fileURL, err := gourl.JoinPath(retrievalURI, strings.TrimPrefix(url, baseURI))
			if err != nil {
				return nil, err
			}

			doc, err = l.fileLoader.Load(fileURL)
			if err != nil {
				return nil, err
			}

			return doc, nil
		}
	}
	return nil, lerr
}
