const fileTree = [
  {
    "stageName": "dev",
    "tree": [
      {
        "fileName": "helm",
        "children": [
          {
            "fileName": "carts",
            "children": [
              {
                "fileName": "templates",
                "children": [
                  {
                    "fileName": "deployment.yaml",
                  },
                  {
                    "fileName": "service.yaml",
                  },
                ],
              },
              {
                "fileName": "Chart.yaml",
              },
              {
                "fileName": "values.yaml",
              },
            ],
          },
        ],
      },
      {
        "fileName": "metadata.yaml",
      },
    ],
  },
  {
    "stageName": "staging",
    "tree": [
      {
        "fileName": "helm",
        "children": [
          {
            "fileName": "carts",
            "children": [
              {
                "fileName": "templates",
                "children": [
                  {
                    "fileName": "deployment.yaml",
                  },
                  {
                    "fileName": "service.yaml",
                  },
                ],
              },
              {
                "fileName": "Chart.yaml",
              },
              {
                "fileName": "values.yaml",
              },
            ],
          },
        ],
      },
      {
        "fileName": "metadata.yaml",
      },
    ],
  },
  {
    "stageName": "production",
    "tree": [
      {
        "fileName": "helm",
        "children": [
          {
            "fileName": "carts",
            "children": [
              {
                "fileName": "templates",
                "children": [
                  {
                    "fileName": "deployment.yaml",
                  },
                  {
                    "fileName": "service.yaml",
                  },
                ],
              },
              {
                "fileName": "Chart.yaml",
              },
              {
                "fileName": "values.yaml",
              },
            ],
          },
        ],
      },
      {
        "fileName": "metadata.yaml",
      },
    ],
  },
]

const FileTreeMock = JSON.parse(JSON.stringify(fileTree));
export { FileTreeMock };
