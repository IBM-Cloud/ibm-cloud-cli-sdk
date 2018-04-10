#!/bin/bash

set -e

. ./project.properties

namespace="${project_namespace}"
description="TODO"
usage="TODO"
name=""
command_package="plugin/commands"
while [[ $# -gt 0 ]]
do
	key="$1"

	case $key in 
		-n)
		namespace="$2"
		shift
		;;
		-d)
		description="$2"
		shift
		;;
		-u)
		usage="$2"
		shift
		;;
		-h)
		echo "Usage: new_command.sh COMMAND_NAME [-n Namespace] [-d DESCRIPTIN] [-u USAGE] [-p COMMAND_PACKAGE] -h"
		exit 0
		;;
		-p)
		command_package="$2"
		shift
		;;
		*)
		name="$key"
		;;
	esac
	shift
done
if [[ -z "${name}" ]]; then
    echo "Please specify the command name"
    echo "Usage: new_command.sh COMMAND_NAME [-n Namespace] [-d DESCRIPTIN] [-u USAGE] -h"
	exit 1
fi


IFS="-", read -r -a names <<< "${name}"
for index in "${!names[@]}"
do
  names[$index]="$(tr '[:lower:]' '[:upper:]' <<< ${names[index]:0:1})${names[index]:1}"
done

readonly command_struct_name=$(printf "%s" "${names[@]}")
readonly package="${command_package##*/}"
mkdir -p ${command_package}
cat <<ENDGO > ${command_package}/${name}.go
package ${package}

import (

	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"${project_path}/plugin/metadata"
	. "${project_path}/plugin/i18n"
)

type ${command_struct_name} struct {
	ui          terminal.UI
	context     plugin.PluginContext
}

func (cmd *${command_struct_name}) GetMetadata() metadata.CommandMetadata {
	return metadata.CommandMetadata{
		Namespace:   "${namespace}",
		Name:        "${name}",
		Description: T("${description}"),
		Usage:       "${usage}",
	}
}

func (cmd *${command_struct_name}) Setup(ui terminal.UI, context plugin.PluginContext) {
	cmd.ui = ui
	cmd.context = context
}

func (cmd *${command_struct_name}) Run(context *cli.Context) error {
	return nil
}
ENDGO

if [ "$(uname)" == "Darwin" ]; then
    sed -i "" "/return \[\]metadata.Command{/a\\
        new(${package}.${command_struct_name})," plugin/plugin.go
else
    sed -i "/return \[\]metadata.Command{/a\\
        new(${package}.${command_struct_name})," plugin/plugin.go
fi


readonly file_import="${project_path}/${command_package}"
if ! grep -q "${file_import}"  plugin/plugin.go; then
	if [ "$(uname)" == "Darwin" ]; then
		sed -i "" "/github\.ibm\.com\/IAM\/test\/plugin\/metadata/a\\
			    \"${file_import}\"\\
			" plugin/plugin.go
	else
		sed -i "/github\.ibm\.com\/IAM\/test\/plugin\/metadata/a\\
			    \"${file_import}\"\\
			" plugin/plugin.go
	fi
fi
