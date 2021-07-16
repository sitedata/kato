#!/bin/sh
pull_number=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")

URL="https://api.github.com/repos/gridworkz/kato/pulls/${pull_number}/files"

# Request GitHub api interface to parse changed files
# Used jq here. Some files are filtered
CHANGED_MARKDOWN_FILES=$(curl -s -X GET -G $URL | jq -r '.[] | select(.status != "removed") | select(.filename | endswith(".go")) | .filename')
for file in ${CHANGED_MARKDOWN_FILES}; do
  echo "golint ${file}"
  golint -set_exit_status=true ${file} || exit 1
done

echo "code golint check success"
