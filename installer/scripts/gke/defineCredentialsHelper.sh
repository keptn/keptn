

# Helper function to replace place holder in creds.json.
function replaceCreds {
    CREDS=creds.json
    rm $CREDS 2> /dev/null
    cat ./gke/creds.sav | sed 's~GITHUB_USER_NAME_PLACEHOLDER~'"$GITU"'~' | \
      sed 's~PERSONAL_ACCESS_TOKEN_PLACEHOLDER~'"$GITAT"'~' | \
      sed 's~CLUSTER_NAME_PLACEHOLDER~'"$CLN"'~' | \
      sed 's~CLUSTER_ZONE_PLACEHOLDER~'"$CLZ"'~' | \
      sed 's~GKE_PROJECT_PLACEHOLDER~'"$PROJ"'~' | \
      sed 's~GITHUB_ORG_PLACEHOLDER~'"$GITO"'~' >> $CREDS
}