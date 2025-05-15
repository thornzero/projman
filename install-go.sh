#!/bin/bash

# Configuration
GO_ROOT="/usr/local/go"
INSTALL_DIR="$HOME/.local/bin"
GO_UPDATE_SCRIPT="$INSTALL_DIR/go-update"

# Function to fetch the latest Go version tag
get_latest_go_version() {
  curl -s https://go.dev/VERSION?m=text | head -n 1 | tr -d '\r\n'
}

# Check installed version
INSTALLED=$(go version 2>/dev/null | awk '{print $3}')
LATEST=$(get_latest_go_version)

echo "Installed: ${INSTALLED:-none}"
echo "Latest:    $LATEST"

if [[ "$INSTALLED" != "$LATEST" ]]; then
  echo "🔄 Updating Go to $LATEST"

  cd /tmp
  FILE="${LATEST}.linux-amd64.tar.gz"
  URL="https://go.dev/dl/${FILE}"

  wget -q --show-progress "$URL"
  sudo rm -rf "$GO_ROOT"
  sudo tar -C /usr/local -xzf "$FILE"
  rm "$FILE"

  echo "✅ Go updated to $LATEST"
else
  echo "✅ Go is already up to date."
fi

# Create full update script
mkdir -p "$INSTALL_DIR"
cat <<'EOF' > "$GO_UPDATE_SCRIPT"
#!/bin/bash

GO_ROOT="/usr/local/go"
LATEST=$(curl -s https://go.dev/VERSION?m=text | head -n1 | tr -d '\r\n')
INSTALLED=$(go version 2>/dev/null | awk '{print $3}')

echo "Installed: ${INSTALLED:-none}"
echo "Latest:    $LATEST"

if [[ "$INSTALLED" != "$LATEST" ]]; then
  echo "󰖷 Updating Go to $LATEST"

  cd /tmp
  FILE="${LATEST}.linux-amd64.tar.gz"
  URL="https://go.dev/dl/${FILE}"

  wget -q --show-progress "$URL"
  sudo rm -rf "$GO_ROOT"
  sudo tar -C /usr/local -xzf "$FILE"
  rm "$FILE"

  echo "✅ Go updated to $LATEST"
else
  echo "✅ Go is already up to date."
fi
EOF

chmod +x "$GO_UPDATE_SCRIPT"
echo "󰖷 'go-update' tool installed at $GO_UPDATE_SCRIPT"
