## Customization options

The env var `LOGGING_COMPONENTS` allows to enable debug logs for specific components. The value for this variable is a list separated by comma of `Component=bool`, where `Component` is the component we want to configure and `bool` is a boolean value.

The following components are currently supported:

- `API`
- `APIService`
- `App`
- `OAuth`
