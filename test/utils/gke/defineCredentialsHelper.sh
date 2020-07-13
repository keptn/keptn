

# Helper function to replace place holder in creds.json.
function replaceCreds {
    CREDS=creds.json
    rm $CREDS 2> /dev/null
    cat ./creds.sav | sed 's~CLUSTER_NAME_PLACEHOLDER~'"$CLN"'~' | \
      sed 's~CLUSTER_ZONE_PLACEHOLDER~'"$CLZ"'~' | \
      sed 's~GKE_PROJECT_PLACEHOLDER~'"$PROJ"'~' >> $CREDS
}
