![Generic badge](https://github.com/konoui/alfred-bookmarks/workflows/test/badge.svg)
## alfred bookmarks
The workflow is a cross-browser bookmark searcher across Firefox, Google Chrome and Safari.

## Install
Download the workflow from [latest release](https://github.com/konoui/alfred-bookmarks/releases).

## Customize Configuration
Please create configuration file as `.alfred-bookmarks` in home directory (`~/`) if you want to customize.
```
firefox:
    enable: true
chrome:
    enable: true
safari:
    enable: false
remove_duplicate: true
```

If the configuration file does not exists, the workflow try to use available bookmark files of web browsers.

## Usage
Please type `bs <query>` and select your desired bookmark to open on default web browser.

## Feature
Supports fuzzy search.   
Supports following web browsers.
- Firefox
- Google Chrome
- Safari

## Limitation
### Firefox 
The workflow reads latest bookmark data from `~/Library/Application Support/Firefox/Profiles/<xxxxx>.default/bookmarkbackups/` directory.
If you register a web site to bookmarks, the workflow does not read and search the web site immediately.

## License
MIT License.
