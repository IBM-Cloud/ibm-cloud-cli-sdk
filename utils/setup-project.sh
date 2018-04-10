#!/bin/bash
set -e

while [[ $# -gt 0 ]]
do
	key="$1"

	case $key in 
		-travis)
		travis="true"
		;;
		-h)
		printf "Usage:\n  setup-project.sh [-travis] [-h]\n\t-travis: create a travis build file\n\t-h print help\n"
		exit 0
		;;
		*)
		;;
	esac
	shift
done

if [[ ! -f project.properties ]]; then
	echo "project.properties is not found"
	exit 1
fi

. ./project.properties

if [[ -z $go_version ]]; then
        go_version=`go version | sed 's/[^0-9.]*\([0-9.]*\).*/\1/'`
fi

IFS="-", read -r -a names <<< "${plugin_name}"
for index in "${!names[@]}"
do
  names[$index]="$(tr '[:lower:]' '[:upper:]' <<< ${names[index]:0:1})${names[index]:1}"
done
readonly plugin_struct_name=$(printf "%s" "${names[@]}")

function check_exist() {
    if ! hash $1 2>/dev/null; then
        echo "'$1' is not found. Please make sure it is installed."
        exit 1
    fi
}

if [[ ${#new_namespaces[@]} -gt 0 ]]; then
	namespace_definition=""

	for namespace in "${new_namespaces[@]}"
	do
		namespace_definition="${namespace_definition}plugin.Namespace{Name:\"${namespace}\",Description:\"TODO\"},"
		namespace_definition="${namespace_definition}"$'\n'
	done
	namespace_definition=$'Namespaces: []plugin.Namespace{\n\t\t\t'"${namespace_definition}"$'\n\t\t\t},'
fi

# check prerequites
check_exist "cucumber"
check_exist "go-bindata"
check_exist "i18n4go"

# set up BDD
mkdir -p bdd
(cd bdd && cucumber --init)

# set up devops scripts
mkdir -p bin
cat <<ENDSCRIPT > bin/build.sh
#!/bin/bash

set -e

readonly ROOT_DIR=\$(cd \$(dirname \$(dirname \$0)) && pwd)

go build -ldflags "-w" -o \$ROOT_DIR/out/${plugin_file_basename} .
ENDSCRIPT

cat <<ENDSCRIPT > bin/build-all.sh
#!/bin/bash

set -e

readonly ROOT_DIR=\$(cd \$(dirname \$(dirname \$0)) && pwd)
readonly OUT_DIR=\$ROOT_DIR/out

build() {
    os=\$1
    arch=\$2
    GOARCH=\$arch GOOS=\$os \$ROOT_DIR/bin/build.sh

    nf="${plugin_file_basename}-\$os-\$arch"
    if [ \$os == "windows" ]; then
        nf="\$nf.exe"
    fi
    mv \$OUT_DIR/${plugin_file_basename} "\$OUT_DIR/\$nf"
}

build windows amd64
build windows 386
#disable CGO for Linux
CGO_ENABLED=0 build linux amd64
CGO_ENABLED=0 build linux 386
build darwin amd64
ENDSCRIPT

cat <<ENDSCRIPT > bin/catch-i18n-mismatch.sh
#!/bin/bash

ROOT_DIR=\$(cd \$(dirname \$(dirname \$0)) && pwd)

(cd \${ROOT_DIR}/plugin && i18n4go -c checkup -v | sed -E 's/(.+) exists in the code, but not in en_US/{"id": \1, "translation": \1},/g')
cd ..
ENDSCRIPT

cat <<ENDSCRIPT > bin/format-translation-files.sh
#!/usr/bin/env ruby

require 'json'

ROOT = File.expand_path(File.join(File.dirname(__FILE__), ".."))
RESOURCES_DIR = File.join(ROOT, "plugin", "i18n", "resources")

SRC_LANG = "en_US"

def run
  srcLangFile = sourceLangFile()
  srcTranslations = loadFromFile(srcLangFile)

  Dir.glob(File.join(RESOURCES_DIR, "*.all.json")) do |file|
    puts "*** Process #{File.basename(file)}"
    if file == srcLangFile
      translations = srcTranslations
    else
      translations = getTranslations(srcTranslations, file)
    end
    normalized = normalize(translations)
    saveToFile(normalized, file)
  end
end

def sourceLangFile
  return File.join(RESOURCES_DIR, SRC_LANG + ".all.json")
end

def loadFromFile(file)
  array = JSON.parse(File.read(file))
  translations = {}
  array.each{ |t| translations[t["id"]] = t["translation"] unless t["id"].to_s.empty? }
  return translations
end

def getTranslations(srcTranslations, file)
  translations = loadFromFile(file)
  result = {}
  srcTranslations.each do |k,v|
    if translations.include? k
      result[k] = translations[k]
    else
      result[k] = ""
    end
  end
  return result
end

def normalize(translations)
  result = {}
  sorted = Hash[translations.sort]
  sorted.each {|k,v| result[k] = v.to_s.empty?? k : v}
  return result
end

def saveToFile(translations, path)
  json = []
  translations.each { |k,v| json << {"id" => k, "translation" => v} }
  File.open(path,"w") do |f|
    f.write(JSON.pretty_generate(json))
  end
end

run
ENDSCRIPT

cat <<ENDSCRIPT > bin/generate-i18n-resources.sh
#!/bin/bash

set -e

echo "Generating i18n resource file ..."
\$GOPATH/bin/go-bindata -pkg resources -o plugin/resources/i18n_resources.go plugin/i18n/resources
echo "Done."
ENDSCRIPT

cat <<ENDSCRIPT > bin/update-translation-files.sh
#!/bin/bash

set -e

ROOT_DIR=\$(cd \$(dirname \$(dirname \$0)) && pwd)
\$ROOT_DIR/bin/format-translation-files.sh
\$ROOT_DIR/bin/generate-i18n-resources.sh
ENDSCRIPT

chmod a+x bin/*

# set up source code folders
mkdir -p plugin/api
mkdir -p plugin/metadata
mkdir -p plugin/commands
mkdir -p plugin/errors
mkdir -p plugin/models
mkdir -p plugin/version
mkdir -p plugin/i18n/resources
mkdir -p plugin/i18n/detection
mkdir -p plugin/resource

cat <<ENDGO > plugin/metadata/command.go
package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
)

type CommandMetadata struct {
	Namespace   string
	Name        string
	Description string
	Usage       string
	Flags       []cli.Flag
}

type Command interface {
	GetMetadata() CommandMetadata
	Setup(ui terminal.UI, context plugin.PluginContext)
	Run(*cli.Context) error
}
ENDGO

cat <<ENDGO > plugin/errors/invalid_usage.go
package errors

import "errors"

var ErrInvalidUsage = errors.New("invalid usage")
ENDGO

cat <<ENDGO > plugin/version/version.go
package version

// plugin version
const PLUGIN_VERSION = "0.1.0"

// plugin major version
const PLUGIN_MAJOR_VERSION = 0

// plugin minor version
const PLUGIN_MINOR_VERSION = 1

// plugin build version
const PLUGIN_BUILD_VERSION = 0
ENDGO

cat <<ENDGO > plugin/i18n/detection/detection.go
package detection

import "github.com/cloudfoundry/jibber_jabber"

type Detector interface {
	DetectLocale() string
	DetectLanguage() string
}

type JibberJabberDetector struct{}

func (d *JibberJabberDetector) DetectLocale() string {
	userLocale, err := jibber_jabber.DetectIETF()
	if err != nil {
		userLocale = ""
	}
	return userLocale
}

func (d *JibberJabberDetector) DetectLanguage() string {
	lang, err := jibber_jabber.DetectLanguage()
	if err != nil {
		lang = ""
	}
	return lang
}
ENDGO

cat <<ENDGO > plugin/i18n/i18n.go
package i18n

import (
	"path/filepath"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/i18n"

	"github.ibm.com/Bluemix/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"${project_path}/plugin/i18n/detection"
	"${project_path}/plugin/resources"
)

const (
	DEFAULT_LOCALE = "en_US"
)

var SUPPORTED_LOCALES = []string{
	"de_DE",
	"en_US",
	"es_ES",
	"fr_FR",
	"it_IT",
	"ja_JP",
	"ko_KR",
	"pt_BR",
	"zh_Hans",
	"zh_Hant",
}

var resourcePath = filepath.Join("plugin", "i18n", "resources")

func GetResourcePath() string {
	return resourcePath
}

func SetResourcePath(path string) {
	resourcePath = path
}

var T goi18n.TranslateFunc = Init(core_config.NewCoreConfig(func(error) {}), new(detection.JibberJabberDetector))

func Init(coreConfig core_config.Reader, detector detection.Detector) goi18n.TranslateFunc {
	userLocale := coreConfig.Locale()
	if userLocale != "" {
		return initWithLocale(userLocale)
	}
	locale := supportedLocale(detector.DetectLocale())
	if locale == "" {
		locale = defaultLocaleForLang(detector.DetectLanguage())
	}
	if locale == "" {
		locale = DEFAULT_LOCALE
	}
	return initWithLocale(locale)
}

func initWithLocale(locale string) goi18n.TranslateFunc {
	err := loadFromAsset(locale)
	if err != nil {
		panic(err)
	}
	return goi18n.MustTfunc(locale)
}

func loadFromAsset(locale string) (err error) {
	assetName := locale + ".all.json"
	assetKey := filepath.Join(resourcePath, assetName)
	bytes, err := resources.Asset(assetKey)
	if err != nil {
		return
	}
	err = goi18n.ParseTranslationFileBytes(assetName, bytes)
	return
}

func supportedLocale(locale string) string {
	locale = normailizeLocale(locale)
	for _, l := range SUPPORTED_LOCALES {
		if strings.EqualFold(locale, l) {
			return l
		}
	}
	switch locale {
	case "zh_cn", "zh_sg":
		return "zh_Hans"
	case "zh_hk", "zh_tw":
		return "zh_Hant"
	}
	return ""
}

func normailizeLocale(locale string) string {
	return strings.ToLower(strings.Replace(locale, "-", "_", 1))
}

func defaultLocaleForLang(lang string) string {
	if lang != "" {
		lang = strings.ToLower(lang)
		for _, l := range SUPPORTED_LOCALES {
			if lang == l[0:2] {
				return l
			}
		}
	}
	return ""
}
ENDGO

cat <<ENDGO > plugin/commands/namespace.go
package commands

const PLUGIN_NAMESPACE = "${default_namespace}"
ENDGO

cat <<ENDJSON > plugin/i18n/resources/en_US.all.json
[
  {"id": "OPTIONS:", "translation": "OPTIONS:"},
  {"id": "Incorrect Usage.", "translation": "Incorrect Usage."},
  {"id": "Say 'Hello, world!'", "translation": "Say 'Hello, world!'"},
  {"id": "NAME:", "translation": "NAME:"},
  {"id": "ALIAS:", "translation": "ALIAS:"},
  {"id": "USAGE:", "translation": "USAGE:"}
]
ENDJSON
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/de_DE.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/fr_FR.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/es_ES.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/ja_JP.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/it_IT.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/pt_BR.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/ko_KR.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/zh_Hans.all.json
cp plugin/i18n/resources/en_US.all.json plugin/i18n/resources/zh_Hant.all.json

bin/update-translation-files.sh

cat <<ENDGO > plugin/plugin.go
package plugin

import (
	"os"
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"
	"${project_path}/plugin/metadata"
	"${project_path}/plugin/commands"
	"${project_path}/plugin/errors"

	. "${project_path}/plugin/i18n"
	"${project_path}/plugin/version"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
)

// plugin name
const PLUGIN_NAME = "${plugin_name}"

var (
	COMMAND_HELP_TEMPLATE = T("NAME:") + \`
{{.Name}} - {{.Description}}{{with .ShortName}}
\` + T("ALIAS:") + \`
   {{.}}{{end}}

\` + T("USAGE:") + \`
   {{.Usage}}
{{with .Flags}}
\` + T("OPTIONS:") + \`
{{range .}}   {{.}}
{{end}}{{end}}
\`
)

type ${plugin_struct_name}Plugin struct {
	ui terminal.UI
}

func (p *${plugin_struct_name}Plugin) GetMetadata() plugin.PluginMetadata {

	metadata := plugin.PluginMetadata{
		${namespace_definition}

		Name: PLUGIN_NAME,

		Version: plugin.VersionType{
			Major: version.PLUGIN_MAJOR_VERSION,
			Minor: version.PLUGIN_MINOR_VERSION,
			Build: version.PLUGIN_BUILD_VERSION,
		},

		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 5,
			Build: 0,
		},
	}

	for _, cmd := range getCommands() {
		cmdMeta := cmd.GetMetadata()
		metadata.Commands = append(metadata.Commands, plugin.Command{
			Namespace:   cmdMeta.Namespace,
			Name:        cmdMeta.Name,
			Description: cmdMeta.Description,
			Usage:       cmdMeta.Usage,
			Flags:       convertToPluginFlags(cmdMeta.Flags),
		})
	}
	return metadata
}

func (p *${plugin_struct_name}Plugin) Run(context plugin.PluginContext, args []string) {

	trace.Logger = trace.NewLogger(context.Trace())

	terminal.UserAskedForColors = context.ColorEnabled()
	terminal.InitColorSupport()
	p.ui = terminal.NewStdUI()

	cli.CommandHelpTemplate = COMMAND_HELP_TEMPLATE

	app := cli.NewApp()
	app.Name = "bluemix ${default_namespace}"
	app.Version = version.PLUGIN_VERSION

	for _, c := range getCommands() {
		cmd := c
		meta := cmd.GetMetadata()

		app.Commands = append(app.Commands, cli.Command{
			Name:      meta.Name,
			Usage:     meta.Description,
			UsageText: meta.Usage,
			Flags:     meta.Flags,
			Action: func(ctx *cli.Context) error {
				cmd.Setup(p.ui, context)
				err := cmd.Run(ctx)
				switch err {
				case nil:
					return nil
				case errors.ErrInvalidUsage:
					p.ui.Failed(T("Incorrect Usage."))
					cli.ShowCommandHelp(ctx, ctx.Command.Name)
					os.Exit(2)
				default:
					p.ui.Failed(err.Error())
					os.Exit(1)
				}

				return nil
			},
		})
	}

	app.Run(append([]string{context.CommandNamespace()}, args...)) // need to append a fake CLI placeholder, since plugin 'args' is stripped of leading 'bx' or 'bluemix'
}

func getCommands() []metadata.Command {
	return []metadata.Command{
		new(commands.Hello),
	}
}

func convertToPluginFlags(flags []cli.Flag) []plugin.Flag {
	var ret []plugin.Flag
	for _, f := range flags {
		ret = append(ret, plugin.Flag{
			Name:        f.GetName(),
			Description: reflect.ValueOf(f).FieldByName("Usage").String(),
		})
	}
	return ret
}
ENDGO

./new-command.sh hello -n "${default_namespace}" -d "Say 'Hello, world!'" -u "\${COMMAND_NAME} ${default_namespace} hello"

cat <<ENDGO > main.go
package main

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	${plugin_name//-/_}_plugin "${project_path}/plugin"
)

func main() {
	plugin.Start(new(${plugin_name//-/_}_plugin.${plugin_struct_name}Plugin))
}

ENDGO


echo "Setting up dependency..."
govendor init
govendor add +external
echo "Dependency added"

if [[ ! -z "${travis}" ]]; then
echo "Creating Travis file..."
cat <<ENDTRAVIS > .travis.yml
language: go
go:
  - ${go_version}
install: true
before_script:
  - go vet \$(go list ./plugin/...)
  - go test \$(go list ./plugin/...)
script:
  - bin/build-all.sh
ENDTRAVIS
fi

echo ""
echo "Done!"
