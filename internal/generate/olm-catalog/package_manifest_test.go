// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package olmcatalog

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/operator-framework/operator-registry/pkg/registry"
	gen "github.com/operator-framework/operator-sdk/internal/generate/gen"

	"github.com/stretchr/testify/assert"
)

// newTestPackageManifestGenerator returns a package manifest Generator populated with test values.
func newTestPackageManifestGenerator() gen.Generator {
	cfg := gen.Config{
		OperatorName: testProjectName,
	}
	g := NewPackageManifest(cfg, csvVersion, "stable", true)

	// TODO: The deploy dir input should be the way to point the generator to test data
	// E.g: Inputs: map[string]string{DeployDirKey: testGoDataDir}

	// Override the generator to point to the package manifest in testGoDataDir
	testPkgManifestDir := filepath.Join(testGoDataDir, OLMCatalogDir, testProjectName)
	pkgGen := g.(pkgGenerator)
	pkgGen.existingPkgManifestDir = testPkgManifestDir
	return pkgGen
}

func TestNewPkgManifestWithInputsAndOutput(t *testing.T) {
	nonStandardTestDataDir := filepath.Join(testDataDir, "non-standard-layout")
	for _, cleanupFunc := range setupTestEnvWithCleanup(t, nonStandardTestDataDir) {
		defer cleanupFunc()
	}

	// Temporary output dir for generating catalog bundle
	outputDir, err := ioutil.TempDir("", t.Name()+"-output-catalog")
	if err != nil {
		log.Fatal(err)
	}
	// Clean up output catalog dir
	defer func() {
		if err := os.RemoveAll(outputDir); err != nil && !os.IsNotExist(err) {
			// Not a test failure since files in /tmp will eventually get deleted
			t.Logf("Failed to remove tmp generated catalog directory (%s): %v", outputDir, err)
		}
	}()

	cfg := gen.Config{
		OperatorName: testProjectName,
		Inputs: map[string]string{
			DeployDirKey: "config",
			APIsDirKey:   "api",
		},
		OutputDir: outputDir,
	}

	g := NewPackageManifest(cfg, csvVersion, "stable", true)
	if err := g.Generate(); err != nil {
		t.Fatalf("Failed to execute package manifest generator: %v", err)
	}

	pkgManFileName := getPkgFileName(testProjectName)

	// Read expected Package Manifest
	expCatalogDir := filepath.Join("expected-catalog", OLMCatalogChildDir)
	pkgManExpBytes, err := ioutil.ReadFile(filepath.Join(expCatalogDir, testProjectName, pkgManFileName))
	if err != nil {
		t.Fatalf("Failed to read expected package manifest file: %v", err)
	}
	pkgManExp := string(pkgManExpBytes)

	// Read generated Package Manifest from OutputDir/olm-catalog
	outputCatalogDir := filepath.Join(cfg.OutputDir, OLMCatalogChildDir)
	pkgManOutputBytes, err := ioutil.ReadFile(filepath.Join(outputCatalogDir, testProjectName, pkgManFileName))
	if err != nil {
		t.Fatalf("Failed to read output package manifest file: %v", err)
	}
	pkgManOutput := string(pkgManOutputBytes)

	assert.Equal(t, pkgManExp, pkgManOutput)

}

func TestGeneratePackageManifest(t *testing.T) {
	g := newTestPackageManifestGenerator()
	fileMap, err := g.(pkgGenerator).generate()
	if err != nil {
		t.Fatalf("Failed to execute package manifest generator: %v", err)
	}

	if b, ok := fileMap[g.(pkgGenerator).fileName]; !ok {
		t.Error("Failed to generate package manifest")
	} else {
		assert.Equal(t, packageManifestExp, string(b))
	}
}

func TestValidatePackageManifest(t *testing.T) {
	g := newTestPackageManifestGenerator()

	// pkg is a basic, valid package manifest.
	pkg, err := g.(pkgGenerator).buildPackageManifest()
	if err != nil {
		t.Fatalf("Failed to execute package manifest generator: %v", err)
	}

	g.(pkgGenerator).setChannels(&pkg)
	sortChannelsByName(&pkg)

	// invalid mock data, pkg with empty channel
	invalidPkgWithEmptyChannels := pkg
	invalidPkgWithEmptyChannels.Channels = []registry.PackageChannel{}

	type args struct {
		pkg *registry.PackageManifest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should work successfully with a valid pkg",
			wantErr: false,
			args: args{
				pkg: &pkg,
			},
		},
		{
			name:    "Should return error when the pkg is not informed",
			wantErr: true,
		},
		{
			name:    "Should return error when the pkg is invalid",
			wantErr: true,
			args: args{
				pkg: &invalidPkgWithEmptyChannels,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePackageManifest(tt.args.pkg); (err != nil) != tt.wantErr {
				t.Errorf("Failed to check package manifest validate: error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPackageManifest(t *testing.T) {
	type args struct {
		operatorName string
		channelName  string
		version      string
	}
	tests := []struct {
		name string
		args args
		want registry.PackageManifest
	}{
		{
			name: "Should return a valid registry.PackageManifest",
			want: registry.PackageManifest{
				PackageName: "memcached-operator",
				Channels: []registry.PackageChannel{
					registry.PackageChannel{
						Name:           "stable",
						CurrentCSVName: "memcached-operator.v0.0.3",
					},
				},
				DefaultChannelName: "stable",
			},
			args: args{
				operatorName: testProjectName,
				channelName:  "stable",
				version:      csvVersion,
			},
		},
		{
			name: "Should return a valid registry.PackageManifest with channel == alpha when it is not informed",
			want: registry.PackageManifest{
				PackageName: "memcached-operator",
				Channels: []registry.PackageChannel{
					registry.PackageChannel{
						Name:           "alpha",
						CurrentCSVName: "memcached-operator.v0.0.3",
					},
				},
				DefaultChannelName: "alpha",
			},
			args: args{
				operatorName: testProjectName,
				version:      csvVersion,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newPackageManifest(tt.args.operatorName, tt.args.channelName, tt.args.version)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPackageManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}

const packageManifestExp = `channels:
- currentCSV: memcached-operator.v0.0.2
  name: alpha
- currentCSV: memcached-operator.v0.0.3
  name: stable
defaultChannel: stable
packageName: memcached-operator
`
