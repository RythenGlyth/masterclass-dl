# masterclass-dl

A command line tool to download masterclass.com classes.

## Features

- **Download entire categories** - Bulk download all classes from any category page
- **Plex-ready output** - Poster (`poster.jpg`) and fanart (`fanart.jpg`) images for media servers
- **Embedded subtitles** - 10+ languages automatically embedded in videos
- **Subtitles-only mode** - Download just subtitles without video files
- **PDF workbooks** - Class guides and supplementary materials
- **Flexible options** - Control what gets downloaded with `--pdfs`, `--posters`, `--limit`
- **SSO support** - Login via Safari session for Google/Apple/company SSO accounts (macOS)

## Prerequisites

- yt-dlp (recommended) or youtube-dl
- ffmpeg

## Installation

Download the executable from the [releases](https://github.com/RythenGlyth/masterclass-dl/releases) page.

Or build from source:
```bash
go build -o masterclass-dl .
```

## Usage

### Login

First, authenticate with your Masterclass account:

```bash
masterclass-dl login <email> <password>
```

You'll be prompted to select a profile if your account has multiple profiles.

### SSO Login (macOS only)

If you sign in to MasterClass via **Google, Apple, or company SSO** (and don't have a password), use `safari-login` instead. It reads your session cookies directly from Safari — no password needed.

**Requirements:**
- macOS only
- Must already be logged in to [masterclass.com](https://www.masterclass.com) in Safari
- Terminal may need **Full Disk Access** in System Settings → Privacy & Security → Full Disk Access

```bash
masterclass-dl safari-login
```

You'll be prompted to select a profile if your account has multiple profiles.

> **Note:** On Windows and Linux, use `masterclass-dl login <email> <password>` instead.

### Check Status

Verify your login and subscription status:

```bash
masterclass-dl status
```

### Download

Download individual classes, specific chapters, or **bulk download entire categories**:

```bash
# Download a single class
masterclass-dl download -o ./downloads "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking"

# Download a specific chapter
masterclass-dl download -o ./downloads "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking/chapters/introduction"

# Download all classes in a category (first 10 by default)
masterclass-dl download -o ./downloads "https://www.masterclass.com/homepage/science-and-tech"

# Download all classes in a category (no limit)
masterclass-dl download -o ./downloads --limit 0 "https://www.masterclass.com/homepage/science-and-tech"
```

#### Download Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | (required) | Output directory |
| `--limit` | `-l` | 10 | Max classes to download from a category (0 = unlimited) |
| `--pdfs` | `-p` | true | Download PDF workbooks |
| `--posters` | | true | Download poster and fanart images |
| `--ytdl-exec` | `-y` | yt-dlp | Path to yt-dlp/youtube-dl executable |
| `--name-files-as-series` | | false | Name files in TV series format for Plex/Jellyfin |
| `--write-nfo` | | false | Generate NFO metadata files alongside downloads |
| `--metadata-only` | | false | Download only metadata (no videos or PDFs) |
| `--subs-only` | `-s` | false | Download only subtitles (no video) |
| `--force` | | false | Force re-download and overwrite existing files |
| `--concurrency` | | 1 | Number of concurrent downloads to run |

#### Examples

```bash
# Download with Plex-friendly TV series naming (s01e01-Title.mp4)
masterclass-dl download -o ./downloads --name-files-as-series "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking"

# Download without PDFs
masterclass-dl download -o ./downloads --pdfs=false "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking"

# Download without poster images
masterclass-dl download -o ./downloads --posters=false "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking"

# Download subtitles only (no video)
masterclass-dl download -o ./downloads --subs-only "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking"

# Download first 5 classes from a category
masterclass-dl download -o ./downloads --limit 5 "https://www.masterclass.com/homepage/science-and-tech"

# Use a specific yt-dlp path
masterclass-dl download -o ./downloads -y /usr/local/bin/yt-dlp "https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking"
```

### Output Structure

Downloads are organized in a Plex/Jellyfin-friendly format:

**Default naming:**
```
downloads/
├── Gordon Ramsay Teaches Cooking/
│   ├── poster.jpg
│   ├── fanart.jpg
│   ├── Class Guide.pdf
│   ├── 001-Introduction.mp4
│   ├── 002-Knives.mp4
│   └── ...
```

**With `--name-files-as-series`:**
```
downloads/
├── Gordon Ramsay Teaches Cooking/
│   ├── poster.jpg
│   ├── fanart.jpg
│   ├── Class Guide.pdf
│   ├── s01e01-Introduction.mp4
│   ├── s01e02-Knives.mp4
│   ├── s01e15-Closing-Extra_trailer.mp4  # Example lessons marked as extras
│   └── ...
```

**Video features:**
- Embedded subtitles in 10+ languages (English, Spanish, French, German, Italian, Japanese, Chinese, Hindi, Polish, Portuguese)
- Best available video/audio quality
- Metadata embedded (title, description, episode number)

**With `--name-files-as-series`, additional metadata:**
- `show` - Series/course title
- `artist` - Instructor name
- `genre` - Category (e.g., "Arts & Entertainment")
- `date` - Episode date
- `episode_id` - s01e01 format
- `network` - "MasterClass"
- `synopsis` - Course overview

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--datDir` | `-d` | ~/.masterclass/ | Directory for cookies and data |
| `--debug` | | false | Enable debug output |

## Commands Reference

```
masterclass-dl [command]

Available Commands:
  completion    Generate the autocompletion script for the specified shell
  download      Download a class, chapter, or category from masterclass.com
  help          Help about any command
  login         Login to masterclass.com
  safari-login  Login using Safari session cookies (macOS only, for SSO users)
  status        Check login status

Use "masterclass-dl [command] --help" for more information about a command.
```
