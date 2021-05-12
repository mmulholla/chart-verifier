# Chart Verifier

`chart-verifier` is a tool that verifies a Helm chart is compliant against a configurable list of checks. 

This tool can be used to help ensure the quality of Helm Charts, from its associated metadatas, formating and 
readiness for distribution. Additionaly, it helps ensure that a  Helm Chart will work seamlessly on Red Hat 
OpenShift and can be submitted as a certified Helm Chart in the [Red Hat Helm Repository](https://github.com/openshift-helm-charts).

## Features

- Helm Chart Verification: Verifies an Helm Chart is compliant with a certain set of checks.
- Red Hat Certified Chart Validation: Verifies an Helm Chart's readiness for being certified and submitted in [Red Hat Helm Repository](https://github.com/openshift-helm-charts).
- Report Generation: Generates a verification report in YAML format.
- Optionable Checks: Defines the checks you want to execute during the verification process.

## Existing Checks

The following checks have been implemented:

| Check | Type | Description
|---|---
| `is-helm-v3` | Mandatory | Checks whether the given `uri` is a Helm v3 chart.
| `has-readme` | Mandatory | Checks whether the Helm chart contains a `README.md` file.
| `contains-test` | Mandatory | Checks whether the Helm chart contains at least one test file.
| `has-minkubeversion` | Mandatory | Checks whether the Helm chart's `Chart.yaml` includes the `minKubeVersion` field.
| `contains-values-schema` | Mandatory | Checks whether the Helm chart contains a values schema.
| `not-contains-crds` | Mandatory | Checks whether the Helm chart does not include CRDs.
| `not-contain-csi-objects` | Mandatory | Checks whether the Helm chart does not include CSI objects.
| `images-are-certified` | Mandatory | Checks whether images referenced by the helm chart are Red Hat certified images.  
| `helm-lint` | Mandatory | Runs the helm lint command to check that the chart is wel formed.
| `contains-values` | Mandatory | Checks whether the Helm chart contains a values file.

Further the checks include installing the chart on an available cluster and running the chart tests. Information on this 
will be provided when this functionality will be added.

All current checks are of type "Mandatory". Mandatory indicates that a check is required for chart submission. 

## Usage

### Pre Requisities

- Docker installed.
- Internet Connection: The check that images are Red Hat Certified requires an internet connection.
- OCP Cluster. 

### Know before you start

- The default set of check covers Red Hat’s submission requirements.
- Each check is independent and execution order is not guaranteed. 

### Basic Usage with Docker

1. To run all checks for a chart available using a uri: 
   ```
   docker run -it --rm quay.io/redhat-certification/chart-verifier verify <chart-uri>
   ```
1. For a chart available locally on disk, from the same directory as the chart: 
   ```
   docker run -v $(pwd):/charts --rm quay.io/redhat-certification/chart-verifier verify /charts/<chart>
   ``` 
1. To get a list of options for the verify request:
   ```
   docker run -it --rm quay.io/redhat-certification/chart-verifier verify help
   ```
   This will produce the following output:
   ```
   Verifies a Helm chart by checking some of its characteristics

   Usage:
     chart-verifier verify <chart-uri> [flags]

   Flags:
      -S, --chart-set strings          set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
      -F, --chart-set-file strings     set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
      -X, --chart-set-string strings   set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
      -f, --chart-values strings       specify values in a YAML file or a URL (can specify multiple)
      -x, --disable strings            all checks will be enabled except the informed ones
      -e, --enable strings             only the informed checks will be enabled
      -h, --help                       help for verify
      -o, --output string              the output format: default, json or yaml
      -s, --set strings                overrides a configuration, e.g: dummy.ok=false

    Global Flags:
          --config string   config file (default is $HOME/.chart-verifier.yaml)
   ```

### Examples 

1. Run a subset of checks:
   ```
   chart-verifier verify -e images-are-certified,helm-lint
   ```
1. Run all checks except a subset:
   ```
   chart-verifier verify -x images-are-certified,helm-lint
   ```
1. Provide chart override values:   
   ```
   chart-verifier verify -S default.port=8080
   ```
1. Provide chart override values in a file:  
   ```
   chart-verifier verify -F overrides.yaml
   ```
   
### Notes on usage

The checks performed include running ```helm lint```, and ```helm template```(for red hat image certification) against 
the chart. As a result if the chart requires additional values for these to succeed the values must be specified using 
the options available. These options are similar to those use by ```helm lint``` and ```helm template```.

Note, for ``helm lint`` the check will pass if there are no error messages - warning and info messages do not cause the check to fail.

Running a subset of checks using the -e and -x flags are provided as a convenience when working to get all checks to pass 
for a chart. A report generated for the submission process must include all checks. 

## Submitting a Chart for inclusion in Red Hat Helm Repository and Certification

Information on the submission process is defined in the [Red Hat Helm Repository](https://github.com/openshift-helm-charts/repo)
Please read this information before submitting a chart.

## Submission notes: 

A verifier report is an integral part of the submission process. There are 3 options for submitting a chart:

| Option | Description
|---|---
| **1. Helm Chart Tarball** | Submit your Chart with its tarball (`chart-verifier`'s report optional).
| **2. Helm Chart extracted Tarball** | Submit your Chart with its extracted tarball (`chart-verifier`'s report optional).
| **3. chart-verifier Report only** | When your Chart will not be hosted in the Red Hat Helm repository, you can just submit the generated report from `chart-verifier` tool.

With the options which do not require and report, a report will be generated as part of the submission process. 


### Notes

If a report is not included it will be generated as part of the submission process. 

When a Chart is submitted a series of checks will be run against the associated Pull Request. The PR will fail
and an exception process will be started if the report contains one or more failures or is missing any mandatory 
tests. For more information on the submission process see: https://github.com/openshift-helm-charts/repo.

If the report is to be submitted without a chart, the report should be run against the chart in its final 
location. This is because the verifier will record the chart-uri specified when the report was run and, 
in the absence of a submitted chart, this uri will be used for publication.

If the report is submitted with a chart it must be run against the chart as submitted. So, for example, if submitting 
a tarball, run the report against the tarball that will be submitted. This is important because the report will calculate 
and record a sha256 value for the chart. The submission process will then re-generate the sha256 value and the process 
will fail if the sha values do not match.

If a successful run of the report requires additional values to be specified the report must be submitted with the chart.
This is because the submission process does not have access to the values and the report generated would inevitably include
failures.

## Trouble shooting check failures

### `is-helm-v3`

Requires the "api-version" attribute of chart.yaml to be set to "v2". Any other value will result in check fail.

### `has-readme`

Requires a "README.md" file to exist in the root directory of the chart. Any other spelling or 
capitialisation of letters with result if check failure.

### `contains-test`

Requires at least one file to exists in the ```templates/test``` subdirectory of the chart. If no such file 
exists this check will fail. Note other checks will require the directory to contain a valid test.

### `has-minkubeversion`

Requires the "kubeVersion" attribute of chart.yaml to be set to a value. If the attribute is not set the check 
will fail. The vaue set is not checked.

### `contains-values`

Requires a ```values.schema``` file to be present in the chart. If the file is not present the check will fail.

### `contains-values-schema`

Requires a ```values.schema.json``` file to be present in the chart. If the file is not present the check will fail. 

### `not-contains-crds`

Requires no crds to be defined in the chart. A crd is a file with an extension of `.yaml`, `.yml` or `.json`
in a `crd` subdirectory of the chart and should be removed if present.

### `not-contain-csi-objects`

Requires no csi objects in a chart. A csi object is a file in the template subdirectory, with an extension of `.yaml`,
and containing an `kind` attribute set to `CSIDriver`. If such a file exists it should be removed.

### `helm-lint`

Requires a `helm lint` of the chart to not result in any `ERROR` messages. If a ERROR does occur the helm lint messages 
will be output. Run `helm lint` on your chart for additional information. If the chart requires specification of additional 
attributes to pass `helm lint` use one of the `chart-set` flags of the verifier tool for this check to pass. If additional 
attributes are required a verifier report mut be included in the chart submission.

### `images-are-certified`

Requires any images referenced in a chart to be Red Hat Certified. 
- The list of image references is found by running `helm template` and if this fails the error output from `helm template` 
  will be output. Run `helm template` on your chart for additional information. If the chart requires specification of additional
  attributes to pass `helm template` use one of the `chart-set` flags of the verifier tool for this check to pass. If additional
  attributes are required a verifier report must be included in the chart submission. 
- Each image reference is then parsed to determine the registry, repository and tag or digest value.
    - registry is the string before the first "/" in the image reference but only if it includes a "." character.
    - the repository is what remains in the image reference, after the registry is removed and before ":" or "@sha" 
    - tag is what is set after the ":" character
    - digest is what is set after the "@" character in "@sha"
- If a registry is not found the pyxis swagger api is used to find the repository and from it, extract the registry
    - `https://catalog.redhat.com/api/containers/v1/repositories?filter=repository==<repository>`
    - if the repository is not found the check will fail.
- The registry and repository are the used to find images:
    - `https://catalog.redhat.com/api/containers/v1/repositories/registry/<registry>/repository/<repository>/images`
    - if the image specified a sha value it is compared with the `parsed_data.docker_image_digest` attribute. If a 
      match is not found the check fails.
    - if the image specified a tag value it is compared with the `repositories.tags.name` attributes. If a match is 
      not found the check fails.
- If the check fails use the point fo failure to determine how to address the issue. 
    


## Suggestions

If you have any suggestions for improving the verifier, for example additional checks to add, please open 
an issue in this repository.

# A deeper dive for developers 

## Architecture

This tool inspects a Helm
chart URI (`file://`, `https?://`, etc)
and returns either a *positive* result indicating the Helm chart has passed all checks, or a *negative* result indicating
which checks have failed and remedial actions.

The application is separated in two pieces: a command line interface and a library. This is handy because the command
line interface is specific to the user interface, and the library can be generic enough to be used to, for example,
inspect Helm chart bytes in flight.

One positive aspect of the command line interface specificity is that its output can be tailored to the methods of
consumption the user expects; in other words, the command line interface can be programmed in such way it can be
represented as either *YAML* or *JSON* formats, in addition to a descriptive representation tailored to human actors.

Primitive functions to manipulate the Helm chart should be provided, since most checks involve inspecting the contents
of the chart itself; for example, whether a `README.md` file exists, or whether `README.md` contains the `values`'
specification, implicating in offering a cache API layer is required to avoid downloading and unpacking the charts for
each test.


## Building chart-verifier

To build `chart-verifier` locally, please execute `hack/build.sh` or its PowerShell alternative.

To build `chart-verifier` container image, please execute `hack/build-image.sh` or its PowerShell alternative:

```text
PS C:\Users\igors\GolandProjects\chart-verifier> .\hack\build-image.ps1
[+] Building 15.1s (15/15) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                                                                                                                 0.0s
 => => transferring dockerfile: 32B                                                                                                                                                                                                                                  0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                                    0.0s
 => => transferring context: 2B                                                                                                                                                                                                                                      0.0s
 => [internal] load metadata for docker.io/library/fedora:31                                                                                                                                                                                                         1.4s
 => [internal] load metadata for docker.io/library/golang:1.15                                                                                                                                                                                                       1.3s
 => [build 1/7] FROM docker.io/library/golang:1.15@sha256:d141a8bca046ade2c96f89e864cd31f5d0ba88d5a71d62d59e0e1f2ecc2451f1                                                                                                                                           0.0s
 => CACHED [stage-1 1/2] FROM docker.io/library/fedora:31@sha256:ba4fe6a3da48addb248a16e8a63599cc5ff5250827e7232d2e3038279a0e467e                                                                                                                                    0.0s
 => [internal] load build context                                                                                                                                                                                                                                    0.5s
 => => transferring context: 43.06MB                                                                                                                                                                                                                                 0.5s
 => CACHED [build 2/7] WORKDIR /tmp/src                                                                                                                                                                                                                              0.0s
 => CACHED [build 3/7] COPY go.mod .                                                                                                                                                                                                                                 0.0s
 => CACHED [build 4/7] COPY go.sum .                                                                                                                                                                                                                                 0.0s
 => CACHED [build 5/7] RUN go mod download                                                                                                                                                                                                                           0.0s
 => [build 6/7] COPY . .                                                                                                                                                                                                                                             0.2s
 => [build 7/7] RUN ./hack/build.sh                                                                                                                                                                                                                                 12.5s
 => [stage-1 2/2] COPY --from=build /tmp/src/out/chart-verifier /app/chart-verifier                                                                                                                                                                                  0.1s
 => exporting to image                                                                                                                                                                                                                                               0.2s
 => => exporting layers                                                                                                                                                                                                                                              0.2s
 => => writing image sha256:7302e88a2805cb4be1b9e130d057bd167381e27f314cbe3c28fbc6cb7ee6f2a1                                                                                                                                                                         0.0s
 => => naming to quay.io/redhat-certification/chart-verifier:0d3706f
```

The container image created by the build program is tagged with the commit ID of the working directory at the time of
the build: `quay.io/redhat-certification/chart-verifier:0d3706f`.

## Running built images

### Local command

To verify a chart against all available checks, for exmaple:

```text
> out/chart-verifier verify ./chart.tgz
> out/chart-verifier verify ~/src/chart
> out/chart-verifier verify https://www.example.com/chart.tgz
```

To apply only the `is-helm-v3` check:

```text
> out/chart-verifier verify --enable is-helm-v3 https://www.example.com/chart.tgz
```

To apply all checks except `is-helm-v3`:

```text
> out/chart-verifier verify --disable is-helm-v3 https://www.example.com/chart.tgz
```

### Container Image

The container image produced in 'Building chart-verifier' can then be executed with the Docker client
as `docker run -it --rm quay.io/redhat-certification/chart-verifier:0d3706f verify`.

If you haven't built a container image, you could still use the Docker client to execute the latest release available in
Quay:

```text
> docker run --rm quay.io/redhat-certification/chart-verifier:latest verify --help
Verifies a Helm chart by checking some of its characteristics

Usage:
  chart-verifier verify <chart-uri> [flags]

Flags:
  -S, --chart-set strings          set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
  -F, --chart-set-file strings     set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
  -X, --chart-set-string strings   set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
  -f, --chart-values strings       specify values in a YAML file or a URL (can specify multiple)
  -x, --disable strings            all checks will be enabled except the informed ones
  -e, --enable strings             only the informed checks will be enabled
  -h, --help                       help for verify
  -o, --output string              the output format: default, json or yaml
  -s, --set strings                overrides a configuration, e.g: dummy.ok=false

Global Flags:
      --config string   config file (default is $HOME/.chart-verifier.yaml)
```

To verify a chart on the host system, the directory containing the chart should be mounted in the container; for http or
https verifications, no mounting is required:

```text
> docker run --rm quay.io/redhat-certification/chart-verifier:latest verify https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz?raw=true
```

Here is another example for a chart on the host system using volume mount. In
the below example, the chart is located in the current directory:

```text
> docker run -v $(pwd):/charts --rm quay.io/redhat-certification/chart-verifier:latest verify /charts/chart-0.1.0-v3.valid.tgz
```
