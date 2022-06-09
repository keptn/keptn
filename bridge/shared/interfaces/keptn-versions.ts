export interface VersionInfo {
  stable: string[];
  prerelease: string[];
}

export interface KeptnVersions {
  cli: VersionInfo;
  bridge: VersionInfo;
  keptn: {
    stable: {
      version: string;
      upgradableVersions: string[];
    }[];
  };
}
