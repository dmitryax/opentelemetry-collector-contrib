# Addon (sub-)package support in the OpAMP Supervisor

Status: proof of concept. This document and the code under
`cmd/opampsupervisor/supervisor/addon` exist to back an upstream proposal with a
working reference implementation. It is not proposed for merge as-is.

## Why

The Collector cannot deliver or supervise artifacts that are not Go code
compiled into its binary. Two current, unserved needs:

**Helper processes.** The JMX Metric Scraper
(`io.opentelemetry.contrib:opentelemetry-jmx-scraper`) runs beside a Collector
and exports OTLP to it. `receiver/jmxreceiver` was removed in **v0.157.0** (code
removal PR #49478, merged 2026-07-10), and with it the only managed deployment
path. Issue #34825 settled that JMX functionality moves to "a direct
installation of the jmx scraper as a standalone replacement". Nobody specified
how direct installation happens across a fleet.

**Data files.** `processor/geoipprocessor` requires a MaxMind `.mmdb` through
`database_path` (`internal/provider/maxmindprovider/config.go`), which is
required and errors if empty. MaxMind publishes updates weekly. Procurement and
refresh are entirely manual.

OpAMP already specifies everything needed to deliver both: `PackageType_Addon`,
`download_url`, `content_hash`, `version` and `signature`. None of it is
implemented anywhere, in any OpAMP implementation.

**No spec change is needed** for anything in this proposal.

## Design: typed addons

An addon's behaviour is decided by its **type**, and each type is implemented
in Supervisor code. The type is carried in the package name:

```
<type>:<name>          jmx:tomcat   file:geoip   splunk_ta:unix
```

The same string is the key of the addon's configuration in the remote config
map, so one name identifies an addon in `PackagesAvailable`, `PackageStatuses`
and `AgentConfigMap`.

Built-in types:

| Type | Package content | What the Supervisor does |
|---|---|---|
| `file` | a `.tar.gz`/`.tgz`, or any single file | extracts it, or stores it under the download URL's file name; runs nothing |
| `jmx` | the `opentelemetry-jmx-scraper` jar, as published | checks for a JVM, runs `java -jar <jar> -config <properties>`, restarts it on config change |

Other types are added by registering a factory. A distribution does that
through a constructor option:

```go
sup, err := supervisor.NewSupervisor(ctx, logger, cfg,
	supervisor.WithAddonFactories(splunkta.NewFactory()))
```

`supervisor.NewSvcHandler(opts...)` takes the same options on the Windows
service path.

### Why types, not a package manifest

An earlier version of this PoC had each package carry a manifest declaring its
type and entrypoint. Typed addons are better on two counts.

**The server cannot make the Supervisor run arbitrary commands.** With a
manifest, a package declares what to execute, so whoever controls the OpAMP
server or the download location controls what runs on the host. With types, the
server can only pick among the types the Supervisor was built with, and each
type owns its launch logic: the `jmx` type will only ever run the scraper jar
under a JVM. This is the property #45740 was asking for, and checksums alone did
not provide.

**Upstream artifacts are usable as-is.** The scraper's release asset is a raw
jar at a stable URL with a published checksum. The `jmx` type accepts it
directly, so a server offers that URL and its SHA-256 and nothing is repackaged.
A manifest design requires wrapping the jar in an archive per release.

There is deliberately no generic `exec` type. It would be the manifest design
again, and the arbitrary-execution objection with it.

### Naming

The shape follows confmap URIs: everything before the first `:` is the type and
the rest is opaque, so names may contain further `:` or `/`.

The grammar of the type is the Collector's **component type** grammar,
`^[a-zA-Z][0-9a-zA-Z_]{0,62}$`, not confmap's scheme grammar. Confmap schemes
follow RFC 3986 and do not allow `_` (`[A-Za-z][A-Za-z0-9+.-]+`), which rules out
`splunk_ta`. Addon types read like component types, so they use that grammar.

Top-level packages are not subject to the naming rule.

### Framework

```go
type Factory interface {
	Type() Type
	// Install lays out verified package content in an empty staging directory.
	// An error is reported to the server as INSTALL_FAILED.
	Install(ctx context.Context, set InstallSettings) error
	// Create returns the runtime object for an installed addon. Starts nothing.
	Create(ctx context.Context, set Settings) (Addon, error)
}

type Addon interface {
	Configure(ctx context.Context, body []byte) (Effect, error) // nil body: no config offered
	Start(ctx context.Context) (Effect, error)                  // idempotent
	Shutdown(ctx context.Context) error                         // files stay
}

type Effect struct {
	RestartAgent bool // the Collector must restart to observe this change
}
```

The Supervisor **reconciles** rather than replaying events. After every package
sync, every remote config change, and once at startup, it calls `Configure` on
each addon whose configuration changed and then `Start` on all of them. That
stays correct after a restart, where nothing is replayed, and means a type never
has to reason about message ordering. The manager tracks the last applied config
per addon, so an unchanged offer is a no-op and a type need not diff.

`addon/process` is an exported process supervisor that types needing a child
process can reuse.

### Configuration and ordering

Addon configuration is a remote config entry whose `AgentConfigObject.role` is
`addon-config` and whose key is the addon's name. The role, not the key's shape,
selects it: config map keys are free-form, and a Collector config entry keyed
`file:base.yaml` must not be taken for addon configuration. `role` was added for
this purpose (opamp-spec #350, resolving #184).

`remote_config` and `packages_available` are each complete when present, but
the server sends each only when its own hash is stale, so they are not
guaranteed to arrive together. Reconciliation makes arrival order irrelevant:
configuration for an addon that is not installed yet is held, and applied by the
reconciliation that follows its package sync.

**One Collector restart per message.** Types the Collector itself reads, such as
a Splunk TA, return `Effect{RestartAgent: true}`. When a message changes the
Collector config and also offers packages, the Supervisor holds that restart
until the package sync completes, then signals once. The Collector restarts a
single time and observes both changes, instead of restarting with the old addons
and again with the new. At startup, addons are reconciled before the Collector
first starts, so their restart requests are moot.

Servers should still send coupled changes, such as a new TA and its state, in
the same `ServerToAgent`. Across messages the agent converges, but a restart can
happen in between.

### Enable/disable without a spec change

A type can carry state in its configuration. The `splunk_ta` test type takes
`enabled` or `disabled` from its config entry, writes `local/app.conf`, and asks
for a Collector restart when the state changes on a loaded TA
(`TestAddonStateChangeRestartsCollectorOnce`). No proto change is involved, and
the content hash is unchanged, so nothing is downloaded.

What this does not give is per-addon reporting. `RemoteConfigStatus` is one
status per `config_hash`, so a failed state change fails the whole remote config.
Reporting each addon's actual state in `EffectiveConfig`, under the same role and
key, would close that without a spec change. Not implemented here.

## Prior art in this repo

`jmxreceiver` already implemented package delivery and process supervision, in
the wrong layer. At v0.156.0, its last release:

- `internal/subprocess/subprocess.go`: a 322-line process supervisor with
  `RestartOnError`, `RestartDelay`, a running/restarting state machine and
  stdout capture.
- `config.go`: `jar_path`, which the operator had to populate by hand, plus
  SHA-256 verification in `hashFile()` and `validateJar()`.
- `supported_jars.go`: **83 pinned jar hashes compiled into the Collector
  binary**, across `jmxMetricsGathererVersions` and `jmxScraperVersions`, plus a
  build-time `MetricsGathererHash` override.

The structural failure was that the allowlist was compile-time, so every new jar
release needed a Collector release. PR #41863 ("fix sha256sum for
opentelemetry-jmx-scraper-1.48.0-alpha.jar") is one such fix; the security
hardening PRs #9721 and #9985 and the orphaned-process bug #23051 are all
symptoms of solving this inside a receiver.

| jmxreceiver | Typed addon |
|---|---|
| `hashFile()` + 83-entry allowlist | `DownloadableFile.content_hash`, supplied per offer |
| `supportedJar.version` | `PackageAvailable.version` |
| `jar_path`, populated manually | `DownloadableFile.download_url` |
| receiver builds the `java` command line | the `jmx` type builds it; the server cannot change it |
| `subprocess.Config` | `addon/process` |

## The objection, and the answer

Issue **#45740** (rogercoll, 2026-01-30) deprecated the receiver and states:

> Requiring the Collector to run Java is problematic, especially since the
> Collector is commonly deployed as a Docker image and none of the official
> images ship a Java runtime [...] More broadly, the Collector should not run
> subprocesses at all; the jmxreceiver is an exception that sets an undesirable
> precedent, even with its checksum-based restrictions.

1. **The Supervisor is not the Collector.** #45740 is scoped to the Collector
   process: no JRE in its image, and subprocess lifecycle entangled with
   receiver start and shutdown (#23051). The Supervisor is already a process
   manager; a second managed child is a quantitative change, not a new
   precedent.
2. **Typed addons are a narrower precedent than checksums.** The objection
   dismisses checksum restrictions. Typed addons do not rely on them: the set of
   things that can run is fixed by the Supervisor build, per type.
3. **The runtime problem does not go away, and this proposal concedes it.** A
   Supervisor-delivered jar fails on a host without Java exactly as a
   Collector-delivered one does, and identically in containers.

On (3), the scope boundary is explicit: **executable addon types target host and
VM deployments, not containers.** The `jmx` type checks for a JVM (`JAVA_HOME`,
then `PATH`) in `Install`, so a host without Java reports `INSTALL_FAILED` with
a reason the server sees, rather than reporting a successful install of an addon
that never runs (`TestJMXAddonWithoutJavaReportsInstallFailure`). File addons
have none of these problems.

## What the PoC implements

```
<storage>/packages/
  state.json                     package records and all_packages_hash
  last_reported_statuses.bin     PackageStatuses, so a restart reports what it has
  <escaped-package-name>/        e.g. jmx%3Atomcat
    content                      the package file exactly as downloaded
    files/                       the installed addon
```

- `packages.go`: a filesystem-backed `types.PackagesStateProvider`. Rejects
  unsupported addon types in `CreatePackage`, which the syncer calls before
  downloading, so they fail without being fetched. Verifies the content hash,
  installs into a staging directory through the type's factory, shuts the old
  addon down, then swaps.
- `addon/`: `ID`, `Factory`, `Addon`, `Effect`, `Manager`, and tar.gz
  extraction with escape protection.
- `addon/fileaddon`, `addon/jmxaddon`: the built-in types.
- `addon/process`: the shared process supervisor.
- `options.go`: `WithAddonFactories`.
- `addons.go`: role-based config routing, reconciliation, and the single
  coordinated Collector restart.

## Findings

**1. Retracted.** An earlier version claimed `accepts_packages: true` was broken.
It is not: `supervisor.go` rejects it during `Start` with a link to #47272,
covered by `TestSupervisorPackageCapabilitiesReturnError`. That guard is step 1
of #47272 and removing it is step 7. This PoC removes it, along with that test
and its fixture, so the addon path is reachable; upstream it belongs in step 7.

**2. The `role` field is specified but ignored.** `createRemoteConfigComposers`
merges every `config_map` entry into the Collector's YAML. A server that sends a
role-tagged entry to a current Supervisor corrupts the Collector config. The
server needs to know the agent honours `role`: a capability bit or a documented
minimum version.

**3. The existing `archive` package cannot extract an addon.** It pulls one
named binary out of a package. `addon.ExtractTarGz` extracts a whole tree and
rejects absolute paths, upward traversal, symlinks, hardlinks and devices
outright, with bounds on entry count and total size.

**4. `Commander` cannot supervise an addon.** It signals `os.Interrupt` and does
not use a process group, which is right for the Collector and wrong for a
runtime or a wrapper script, which ignores SIGINT while waiting on its child.
Measured against `Commander` alone:

| Delay between Start and Stop | `Commander.Stop` latency |
|---|---|
| 0 ms | 408 µs |
| 50 ms | 556 µs |
| 200 ms | **10.005 s** (the full grace period, then kill) |

`addon/process` uses a process group and SIGTERM; the lifecycle test went from
20.17 s to 0.20 s. `Commander` and `addon/process` should converge.

**5. Package names cannot be used as directory names.** Addon names always
contain `:`. `url.PathEscape` leaves it alone, so the PoC first stored the jar
under `jmx:tomcat/`. Two failures follow. Windows forbids `:` in file names. And
`java -jar` sets `java.class.path` to the jar's path, which the JVM splits on
`:`, so the scraper fails with `ClassNotFoundException` for its own main class.
Only the run against the real JVM found this; the stand-in JVM in the unit tests
did not. Names are now escaped with `url.QueryEscape`, and
`TestPackageDirIsSafeFileName` pins it.

## Validation

All pass on macOS, with `go vet` clean for Linux and Windows, and the addon
tests clean under the race detector.

End to end against the real artifact (`TestJMXScraperAddonEndToEnd`, needs
`OPAMP_JMX_SCRAPER_JAR` and `java`): the published
`opentelemetry-jmx-scraper.jar` is offered as `jmx:tomcat` with no wrapper,
installed through the provider call sequence opamp-go's syncer uses, configured
by a `role`-tagged remote config entry under the same key, run against a live
JVM, and its OTLP export is received. Shutting it down leaves the jar and stops
the exports; starting it again resumes them without a download.

Framework: ID grammar; duplicate and invalid type registration; unsupported types
rejected before download; configure only on change, nil on withdrawal, start on
every pass, failures isolated per addon; upgrade shutting down the old instance;
the constructor wiring the option; a distribution-defined `splunk_ta` type;
exactly one Collector restart per message, after the TA is on disk; no restart
when nothing changed.

Types: `file` raw and archived content and rejected config; `jmx` with a
stand-in JVM covering idle-until-configured, the fixed command line, restart on
config change, stop on withdrawal, missing JVM at install and at start, and
non-jar content rejected.

Lifecycle: content-hash mismatch rejected before install; a failed install
leaving the previous version; state, content hash and running addons surviving a
Supervisor restart; shutdown leaving no orphaned processes; addon config never
reaching the Collector's YAML.

## What this does not do

- No top-level package installation; that is #47272's binary update path.
- No signature verification. `InstallSettings` carries the signature and the
  Supervisor has a cosign verifier; wiring them is out of scope.
- No operator policy over which types may run, for example an allowlist in the
  Supervisor config. Worth adding before executable types are enabled by
  default.
- No advertisement of supported types to the server. A non-identifying
  `AgentDescription` attribute listing them would let a server avoid offering an
  unsupported type; today it learns from `INSTALL_FAILED`.
- No per-addon operational status (running, crash-looping) beyond
  `PackageStatuses`, which describes installation only.
- Windows graceful stop ends the process rather than sending a console control
  event, matching what the Supervisor does for the Collector.

## Relationship to #47272

#47272 implements packages for the **top-level** package, to upgrade the
Collector binary. Four of its seven steps are merged, most recently the Cosign
verifier (#50707, 2026-10-28). The remaining three are unclaimed: **package
manager**, **binary update orchestration** and **e2e and finalization**.

This PoC's `packages.go` overlaps step 5 and is not proposed as step 5. The
intent is a parallel addon track that consumes step 5's package manager:

- Step 5 implements the `types.PackagesStateProvider` stubs. The version here is
  a stand-in so the addon path runs today and should give way to step 5's. The
  ask of step 5 is that it not assume a single top-level package, and that it
  delegate `PackageType_Addon` to the addon manager.
- Addons need no part of step 6: no binary to swap atomically and no rollback of
  a running Collector, so the addon track can land before it.
- What the addon track adds is additive and in separate files: the type
  framework, the built-in types, the process supervisor, role routing and the
  coordinated restart.

Findings 3 and 4 matter to #47272 directly: step 6 will supervise a swapped
binary with `Commander`, and the step 4 extractor handles only a single named
binary.

## Suggested sequencing

1. Raise the parallel addon track on #47272 before writing more code, so step 5
   and the addon track agree on the provider boundary.
2. Open the contrib issue for typed addon support with `file` and `jmx`, the
   runtime boundary, and the no-generic-exec rule stated up front. No spec
   change.
3. Generalise process supervision (finding 4), usable by both tracks.
4. Land `role` handling (finding 2) with the interoperability question answered.
5. Distributions implement their own types against the framework; for Splunk,
   `splunk_ta` with state carried in configuration.
6. Revisit per-addon status reporting and, only if configuration-carried state
   proves insufficient, a protocol-level enable/disable.
