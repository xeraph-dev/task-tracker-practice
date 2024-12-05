# Task Tracker

[Project Details](https://roadmap.sh/projects/task-tracker)

## Features

- Add, update, and delete tasks.
- Mark tasks as **to-do**, **in-progress**, or **done**.
- List tasks by status: all, done, to-do, or in-progress.
- JSON storage: Task data is stored persistently in a JSON file.

### Task Properties

Each task contains:

- `id`: Unique identifier.
- `description`: Brief description of the task.
- `status`: Current status (`todo`, `in-progress`, or `done`).
- `createdAt`: Timestamp for creation.
- `updatedAt`: Timestamp for the last update.

## Getting Started

### Precompiled Binaries

Download precompiled binaries for Windows, macOS, and Linux from the [Releases](https://github.com/xeraph-dev/task-tracker-practice/releases) page.

### Building from Source

If you'd like to build the project yourself, a `Makefile` is provided to simplify the process.

#### Prerequisites

- [Go](https://golang.org/dl/)
- Make (optional, for using the `Makefile`)

#### Using the Makefile

1. Clone the repository:

   ```bash
   git clone https://github.com/xeraph-dev/task-tracker-practice.git
   cd task-tracker-practice
   ```

2. Build the project for your platform:

   ```bash
   make build
   ```

_This generates the binary in the build/ directory._

3. Run the project:

```bash
go run . [arguments]
```

4. Clean build artifacts:

```bash
make clean
```

##### Cross-Platform Builds

The release target in the Makefile builds binaries for Windows, macOS, and Linux:

```bash
make release
```

### Example Usage

- Add a Task:

  ```bash
  task add "Buy groceries"
  ```

- List Tasks:
  ```bash
  task list
  ```

## License

This project is licensed under the BSD License. See the [LICENSE](./LICENSE) file for more details.
