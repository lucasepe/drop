#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/.env

ARCHS=(
    amd64
    arm64
)

# Costruzione per linux architetture amd64 e arm64
for ARCH in ${ARCHS[@]}; do
  if [[ "$ARCH" == "amd64" ]]; then
    OUTPUT_DIR="${DIST_DIR}/linux-x86_64"
    ARCHIVE_FILE=${DIST_DIR}/${PROJECT_NAME}-linux-x86_64
  else
    OUTPUT_DIR="${DIST_DIR}/linux-${ARCH}"
    ARCHIVE_FILE=${DIST_DIR}/${PROJECT_NAME}-linux-${ARCH}
  fi

  # Crea la directory per l'architettura specifica
  mkdir -p "$OUTPUT_DIR"

  # Imposta le variabili d'ambiente
  CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH \
    go build -ldflags="-s -w -X main.Version=$VERSION -X main.Build=${VERSION} -a -extldflags '-static'" \
    -o "${OUTPUT_DIR}/${PROJECT_NAME}"
  
  #if [[ "$ARCH" == "amd64" ]]; then
  #  strip -arch x86_64 "${OUTPUT_DIR}/${PROJECT_NAME}"
  #else
  #  strip -arch ${ARCH} "${OUTPUT_DIR}/${PROJECT_NAME}"
  #fi

  upx --best --lzma --brute "${OUTPUT_DIR}/${PROJECT_NAME}"

  # Costruisci la lista di file esistenti da aggiungere all'archivio
  FILES_TO_ARCHIVE=("${OUTPUT_DIR}/${PROJECT_NAME}")
  for FILE in "${EXTRA_FILES[@]}"; do
    if [[ -f "$FILE" ]]; then
      FILES_TO_ARCHIVE+=("$FILE")
    fi
  done
  

  zip -j -r "${ARCHIVE_FILE}.zip" "${FILES_TO_ARCHIVE[@]}"

  shasum -a 256 "${ARCHIVE_FILE}.zip" >> ${DIST_DIR}/digests.txt

  echo "Built and archived ${ARCHIVE_FILE}.zip for linux ($ARCH)"

done
