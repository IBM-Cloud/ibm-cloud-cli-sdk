# IBM Cloud CLI SDK

IBM Cloud CLI plugin SDK provides predefined plugin interface, utilities and libraries to develop plugins for [IBM Cloud cli](https://clis.cloud.ibm.com).

# Get started

You firstly need [Go](http://www.golang.org) installed on your machine. Then clone this repository into `$GOPATH/src/github.com/IBM-Cloud/ibm-cloud-cli-sdk`. 

This project uses [go modules](https://go.dev/blog/using-go-modules) to manage dependencies. Go to the project directory and run the following command to restore the dependencies into vendor folder:

```bash
$ go mod vendor
```

and then run tests:

```bash
$ go test ./...
```

# Build and run plugin

Download and install the IBM Cloud CLI. See instructions [here](https://clis.cloud.ibm.com).

Compile the plugin source code with `go build` command, for example

```bash
$ go build plugin_examples/hello.go
```

Install the plugin:

```bash
$ ibmcloud plugin install ./hello
```

List installed plugins:

```bash
$ ibmcloud plugin list

# list plugin commands with '-c' flag
$ ibmcloud plugin list -c
```

Uninstall the plugin:

```bash
$ ibmcloud plugin uninstall SayHello # SayHello is the plugin name
```

For more usage of 
plugin management, run `ibmcloud help plugin`

# Develop plugins

Refer to [plugin developer guide](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md) for how to develop a plugin.

See plugin examples [here](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/tree/master/plugin_examples)

# Publish plugins

IBM Cloud has a public plugin repository by default installed in IBM Cloud CLI. Run `ibmcloud plugin`, you can see a repository named `IBM Cloud` (`https://plugins.cloud.ibm.com`). The repository support multiple version of plugin. You can list all plugins in the repository by using `ibmcloud plugin repo-plugins -r 'IBM Cloud'`.

To publish, update or remove your plugin in IBM Cloud plugin repository, you can simply [create an issue on GitHub](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/issues/new) following below samples:

**Example to publish a new plugin**:

Title: [plugin-publish] Request to publish a new plugin 'SayHello'

Content:

```yaml

- name: SayHello
  description: Say hello
  company: YYY
  authors:
  - name: xxx
    contact: xxx@example.com
  homepage: http://www.example.com/hello
  version: 0.0.1
  binaries:
  - platform: osx
    url: http://www.example.com/downloads/hello/hello-darwin-amd64-0.0.1
    checksum: xxxxx
  - platform: win32
    url: http://www.example.com/downloads/hello/hello-windows-386-0.0.1.exe
    checksum: xxxxx
  - platform: win64
    url: http://www.example.com/downloads/hello/hello-windows-amd64-0.0.1.exe
    checksum: xxxxx
  - platform: linux32
    url: http://www.example.com/downloads/hello/hello-linux-386-0.0.1.exe
    checksum: xxxxx
  - platform: linux64
    url: http://www.example.com/downloads/hello/hello-linux-amd64-0.0.1.exe
    checksum: xxxxx
```

The following descibes each field's usage.

Field | Description
------ | ---------
name | Name of your plugin, must not conflict with other existing plugins in the repo.
description | Describe your plugin in a line or two. This description will show up when your plugin is listed on the command line. Avoid saying "A plugin to ..." as it's redundant. Just briefly describe what the plugin provides.
company | *Optional*
authors | authors of the plugin: `name`: name of author; `homepage`: *Optional* link to the homepage of the author; `contact`: *Optional* ways to contact author, email, twitter, phone etc ...
homepage | Link to the homepage
version | Version number of your plugin, in [major].[minor].[build] form
binaries | This section has fields detailing the various binary versions of your plugin. To reach as large an audience as possible, we encourage contributors to cross-compile their plugins on as many platforms as possible. Go provides everything you need to cross-compile for different platforms. `platform`: The os for this binary. Supports `osx`, `linux32`, `linux64`, `win32`, `win64`; `url`: Link to the binary file itself; `checksum`: SHA-1 of the binary file for verification.

**Example to update a plugin**:

Title: [plugin-update] Request to update plugin 'SayHello'

Content:

```yaml

- name: SayHello
  description: Updated description of plugin Hello
  company: YYY
  authors:
  - name: xxx
    contact: xxx@example.com
  homepage: http://www.example.com/hello
```

**Example to remove a plugin**:

Title: [plugin-remove] Request to remove plugin 'SayHello'


**Example to submit/update a version**:

Title: [plugin-version-update] Request to submit a new version of plugin 'SayHello'

Content:

```yaml

- name: SayHello
  version: 0.0.2
  binaries:
  - platform: osx
    url: http://www.example.com/downloads/hello/hello-darwin-amd64-0.0.2
    checksum: xxxxx
  - platform: win32
    url: http://www.example.com/downloads/hello/hello-windows-386-0.0.2.exe
    checksum: xxxxx
  - platform: win64
    url: http://www.example.com/downloads/hello/hello-windows-amd64-0.0.2.exe
    checksum: xxxxx
  - platform: linux32
    url: http://www.example.com/downloads/hello/hello-linux-386-0.0.2.exe
    checksum: xxxxx
  - platform: linux64
    url: http://www.example.com/downloads/hello/hello-linux-amd64-0.0.2.exe
    checksum: xxxxx
```

**Example to remove plugin versions**:

Title: [plugin-remove] Request to remove plugin 'SayHello'

Content:

```yaml

- name: SayHello
  versions:
    - 0.0.1
      0.0.2
```

# Issues

Report problems by [adding an issue on GitHub](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/issues/new).

# License

This project is released under version 2.0 of the [Apache License](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/LICENSE)





