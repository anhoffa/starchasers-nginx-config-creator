package main

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata")

// TestGenerateNginxConfigGolden renders every testdata/*.json fixture and compares the
// result byte for byte with the matching .golden.conf file.
func TestGenerateNginxConfigGolden(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("testdata", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no fixtures found in testdata")
	}

	for _, fixture := range fixtures {
		name := strings.TrimSuffix(filepath.Base(fixture), ".json")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}

			got, err := generateNginxConfig(input)
			if err != nil {
				t.Fatalf("generateNginxConfig returned an error: %v", err)
			}

			goldenPath := filepath.Join("testdata", name+".golden.conf")
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatal(err)
			}
			if got != string(want) {
				t.Errorf("generated config differs from %s\n--- got ---\n%s\n--- want ---\n%s", goldenPath, got, want)
			}
		})
	}
}

func TestGenerateNginxConfigRejectsMalformedJSON(t *testing.T) {
	for _, input := range []string{"", "{", "not json", `{"domains": "nope"}`} {
		if _, err := generateNginxConfig([]byte(input)); err == nil {
			t.Errorf("expected an error for input %q", input)
		}
	}
}

func TestPrepareTemplateParamsDeduplicatesContainers(t *testing.T) {
	config := &Config{Domains: []Domain{
		{Name: "a.example.com", ContainerName: "app", Ip: "10.0.0.1"},
		{Name: "b.example.com", ContainerName: "app", Ip: "10.0.0.1"},
		{Name: "c.example.com", ContainerName: "other", Ip: "10.0.0.2"},
	}}

	params := prepareTemplateParams(config)

	wantContainers := []Container{
		{ContainerName: "app", Ip: "10.0.0.1"},
		{ContainerName: "other", Ip: "10.0.0.2"},
	}
	if !reflect.DeepEqual(params.Containers, wantContainers) {
		t.Errorf("containers = %+v, want %+v", params.Containers, wantContainers)
	}
	if !reflect.DeepEqual(params.Domains, config.Domains) {
		t.Errorf("domains = %+v, want %+v", params.Domains, config.Domains)
	}
}
