# cloud-spanner

Authors :

+ Matthieu Burguburu
+ David Dupraz
+ Clarisse Fleurimont
+ Alexandre Piveteau
+ Guy-Laurent Subri

## Setting it up locally

1. Run `git clone git@github.com:heig-CLD/cloud-spanner.git` or `git clone https://github.com/heig-CLD/cloud-spanner.git` to download our code.
2. Run `cd cloud-spanner`
3. To run the server, run `go run main.go server`
4. To run the client, run `go run main.go client`

The server will use your local [`gcloud`](https://cloud.google.com/sdk/gcloud/reference/config/) configuration to
connect to Cloud Spanner. Make sure to update your Google Cloud project configuration in `LocalConfig()` inside
`shared/shared.go`.

Additionally, it's also possible to use the Cloud Spanner emulator to run the server and client locally. After having
started the emulator with `emulator.sh`, make sure to configure your environment locally with `config_env.sh` and
set up the database schema with `config_db.sh --no-checks`.

## References

+ [1] Getting started with Google Cloud : https://cloud.google.com/spanner/docs/getting-started/gcloud
+ [2] `gcloud` reference : https://cloud.google.com/sdk/gcloud/reference/config/
