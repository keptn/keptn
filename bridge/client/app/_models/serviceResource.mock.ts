const resources = [
  {"resourceURI": "/helm/carts/Chart.yaml", "stageName": "dev"},
  {"resourceURI": "/helm/carts/templates/deployment.yaml", "stageName": "dev"},
  {"resourceURI": "/helm/carts/templates/service.yaml", "stageName": "dev"},
  {"resourceURI": "/helm/carts/values.yaml", "stageName": "dev"},
  {"resourceURI": "/metadata.yaml", "stageName": "dev"},
  {"resourceURI": "/helm/carts/Chart.yaml", "stageName": "staging"},
  {"resourceURI": "/helm/carts/templates/deployment.yaml", "stageName": "staging"},
  {"resourceURI": "/helm/carts/templates/service.yaml", "stageName": "staging"},
  {"resourceURI": "/helm/carts/values.yaml", "stageName": "staging"},
  {"resourceURI": "/metadata.yaml", "stageName": "staging"},
  {"resourceURI": "/helm/carts/Chart.yaml", "stageName": "production"},
  {"resourceURI": "/helm/carts/templates/deployment.yaml", "stageName": "production"},
  {"resourceURI": "/helm/carts/templates/service.yaml", "stageName": "production"},
  {"resourceURI": "/helm/carts/values.yaml", "stageName": "production"},
  {"resourceURI": "/metadata.yaml", "stageName": "production"},
]

const ServiceResourceMock = JSON.parse(JSON.stringify(resources));
export { ServiceResourceMock };
