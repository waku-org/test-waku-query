
## Summary

The [jmeter](https://jmeter.apache.org/) is used to force concurrent Store
requests.

_jmeter_ launches concurrent users that perform http requests to two _Store_
node clients. And the _Test Plan_ is defined in _http_store_requests.jmx_.

_jmeter_ is thought to be used in combination of the `docker compose up`.

## jmeter

### Prerequisite
- `docker compose up` to be running.
- Java should be installed (e.g. "1.8.0_382")

### How to use jmeter
- From CLI and Linux, go to `apache-jmeter-5.6.2/bin` and run `bash jmeter.sh`.
- Load the _http_store_requests.jmx_ test plan file.
- Hit the "play" button and check the "Summary Report", for example.

