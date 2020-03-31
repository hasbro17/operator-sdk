## operator-sdk generate csv

Generates a ClusterServiceVersion YAML file for the operator

### Synopsis

The 'generate csv' command generates a ClusterServiceVersion (CSV) YAML manifest
for the operator. This file is used to publish the operator to the OLM Catalog.

A CSV semantic version is supplied via the --csv-version flag. If your operator
has already generated a CSV manifest you want to use as a base, supply its
version to --from-version. Otherwise the SDK will scaffold a new CSV manifest.

```
operator-sdk generate csv [flags]
```

### Options

```
      --csv-channel string      Channel the CSV should be registered under in the package manifest
      --csv-version string      Semantic version of the CSV
      --default-channel         Use the channel passed to --csv-channel as the package manifests' default channel. Only valid when --csv-channel is set
      --from-version string     Semantic version of an existing CSV to use as a base
  -h, --help                    help for csv
      --inputs stringToString   Key value input paths used in CSV generation.
                                Use this to set custom paths for operator manifests and API type definitions
                                E.g: --inputs deploy=config/production,apis=pkg/myapp/apis 
                                Supported input keys:
                                	- deploy=<project-relative path to root directory for operator manifests (Deployment, RBAC, CRDs)>
                                	- apis=<project-relative path to root directory for API type defintions>
                                 (default [apis=pkg/apis,deploy=deploy])
      --operator-name string    Operator name to use while generating CSV
      --output-dir string       Base directory to output generated CSV. The resulting CSV bundle directorywill be "<output-dir>/olm-catalog/<operator-name>/<csv-version>" (default "deploy")
      --update-crds             Update CRD manifests in deploy/{operator-name}/{csv-version} the using latest API's
```

### SEE ALSO

* [operator-sdk generate](operator-sdk_generate.md)	 - Invokes a specific generator

