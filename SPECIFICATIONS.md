# Main
The application will be a dockerized http server that hits configured external
endpoints to ensure GDPR compliance. The current GDPR scope:

- Removing records for churned customers.
...

#### Integrations
- Secrets will be stored and accessed as environment variables
and compiled into the binary.
- Each integration will likely require a bespoke solution for secure connectivity.

#### Storage
There are my current thoughts on the different storage needs.
1. Record database that holds all transactions.

The record database will be a SQL database.
- Possibly provide options for a desired third pary service/self hosting/docker volume.

2. Configurable application settings.

The application settings can be stored and read from a docker volume (JSON).
- Need to give customers the ability to export/backup settings.


#### Authorization/Authentication
- Simple password solution.
- Sessions vs JWTs to be decided


