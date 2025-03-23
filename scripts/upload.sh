#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/.env

FILES=(
    ${DIST_DIR}/${PROJECT_NAME}-macos-x86_64.zip
    ${DIST_DIR}/${PROJECT_NAME}-macos-arm64.zip
    ${DIST_DIR}/${PROJECT_NAME}-linux-x86_64.zip
    ${DIST_DIR}/${PROJECT_NAME}-linux-arm64.zip
    ${DIST_DIR}/${PROJECT_NAME}-windows-x86_64.zip
    ${DIST_DIR}/digests.txt
)

TEMP_DIR=$(mktemp -d)

response=$(curl -X POST "https://api.github.com/repos/$GITHUB_USER/$PROJECT_NAME/releases" \
    -s -w "%{http_code}" -o "$TEMP_DIR/tmp.json" \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    -H 'Content-Type: application/json' \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    -d @- <<EOF
{
  "tag_name": "$VERSION"
}
EOF
)

if [[ "$response" -ne 201 ]]; then
    echo "Unable to create release:"
    cat "$TEMP_DIR/tmp.json"
    exit 1
fi

RELEASE_ID=$(jq -r '.id' "$TEMP_DIR/tmp.json")
rm "$TEMP_DIR/tmp.json"

if [[ -z "$RELEASE_ID" ]]; then
    echo "RELEASE_ID not found"
    exit 1
fi


for file in "${FILES[@]}"; do
    echo "Uploading $file..."
    
    response=$(curl -X POST "https://uploads.github.com/repos/$GITHUB_USER/$PROJECT_NAME/releases/$RELEASE_ID/assets?name=$(basename "$file")" \
        -s -w "%{http_code}" -o "$TEMP_DIR/tmp.json" -L \
        -H "Authorization: Bearer $GITHUB_TOKEN" \
        -H "Accept: application/vnd.github+json" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        -H "Content-Type: application/octet-stream" \
        --data-binary @"$file")

    if [[ "$response" -ne 201 ]]; then
        echo "Unable to upload: $file:"
        cat "$TEMP_DIR/tmp.json"
        exit 1
    fi

    rm "$TEMP_DIR/tmp.json"
    echo "$file successfully uploaded"
done

echo "Done!"
