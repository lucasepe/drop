#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

DIST_DIR=dist
rm -rf ${DIST_DIR}

go mod tidy

${SCRIPT_DIR}/build-darwin.sh
${SCRIPT_DIR}/build-linux.sh
${SCRIPT_DIR}/build-windows.sh
