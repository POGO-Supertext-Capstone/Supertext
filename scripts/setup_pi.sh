#!/usr/bin/env bash
# provision a Raspberry Pi 5 running Debian 13 as a Supertext server
#
# manual use:      sudo bash scripts/setup_pi.sh
# first-boot use:  scripts/first_boot.sh runs this with POGO_APP_USER and POGO_AP_PASSWORD set
#
# safe to re-run, if the hotspot profile already exists and no password is
# supplied, the existing hotspot password is kept
set -Eeuo pipefail

SSID="${POGO_SSID:-POGOSupertext}"
AP_CONNECTION="${POGO_AP_CONNECTION:-POGOSupertext}"
WIFI_IF="${POGO_WIFI_IF:-wlan0}"
AP_ADDRESS="10.42.0.1/24"
APP_PORT=8081
BINARY_NAME="supertext_server"
DNSMASQ_SOURCE="config/dnsmasq/captive.conf"
SERVICE_TEMPLATE="systemd/supertext.service.in"

log()  { printf '\n==> %s\n' "$*"; }
warn() { printf 'WARNING: %s\n' "$*" >&2; }
die()  { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

[[ "$EUID" -eq 0 ]] || die "Run as root, normally with: sudo bash scripts/setup_pi.sh"

# a normal user owns the checkout and runs the Go process. first_boot.sh passes
# POGO_APP_USER because cloud-init runs as root
APP_USER="${POGO_APP_USER:-${SUDO_USER:-}}"
[[ -n "$APP_USER" && "$APP_USER" != root ]] || \
    die "Could not identify the normal project user. Use sudo from that user's account, or set POGO_APP_USER."
APP_HOME="$(getent passwd "$APP_USER" | cut -d: -f6)"
[[ -n "$APP_HOME" && -d "$APP_HOME" ]] || die "Could not determine the home directory for $APP_USER."
APP_GROUP="$(id -gn "$APP_USER")"

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd -- "$SCRIPT_DIR/.." && pwd)"
[[ -f "$PROJECT_DIR/go.mod" ]] || die "go.mod not found at $PROJECT_DIR."
[[ -f "$PROJECT_DIR/cmd/main.go" ]] || die "Expected Go entry point cmd/main.go was not found."
[[ -f "$PROJECT_DIR/$DNSMASQ_SOURCE" ]] || die "Missing $DNSMASQ_SOURCE."
[[ -f "$PROJECT_DIR/$SERVICE_TEMPLATE" ]] || die "Missing $SERVICE_TEMPLATE."

[[ -r /etc/os-release ]] || die "Cannot identify the operating system."
. /etc/os-release
[[ "${ID:-}" == debian && "${VERSION_ID:-}" == 13 ]] || \
    die "This script targets Debian 13 (found ${PRETTY_NAME:-unknown})."

