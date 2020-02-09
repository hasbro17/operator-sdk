package olmcatalog

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	gen "github.com/operator-framework/operator-sdk/internal/generate/gen"
	"github.com/operator-framework/operator-sdk/internal/scaffold"

	"github.com/stretchr/testify/assert"
)

var testHelmDataDir = filepath.Join(testDataDir, "helm")

func TestHelmCSVNew(t *testing.T) {
	for _, cleanupFunc := range setupTestEnvWithCleanup(t, testHelmDataDir) {
		defer cleanupFunc()
	}

	cfg := gen.Config{
		OperatorName: testProjectName,
		// Only include one CRD and all deploy manifests in this run.
		Filters: gen.MakeFilters(
			filepath.Join(scaffold.CRDsDir, getTestCRDFile(testGroup, testKind1)),
			filepath.Join(scaffold.CRDsDir, getTestCRFile(testGroup, testVersion1, testKind1)),
			filepath.Join(scaffold.DeployDir, "operator.yaml"),
			filepath.Join(scaffold.DeployDir, "role_binding.yaml"),
			filepath.Join(scaffold.DeployDir, "role.yaml"),
			filepath.Join(scaffold.DeployDir, "service_account.yaml"),
		),
	}
	// Create new CSV from scratch and compare against an existing CSV.
	g := NewCSV(cfg, notExistVersion, "")
	fileMap, err := g.(csvGenerator).generate()
	if err != nil {
		t.Fatalf("Failed to execute CSV generator: %v", err)
	}
	// Get an existing CSV created from scratch.
	csvExpFile := getCSVFileName(testProjectName, notExistVersion)
	csvExpBytes, err := ioutil.ReadFile(filepath.Join(OLMCatalogDir, testProjectName, scratchBundleDir, csvExpFile))
	if err != nil {
		t.Fatalf("Failed to read expected CSV file: %v", err)
	}
	// Compare scratch CSV to existing CSV.
	csvExp := string(csvExpBytes)
	genCSVFile := getCSVFileName(testProjectName, notExistVersion)
	if b, ok := fileMap[genCSVFile]; !ok {
		t.Errorf("Failed to generate CSV for version %s: file %s not found", notExistVersion, genCSVFile)
	} else {
		assert.Equal(t, csvExp, string(b))
	}

	// Include all CRDs.
	cfg.Filters = gen.MakeFilters(scaffold.DeployDir)
	g = NewCSV(cfg, csvVersion, "")
	fileMap, err = g.(csvGenerator).generate()
	if err != nil {
		t.Fatalf("Failed to execute CSV generator: %v", err)
	}
	csvExpFile = getCSVFileName(testProjectName, csvVersion)
	csvExpBytes, err = ioutil.ReadFile(filepath.Join(OLMCatalogDir, testProjectName, csvVersion, csvExpFile))
	if err != nil {
		t.Fatalf("Failed to read expected CSV file: %v", err)
	}
	csvExp = string(csvExpBytes)
	if b, ok := fileMap[csvExpFile]; !ok {
		t.Errorf("Failed to generate CSV for version %s", csvVersion)
	} else {
		assert.Equal(t, csvExp, string(b))
	}
}
