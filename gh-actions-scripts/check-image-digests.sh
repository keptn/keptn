#!/usr/bin/env bash
DIGEST_FILE="digests.csv"

readarray -t tags <<< "$(gh release list --limit 1000 | awk '{print $1}')"

for tag in "${tags[@]}"; do
  readarray -t release_assets <<< "$(gh release view "$tag" --json assets --jq ".assets[].name")"

  echo "Checking if release $tag has docker image digests attached..."
  if [[ "${release_assets[*]}" =~ $DIGEST_FILE ]]; then
    echo "Release $tag has docker image digests attached. Downloading digests file..."
    gh release download --pattern="$DIGEST_FILE"
    echo "Digests from release:"
    cat "$DIGEST_FILE"

    while IFS=, read -r artifact release_digest; do
      echo "Pulling docker image $artifact:$tag..."
      docker pull "keptn/$artifact:$tag" --quiet
      docker_digest=$(docker inspect "keptn/$artifact:$tag" | jq -r '.[0].RepoDigests[0]' | cut -d'@' -f2)

      echo "Checking image digest..."
      if [[ "$docker_digest" != "$release_digest" ]]; then
        echo "CAUTION: Docker image digest does not match released digest!"
        echo "Digest from dockerhub: $docker_digest"
        echo "Digest from release:   $release_digest"
        exit 1
      else
        echo "Image digest for $artifact:$tag matches released digest. Continuing..."
      fi
    done < $DIGEST_FILE
    rm "$DIGEST_FILE"
  fi
done
