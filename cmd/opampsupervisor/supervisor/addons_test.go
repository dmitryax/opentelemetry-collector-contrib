// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package supervisor

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/supervisor/addon"
	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/supervisor/config"
)

// addonFile is one file to place in a test addon package.
type addonFile struct {
	name string
	body string
	exec bool
}

// buildAddonPackage builds an addon package as a gzipped tarball.
func buildAddonPackage(t *testing.T, files []addonFile) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	for _, f := range files {
		mode := int64(0o600)
		if f.exec {
			mode = 0o755
		}
		require.NoError(t, tw.WriteHeader(&tar.Header{
			Name:     f.name,
			Mode:     mode,
			Size:     int64(len(f.body)),
			Typeflag: tar.TypeReg,
		}))
		_, err := tw.Write([]byte(f.body))
		require.NoError(t, err)
	}

	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())
	return buf.Bytes()
}

// newAddonTestSupervisor returns a Supervisor with package storage and an addon
// registry wired up, as NewSupervisor does when accepts_packages is enabled.
func newAddonTestSupervisor(t *testing.T) *Supervisor {
	t.Helper()

	cfg := config.DefaultSupervisor()
	cfg.Storage.Directory = t.TempDir()
	cfg.Capabilities.AcceptsPackages = true

	s := &Supervisor{
		config:            cfg,
		telemetrySettings: newNopTelemetrySettings(),
		persistentState:   &persistentState{InstanceID: uuid.New()},
	}
	s.addons = addon.NewRegistry(zap.NewNop())

	var err error
	s.packages, err = newPackageManager(
		zap.NewNop(),
		filepath.Join(cfg.Storage.Directory, packagesDirName),
		s.addons,
	)
	require.NoError(t, err)

	return s
}

// installPackage drives the packageManager through the same call sequence the
// opamp-go PackagesSyncer uses when it installs a package, so the PoC is
// exercised the way the real client drives it.
func installPackage(t *testing.T, p *packageManager, name string, typ protobufs.PackageType, content []byte, version string) {
	t.Helper()

	state, err := p.PackageState(name)
	require.NoError(t, err)

	if !state.Exists {
		require.NoError(t, p.CreatePackage(name, typ))
	}

	sum := sha256.Sum256(content)

	// The syncer skips the download when the stored file already matches.
	existing, err := p.FileContentHash(name)
	require.NoError(t, err)
	if bytes.Equal(existing, sum[:]) {
		return
	}

	require.NoError(t, p.UpdateContent(
		context.Background(), name, "https://example.invalid/"+name+".tar.gz",
		bytes.NewReader(content), sum[:], nil,
	))

	require.NoError(t, p.SetPackageState(name, types.PackageState{
		Exists:  true,
		Type:    typ,
		Hash:    []byte("package-hash-" + version),
		Version: version,
	}))
}

