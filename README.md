# Moon Patrol

A Go + Ebitengine remake of the classic arcade game Moon Patrol: drive a lunar buggy across
the course, jump craters, shoot rocks and enemies. All art and sound are original and
generated procedurally; nothing is copied from the original game.

## Requirements

System packages (Debian/Ubuntu; cgo is required):

    sudo apt install gcc libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

Tools (Go, golangci-lint, gremlins, govulncheck, gitleaks, gocover-cobertura) are pinned in
`mise.toml`:

    mise install

## Build, run, test

    make build   # creates bin/moonpatrol
    make run     # builds and starts the game in a window
    make test    # go test -race ./...
    make lint    # golangci-lint

Flags: `-windowed`, `-seed <n>` (default 1), `-version`.

## Controls (defaults)

| Action | Keys |
|---|---|
| Move forward / back | Right / Left arrows |
| Jump | Space |
| Fire | Ctrl |
| Pause | P |
| Quit | Esc |
