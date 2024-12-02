#!/usr/bin/env bash
# Copyright (C) 2023, Your Name/Organization
# All rights reserved.

set -o errexit
set -o nounset
set -o pipefail

# Set the CGO flags for building with portable configurations if applicable
export CGO_CFLAGS="-O -D__PORTABLE__"

# Ensure the script is run from the repository root
if ! [[ "$0" =~ scripts/build.sh ]]; then
  echo "This script must be run from the repository root"
  exit 255
fi

# Create the build directory
mkdir -p ./build

# Build the VM binary
vm_name="vm_binary"
echo "Building VM in ./build/$vm_name"
go build -o ./build/$name ./vm

# Build the CLI binary
cli_name="cli_binary"
echo "Building CLI in ./build/$cli_name"
go build -o ./build/$cli_name ./cli

# Optional: Build any other specific binaries or modules
# Uncomment or add as needed
# echo "Building additional components..."
# go build -o ./build/some_binary ./path/to/module

echo "Build completed successfully!"