// TestAddonLifecycle walks an executable addon through install, configure, run,
// reconfigure and delete. It is the end-to-end proof for the PoC: a package
// offered over OpAMP becomes a supervised process that reads configuration
// delivered as remote config, and stops when the package goes away.
func TestAddonLifecycle(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test addon is a shell script")
	}

	s := newAddonTestSupervisor(t)
	ctx := context.Background()

	// The addon writes the config it was given to a marker file and then idles,
	// which lets the test observe both that it ran and what config it saw.
	script := `#!/bin/sh
cat "$OPAMP_ADDON_CONFIG" > "$OPAMP_ADDON_DIR/observed-config"
while true; do sleep 0.05; done
`
	pkg := buildAddonPackage(t, []addonFile{
		{name: addon.ManifestName, body: "type: executable\nentrypoint: [\"./run.sh\"]\nconfig_file: addon.properties\n"},
		{name: "run.sh", body: script, exec: true},
	})

	installPackage(t, s.packages, "jmx-scraper", protobufs.PackageType_PackageType_Addon, pkg, "1.0.0")

	// The package is recorded and the addon is registered.
	state, err := s.packages.PackageState("jmx-scraper")
	require.NoError(t, err)
	assert.True(t, state.Exists)
	assert.Equal(t, "1.0.0", state.Version)
	assert.Equal(t, []string{"jmx-scraper"}, s.addons.Names())

	installed, ok := s.addons.Get("jmx-scraper")
	require.True(t, ok)
	assert.True(t, installed.IsExecutable())

	// Configuration arrives as a remote config entry whose role marks it as addon
	// config and whose key names the addon.
	require.NoError(t, s.applyAddonConfigs(ctx, addonRemoteConfig(map[string]string{
		"jmx-scraper": "otel.jmx.target.system=tomcat\n",
	})))

	assert.FileExists(t, installed.ConfigPath())
	body, err := os.ReadFile(installed.ConfigPath())
	require.NoError(t, err)
	assert.Equal(t, "otel.jmx.target.system=tomcat\n", string(body))

	// Starting the addon runs the entrypoint, which observes its config.
	require.NoError(t, s.addons.StartAll(ctx))
	assert.True(t, installed.IsRunning())
	assert.NotZero(t, installed.Pid())

	observed := filepath.Join(installed.Dir(), "observed-config")
	requireEventuallyFileContains(t, observed, "otel.jmx.target.system=tomcat")

	firstPid := installed.Pid()

	// A configuration change restarts a running addon so it rereads its config.
	require.NoError(t, os.Remove(observed))
	require.NoError(t, s.applyAddonConfigs(ctx, addonRemoteConfig(map[string]string{
		"jmx-scraper": "otel.jmx.target.system=kafka\n",
	})))

	assert.True(t, installed.IsRunning())
	assert.NotEqual(t, firstPid, installed.Pid(), "addon should have been restarted")
	requireEventuallyFileContains(t, observed, "otel.jmx.target.system=kafka")

	// Unchanged configuration must not restart it: a no-op offer should not
	// interrupt data collection.
	stablePid := installed.Pid()
	require.NoError(t, s.applyAddonConfigs(ctx, addonRemoteConfig(map[string]string{
		"jmx-scraper": "otel.jmx.target.system=kafka\n",
	})))
	assert.Equal(t, stablePid, installed.Pid(), "unchanged config should not restart the addon")

	// Deleting the package stops the process and removes its files, which is what
	// the syncer does for a package the server no longer offers.
	addonDir := installed.Dir()
	require.NoError(t, s.packages.DeletePackage("jmx-scraper"))

	assert.False(t, installed.IsRunning())
	assert.Empty(t, s.addons.Names())
	assert.NoDirExists(t, addonDir)

	names, err := s.packages.Packages()
	require.NoError(t, err)
	assert.Empty(t, names)
}

// TestAddonUpgradeReplacesProcess checks that offering a new version of an
// installed addon stops the old process and starts the new content, which is
// the case the staged-upgrade discussions turn on.
func TestAddonUpgradeReplacesProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test addon is a shell script")
	}

	s := newAddonTestSupervisor(t)
	ctx := context.Background()

	makePkg := func(version string) []byte {
		return buildAddonPackage(t, []addonFile{
			{name: addon.ManifestName, body: "type: executable\nentrypoint: [\"./run.sh\"]\n"},
			{name: "run.sh", body: "#!/bin/sh\necho " + version + " > \"$OPAMP_ADDON_DIR/version\"\nwhile true; do sleep 0.05; done\n", exec: true},
		})
	}

	installPackage(t, s.packages, "scraper", protobufs.PackageType_PackageType_Addon, makePkg("v1"), "1.0.0")
	require.NoError(t, s.addons.StartAll(ctx))

	first, ok := s.addons.Get("scraper")
	require.True(t, ok)
	requireEventuallyFileContains(t, filepath.Join(first.Dir(), "version"), "v1")
	firstPid := first.Pid()
	require.NotZero(t, firstPid)

	// Upgrade.
	installPackage(t, s.packages, "scraper", protobufs.PackageType_PackageType_Addon, makePkg("v2"), "2.0.0")

	// The old process must be gone, and the registry must now hold the new addon.
	assert.False(t, first.IsRunning(), "previous addon process should be stopped on upgrade")

	second, ok := s.addons.Get("scraper")
	require.True(t, ok)
	assert.NotSame(t, first, second, "registry should hold the upgraded addon")

	require.NoError(t, s.addons.StartAll(ctx))
	requireEventuallyFileContains(t, filepath.Join(second.Dir(), "version"), "v2")
	assert.NotEqual(t, firstPid, second.Pid())

	state, err := s.packages.PackageState("scraper")
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", state.Version)

	require.NoError(t, s.addons.StopAll(ctx))
}

// TestAddonContentHashMismatchRejected checks that content failing its offered
// hash is never extracted or run.
func TestAddonContentHashMismatchRejected(t *testing.T) {
	s := newAddonTestSupervisor(t)

	pkg := buildAddonPackage(t, []addonFile{
		{name: addon.ManifestName, body: "type: file\n"},
	})

	require.NoError(t, s.packages.CreatePackage("bad", protobufs.PackageType_PackageType_Addon))

	err := s.packages.UpdateContent(
		context.Background(), "bad", "https://example.invalid/bad.tar.gz",
		bytes.NewReader(pkg), []byte("not-the-right-hash"), nil,
	)
	require.ErrorContains(t, err, "content hash mismatch")

	assert.Empty(t, s.addons.Names(), "a package failing verification must not be registered")
	assert.NoDirExists(t, s.packages.AddonDir("bad"))
}

