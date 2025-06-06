---
title: "ZAP"
category: "scanner"
type: "WebApplication"
state: "released"
appVersion: "2.16.1"
usecase: "WebApp & OpenAPI Vulnerability Scanner"
---

![zap logo](https://raw.githubusercontent.com/wiki/zaproxy/zaproxy/images/zap32x32.png)

<!--
SPDX-FileCopyrightText: the secureCodeBox authors

SPDX-License-Identifier: Apache-2.0
-->
<!--
.: IMPORTANT! :.
--------------------------
This file is generated automatically with `helm-docs` based on the following template files:
- ./.helm-docs/templates.gotmpl (general template data for all charts)
- ./chart-folder/.helm-docs.gotmpl (chart specific template data)

Please be aware of that and apply your changes only within those template files instead of this file.
Otherwise your changes will be reverted/overwritten automatically due to the build process `./.github/workflows/helm-docs.yaml`
--------------------------
-->

<p align="center">
  <a href="https://opensource.org/licenses/Apache-2.0"><img alt="License Apache-2.0" src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"/></a>
  <a href="https://github.com/secureCodeBox/secureCodeBox/releases/latest"><img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/secureCodeBox/secureCodeBox?sort=semver"/></a>
  <a href="https://owasp.org/www-project-securecodebox/"><img alt="OWASP Lab Project" src="https://img.shields.io/badge/OWASP-Lab%20Project-yellow"/></a>
  <a href="https://artifacthub.io/packages/search?repo=securecodebox"><img alt="Artifact HUB" src="https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/securecodebox"/></a>
  <a href="https://github.com/secureCodeBox/secureCodeBox/"><img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/secureCodeBox/secureCodeBox?logo=GitHub"/></a>
  <a href="https://infosec.exchange/@secureCodeBox"><img alt="Mastodon Follower" src="https://img.shields.io/mastodon/follow/111902499714281911?domain=https%3A%2F%2Finfosec.exchange%2F"/></a>
</p>

## What is ZAP?

:::caution Deprecation Notice
The `zap-advanced` and `zap` ScanType are being deprecated in favor of the `zap-automation-framework`, which encompasses all functionalities of the previous ScanTypes. We recommend transitioning to "zap-automation-framework". This change will take effect in the upcoming release cycle. For guidance on migrating to "zap-automation-framework," please refer to [here](/docs/scanners/zap-automation-framework#migration-to-zap-automation-framework).
:::

The Zed Attack Proxy (ZAP) is one of the world’s most popular free security tools and is actively maintained by hundreds of international volunteers*. It can help you automatically find security vulnerabilities in your web applications while you are developing and testing your applications. It's also a great tool for experienced pentesters to use for manual security testing.

To learn more about the ZAP scanner itself visit [https://www.zaproxy.org/](https://www.zaproxy.org/).
To learn more about the ZAP Automation Framework itself visit [https://www.zaproxy.org/docs/desktop/addons/automation-framework/](https://www.zaproxy.org/docs/desktop/addons/automation-framework/).

## Deployment
The zap chart can be deployed via helm:

```bash
# Install HelmChart (use -n to configure another namespace)
helm upgrade --install zap oci://ghcr.io/securecodebox/helm/zap
```

## Scanner Configuration

The following security scan configuration example are based on the ZAP Docker Scan Scripts. By default, the secureCodeBox ZAP Helm Chart installs all four ZAP scripts: `zap-baseline`, `zap-full-scan` , `zap-api-scan` & `zap-automation-scan`. Listed below are the arguments supported by the `zap-baseline` script, which are mostly interchangeable with the other ZAP scripts (except for `zap-automation-scan`). For a more complete reference check out the [ZAP Documentation](https://www.zaproxy.org/docs/docker/) and the secureCodeBox based ZAP examples listed below.

The command line interface can be used to easily run server scans: `-t www.example.com`

```bash
Usage: zap-baseline.py -t <target> [options]
    -t target         target URL including the protocol, eg https://www.example.com
Options:
    -h                print this help message
    -c config_file    config file to use to INFO, IGNORE or FAIL warnings
    -u config_url     URL of config file to use to INFO, IGNORE or FAIL warnings
    -g gen_file       generate default config file (all rules set to WARN)
    -m mins           the number of minutes to spider for (default 1)
    -r report_html    file to write the full ZAP HTML report
    -w report_md      file to write the full ZAP Wiki (Markdown) report
    -x report_xml     file to write the full ZAP XML report
    -J report_json    file to write the full ZAP JSON document
    -a                include the alpha passive scan rules as well
    -d                show debug messages
    -P                specify listen port
    -D                delay in seconds to wait for passive scanning
    -i                default rules not in the config file to INFO
    -I                do not return failure on warning
    -j                use the Ajax spider in addition to the traditional one
    -l level          minimum level to show: PASS, IGNORE, INFO, WARN or FAIL, use with -s to hide example URLs
    -n context_file   context file which will be loaded prior to spidering the target
    -p progress_file  progress file which specifies issues that are being addressed
    -s                short output format - dont show PASSes or example URLs
    -T                max time in minutes to wait for ZAP to start and the passive scan to run
    -z zap_options    ZAP command line options e.g. -z "-config aaa=bbb -config ccc=ddd"
    --hook            path to python file that define your custom hooks
```

## ZAP Automation Scanner Configuration

The Automation Framework allows for higher flexibility in configuring ZAP scans. Its goal is the automation of the full functionality of ZAP's GUI. The configuration of the Automation Framework differs from the other three ZAP scan types. The following security scan configuration example highlights the differences for running a `zap-automation-scan`.
Of particular interest for us will be the -autorun option. `zap-automation-scan` allows for providing an automation file as a ConfigMap that defines the details of the scan. See the secureCodeBox based ZAP Automation example listed below for what such a ConfigMap would look like.

```bash
Usage: zap.sh -cmd -host <target> [options]
    -t target         target URL including the protocol, eg https://www.example.com
Add-on options:
    -script <script>          Run the specified script from commandline or load in GUI
    -addoninstall <addOnId>   Installs the add-on with specified ID from the ZAP Marketplace
    -addoninstallall          Install all available add-ons from the ZAP Marketplace
    -addonuninstall <addOnId> Uninstalls the Add-on with specified ID
    -addonupdate              Update all changed add-ons from the ZAP Marketplace
    -addonlist                List all of the installed add-ons
    -certload <path>          Loads the Root CA certificate from the specified file name
    -certpubdump <path>       Dumps the Root CA public certificate into the specified file name, this is suitable for importing into browsers
    -certfulldump <path>      Dumps the Root CA full certificate (including the private key) into the specified file name, this is suitable for importing into ZAP
    -notel                    Turns off telemetry calls
    -hud                      Launches a browser configured to proxy through ZAP with the HUD enabled, for use in daemon mode
    -hudurl <url>             Launches a browser as per the -hud option with the specified URL
    -hudbrowser <browser>     Launches a browser as per the -hud option with the specified browser, supported options: Chrome, Firefox by default 'Firefox'
    -openapifile <path>       Imports an OpenAPI definition from the specified file name
    -openapiurl <url>         Imports an OpenAPI definition from the specified URL
    -openapitargeturl <url>   The Target URL, to override the server URL present in the OpenAPI definition. Refer to the help for supported format.
    -quickurl <target url>    The URL to attack, e.g. http://www.example.com
    -quickout <filename>      The file to write the HTML/JSON/MD/XML results to (based on the file extension)
    -autorun <filename>       Run the automation jobs specified in the file.
    -autogenmin <filename>    Generate template automation file with the key parameters.
    -autogenmax <filename>    Generate template automation file with all parameters.
    -autogenconf <filename>   Generate template automation file using the current configuration.
    -graphqlfile <path>       Imports a GraphQL Schema from a File
    -graphqlurl <url>         Imports a GraphQL Schema from a URL
    -graphqlendurl <url>      Sets the Endpoint URL
```

## Requirements

Kubernetes: `>=v1.11.0-0`

The secureCodeBox provides two different scanner charts (`zap`, `zap-advanced`) to automate ZAP WebApplication security scans. The first one `zap` comes with four scanTypes:
- `zap-baseline-scan`
- `zap-full-scan`
- `zap-api-scan`
- `zap-automation-scan`

The scanTypes `zap-baseline-scan`, `zap-full-scan` & `zap-api-scan` can be configured via CLI arguments which are somehow a bit limited for some advanced usecases, e.g. using custom zap scripts or configuring complex authentication settings.

That's why we introduced this `zap-advanced` scanner chart, which introduces extensive YAML configuration options for ZAP. The YAML configuration can be split in multiple files and will be merged at start.
ZAP's own Automation Framework provides similar functionality to the `zap-advanced` scanner chart and is set to displace it in the future.

## ZAP Automation Configuration

The ZAP Automation Scanner supports the use of secrets, as to not have hardcoded credentials in the scan definition.
Generate secrets using the credentials that will later be used in the scan for authentication. Supported authentication methods for the ZAP Authentication scanner are Manual, HTTP / NTLM, Form-based, JSON-based, and Script-based.

```bash
kubectl create secret generic unamesecret --from-literal='username=<USERNAME>'
kubectl create secret generic pwordsecret --from-literal='password=<PASSWORD>'
```

You can now include the secrets in the scan definition and reference them in the ConfigMap that defines the scan options.
A ZAP Automation scan using JSON-based authentication may look like this:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: "zap-automation-scan-config"
data:
  1-automation.yaml: |-

    env:                                   # The environment, mandatory
      contexts:                           # List of 1 or more contexts, mandatory
        - name: test-config                  # Name to be used to refer to this context in other jobs, mandatory
          urls: ["http://juiceshop.demo-targets.svc:3000"]                           # A mandatory list of top level urls, everything under each url will be included
          includePaths:
            - "http://juiceshop.demo-targets.svc:3000/.*"                   # An optional list of regexes to include
          excludePaths:
            - ".*socket\\.io.*"
            - ".*\\.png"
            - ".*\\.jpeg"
            - ".*\\.jpg"
            - ".*\\.woff"
            - ".*\\.woff2"
            - ".*\\.ttf"
            - ".*\\.ico"                  
          authentication:
            method: "json"
            parameters:
              loginPageUrl: "http://juiceshop.demo-targets.svc:3000/rest/user"
              loginRequestUrl: "http://juiceshop.demo-targets.svc:3000/rest/user/login"
              loginRequestBody: '{"email":"${EMAIL}","password":"${PASS}"}'
            verification:
              method: "response"
              loggedOutRegex: '\Q{"user":{}}\E'
              loggedInRegex: '\Q<a href="password.jsp">\E'
          users:
          - name: "juiceshop-user-1"
            credentials:
              username: "${EMAIL}"
              password: "${PASS}"
      parameters:
        failOnError: true                  # If set exit on an error        
        failOnWarning: false               # If set exit on a warning
        progressToStdout: true             # If set will write job progress to stdout

    jobs:
      - type: passiveScan-config           # Passive scan configuration
        parameters:
          maxAlertsPerRule: 10             # Int: Maximum number of alerts to raise per rule
          scanOnlyInScope: true            # Bool: Only scan URLs in scope (recommended)
      - type: spider                       # The traditional spider - fast but doesnt handle modern apps so well
        parameters:
          context: test-config                        # String: Name of the context to spider, default: first context
          user: juiceshop-user-1                           # String: An optional user to use for authentication, must be defined in the env
          maxDuration: 2                     # Int: The max time in minutes the spider will be allowed to run for, default: 0 unlimited
      - type: spiderAjax                   # The ajax spider - slower than the spider but handles modern apps well
        parameters:
          context: test-config                         # String: Name of the context to spider, default: first context
          maxDuration: 2                     # Int: The max time in minutes the ajax spider will be allowed to run for, default: 0 unlimited
      - type: passiveScan-wait             # Passive scan wait for the passive scanner to finish
        parameters:
          maxDuration: 10                   # Int: The max time to wait for the passive scanner, default: 0 unlimited
      - type: report                       # Report generation
        parameters:
          template: traditional-xml                        # String: The template id, default : modern
          reportDir: /home/securecodebox/               # String: The directory into which the report will be written
          reportFile: zap-results                     # String: The report file name pattern, default: [[yyyy-MM-dd]]-ZAP-Report-[[site]]
        risks:                             # List: The risks to include in this report, default all
          - high
          - medium
          - low
---
apiVersion: "execution.securecodebox.io/v1"
kind: Scan
metadata:
  name: "zap-example-scan"
spec:
  scanType: "zap-automation-scan"
  parameters:
    - "-autorun"
    - "/home/securecodebox/scb-automation/1-automation.yaml"
  volumeMounts:
    - mountPath: /home/securecodebox/scb-automation/1-automation.yaml
      name: zap-automation
      subPath: 1-automation.yaml
  volumes:
    - name: zap-automation
      configMap:
        name: zap-automation-scan-config
  env:
    - name: EMAIL
      valueFrom:
        secretKeyRef:
          name: unamesecret
          key: username
    - name: PASS
      valueFrom:
        secretKeyRef:
          name: pwordsecret
          key: password
```

For a complete overview of all the possible options you have for configuring a ZAP Automation scan, run
```bash
./zap.sh -cmd -autogenmax zap.yaml
```
For an overview of all required configuration options, run
```
bash ./zap.sh -cmd -autogenmin zap.yaml
```
Alternatively, have a look at the [official documentation](https://www.zaproxy.org/docs/desktop/addons/automation-framework/).

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| cascadingRules.enabled | bool | `false` | Enables or disables the installation of the default cascading rules for this scanner |
| imagePullSecrets | list | `[]` | Define imagePullSecrets when a private registry is used (see: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) |
| parser.affinity | object | `{}` | Optional affinity settings that control how the parser job is scheduled (see: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/) |
| parser.env | list | `[]` | Optional environment variables mapped into each parseJob (see: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/) |
| parser.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images |
| parser.image.repository | string | `"docker.io/securecodebox/parser-zap"` | Parser image repository |
| parser.image.tag | string | defaults to the charts version | Parser image tag |
| parser.nodeSelector | object | `{}` | Optional nodeSelector settings that control how the scanner job is scheduled (see: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/) |
| parser.resources | object | `{ requests: { cpu: "200m", memory: "100Mi" }, limits: { cpu: "400m", memory: "200Mi" } }` | Optional resources lets you control resource limits and requests for the parser container. See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |
| parser.scopeLimiterAliases | object | `{}` | Optional finding aliases to be used in the scopeLimiter. |
| parser.tolerations | list | `[]` | Optional tolerations settings that control how the parser job is scheduled (see: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) |
| parser.ttlSecondsAfterFinished | string | `nil` | seconds after which the Kubernetes job for the parser will be deleted. Requires the Kubernetes TTLAfterFinished controller: https://kubernetes.io/docs/concepts/workloads/controllers/ttlafterfinished/ |
| scanner.activeDeadlineSeconds | string | `nil` | There are situations where you want to fail a scan Job after some amount of time. To do so, set activeDeadlineSeconds to define an active deadline (in seconds) when considering a scan Job as failed. (see: https://kubernetes.io/docs/concepts/workloads/controllers/job/#job-termination-and-cleanup) |
| scanner.affinity | object | `{}` | Optional affinity settings that control how the scanner job is scheduled (see: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/) |
| scanner.backoffLimit | int | 3 | There are situations where you want to fail a scan Job after some amount of retries due to a logical error in configuration etc. To do so, set backoffLimit to specify the number of retries before considering a scan Job as failed. (see: https://kubernetes.io/docs/concepts/workloads/controllers/job/#pod-backoff-failure-policy) |
| scanner.env | list | `[]` | Optional environment variables mapped into each scanJob (see: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/) |
| scanner.envFrom | list | `[]` | Optional mount environment variables from configMaps or secrets (see: https://kubernetes.io/docs/tasks/inject-data-application/distribute-credentials-secure/#configure-all-key-value-pairs-in-a-secret-as-container-environment-variables) |
| scanner.extraContainers | list | `[]` | Optional additional Containers started with each scanJob (see: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) |
| scanner.extraVolumeMounts | list | `[{"mountPath":"/zap/wrk","name":"zap-workdir"}]` | Optional VolumeMounts mapped into each scanJob (see: https://kubernetes.io/docs/concepts/storage/volumes/) |
| scanner.extraVolumes | list | `[{"emptyDir":{},"name":"zap-workdir"}]` | Optional Volumes mapped into each scanJob (see: https://kubernetes.io/docs/concepts/storage/volumes/) |
| scanner.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images |
| scanner.image.repository | string | `"docker.io/zaproxy/zap-stable"` | Container Image to run the scan |
| scanner.image.tag | string | `nil` | defaults to the charts appVersion |
| scanner.nameAppend | string | `nil` | append a string to the default scantype name. |
| scanner.nodeSelector | object | `{}` | Optional nodeSelector settings that control how the scanner job is scheduled (see: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/) |
| scanner.podSecurityContext | object | `{}` | Optional securityContext set on scanner pod (see: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) |
| scanner.resources | object | `{}` | CPU/memory resource requests/limits (see: https://kubernetes.io/docs/tasks/configure-pod-container/assign-memory-resource/, https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/) |
| scanner.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["all"]},"privileged":false,"readOnlyRootFilesystem":false,"runAsNonRoot":false}` | Optional securityContext set on scanner container (see: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) |
| scanner.securityContext.allowPrivilegeEscalation | bool | `false` | Ensure that users privileges cannot be escalated |
| scanner.securityContext.capabilities.drop[0] | string | `"all"` | This drops all linux privileges from the container. |
| scanner.securityContext.privileged | bool | `false` | Ensures that the scanner container is not run in privileged mode |
| scanner.securityContext.readOnlyRootFilesystem | bool | `false` | Prevents write access to the containers file system |
| scanner.securityContext.runAsNonRoot | bool | `false` | Enforces that the scanner image is run as a non root user |
| scanner.suspend | bool | `false` | if set to true the scan job will be suspended after creation. You can then resume the job using `kubectl resume <jobname>` or using a job scheduler like kueue |
| scanner.tolerations | list | `[]` | Optional tolerations settings that control how the scanner job is scheduled (see: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) |
| scanner.ttlSecondsAfterFinished | string | `nil` | seconds after which the Kubernetes job for the scanner will be deleted. Requires the Kubernetes TTLAfterFinished controller: https://kubernetes.io/docs/concepts/workloads/controllers/ttlafterfinished/ |

## License
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Code of secureCodeBox is licensed under the [Apache License 2.0][scb-license].

[scb-owasp]:    https://www.owasp.org/index.php/OWASP_secureCodeBox
[scb-docs]:     https://www.securecodebox.io/
[scb-site]:     https://www.securecodebox.io/
[scb-github]:   https://github.com/secureCodeBox/
[scb-mastodon]: https://infosec.exchange/@secureCodeBox
[scb-slack]:    https://owasp.org/slack/invite
[scb-license]:  https://github.com/secureCodeBox/secureCodeBox/blob/master/LICENSE
[zap github]: https://github.com/zaproxy/zaproxy/
[zap user guide]: https://www.zaproxy.org/docs/
[zap automation framework]: https://www.zaproxy.org/docs/desktop/addons/automation-framework/
