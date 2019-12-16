## alfred-bookmarks
Alfred workflow to search bookmarks in Firefox and Google Chrome.

## Install
Download the workflow form [latest release](https://github.com/konoui/alfred-bookmarks/releases).

## Configuration
Please create configuration file as `.alfred-bookmarks` in home directory (`~/`).
```
firefox:
    enable: true
chrome:
    enable: true
remove_duplicate: true
```

## Usage
Please type `bs <query>` and select your desired bookmark to open on default web browser.

## Feature
Supports fuzzy search.   
Supports following web browsers.
- Firefox
- Google Chrome

## Limitation
### Firefox 
The workflow reads latest bookmark data from `~/Library/Application Support/Firefox/Profiles/<xxxxx>.default/bookmarkbackups/` directory.
If you register a web site to bookmarks, the workflow does not read and search the web site immediately.
