import { KeptnInfo as ki } from '../../_models/keptn-info';

export const KeptnInfo = {
  "availableVersions": {
    "cli": {
      "stable": [
        "0.7.0",
        "0.7.1",
        "0.7.2",
        "0.7.3",
        "0.8.0",
        "0.8.1"
      ],
      "prerelease": []
    },
    "bridge": {
      "stable": [
        "0.7.0",
        "0.7.1",
        "0.7.2",
        "0.7.3",
        "0.8.0",
        "0.8.1"
      ],
      "prerelease": []
    },
    "keptn": {
      "stable": [
        {
          "version": "0.8.1",
          "upgradableVersions": [
            "0.8.0"
          ]
        },
        {
          "version": "0.8.0",
          "upgradableVersions": [
            "0.7.1",
            "0.7.2",
            "0.7.3"
          ]
        },
        {
          "version": "0.7.3",
          "upgradableVersions": [
            "0.7.0",
            "0.7.1",
            "0.7.2"
          ]
        },
        {
          "version": "0.7.2",
          "upgradableVersions": [
            "0.7.0",
            "0.7.1"
          ]
        },
        {
          "version": "0.7.1",
          "upgradableVersions": [
            "0.7.0"
          ]
        }
      ]
    }
  },
  "bridgeInfo": {
    "apiUrl": "http://localhost:3000/api",
    "apiToken": "TOKEN",
    "cliDownloadLink": "https://github.com/keptn/keptn/releases",
    "enableVersionCheckFeature": true,
    "showApiToken": true,
    "authType": "BASIC"
  },
  "authCommand": "keptn auth --endpoint=http://localhost:3000/api --api-token=TOKEN",
  "keptnVersion": "0.8.1-dev-PR-3529",
  "versionCheckEnabled": true,
  "metadata": {
    "bridgeversion": "docker.io/keptn/bridge2:0.8.1-dev-PR-3611.202103251226",
    "keptnlabel": "keptn",
    "keptnversion": "0.8.1-dev-PR-3529",
    "shipyardversion": "0.2.0",
    "namespace": "keptn"
  },
  "keptnVersionInvalid": false,
  "bridgeVersionInvalid": true
} as ki;
