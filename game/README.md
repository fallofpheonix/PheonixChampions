# Phoenix Game Layer (Champions)

> **Status**: v1.1 Implementation Verified | **Foundation**: [SACRED]

Phoenix Champions is a deterministic, grid-based puzzle simulation designed to train autonomous AI agents. Players take the role of an **AI Architect**, writing Python, C++, or Java scripts to navigate high-fidelity security challenges.

## 🚀 Key Features

- **Multi-Language IDE**: In-game `CodeEdit` editor with syntax highlighting and hot-reloading for Python, C++, and Java.
- **Unified IPC Engine**: All languages run in persistent subprocesses, communicating with the core simulation via an atomic stdin/stdout bridge.
- **Deterministic Foundation**: Simulation logic is backed by fixed-point `int64` math to ensure bit-perfect consistency across platforms.
- **Hacker Control Room**: A dark-themed terminal UI featuring real-time telemetry, agent path trails, and dynamic level loading.
- **Integrated Error Console**: Compiler and runtime errors are routed directly from the backend into the visual client.

## 🛠️ Project Structure

- `core/`: Multi-language simulation engine (Python) and P-Script subsystem (Go).
- `client/`: Godot 4 visualizer and interactive GUI.
- `levels/`: Data-driven level definitions in JSON format.
- `docs/`: Detailed architectural specifications and development roadmap.

## 🏁 Getting Started

### Prerequisites
- Python 3.x
- Godot 4.x
- `g++` (for C++ support)
- `javac` (for Java support)

### Run the Game
```bash
cd pheonix-champions/game
make run-all
```

## 📜 Development Standards
- **Determinism is Sacred**: No floating-point math in the simulation core. Use `int64`.
- **Atomic Reliability**: All script deployments use atomic filesystem swaps to prevent race conditions.
- **Stateless Intelligence**: The interpreter/runners are stateless; the authoritative state is passed per-tick.

---
*Part of the Phoenix Matrix OS Monorepo.*
