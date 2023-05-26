# mercury

⚠ Only works on Windows!

A simple executable for converting docx documents to Markdown and uploading them to an Astro website.

This tool gives users access to the powerful git source control toolbox in a transparent and seamless manner. Non technically savy authors will only have to open the "mercury.exe" executable to automatically sync their work with their Astro website.

This was built for the website of Dr. [Paolo Sernini](https://github.com/ornato-t/paolo-sernini).

## Requirements (for users)
- git (credentials or SSH keys are not required)

## How it works
The app automatically extracts a GitHub token (contained in a "env" file bundled along with the exe and not included in this repo) and uses it to authenticate with GitHub. The token needs to have access to the "repo" scope for the repository containing the Astro website.

Documents are converted from .docx to Markdown with [Pandoc](https://pandoc.org/)

## Requirements (devs)
Developers looking to adopt this tool will have to compile and bundle it on their own.

In order to work the file requires:
- a file named `env` containing the aforementioned GitHub token
- an executable for pandoc. The latest release of the [pandoc-portable](https://github.com/pandoc-extras/pandoc-portable/releases) project is the reccomended option

These files should be placed in the root directory. Once this is done run `go build -o mercury.exe`.

## Optional
Before building the app, it is possible to add meta data and an icon with go-winres. Install it with `go install github.com/tc-hib/go-winres@latest`. Then run `go-winres make` and lastly `go build -o mercury.exe`.

## Conclusion
The resulting executable will include the GitHub token in its bundle. Because of this, for security reasons, there are no releases available. Users are invited to build their own version of the tool, with their own tokens.