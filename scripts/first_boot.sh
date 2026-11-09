#!/usr/bin/env bash
# supertext first-boot installer for a freshly flashed Raspberry Pi OS
# cloud-init downloads and runs this as root (see cloud_init/user_data.example)
#
# waits for internet, reads the hotspot password from supertext.env on the
# SD card's boot partition, clones the main branch, runs scripts/setup_pi.sh,
# and reboots so wlan0 becomes the hotspot
#
# if it fails, fix the cause and re-run:  sudo bash /root/supertext-first-boot.sh
# Details are in /var/log/cloud-init-output.log
set -Eeuo pipefail

REPO_URL="https://github.com/POGO-Supertext-Capstone/Supertext.git"
REPO_BRANCH="main"
APP_DIR_NAME="Supertext"
SECRETS_ON_BOOT="/boot/firmware/supertext.env"
SECRETS_FILE="/etc/supertext/provisioning.env"
DONE_MARKER="/var/lib/supertext/first-boot-complete"

log() { printf '\n==> %s\n' "$*"; }
die() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

[[ "$EUID" -eq 0 ]] || die "first_boot.sh must run as root."

if [[ -e "$DONE_MARKER" ]]; then
    log "First-boot provisioning has already completed; nothing to do."
    exit 0
fi

# read KEY=value from an env file WITHOUT running it as shell code
# handles windows (CRLF) line endings and optional surrounding quotes, and doesn't
# expand $ or backticks in the value
env_get() {
    local file="$1" key="$2" line value
    line="$(grep -E "^[[:space:]]*${key}=" "$file" | tail -n 1 || true)"
    value="${line#*=}"
    value="${value%$'\r'}"
    if [[ "$value" =~ ^\"(.*)\"$ ]] || [[ "$value" =~ ^\'(.*)\'$ ]]; then
        value="${BASH_REMATCH[1]}"
    fi
    printf '%s' "$value"
}

# user copies cloud_init/supertext.env.example to the boot partition as
# supertext.env and sets the hotspot password
# keep a root-only copy until revisioning succeeds, 
# and remove the world-readable one from the FAT partition
if [[ -f "$SECRETS_ON_BOOT" ]]; then
    install -d -m 0700 /etc/supertext
    install -o root -g root -m 0600 "$SECRETS_ON_BOOT" "$SECRETS_FILE"
    rm -f "$SECRETS_ON_BOOT"
fi
[[ -f "$SECRETS_FILE" ]] || \
    die "Missing $SECRETS_FILE. Put supertext.env on the SD card's boot partition before first boot."

POGO_AP_PASSWORD="$(env_get "$SECRETS_FILE" POGO_AP_PASSWORD)"
POGO_APP_USER="$(env_get "$SECRETS_FILE" POGO_APP_USER)"
POGO_SSID="$(env_get "$SECRETS_FILE" POGO_SSID)"
POGO_AP_CONNECTION="$(env_get "$SECRETS_FILE" POGO_AP_CONNECTION)"
POGO_WIFI_IF="$(env_get "$SECRETS_FILE" POGO_WIFI_IF)"

[[ -n "$POGO_AP_PASSWORD" ]] || die "POGO_AP_PASSWORD is not set in $SECRETS_FILE."
[[ "$POGO_AP_PASSWORD" != "CHANGE_THIS_TO_A_UNIQUE_WIFI_PASSWORD" ]] || \
    die "Replace the example hotspot password in supertext.env before booting the Pi."
(( ${#POGO_AP_PASSWORD} >= 8 && ${#POGO_AP_PASSWORD} <= 63 )) || \
    die "POGO_AP_PASSWORD must be between 8 and 63 characters."

log "Identifying the first normal user created in Raspberry Pi Imager"
APP_USER="${POGO_APP_USER:-$(getent passwd 1000 | cut -d: -f1)}"
[[ -n "$APP_USER" && "$APP_USER" != root ]] || \
    die "Could not find the initial non-root user (UID 1000). Create a user in Raspberry Pi Imager."
APP_HOME="$(getent passwd "$APP_USER" | cut -d: -f6)"
[[ -n "$APP_HOME" && -d "$APP_HOME" ]] || die "Home directory for $APP_USER does not exist."

log "Waiting for internet access (Wi-Fi can take a minute to connect)"
online=0
for _ in $(seq 1 60); do
    if getent hosts github.com >/dev/null 2>&1; then online=1; break; fi
    sleep 5
done
[[ "$online" -eq 1 ]] || die "No internet after 5 minutes. Check the Wi-Fi settings entered in Raspberry Pi Imager (or plug in Ethernet)."

log "Installing tools needed to clone the repository"
export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y ca-certificates git

run_as_app() { runuser -u "$APP_USER" -- env HOME="$APP_HOME" "$@"; }

PROJECT_DIR="$APP_HOME/$APP_DIR_NAME"
if [[ -d "$PROJECT_DIR/.git" ]]; then
    log "Reusing the existing checkout at $PROJECT_DIR"
    run_as_app git -C "$PROJECT_DIR" fetch origin "$REPO_BRANCH"
    run_as_app git -C "$PROJECT_DIR" checkout "$REPO_BRANCH"
    run_as_app git -C "$PROJECT_DIR" merge --ff-only "origin/$REPO_BRANCH"
elif [[ -e "$PROJECT_DIR" ]]; then
    die "$PROJECT_DIR exists but is not a Git checkout. Move it aside and retry."
else
    log "Cloning $REPO_BRANCH from GitHub into $PROJECT_DIR"
    run_as_app git clone --branch "$REPO_BRANCH" --single-branch "$REPO_URL" "$PROJECT_DIR"
fi
[[ -f "$PROJECT_DIR/scripts/setup_pi.sh" ]] || \
    die "The $REPO_BRANCH branch does not contain scripts/setup_pi.sh."

log "Running the Supertext provisioning script"
POGO_APP_USER="$APP_USER" \
POGO_AP_PASSWORD="$POGO_AP_PASSWORD" \
POGO_SSID="${POGO_SSID:-POGOSupertext}" \
POGO_AP_CONNECTION="${POGO_AP_CONNECTION:-POGOSupertext}" \
POGO_WIFI_IF="${POGO_WIFI_IF:-wlan0}" \
    bash "$PROJECT_DIR/scripts/setup_pi.sh"

# NetworkManager now holds the passphrase in its root-only connection profile
rm -f "$SECRETS_FILE"
unset POGO_AP_PASSWORD
install -d -m 0755 /var/lib/supertext
touch "$DONE_MARKER"

log "Provisioning succeeded. Rebooting so wlan0 can become the hotspot."
sync
systemctl reboot
