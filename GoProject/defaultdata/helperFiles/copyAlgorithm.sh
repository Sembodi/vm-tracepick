#!/bin/bash

SOURCE=$1
DEST=$2
PROTECTED_PATHS_FILE=$3

mapfile -t protectedFolders < "$PROTECTED_PATHS_FILE"

cd "$SOURCE" || exit 1

for item in *; do
  protected=false
  for protectedPath in "${protectedFolders[@]}"; do
    if [[ "$DEST/$item" == "$protectedPath" ]]; then
      protected=true
      break
    fi
  done

  if [[ -d "$item" && "$protected" == true ]]; then
    # # DEBUG:
    # echo "Protected: merging contents of $item"
    # sleep 3
    /helperFiles/copyAlgorithm.sh "$SOURCE/$item" "$DEST/$item" $PROTECTED_PATHS_FILE
  elif [[ -e "$SOURCE/$item" ]]; then
    # # DEBUG:
    # echo "Replacing or copying $item"
    # sleep 3
    rm -r "$DEST/$item" && cp -r "$SOURCE/$item" "$DEST/$item" || cp -r "$SOURCE/$item" "$DEST/$item" # && echo "rm failed, source copied to dest"
  fi
done
