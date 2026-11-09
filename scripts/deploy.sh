#!/usr/bin/env bash
# deploy the latest origin/main to this Raspberry Pi
# use as the normal project user (NOT with sudo):   bash scripts/deploy.sh
# SKIP_TESTS=1 bash scripts/deploy.sh    skips 'go test' (e.g. a last-minute fix)
#
# while wlan0 is the Supertext hotspot the Pi has no internet, so use
# ethernet or a second Wi-Fi adapter for updates
#
# switches to main if needed, fast-forwards it, syncs the dnsmasq and
# systemd files if they changed, rebuilds build/supertext_server, restarts
# supertext.service, and rolls back to the previous executable if the new one
# does not come up
set -Eeuo pipefail

BINARY_NAME="supertext_server"
SERVICE_NAME="supertext.service"
APP_PORT=8081
DNSMASQ_SOURCE="config/dnsmasq/captive.conf"
SERVICE_TEMPLATE="systemd/supertext.service.in"

log()  { printf '\n==> %s\n' "$*"; }
warn() { printf 'WARNING: %s\n' "$*" >&2; }
die()  { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

[[ "$EUID" -ne 0 ]] || die "Run as the normal project user, without sudo."

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd -- "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_DIR"

[[ -d .git ]] || die "This directory is not the root of a Git checkout."
[[ -f go.mod && -f cmd/main.go ]] || die "Expected go.mod and cmd/main.go were not found."
[[ -f "$DNSMASQ_SOURCE" && -f "$SERVICE_TEMPLATE" ]] || die "Missing $DNSMASQ_SOURCE or $SERVICE_TEMPLATE."
for tool in go git sudo curl; do
    command -v "$tool" >/dev/null 2>&1 || die "$tool is not installed. Run scripts/setup_pi.sh first."
done
git remote get-url origin >/dev/null 2>&1 || die "The checkout has no 'origin' remote."

# Go 1.27+ is required but debian packages an older Go so fetch the right toolchain
export GOTOOLCHAIN=auto

APP_USER="$(id -un)"
APP_HOME="$(getent passwd "$APP_USER" | cut -d: -f6)"

BUILD_TMP=""
SERVICE_NEW=""
cleanup() { rm -f "${BUILD_TMP:-}" "${SERVICE_NEW:-}"; }
trap cleanup EXIT

# don't discard work done directly on the Pi
if ! git diff --quiet || ! git diff --cached --quiet; then
    die "Tracked files have local modifications. Commit, stash, or save them before deploying."
fi
UNTRACKED="$(git ls-files --others --exclude-standard)"
if [[ -n "$UNTRACKED" ]]; then
    warn "Untracked files will be left alone:"
    printf '%s\n' "$UNTRACKED" >&2
fi

log "Fetching origin/main"
if ! git fetch origin +refs/heads/main:refs/remotes/origin/main; then
    die "Could not reach GitHub. While wlan0 is the Supertext hotspot the Pi has no internet; connect Ethernet or a second Wi-Fi adapter and retry."
fi

CURRENT="$(git branch --show-current)"
if [[ "$CURRENT" != main ]]; then
    log "Switching from '${CURRENT:-detached HEAD}' to main"
    if git show-ref --verify --quiet refs/heads/main; then
        git checkout main
    else
        git checkout -b main --track origin/main
    fi
fi
git branch --set-upstream-to=origin/main main >/dev/null 2>&1 || true

# fast-forward, never overwrite local-only commits or merge a divergent branch
git merge-base --is-ancestor HEAD origin/main || \
    die "Local main has commits that are not on origin/main (or has diverged). Resolve then retry."
log "Fast-forwarding local main to origin/main"
git merge --ff-only origin/main

# keep the system files in step with the repository
log "Checking system configuration files"
DNSMASQ_LIVE="/etc/dnsmasq.d/captive.conf"
if ! cmp -s "$DNSMASQ_SOURCE" "$DNSMASQ_LIVE"; then
    if ! DNSMASQ_CHECK="$(sudo dnsmasq --test --conf-file="$PROJECT_DIR/$DNSMASQ_SOURCE" 2>&1)"; then
        die "New dnsmasq configuration failed its syntax check; not installed: $DNSMASQ_CHECK"
    fi
    log "Installing updated dnsmasq configuration"
    sudo install -m 0644 "$DNSMASQ_SOURCE" "$DNSMASQ_LIVE"
    sudo systemctl restart dnsmasq
fi

SERVICE_LIVE="/etc/systemd/system/$SERVICE_NAME"
SERVICE_NEW="$(mktemp)"
sed \
    -e "s|@APP_USER@|$APP_USER|g" \
    -e "s|@PROJECT_DIR@|$PROJECT_DIR|g" \
    -e "s|@APP_HOME@|$APP_HOME|g" \
    "$SERVICE_TEMPLATE" > "$SERVICE_NEW"
if ! cmp -s "$SERVICE_NEW" "$SERVICE_LIVE"; then
    log "Installing updated $SERVICE_NAME"
    sudo install -m 0644 "$SERVICE_NEW" "$SERVICE_LIVE"
    sudo systemctl daemon-reload
    sudo systemctl enable "$SERVICE_NAME"
fi

log "Downloading Go modules"
go mod download

if [[ "${SKIP_TESTS:-0}" == 1 ]]; then
    warn "SKIP_TESTS=1: not running go test"
else
    log "Running tests"
    go test ./...
fi

log "Building build/$BINARY_NAME"
install -d -m 0755 build
BUILD_TMP="$PROJECT_DIR/build/$BINARY_NAME.new.$$"
go build -o "$BUILD_TMP" ./cmd
chmod 0755 "$BUILD_TMP"

LIVE="$PROJECT_DIR/build/$BINARY_NAME"
PREV="$LIVE.prev"
rm -f "$PREV"
[[ ! -f "$LIVE" ]] || cp -p "$LIVE" "$PREV"
mv -f "$BUILD_TMP" "$LIVE"
BUILD_TMP=""

wait_until_healthy() {
    local i
    for i in 1 2 3 4 5 6 7 8 9 10; do
        if systemctl is-active --quiet "$SERVICE_NAME" && \
           curl -fsS -o /dev/null --max-time 3 "http://127.0.0.1:$APP_PORT/" 2>/dev/null; then
            return 0
        fi
        sleep 1
    done
    return 1
}

log "Restarting $SERVICE_NAME"
RESTART_OK=1
sudo systemctl restart "$SERVICE_NAME" || RESTART_OK=0
if [[ "$RESTART_OK" -ne 1 ]] || ! wait_until_healthy; then
    warn "$SERVICE_NAME is not healthy after the restart."
    sudo systemctl --no-pager --full status "$SERVICE_NAME" || true
    sudo journalctl -u "$SERVICE_NAME" -n 60 --no-pager || true
    if [[ -f "$PREV" ]]; then
        warn "Rolling back to the previous executable."
        mv -f "$PREV" "$LIVE"
        sudo systemctl restart "$SERVICE_NAME" || true
    fi
    die "Deployment failed; see the log above."
fi
rm -f "$PREV"

log "Deployment successful"
printf 'Branch:      %s\nCommit:      %s\nExecutable:  %s\n' \
    "$(git branch --show-current)" "$(git rev-parse --short HEAD)" "$LIVE"
