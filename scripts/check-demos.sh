#!/bin/sh

set -eu

video_extensions="mov mp4 webm avi mkv"
allowed_video_hosts="https://user-images.githubusercontent.com/ https://github.com/user-attachments/"

report="$(mktemp)"
trap 'rm -f "$report"' EXIT

file_ext() {
  printf '%s\n' "$1" | sed -E 's/.*\.([A-Za-z0-9]+)$/\1/' | tr '[:upper:]' '[:lower:]'
}

is_video_ext() {
  ext="$1"
  for vext in $video_extensions; do
    [ "$ext" = "$vext" ] && return 0
  done
  return 1
}

check_video_tags() {
  file="$1"
  matches="$(grep -noE '<video[^>]*src="[^"]*"' "$file" || true)"
  [ -z "$matches" ] && return 0
  printf '%s\n' "$matches" | while IFS=: read -r line rest; do
    src="$(printf '%s' "$rest" | sed -E 's/.*src="([^"]*)".*/\1/')"
    allowed=0
    for host in $allowed_video_hosts; do
      case "$src" in
      "$host"*) allowed=1 ;;
      esac
    done
    if [ "$allowed" -eq 0 ]; then
      echo "$file:$line: <video> src '$src' won't render on GitHub (not hosted on GitHub's asset CDN)" >>"$report"
    fi
  done
}

check_image_refs() {
  file="$1"
  matches="$(grep -noE '!\[[^]]*\]\([^)]+\)|<img[^>]*src="[^"]*"' "$file" || true)"
  [ -z "$matches" ] && return 0
  printf '%s\n' "$matches" | while IFS=: read -r line rest; do
    case "$rest" in
    '!['*) path="$(printf '%s' "$rest" | sed -E 's/.*\(([^)]+)\).*/\1/')" ;;
    *) path="$(printf '%s' "$rest" | sed -E 's/.*src="([^"]*)".*/\1/')" ;;
    esac
    ext="$(file_ext "$path")"
    if is_video_ext "$ext"; then
      echo "$file:$line: image reference points at '$path', a video file — <img> can't render video" >>"$report"
    fi
  done
}

readmes="$(git ls-files | grep -E '(^|/)README\.md$')"

for file in $readmes; do
  check_video_tags "$file"
  check_image_refs "$file"
done

if [ -s "$report" ]; then
  echo "Demo media references that won't render on GitHub:"
  echo
  cat "$report"
  echo
  echo "Convert the recording to a GIF and embed it as a markdown image instead:"
  echo "  ffmpeg -i <recording> -vf \"fps=12,scale=960:-1:flags=lanczos,palettegen\" palette.png"
  echo "  ffmpeg -i <recording> -i palette.png -lavfi \"fps=12,scale=960:-1:flags=lanczos[x];[x][1:v]paletteuse\" demo.gif"
  echo "  rm palette.png"
  echo "Then reference it with: ![demo](path/to/demo.gif)"
  exit 1
fi
