#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/.env

ARCHS=(
    amd64
)

for ARCH in ${ARCHS[@]}; do
  if [[ "$ARCH" == "amd64" ]]; then
    OUTPUT_DIR="${DIST_DIR}/windows-x86_64"
    ARCHIVE_FILE=${DIST_DIR}/${PROJECT_NAME}-windows-x86_64
  else
    OUTPUT_DIR="${DIST_DIR}/windows-${ARCH}"
    ARCHIVE_FILE=${DIST_DIR}/${PROJECT_NAME}-windows-${ARCH}
  fi

  # Crea la directory per l'architettura specifica
  mkdir -p "$OUTPUT_DIR"

  # Imposta le variabili d'ambiente
  CGO_ENABLED=0 GOOS=windows GOARCH=$ARCH \
    go build -ldflags="-s -w -X main.Version=$VERSION -X main.Build=${VERSION} -a -extldflags '-static'" \
    -o "${OUTPUT_DIR}/${PROJECT_NAME}.exe"

  # Compress con UPX
  upx --best --lzma --brute "${OUTPUT_DIR}/${PROJECT_NAME}.exe"

  # Costruisci la lista di file esistenti da aggiungere all'archivio
  FILES_TO_ARCHIVE=("${OUTPUT_DIR}/${PROJECT_NAME}.exe")
  for FILE in "${EXTRA_FILES[@]}"; do
    if [[ -f "$FILE" ]]; then
      FILES_TO_ARCHIVE+=("$FILE")
    fi
  done

  zip -j -r "${ARCHIVE_FILE}.zip" "${FILES_TO_ARCHIVE[@]}"

  # Digest
  shasum -a 256 "${ARCHIVE_FILE}.zip" >> ${DIST_DIR}/digests.txt

done






# Crea un archivio zip per Windows
#zip "${OUTPUT_FILE%.exe}.zip" "$OUTPUT_FILE" README_it.txt HOWTO_it.txt README_en.txt HOWTO_en.txt

