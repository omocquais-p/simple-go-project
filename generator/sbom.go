package generator

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/v2/sbom"
)

const (
	CycloneDXFormat = "application/vnd.cyclonedx+json"
	SPDXFormat      = "application/spdx+json"
	SyftFormat      = "application/vnd.syft+json"
)

func GenerateSBOM(path string) []error {
	log.Printf("Generate SBOM from this path: %v", path)

	isDirectory, _ := isDirectory(path)

	if isDirectory {
		return []error{errors.New("path is a directory, it must be the path to a binary file")}
	}

	generatedSbom, _ := sbom.Generate(path)

	formatter, _ := generatedSbom.InFormats(CycloneDXFormat, SPDXFormat, SyftFormat)
	parentDirectory := filepath.Dir(path)

	var errs []error

	filename := GenerateFilename(path)

	errs = append(errs, writeSbom(filepath.Join(parentDirectory, filename+".sbom.cdx.json"), formatter.Formats()[0].Content))
	errs = append(errs, writeSbom(filepath.Join(parentDirectory, filename+".sbom.spdx.json"), formatter.Formats()[1].Content))
	errs = append(errs, writeSbom(filepath.Join(parentDirectory, filename+".sbom.syft.json"), formatter.Formats()[2].Content))

	return errs
}

func writeSbom(sbomFullPath string, content io.Reader) error {
	syft := bytes.NewBuffer(nil)
	_, _ = io.Copy(syft, content)
	return ioutil.WriteFile(sbomFullPath, syft.Bytes(), 0644)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func GenerateFilename(absolutePath string) string {
	return filepath.Base(strings.TrimSuffix(absolutePath, filepath.Ext(absolutePath)))
}
