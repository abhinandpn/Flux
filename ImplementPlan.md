# Flux --- Implementation Plan

**Tagline:** Your network. At a glance.

## 1. Project Goal

Build Flux as a local-first, cross-platform network traffic monitor for
macOS, Ubuntu/Linux, and Android, with a future path to iOS/iPadOS.

Flux continuously displays upload and download traffic using
operating-system network counters. It does not run repeated internet
speed tests.

## 2. Engineering Principles

-   **Zero third-party Go modules:** use only the Go standard library in
    Go code.
-   **Native platform APIs:** use documented operating-system APIs to
    collect network counters.
-   **Shared core:** keep traffic calculations and common data models in
    pure Go.
-   **Platform adapters:** isolate OS-specific collection behind a
    common interface.
-   **Local-first:** no account, cloud backend, or database is required
    for the initial MVP.
-   **Test-first milestones:** implement, test, review, and commit each
    milestone before proceeding.
-   **No required shell utilities:** do not depend on tools such as
    `netstat` or `ip` at runtime.

Official platform SDKs and native build toolchains may still be
required. "Zero third-party Go modules" does not mean every platform SDK
can be avoided.

## 3. Target Repository Structure

``` text
flux/
├── cmd/flux/main.go
├── internal/
│   ├── core/
│   │   ├── types.go
│   │   ├── calculator.go
│   │   ├── units.go
│   │   └── calculator_test.go
│   ├── monitor/
│   │   ├── monitor.go
│   │   ├── sampler.go
│   │   └── monitor_test.go
│   ├── network/
│   │   ├── interfaces.go
│   │   └── interfaces_test.go
│   └── platform/
│       ├── darwin/
│       │   ├── counters.go
│       │   └── interfaces.go
│       ├── linux/
│       │   ├── counters.go
│       │   └── interfaces.go
│       └── android/README.md
├── apps/
│   ├── android/
│   ├── macos/Sources/
│   └── ios/
├── tests/
│   ├── fixtures/
│   └── integration/
├── docs/
├── scripts/
├── go.mod
├── README.md
├── IMPLEMENTATION_PLAN.md
├── LICENSE
└── .gitignore
```

This is the target structure. Files should be implemented incrementally;
empty placeholders do not count as completed functionality.

## 4. Implementation Roadmap

### Phase 1 --- Pure Go Core

1.  **`internal/core/types.go`**
    -   Define network counter samples and calculated traffic rates.
    -   Include interface identity, byte counters, and timestamps.
2.  **`internal/core/calculator.go`**
    -   Calculate download and upload bytes per second from counter
        deltas.
    -   Reject missing timestamps, invalid intervals, interface changes,
        and decreasing counters.
    -   Return explicit errors so the monitor can recover with a fresh
        baseline.
3.  **`internal/core/units.go`**
    -   Format rates into readable binary units such as B/s, KiB/s,
        MiB/s, and GiB/s.
    -   Handle negative or invalid values safely.
4.  **Core tests**
    -   Test normal calculations, zero traffic, invalid intervals,
        interface changes, counter resets, and unit formatting.

**Exit criteria:** `go test ./...`, `go test -race ./...`, and
`go vet ./...` pass; no third-party Go modules are required.

### Phase 2 --- Monitoring Engine

5.  **`internal/monitor/sampler.go`**
    -   Define sampling intervals and a cancellable sampling loop.
    -   Avoid blocking the consumer of traffic updates.
6.  **`internal/monitor/monitor.go`**
    -   Coordinate the counter provider and rate calculator.
    -   Emit timestamped samples and recover from temporary provider
        errors.
    -   Reset the baseline after counter resets or interface changes.
7.  **Monitor tests**
    -   Use fake counter providers to test timing, cancellation, errors,
        recovery, and repeated samples without depending on a real
        network.

**Exit criteria:** the monitoring lifecycle is testable with simulated
counters and has no OS-specific dependencies.

### Phase 3 --- Platform Abstraction and Native Counters

8.  **`internal/network/interfaces.go`**
    -   Define shared network interface models and the counter-provider
        contract.
    -   Keep OS-specific structures out of the shared core.
9.  **macOS adapter**
    -   Enumerate interfaces and read cumulative byte counters through
        native macOS APIs.
    -   Handle Wi-Fi, Ethernet, virtual interfaces, VPNs, and interface
        changes.
    -   Avoid parsing command output for runtime monitoring.
10. **Linux adapter**
    -   Read kernel-provided interface counters using an appropriate
        native Linux interface.
    -   Handle interface discovery, virtual interfaces, VPNs, and
        changing connections.
11. **Platform integration tests**
    -   Validate counter direction, monotonicity, interface identity,
        and behavior during disconnects.
    -   Document differences in what each OS exposes.

