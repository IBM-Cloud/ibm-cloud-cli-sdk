# Prerequisites
We assume that the following tools are installed:

- [Go](https://golang.org/)
- [go-bindata](https://github.com/jteeuwen/go-bindata)
- [i18n4go](https://github.com/maximilien/i18n4go)
- [cucumber](https://github.com/cucumber/cucumber)
  - [aruba](https://github.com/cucumber/aruba)

# How to create plugin project

1. Create the project folder
2. Copy `setup-project.sh` and `new-command.sh` under the project root folder
3. Compose `project.properties` under the project root folder
4. Run `./setup-project.sh`

## setup-project.sh
**Syntax**
```setup-project.sh [-travis] [-h]```
**Flags:**

- -travis: create a Travis build file
- -h: show help

## project.properties
To generate the project, you need to provide a `project.properties` file with the following contents.
```properties
# name of the plugin, will be shown in base CLI when `bx plugin list`
plugin_name=test-plugin

# base name of the plugin binary files, i.e. ${base_name}-${os}-${arch}
plugin_file_basename=test-plugin

# major namespace used in this plugin, this will be used as the namespace of "hello" command
# You can always change it later
default_namespace=catalog

# new namespaces used by this plugin
#new_namespaces=("test1" "test2")

# Go version, if not specified, then deduce from the current context
#go_version=1.8.1

# go package path, please align with the project
project_path=github.ibm.com/Bluemix/test-plugin
```

# How to add a new command
Run `new-command.sh`
**Syntax**
```new-command.sh NAME [-d DESCRIPTION] [-u USAGE] [-n NAMESPACE] [-p PACKAGE] [-h]```

**Flags**

- -d DESCRIPTION: description of the command
- -u USAGE: usage of the command, e.g. command syntax
- -n NAMESPACE: namespace of the command
- -p PACKAGE: package of the command **under the project**, i.e. exclude the domain of the project.
- -h: show help