// TestAddonMalformedPackageLeavesPreviousInstalled checks that a bad offer
// cannot take down a working addon.
func TestAddonMalformedPackageLeavesPreviousInstalled(t *testing.T) {
	s := newAddonTestSupervisor(t)

	good := buildAddonPackage(t, []addonFile{
		{name: addon.ManifestName, body: "type: file\n"},
		{name: "data.txt", body: "original"},
	})
	installPackage(t, s.packages, "files", protobufs.PackageType_PackageType_Addon, good, "1.0.0")
	require.Len(t, s.addons.Names(), 1)

	// A package with no manifest is not a usable addon.
	bad := buildAddonPackage(t, []addonFile{{name: "data.txt", body: "replaced"}})
	sum := sha256.Sum256(bad)
	err := s.packages.UpdateContent(
		context.Background(), "files", "https://example.invalid/files.tar.gz",
		bytes.NewReader(bad), sum[:], nil,
	)
	require.ErrorContains(t, err, addon.ManifestName)

	// The previously installed content must still be there.
	body, readErr := os.ReadFile(filepath.Join(s.packages.AddonDir("files"), "data.txt"))
	require.NoError(t, readErr)
	assert.Equal(t, "original", string(body))
}

// TestExecutableAddonMissingRuntimeFailsDeterministically covers the objection
// that moving a Java addon to the Supervisor does not make a JVM appear: a host
// without the runtime must fail with a clear, non-retryable reason rather than
// looking broken.
func TestExecutableAddonMissingRuntimeFailsDeterministically(t *testing.T) {
	s := newAddonTestSupervisor(t)
	ctx := context.Background()

	pkg := buildAddonPackage(t, []addonFile{
		{name: addon.ManifestName, body: "type: executable\nentrypoint: [definitely-not-installed-runtime, -jar, scraper.jar]\n"},
	})
	installPackage(t, s.packages, "needs-runtime", protobufs.PackageType_PackageType_Addon, pkg, "1.0.0")

	installed, ok := s.addons.Get("needs-runtime")
	require.True(t, ok)

	err := installed.Start(ctx)
	require.ErrorIs(t, err, addon.ErrRuntimeNotFound)
	require.ErrorContains(t, err, "definitely-not-installed-runtime")
	require.ErrorContains(t, err, "must be installed on the host")
	assert.False(t, installed.IsRunning())
}

// TestFileAddonIsNotRun checks that a file addon occupies a directory and
// nothing more. This is the geoipprocessor case: the Collector reads the file,
// the Supervisor runs nothing.
func TestFileAddonIsNotRun(t *testing.T) {
	s := newAddonTestSupervisor(t)
	ctx := context.Background()

	pkg := buildAddonPackage(t, []addonFile{
		{name: addon.ManifestName, body: "type: file\n"},
		{name: "GeoLite2-City.mmdb", body: "pretend-database"},
	})
	installPackage(t, s.packages, "geoip-db", protobufs.PackageType_PackageType_Addon, pkg, "2026.10.28")

	installed, ok := s.addons.Get("geoip-db")
	require.True(t, ok)
	assert.False(t, installed.IsExecutable())

	require.NoError(t, s.addons.StartAll(ctx))
	assert.False(t, installed.IsRunning(), "a file addon must not start a process")
	assert.Zero(t, installed.Pid())

	// The extracted path is what a Collector config would point at.
	assert.FileExists(t, filepath.Join(installed.Dir(), "GeoLite2-City.mmdb"))
}

