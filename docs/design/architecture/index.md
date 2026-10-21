# Architecture

GoGoogle is a set of independent packages over the official Google API Go
Client. Each package is named `<name>util` and carries its API version in the
import path, for example `gmailutil/v1` or `searchconsoleutil/v1`.

- Each package takes an authorized `*http.Client` or builds one from
  credentials, so the caller chooses how people or services sign in.
- Packages wrap the generated client in types and functions for common tasks
  and return errors that carry the underlying Google error.
- Packages do not depend on each other.

Package guides are under [Guides](../../guides/getting-started.md).
