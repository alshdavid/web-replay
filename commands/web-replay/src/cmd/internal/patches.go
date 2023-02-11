package internal_serve

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type Patch struct {
	Name          string
	Match         PatchMatch          `yaml:"match,omitempty"`
	ReplaceText   *PatchReplaceText   `yaml:"replaceText,omitempty"`
	StreamLatency *PatchStreamLatency `yaml:"streamLatency,omitempty"`
	AddLatency    *int64              `yaml:"addLatency,omitempty"`
}

type PatchMatch struct {
	MimeTypes    []string `yaml:"mimeTypes,omitempty"`
	Origins      []string `yaml:"origins,omitempty"`
	UrlPaths     []string `yaml:"urlPaths,omitempty"`
	RequestTypes []string `yaml:"requestTypes,omitempty"`
}

type PatchReplaceText struct {
	FindString  string `yaml:"findString,omitempty"`
	ReplaceWith string `yaml:"replaceWith,omitempty"`
}

type PatchStreamLatency struct {
	FindString string `yaml:"findString,omitempty"`
	Latency    *int64 `yaml:"latency,omitempty"`
}

type configFile struct {
	Match         PatchMatch          `yaml:"match,omitempty"`
	ReplaceText   *PatchReplaceText   `yaml:"replaceText,omitempty"`
	StreamLatency *PatchStreamLatency `yaml:"streamLatency,omitempty"`
	AddLatency    *int64              `yaml:"addLatency,omitempty"`
	Patches       []Patch             `yaml:"patches,omitempty"`
}

func LoadPatchesFromDir(pathToDir string) ([]Patch, []fs.DirEntry, error) {
	dirs, _ := os.ReadDir(pathToDir)
	patches := []Patch{}

	for _, file := range dirs {
		fp := filepath.Join(pathToDir, file.Name())
		yamlFile, _ := os.ReadFile(fp)
		config := configFile{}
		err := yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return nil, nil, err
		}
		if config.Patches != nil {
			for i, patch := range config.Patches {
				patch.Name = fmt.Sprintf("%s[%d]", file.Name(), i)
				patches = append(patches, patch)
			}
		} else {
			patch := Patch{}
			patch.Name = file.Name()
			patch.Match = config.Match
			patch.ReplaceText = config.ReplaceText
			patch.StreamLatency = config.StreamLatency
			patch.AddLatency = config.AddLatency
			patches = append(patches, patch)
		}
	}

	return patches, dirs, nil
}
