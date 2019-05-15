
function replaceCreds {
    CREDS=./creds.json
    rm $CREDS 2> /dev/null
    cat ./creds.sav | sed 's~GITHUB_USER_NAME_PLACEHOLDER~'"$GITU"'~' | \
      sed 's~PERSONAL_ACCESS_TOKEN_PLACEHOLDER~'"$GITAT"'~' | \
      sed 's~GITHUB_USER_EMAIL_PLACEHOLDER~'"$GITE"'~' | \
      sed 's~CLUSTER_NAME_PLACEHOLDER~'"$CLN"'~' | \
      sed 's~CLUSTER_ZONE_PLACEHOLDER~'"$CLZ"'~' | \
      sed 's~GKE_PROJECT_PLACEHOLDER~'"$PROJ"'~' | \
      sed 's~GITHUB_ORG_PLACEHOLDER~'"$GITO"'~' >> $CREDS
}