**Exit criteria:** real counter samples are collected on macOS and Linux
and consumed by the shared monitor.

### Phase 4 --- CLI MVP

12. **`cmd/flux/main.go`**
    -   Wire the core, monitor, and current platform adapter together.
13. **`flux watch`**
    -   Display live download and upload rates.
    -   Support clean terminal output and Ctrl+C shutdown.
14. **`flux status`**
    -   Show whether monitoring is running and which interface is
        selected.
15. **`flux interfaces`**
    -   List detected interfaces and their state.
16. **CLI reliability**
    -   Report actionable errors and avoid panics for expected runtime
        failures.
    -   Keep command behavior documented.

**Exit criteria:** Flux monitors real traffic from Terminal without
requiring a UI or external helper command.

### Phase 5 --- Ubuntu/Linux Validation

17. Build and run the CLI on supported Linux architectures.
18. Test multiple interfaces, VPNs, virtual interfaces, disconnects, and
    permission differences.
19. Verify build tags and platform-specific files compile only where
    intended.
20. Document Linux-specific behavior and limitations.

**Exit criteria:** the Linux adapter and CLI pass the agreed tests on a
real Linux environment.

### Phase 6 --- Android Application

21. Create the native Android application using Kotlin and the official
    Android SDK.
22. Implement traffic collection using supported Android APIs and
    document their limitations.
23. Connect shared Go logic only through a deliberate, tested
    integration boundary.
24. Implement the floating overlay and permission flow.
25. Build the compact speed pill and expanded statistics panel.
26. Handle app lifecycle, background execution, battery use, and overlay
    removal.
27. Test on multiple Android versions and devices.

**Exit criteria:** Android displays traffic data reliably while
respecting platform permissions and background restrictions.

### Phase 7 --- macOS Desktop Application

28. Create the native macOS app using SwiftUI and the official Apple
    toolchain.
29. Implement the menu-bar indicator.
30. Build the floating panel with live rates and selected-interface
    information.
31. Connect the UI to monitoring through a defined local integration
    boundary.
32. Handle launch, quit, sleep/wake, network changes, and error states.

**Exit criteria:** the desktop UI works without requiring Terminal to
remain open.

### Phase 8 --- Reliability, Security, and Performance

33. Test counter resets, counter wraparound where applicable, invalid
    samples, and elapsed-time anomalies.
34. Test sleep/wake, network changes, VPNs, multiple interfaces, and
    long-running sessions.
35. Measure CPU and memory use at different sampling intervals.
36. Minimize privileges and request only permissions needed for each
    platform.
37. Verify no traffic data leaves the device in the local-first MVP.
38. Run tests and static analysis for each supported build target.

**Exit criteria:** known limitations are documented, failure modes
recover cleanly, and performance is acceptable.

### Phase 9 --- Documentation and Release

39. Update the README with supported platforms, build instructions, and
    feature status.
40. Document architecture, native API choices, testing, troubleshooting,
    and limitations.
41. Add reproducible build/test scripts where useful.
42. Prepare platform-specific installation and upgrade instructions.
43. Tag and publish the first tested MVP.

**Exit criteria:** another developer can build, test, and run Flux using
the documentation.

### Phase 10 --- Future Enhancements

-   Local traffic history and usage summaries.
-   Configurable display units and sampling intervals.
-   Alerts for user-defined traffic thresholds.
-   iOS/iPadOS widgets or Live Activities where platform APIs allow
    them.
-   Additional OS adapters.
-   Optional cloud sync only if a concrete need emerges.

These are not MVP requirements.

## 5. Recommended Development Order

1.  Pure Go core and unit tests.
2.  Monitoring engine with fake providers.
3.  macOS native counter adapter.
4.  CLI MVP on macOS.
5.  Linux adapter and real Linux validation.
6.  Android monitoring integration and floating UI.
7.  macOS desktop UI.
8.  Reliability, documentation, and release.
9.  Future enhancements.

This order validates the core and real data collection before investing
in multiple UIs.

## 6. Definition of Done for Every Milestone

A milestone is complete only when:

-   The code is implemented and formatted with `gofmt`.
-   Relevant tests pass.
-   `go vet ./...` reports no unresolved issues.
-   No unapproved third-party Go dependencies are introduced.
-   Errors and platform limitations are documented.
-   The milestone is reviewed and committed to version control.

## 7. Initial Commands

From the repository root:

``` bash
go version
cat go.mod
go test ./...
go test -race ./...
go vet ./...
```

The first implementation task is to finish the pure Go core in
`internal/core/`, test it with simulated samples, and only then begin
the native macOS counter adapter.
