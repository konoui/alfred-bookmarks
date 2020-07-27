#!/bin/bash -eu

PROFILE=default
FIREFOX_DIR="${HOME}/Library/Application Support/Firefox/Profiles/xxxxx.${PROFILE}/bookmarkbackups"
CHROME_DIR="${HOME}/Library/Application Support/Google/Chrome/${PROFILE}"
CHROME_BOOKMARK_FILE="Bookmarks"
TEST_DIR="$(pwd)/pkg/bookmarker/testdata"
FIREFOX_TEST_FILE="${TEST_DIR}/test-firefox-bookmarks.jsonlz4"
CHROME_TEST_FILE="${TEST_DIR}/test-chrome-bookmarks.json"

mkdir -p "${FIREFOX_DIR}"
mkdir -p "${CHROME_DIR}"

cp "${FIREFOX_TEST_FILE}" "${FIREFOX_DIR}/"
cp "${CHROME_TEST_FILE}" "${CHROME_DIR}/${CHROME_BOOKMARK_FILE}"