// TestAddonStateSurvivesRestart checks that a restarted Supervisor adopts the
// addons it installed previously instead of waiting to be re-offered them.
func TestAddonStateSurvivesRestart(t *testing.T) {
	s := newAddonTestSupervisor(t)
	ctx := context.Background()

	pkg := buildAddonPackage(t, []addonFile{
		{name: addon.ManifestName, body: "type: file\n"},
		{name: "data.txt", body: "persisted"},
	})
	installPackage(t, s.packages, "files", protobufs.PackageType_PackageType_Addon, pkg, "1.0.0")
	require.NoError(t, s.packages.SetAllPackagesHash([]byte("all-hash")))

	// A new packageManager over the same directory is what a restart looks like.
	restartedRegistry := addon.NewRegistry(zap.NewNop())
	restarted, err := newPackageManager(
		zap.NewNop(),
		filepath.Join(s.config.Storage.Directory, packagesDirName),
		restartedRegistry,
	)
	require.NoError(t, err)

	hash, err := restarted.AllPackagesHash()
	require.NoError(t, err)
	assert.Equal(t, []byte("all-hash"), hash)

	state, err := restarted.PackageState("files")
	require.NoError(t, err)
	assert.True(t, state.Exists)
	assert.Equal(t, "1.0.0", state.Version)

	// The stored content hash must still match, so a restarted Supervisor does
	// not re-download what it already has.
	sum := sha256.Sum256(pkg)
	storedHash, err := restarted.FileContentHash("files")
	require.NoError(t, err)
	assert.Equal(t, sum[:], storedHash)

	require.NoError(t, restarted.loadInstalledAddons(ctx))
	assert.Equal(t, []string{"files"}, restartedRegistry.Names())
}

// TestLastReportedStatusesRoundTrip checks the statuses a restarted Supervisor
// reports are the ones it saved.
func TestLastReportedStatusesRoundTrip(t *testing.T) {
	s := newAddonTestSupervisor(t)

	// Nothing saved yet must not be an error: it is the first run.
	statuses, err := s.packages.LastReportedStatuses()
	require.NoError(t, err)
	assert.Empty(t, statuses.GetPackages())

	require.NoError(t, s.packages.SetLastReportedStatuses(&protobufs.PackageStatuses{
		ServerProvidedAllPackagesHash: []byte("offer-hash"),
		Packages: map[string]*protobufs.PackageStatus{
			"jmx-scraper": {
				Name:            "jmx-scraper",
				AgentHasVersion: "1.0.0",
				Status:          protobufs.PackageStatusEnum_PackageStatusEnum_Installed,
			},
		},
	}))

	statuses, err = s.packages.LastReportedStatuses()
	require.NoError(t, err)
	assert.Equal(t, []byte("offer-hash"), statuses.GetServerProvidedAllPackagesHash())
	require.Contains(t, statuses.GetPackages(), "jmx-scraper")
	assert.Equal(t,
		protobufs.PackageStatusEnum_PackageStatusEnum_Installed,
		statuses.GetPackages()["jmx-scraper"].GetStatus())
}

// TestAddonConfigIsNotMergedIntoCollectorConfig is the core of the config
// routing question: an addon config entry must not reach the Collector's
// configuration, which it would corrupt.
func TestAddonConfigIsNotMergedIntoCollectorConfig(t *testing.T) {
	s := newAddonTestSupervisor(t)
	s.config.Capabilities.AcceptsRemoteConfig = true

	collectorConfig := "receivers:\n  nop:\n"
	addonConfig := "otel.jmx.target.system=tomcat\n"

	remoteConfig := &protobufs.AgentRemoteConfig{
		Config: &protobufs.AgentConfigMap{
			ConfigMap: map[string]*protobufs.AgentConfigObject{
				"collector.yaml": {Body: []byte(collectorConfig), ContentType: "text/yaml"},
				"jmx-scraper": {
					Body:        []byte(addonConfig),
					ContentType: "text/plain",
					Role:        addonConfigRole,
				},
			},
		},
	}

	var composed []byte
	for _, compose := range s.createRemoteConfigComposers(remoteConfig) {
		composed = append(composed, compose()...)
	}

	assert.Contains(t, string(composed), "nop", "collector config should be merged")
	assert.NotContains(t, string(composed), "otel.jmx.target.system",
		"addon config must not be merged into the collector config")
}

// addonRemoteConfig builds a remote config carrying only addon configuration.
func addonRemoteConfig(configs map[string]string) *protobufs.AgentRemoteConfig {
	configMap := map[string]*protobufs.AgentConfigObject{}
	for name, body := range configs {
		configMap[name] = &protobufs.AgentConfigObject{
			Body:        []byte(body),
			ContentType: "text/plain",
			Role:        addonConfigRole,
		}
	}
	return &protobufs.AgentRemoteConfig{
		Config: &protobufs.AgentConfigMap{ConfigMap: configMap},
	}
}

// requireEventuallyFileContains waits for path to exist and contain want.
func requireEventuallyFileContains(t *testing.T, path, want string) {
	t.Helper()
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		body, err := os.ReadFile(path)
		if !assert.NoError(c, err) {
			return
		}
		assert.Contains(c, string(body), want)
	}, 10*time.Second, 25*time.Millisecond, "waiting for %s to contain %q", path, want)
}
