# Phoenix Game Layer — Current Implementation State

Phoenix Champions is a deterministic, multi-language simulation engine where players program autonomous agents to solve security-themed puzzles.

---

## 1. Project Directory Structure

- **`core/`**: The simulation core (The "Brains").
  - **`main.py`**: The primary multi-language simulation engine (Python, C++, Java).
  - **`main.go`**: The Go-based P-Script interpreter (Lexer/Parser/AST verified).
  - **`levels/`**: JSON definitions for levels (1-4), goals, and obstacle blocks.
  - **`scripts/`**: User-generated scripts (`agent.py`, `agent.cpp`, `agent.java`).
  - **`state/`**: Shared state definitions using fixed-point `int64` math.
- **`client/`**: The Godot 4 visual client (The "Eyes").
  - **`main.gd`**: Manages the IPC bridge, editor, and telemetry.
  - **`level_select.gd`**: Dynamic level scanner and UI card generator.
- **`Makefile`**: Orchestration for `run-all`, `run-core`, and `test-core`.
- **`bridge_schema.json`**: The JSON contract for Go-Godot state exchange.

---

## 2. Multi-Language IPC Architecture

The simulation uses a **Unified Persistent Process** model for all languages:
1. **Runner Wrappers**: When a script is deployed, the core generates a wrapper (e.g., `AgentRunner.java` or `agent_wrapper.py`) containing the API stubs.
2. **IPC Bridge**: The core communicates with the subprocess via `stdin/stdout`. Commands like `move_forward()` block until the core sends an acknowledgment.
3. **Error Routing**: Compilation errors (C++/Java) and Runtime errors (Python) are captured by the core and emitted as a `compile_error` status in the JSON bridge.

---

## 3. UI/UX: Hacker Control Room

The visual client has been overhauled with a high-fidelity "Control Room" aesthetic:
- **`CodeEdit` Integration**: Proper IDE feel with syntax highlighting for multiple languages.
- **Real-Time Telemetry**: Cyan-on-black dashboard tracking Ticks, Coordinates, and Goal status.
- **Agent Trail**: Visual feedback showing the agent's path over the last 5 ticks.
- **Data-Driven UX**: Level title, description, and API reference list are populated dynamically from the level JSON.
- **Window Management**: F11 toggle for fullscreen mode.

---

## 4. Stability & Integrity (SACRED)

- **Determinism**: All simulation math uses `int64` (100 units = 1 pixel) to ensure consistency.
- **Resilience**: 
  - **Atomic Swaps**: Godot writes scripts to `.tmp` files then renames them to prevent partial reads by the core.
  - **Pipe Lifecycle**: The core self-terminates on `BrokenPipeError` to prevent memory leaks.
  - **Process Cleanup**: Explicit `OS.kill(pid)` on scene changes or restarts.

---

## 5. Active Levels

1. **Level 1 (Diagonal Dash)**: Basic 2D movement and goal reaching.
2. **Level 2 (The L-Turn)**: Reactive programming using `get_x()` and `get_y()`.
3. **Level 3 (The Staircase)**: Obstacle avoidance and function abstraction.
4. **Level 4 (The Fork)**: Conditional logic based on world-sensing (`get_goal_y()`).

---

## 6. Commands

```bash
make run-all   # Launch the full Godot game
make test-core # Run Go subsystem tests
```