validate_password() {
    local p="$1"
    (( ${#p} >= 8 && ${#p} <= 63 )) || die "The Wi-Fi passphrase must be between 8 and 63 characters."
    LC_ALL=C grep -q '^[ -~]*$' <<<"$p" || die "The Wi-Fi passphrase may only contain printable ASCII characters."
}

# hotspot password comes from first_boot.sh or a prompt
# it is never stored in the Git repository 
# if the hotspot profile already exists, no password is needed and the existing one is kept
AP_PASSWORD="${POGO_AP_PASSWORD:-}"
HAVE_AP_PROFILE=0
if command -v nmcli >/dev/null 2>&1 && nmcli connection show "$AP_CONNECTION" >/dev/null 2>&1; then
    HAVE_AP_PROFILE=1
fi
if [[ -z "$AP_PASSWORD" && "$HAVE_AP_PROFILE" -eq 0 ]]; then
    [[ -t 0 ]] || die "No hotspot password supplied and no terminal to ask on. Set POGO_AP_PASSWORD."
    read -r -s -p "Enter the WPA2 passphrase for $SSID (8-63 characters): " AP_PASSWORD
    printf '\n'
    read -r -s -p "Confirm the Wi-Fi passphrase: " AP_PASSWORD_CONFIRM
    printf '\n'
    [[ "$AP_PASSWORD" == "$AP_PASSWORD_CONFIRM" ]] || die "The passphrases did not match."
    unset AP_PASSWORD_CONFIRM
fi
[[ -z "$AP_PASSWORD" ]] || validate_password "$AP_PASSWORD"

log "Checking internet access"
getent hosts github.com >/dev/null 2>&1 || \
    die "Cannot resolve github.com. Connect this Pi to the internet (Wi-Fi client or Ethernet). If wlan0 is already the Supertext hotspot, the Pi has no internet."

log "Installing required packages"
export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y \
    ca-certificates \
    curl \
    dnsmasq \
    git \
    golang-go \
    iptables \
    iptables-persistent \
    iw \
    network-manager

log "Checking Wi-Fi hardware and access-point support"
[[ -d "/sys/class/net/$WIFI_IF" ]] || die "Wireless interface $WIFI_IF does not exist."
command -v nmcli >/dev/null 2>&1 || die "nmcli is not installed."
command -v iw >/dev/null 2>&1 || die "iw is not installed."
WIPHY="$(iw dev "$WIFI_IF" info 2>/dev/null | awk '$1 == "wiphy" {print $2}')" || true
[[ -n "$WIPHY" ]] || die "Could not find the Wi-Fi device behind $WIFI_IF."
PHY_INFO="$(iw "phy$WIPHY" info)"
grep -Eq '^[[:space:]]*\*[[:space:]]+AP([[:space:]]|$)' <<<"$PHY_INFO" || \
    die "The Wi-Fi adapter behind $WIFI_IF does not report AP mode support."

as_app() { runuser -u "$APP_USER" -- env HOME="$APP_HOME" GOTOOLCHAIN=auto "$@"; }

log "Building $BINARY_NAME as $APP_USER (the first build may download Go 1.27 and modules)"
install -d -o "$APP_USER" -g "$APP_GROUP" "$PROJECT_DIR/build"
install -d -o "$APP_USER" -g "$APP_GROUP" "$APP_HOME/.local/share"
cd "$PROJECT_DIR"
as_app go mod download
PACKAGE_NAME="$(as_app go list -f '{{.Name}}' ./cmd)"
[[ "$PACKAGE_NAME" == main ]] || die "./cmd is not a Go main package (found '$PACKAGE_NAME')."
BUILD_TMP="$PROJECT_DIR/build/$BINARY_NAME.new.$$"
trap 'rm -f "${BUILD_TMP:-}"' EXIT
as_app go build -o "$BUILD_TMP" ./cmd
chmod 0755 "$BUILD_TMP"
mv -f "$BUILD_TMP" "$PROJECT_DIR/build/$BINARY_NAME"

log "Configuring the NetworkManager hotspot profile '$AP_CONNECTION'"
systemctl enable NetworkManager
systemctl start NetworkManager
nmcli device set "$WIFI_IF" managed yes

# new profile starts with autoconnect off so the hotspot cannot take wlan0
# away from the install connection before everything else has succeeded
if ! nmcli connection show "$AP_CONNECTION" >/dev/null 2>&1; then
    nmcli connection add type wifi ifname "$WIFI_IF" \
        con-name "$AP_CONNECTION" ssid "$SSID" connection.autoconnect no
fi

PSK_ARGS=()
if [[ -n "$AP_PASSWORD" ]]; then
    PSK_ARGS=(802-11-wireless-security.psk "$AP_PASSWORD")
fi

nmcli connection modify "$AP_CONNECTION" \
    connection.interface-name "$WIFI_IF" \
    connection.autoconnect-priority 100 \
    802-11-wireless.ssid "$SSID" \
    802-11-wireless.mode ap \
    802-11-wireless-security.key-mgmt wpa-psk \
    802-11-wireless-security.proto rsn \
    802-11-wireless-security.pairwise ccmp \
    802-11-wireless-security.group ccmp \
    802-11-wireless-security.pmf 0 \
    "${PSK_ARGS[@]}" \
    ipv4.method manual \
    ipv4.addresses "$AP_ADDRESS" \
    ipv4.gateway "" \
    ipv4.dns "" \
    ipv6.method ignore
unset AP_PASSWORD POGO_AP_PASSWORD PSK_ARGS || true

log "Installing the dnsmasq configuration"
install -d -m 0755 /etc/dnsmasq.d
install -m 0644 "$PROJECT_DIR/$DNSMASQ_SOURCE" /etc/dnsmasq.d/captive.conf
rm -f /etc/dnsmasq.d/supertext.conf

# add an include only when the main config has none
if ! grep -Eq '^[[:space:]]*conf-dir=/etc/dnsmasq\.d' /etc/dnsmasq.conf; then
    printf '\n# Load local dnsmasq configuration snippets.\nconf-dir=/etc/dnsmasq.d/,*.conf\n' >> /etc/dnsmasq.conf
fi

dnsmasq --test
systemctl enable dnsmasq
systemctl restart dnsmasq

log "Installing the persistent HTTP redirect (port 80 -> $APP_PORT on $WIFI_IF)"
iptables -w -t nat -C PREROUTING -i "$WIFI_IF" -p tcp --dport 80 \
    -j REDIRECT --to-ports "$APP_PORT" 2>/dev/null || \
iptables -w -t nat -A PREROUTING -i "$WIFI_IF" -p tcp --dport 80 \
    -j REDIRECT --to-ports "$APP_PORT"
systemctl enable netfilter-persistent
netfilter-persistent save

log "Installing the systemd service"
SERVICE_TMP="$(mktemp)"
trap 'rm -f "${BUILD_TMP:-}" "${SERVICE_TMP:-}"' EXIT
sed \
    -e "s|@APP_USER@|$APP_USER|g" \
    -e "s|@PROJECT_DIR@|$PROJECT_DIR|g" \
    -e "s|@APP_HOME@|$APP_HOME|g" \
    "$PROJECT_DIR/$SERVICE_TEMPLATE" > "$SERVICE_TMP"
install -m 0644 "$SERVICE_TMP" /etc/systemd/system/supertext.service
systemctl daemon-reload
systemctl enable supertext.service
systemctl restart supertext.service

log "Verifying Supertext"
sleep 3
if ! systemctl is-active --quiet supertext.service; then
    systemctl --no-pager --full status supertext.service || true
    journalctl -u supertext.service -n 60 --no-pager || true
    die "Supertext did not remain active after setup."
fi
systemctl --no-pager --full status supertext.service || true
if ! curl -fsS -o /dev/null --max-time 5 "http://127.0.0.1:$APP_PORT/" 2>/dev/null; then
    warn "Supertext is running but did not answer on port $APP_PORT. Check: journalctl -u supertext.service"
fi

rm -f "$PROJECT_DIR/build/supertext"

# enable automatic hotspot activation
# may activate profile immediately if wlan0 is available, so run provision
# over ethernet or local console when possible
nmcli connection modify "$AP_CONNECTION" connection.autoconnect yes

cat <<SUMMARY

Provisioning completed.

Hotspot SSID:           $SSID
Hotspot interface:      $WIFI_IF
Hotspot address:        $AP_ADDRESS
Application executable: $PROJECT_DIR/build/$BINARY_NAME
Application HTTP port:  $APP_PORT
Systemd service:        supertext.service

The hotspot profile is set to autoconnect at boot. May start hotspot immediately
if wlan0 is idle.
Reboot from a local console or over Ethernet:
    sudo reboot

After reboot, check:
    nmcli connection show --active
    systemctl is-active dnsmasq supertext netfilter-persistent
    sudo iptables -t nat -S PREROUTING
SUMMARY
