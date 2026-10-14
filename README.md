# Flux

**Your network. At a glance.**

Flux is a planned cross-platform network speed monitor that keeps
download and upload activity easy to see through a tiny floating pill
and an expandable statistics panel.

## Product Showcase

![Flux cross-platform UI
showcase](assets/Cross-Platform%20Network%20Speed%20Monitor%20Showcase.png)

## System Architecture

![Flux architecture
overview](assets/Floating%20Internet%20Speed%20Monitor%20Architecture.png)

## Planned Features

-   Minimal floating indicator for download and upload traffic.
-   Expandable panel with live traffic graphs and network statistics.
-   Draggable, customizable overlay on platforms that support it.
-   Local-first monitoring without requiring an account.
-   Platform-appropriate widgets for iOS and iPadOS, where free-floating
    overlays are not available.
-   Optional cloud features in the future, if needed.

## Planned Technology

-   **Android:** React Native UI with native Kotlin integration for
    overlay and network monitoring.
-   **Desktop:** React + TypeScript with Tauri and platform-specific
    integrations.
-   **iOS / iPadOS:** React Native app with native SwiftUI / WidgetKit
    extensions.
-   **Backend:** Go, optional for future cloud synchronization or
    account-based features.
-   **Local storage:** SQLite or platform-native storage, as
    appropriate.

> The stack and features above describe the current plan; implementation
> details may change as development progresses.

## Repository Structure

``` text
Flux/
├── apps/
│   ├── mobile/          # React Native application
│   └── desktop/         # React + Tauri application
├── packages/
│   └── shared/          # Shared TypeScript types and utilities
├── backend/             # Optional Go API
├── assets/
│   ├── Cross-Platform Network Speed Monitor Showcase.png
│   └── Floating Internet Speed Monitor Architecture.png
├── docs/
├── LICENSE
└── README.md
```

## Development Status

**Planning / early development.** The images are product concept
mockups, not screenshots of a finished implementation.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE)
for details.
