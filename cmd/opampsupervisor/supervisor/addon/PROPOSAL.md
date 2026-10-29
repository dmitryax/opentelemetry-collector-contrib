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

Each piece maps onto a protocol field that already exists:

| jmxreceiver | OpAMP Packages |
|---|---|
| `hashFile()` + 83-entry allowlist | `DownloadableFile.content_hash`, supplied per offer |
| `supportedJar.version` | `PackageAvailable.version` |
| `jar_path`, populated manually | `DownloadableFile.download_url` |
| compile-time hash map, needs a release | server computes the hash per offer |
| `subprocess.Config` | Supervisor-managed executable addon |

## The objection, and the answer

Issue **#45740** (rogercoll, 2026-01-30) deprecated the receiver and states:

> Requiring the Collector to run Java is problematic, especially since the
> Collector is commonly deployed as a Docker image and none of the official
> images ship a Java runtime [...] More broadly, the Collector should not run
> subprocesses at all; the jmxreceiver is an exception that sets an undesirable
> precedent, even with its checksum-based restrictions.

That last clause pre-rejects "we would hash it properly". Three responses:

1. **The Supervisor is not the Collector.** #45740 is scoped to the Collector
   process: no JRE in its image, and subprocess lifecycle entangled with
   receiver start and shutdown (#23051). The Supervisor is already a process
   manager; starting, stopping, restarting and configuring the Collector is its
   entire job. A second managed child is a quantitative change, not a new
   architectural precedent.
2. **"Traditional mechanisms" is doing a lot of work.** In Kubernetes that means
   a sidecar and this proposal is unnecessary. On VM and bare-metal fleets it
   means hand-rolled unit files plus a private distribution channel per
   operator, which is the problem OpAMP exists to solve.
3. **The runtime problem does not go away, and this proposal concedes it.** A
   Supervisor-delivered jar fails on a host without Java exactly as a
   Collector-delivered one does, and in containers it fails identically because
   the Supervisor is in the container too.

On (3), the scope boundary is explicit: **executable addons target host and VM
deployments, not containers.** An addon's runtime prerequisite is its
entrypoint, resolved before launch, and a missing runtime is a deterministic,
non-retryable failure naming what is absent rather than an agent that merely
looks broken. `ErrRuntimeNotFound` and
`TestExecutableAddonMissingRuntimeFailsDeterministically` implement and cover
this. File addons have none of these problems and work everywhere.

## The design question this PoC was built to answer: who configures an addon?

The JMX scraper needs `otel.jmx.service.url` per target. Repackaging the addon
per target was ruled out immediately: it destroys artifact reuse and makes the
content hash per-host. That leaves two options.

### Option B: remote config, routed by `role`

`AgentConfigMap` is already a map of named objects, and `AgentConfigObject` has
a `role` field (added in opamp-spec PR #350, 2026-08-03, resolving #184). The
spec's own description of `role` is:

> an Agent type may use this field to distinguish top-level configuration from
> supplementary content such as certificates or other artifacts referenced by
> configuration

So this needs **no spec change**. An entry whose `role` is `addon-config` and
whose key names an addon package configures that addon.

Worth noting: #184 explicitly considered and rejected delivering supplementary
content as packages, calling it "out-of-band for configuration and [...] likely
a complicated flow". That supports the split this PoC implements: the addon
*artifact* travels as a package, its *configuration* as remote config.

**For:** no spec change; central fleet management; one offer configures many
addons; per-target values without repackaging; reuses the server's existing
config delivery, hashing and `RemoteConfigStatus` machinery.

**Against:** addon config is hashed and versioned together with Collector
config, so changing a scraper property bumps `config_hash` and travels the
remote-config apply path, mixing two blast radii; a change to addon config
restarts the addon, not the Collector, which is a different failure domain from
what `config_hash` implies. There is also an interoperability hazard, below.

### Option C: Supervisor-injected from its own local configuration

**For:** zero protocol involvement; the natural home for genuinely host-local
values.

**Against:** it defeats central management. Configuring ten thousand hosts by
editing their Supervisor configs is what OpAMP exists to avoid. There is no
status reporting path and no way to make a fleet-wide change.

### Recommendation

**Option B for delivery, with Option C as a declared escape hatch for host-local
values.** The PoC implements B.

One honest limitation: **B alone does not solve per-host values.** The PoC
delivers one config body to one addon on every host. That is fine where the
target is uniform (`localhost:<fixed port>`, which covers the common JMX case)
and insufficient where it is not. Nothing in OpAMP substitutes per-host values
into a config body. If per-host values are needed, the Supervisor's own config
has to supply them and the remote body has to be a template. That is a real
design question still open, and it is not specific to addons.

## What the PoC implements

```
<storage>/packages/
  state.json                     package records and all_packages_hash
  last_reported_statuses.bin     PackageStatuses, so a restart reports what it has
  <escaped-package-name>/
    content                      the package file exactly as downloaded
    files/                       the extracted addon tree
      opamp-addon.yaml           the addon manifest
      addon.log                  the addon process's stdout and stderr
```

`packages.go` is a real `types.PackagesStateProvider`. `supervisor.go` supplies
it to the OpAMP client and routes `role`-tagged config away from the Collector's
YAML. `addon/` holds the manifest, extraction, process supervision and registry.

An addon describes itself, so one offer serves any Supervisor without the server
knowing how to launch it:

```yaml
type: executable                 # or "file"
entrypoint: ["java", "-jar", "opentelemetry-jmx-scraper.jar", "-config", "scraper.properties"]
config_file: scraper.properties  # where the Supervisor writes delivered config
env: {}                          # optional, merged into the process environment
```

`entrypoint[0]` is resolved on PATH when it has no path separator, and inside
the addon directory otherwise. The process runs with the addon directory as its
working directory and receives `OPAMP_ADDON_DIR` and, when a config file is
declared, `OPAMP_ADDON_CONFIG`.

Reconciliation rather than events: the registry is told which addons should be
present and starts or stops accordingly. That stays correct after a restart,
where no messages are replayed, and avoids ordering install, configure and start
against each other.

## Findings

Four things the PoC established that were not visible from reading the code.

**1. `accepts_packages: true` is broken today, not merely incomplete.** The
Supervisor declared the capability but never supplied a
`PackagesStateProvider`, and opamp-go's `validateCapabilities` returns
`ErrPackagesStateProviderNotSet` from `PrepareStart`, so the OpAMP client
refuses to start. Reproduced in `packages_gap_test.go`. Worth reporting
separately from this proposal.

**2. The `role` field is specified but ignored.** `createRemoteConfigComposers`
merges every `config_map` entry into the Collector's YAML. A server that sends a
`role`-tagged entry to a current Supervisor corrupts the Collector config. So
delivering addon config over remote config needs no spec change, but it does
need a way for the server to know the agent honours `role`. A capability bit or
a documented minimum version is required; silent merging is the current
behaviour and it fails closed-ish at best.

**3. The existing `archive` package cannot extract an addon.** It pulls a single
named binary out of a package. An addon is a tree, so `addon/extract.go` adds
whole-archive extraction with escape protection: entries that are absolute, that
traverse upward, or that are symlinks, hardlinks or devices are rejected
outright rather than skipped, and entry count and total size are bounded.
Covered by `TestExtractTarGzRejectsHostileArchives`.

**4. `Commander` cannot supervise an addon, and the reason is measurable.** It
signals with `os.Interrupt` (SIGINT) and does not place the child in its own
process group, because the Collector is a Go binary that handles SIGINT and has
no children. A shell-wrapped addon ignores SIGINT while waiting on its own
child, so `Stop` burns the entire grace period and then SIGKILLs. Measured
directly against `Commander` with no addon code involved:

| Delay between Start and Stop | `Commander.Stop` latency |
|---|---|
| 0 ms | 408 µs |
| 50 ms | 556 µs |
| 200 ms | **10.005 s** (the full grace period, then kill) |

`addon/process.go` runs addons in their own process group and signals the group
with SIGTERM. `TestAddonLifecycle` went from **20.17 s to 0.20 s**.

This is the concrete form of a general point: the Supervisor already has process
supervision, but its semantics are Collector-specific. `Commander` and
`addon.process` should converge on a shared process supervisor parameterised by
signal, process-group and grace-period policy, instead of `Commander` being
coupled to `config.Agent`.

## Validation

Everything below passes, with no regressions in the existing Supervisor suite.

End-to-end against the real artifact (`TestJMXScraperAddonEndToEnd`, requires
`OPAMP_JMX_SCRAPER_JAR` and `java`): a JVM is started with the JMX remote agent
enabled; the real `opentelemetry-jmx-scraper.jar` is packaged as an addon and
installed through the same provider call sequence opamp-go's syncer uses;
configuration is delivered as a `role`-tagged remote config entry; the addon
runs as a managed child, connects to the JVM and exports OTLP to a local
endpoint. Stopping it leaves the jar installed and stops the exports; starting it
again resumes them without a download.

Lifecycle and failure modes (`addons_test.go`): install, configure, run,
reconfigure with restart, no restart on unchanged config, upgrade replacing a
running process, delete stopping the process and removing files, state and
content hash surviving a restart, `PackageStatuses` round-tripping, content-hash
mismatch rejected before extraction, a malformed package leaving the previous
version installed, a file addon never starting a process, a missing runtime
failing deterministically, and addon config never reaching the Collector's YAML.

## What this does not do

- No top-level package installation; that is #47272's binary update path.
- No signature verification. The Supervisor already has a cosign verifier in
  `supervisor/verifier`, and `UpdateContent` receives the signature; wiring them
  together is deliberately left out of scope.
- No new Supervisor configuration surface for addon directories or policy.
- No enable/disable state. See below.
- Windows graceful stop ends the process rather than sending a console control
  event, which matches what the Supervisor does for the Collector today.

## Path to enable/disable

`Addon.Stop` and `Addon.Start` already give the mechanism: an installed addon
whose process is not running. The protocol has no way to ask for it. Removing a
package from `packages_available` means delete, per Downloading Packages step 2,
so the only way to stop a scraper today is to delete its jar and re-download it
to resume. The PoC demonstrates that stopping and restarting are cheap and that
the artifact survives, which is exactly what makes delete-to-disable the wrong
mechanism.

Three scenarios, none vendor-specific: a monitored JVM goes into maintenance and
its scraper should stop rather than be deleted; an addon is staged across a
fleet and started later, which is impossible today because offering a package
activates it; a misbehaving addon needs stopping without a round trip through
the download server to restart it.

This should be proposed to opamp-spec only once addon support exists upstream,
so the field has a reference implementation behind it. Note that the state is
meaningful for executable addons and fuzzy for file addons, where "disabled" is
hard to distinguish from configuration.

## Suggested sequencing

1. Report finding 1 as a bug; it is independent of everything else.
2. Open the contrib issue for addon support, scoped to `file` and `executable`,
   with the runtime boundary stated up front. No spec change required.
3. Generalise process supervision (finding 4) either before or alongside.
4. Land `role` handling (finding 2) with the interoperability question answered.
5. Only then propose enable/disable to opamp-spec.
