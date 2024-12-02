#!/usr/bin/env bash
# Custom Run Script for Your Project

set -o errexit
set -o nounset
set -o pipefail

if ! [[ "$0" =~ scripts/run.sh ]]; then
  echo "This script must be run from the repository root"
  exit 255
fi

VERSION=0.1.0
MODE=${MODE:-run}
LOGLEVEL=${LOGLEVEL:-info}

echo "Running with:"
echo "VERSION: ${VERSION}"
echo "MODE: ${MODE}"

BUILD_DIR=/tmp/project-${VERSION}
VM_BINARY=${BUILD_DIR}/vm
CLI_BINARY=${BUILD_DIR}/cli
GENESIS_FILE=${BUILD_DIR}/genesis.json
VM_CONFIG_FILE=${BUILD_DIR}/vm.config

############################
# Build VM and CLI binaries
############################
mkdir -p ${BUILD_DIR}

echo "Building VM..."
go build -o ${VM_BINARY} ./vm

echo "Building CLI..."
go build -o ${CLI_BINARY} ./cli

############################
# Create Genesis File
############################
echo "Creating genesis file..."
cat <<EOF > ${GENESIS_FILE}
{
  "allocations": [
    {"address": "addr1xyz", "balance": 100000000}
  ]
}
EOF

############################
# Create VM Configuration File
############################
echo "Creating VM configuration file..."
cat <<EOF > ${VM_CONFIG_FILE}
{
  "mempoolSize": 100000,
  "logLevel": "${LOGLEVEL}",
  "parallelism": 4
}
EOF

############################
# Run VM
############################
if [[ ${MODE} == "run" ]]; then
  echo "Starting VM..."
  ${VM_BINARY} --genesis ${GENESIS_FILE} --config ${VM_CONFIG_FILE}
elif [[ ${MODE} == "test" ]]; then
  echo "Running tests..."
  go test ./... -v
else
  echo "Invalid mode. Use 'run' or 'test'."
  exit 1
fi
