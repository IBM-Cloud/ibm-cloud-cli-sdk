// Code generated for package resources by go-bindata DO NOT EDIT. (@generated)
// sources:
// i18n/resources/all.de_DE.json
// i18n/resources/all.en_US.json
// i18n/resources/all.es_ES.json
// i18n/resources/all.fr_FR.json
// i18n/resources/all.it_IT.json
// i18n/resources/all.ja_JP.json
// i18n/resources/all.ko_KR.json
// i18n/resources/all.pt_BR.json
// i18n/resources/all.zh_Hans.json
// i18n/resources/all.zh_Hant.json
package resources

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _i18nResourcesAllDe_deJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nGeben Sie eine Zahl ein."
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "Füge einen Satz ohne Subjekt hinzu, der beschreibt, was der Befehl bewirkt."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "Fügen Sie einen Hilfetext hinzu, der die Befehlssyntax mit Parametern und Optionen enthält."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "Beim Erstellen der Protokolldatei '{{.Path}}' ist ein Fehler aufgetreten:\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "Bei der Anforderung zum Erstellen eines Speicherauszugs ist ein Fehler aufgetreten:\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "Bei der Antwort bezüglich der Erstellung eines Speicherauszugs ist ein Fehler aufgetreten:\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "Schreibe den ersten Buchstaben der Beschreibung groß."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "Der Befehl „ '{{.Name}}' “ enthält ein Segment „ '{{.Segment}}' “, das weniger als „ {{.Count}} “ Zeichen umfasst. Jedes Wort in einem Befehl sollte mindestens {{.Count}} Zeichen lang sein."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "Der Befehl „ '{{.Name}}' “ verwendet weder ein gängiges Verb noch eine Pluralform. Verwenden Sie beispielsweise Verben wie: auflisten, erstellen, aktualisieren, löschen, anzeigen, abrufen, festlegen... oder Substantive im Plural wie „Instanzen“, „Dienste“"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "Der Befehl „ '{{.Name}}' “ hat keine Beschreibung. Alle Befehle müssen eine klare Beschreibung enthalten."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "Für den Befehl „ '{{.Name}}' “ liegen keine Verwendungshinweise vor."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "Der Befehl „ '{{.Name}}' “ weist {{.Level}} Ebenen auf, was das Maximum von {{.Level}} überschreitet. Tief verschachtelte Befehlshierarchien sind für Benutzer schwer zu merken und zu finden."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "Die Verwendung des Befehls „ '{{.Name}}' “ enthält Argumentwerte in Kleinbuchstaben: {{.Args}}. Die vom Benutzer eingegebenen Werte sollten in Großbuchstaben erfolgen."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "Der Befehl „ '{{.Name}}' “ verwendet einen reservierten Flaggennamen. Diese werden vom CLI-Framework verwaltet."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "Der Befehl „ '{{.Name}}' “ verwendet die Pluralform, doch aus der Beschreibung geht nicht eindeutig hervor, dass er eine Liste oder eine Gruppe von Elementen zurückgibt."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "Werte der Argumente in Großbuchstaben umwandeln (z. B. NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "Lesen der Eingabedaten nicht möglich: "
  },
  {
    "id": "Description",
    "translation": "Beschreibung"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "Die Beschreibung zu „ '{{.Name}}' “ umfasst {{.WordCount}} Wörter. Erwägen Sie, den Text auf weniger als {{.MaxWordCount}} Wörter zu beschränken, um eine bessere Darstellung zu erzielen."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "Die Beschreibung für „ '{{.Name}}' “ sollte mit einem Großbuchstaben beginnen."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "Die Beschreibung für „ '{{.Name}}' “ beginnt mit „ '{{.Bad}}' “. Verwende einen Satz ohne Subjekt."
  },
  {
    "id": "Elapsed:",
    "translation": "Verstrichen:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Die externe Authentifizierung ist fehlgeschlagen. Fehlercode: {{.ErrorCode}}, Meldung: {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "FEHLGESCHLAGEN"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "Fehlgeschlagen, Header konnte nicht in das csv-Format konvertiert werden"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "Fehlgeschlagen, Zeilen konnten nicht in das csv-Format konvertiert werden"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "Ungültiger Grant-Typ: "
  },
  {
    "id": "Invalid token: ",
    "translation": "Ungültiges Token: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) niedriger ist als das zulässige Minimum {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "Geben Sie 'j', 'n', 'ja' oder 'nein' ein."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "Geben Sie eine Zahl zwischen 1 und {{.Count}} ein."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "Geben Sie eine gültige Gleitkommazahl ein."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "Geben Sie eine gültige Zahl ein."
  },
  {
    "id": "Please enter value.",
    "translation": "Geben Sie einen Wert ein."
  },
  {
    "id": "REQUEST:",
    "translation": "ANFORDERUNG:"
  },
  {
    "id": "RESPONSE:",
    "translation": "ANTWORT:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Fehler auf dem fernen Server. Statuscode: {{.StatusCode}}, Fehlercode: {{.ErrorCode}}, Nachricht: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "Entfernen Sie „ '{{.Bad}}' “ und beginnen Sie direkt mit der Aktion. Beispiel: „Alle Instanzen auflisten“ statt „Dieser Befehl listet alle Instanzen auf“."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "Platzhalter aus dem Text zur Befehlsanwendung entfernen"
  },
  {
    "id": "Session inactive: ",
    "translation": "Sitzung inaktiv: "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "Setzen Sie „ MinCliVersion “ auf „ {{.Version}} “ oder höher, um die Kompatibilität mit den unterstützten CLI-Versionen sicherzustellen."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "Fassen Sie die Beschreibung prägnanter zusammen."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Speichern der Plug-in-Konfiguration nicht möglich: "
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "Ergänze die Beschreibung um Begriffe wie „Liste“, „anzeigen“, „darstellen“, „anzeigen“, „alle“ oder „mehrere“, um deutlich zu machen, dass mehrere Elemente zurückgegeben werden."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "Verwendung enthält Platzhalterargumente/Flags"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "Die Verwendung enthält nicht geschlossene {{.UnclosedGroup}} zwischen den Indikatoren {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "Die Verwendung enthält außer „COMMAND“ keinen weiteren Text"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "Verwenden Sie in Befehlsnamen gängige Verben wie „list“, „create“, „update“ und „delete“ oder verwenden Sie Pluralformen, um Auflistungsvorgänge zu kennzeichnen."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "Verwenden Sie aussagekräftigere Formulierungen mit mindestens {{.Count}} Zeichen pro Abschnitt."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} enthält die folgenden verbotenen Zeichen: {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} ist erforderlich"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} muss mindestens das Element {{.Param}} enthalten"
  }
]`)

func i18nResourcesAllDe_deJsonBytes() ([]byte, error) {
	return _i18nResourcesAllDe_deJson, nil
}

func i18nResourcesAllDe_deJson() (*asset, error) {
	bytes, err := i18nResourcesAllDe_deJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.de_DE.json", size: 10174, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllEn_usJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nEnter a number"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "Add a sentence without subject describing what the command does."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "Add usage text showing command syntax with parameters and options."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "An error occurred while dumping request:\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "An error occurred while dumping response:\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "Capitalize the first letter of the description."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "Command '{{.Name}}' has no description. All commands must have a clear description."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "Command '{{.Name}}' has no usage information."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "Could not read from input: "
  },
  {
    "id": "Description",
    "translation": "Description"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "Description for '{{.Name}}' should start with a capital letter."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject."
  },
  {
    "id": "Elapsed:",
    "translation": "Elapsed:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "FAILED"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "Failed, header could not convert to csv format"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "Failed, rows could not convert to csv format"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "Invalid grant type: "
  },
  {
    "id": "Invalid token: ",
    "translation": "Invalid token: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "Please enter 'y', 'n', 'yes' or 'no'."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "Please enter a number between 1 to {{.Count}}."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "Please enter a valid floating number."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "Please enter a valid number."
  },
  {
    "id": "Please enter value.",
    "translation": "Please enter value."
  },
  {
    "id": "REQUEST:",
    "translation": "REQUEST:"
  },
  {
    "id": "RESPONSE:",
    "translation": "RESPONSE:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "Remove placeholders from command usage text"
  },
  {
    "id": "Session inactive: ",
    "translation": "Session inactive: "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "Shorten the description to be more concise."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Unable to save plugin config: "
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "Usage contains placeholder arguments/flags"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "Usage does not have any usage text besides COMMAND"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "Use more descriptive words with at least {{.Count}} characters for each segment."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} contains the following forbidden characters: {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} is required"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} must contain at least {{.Param}} element"
  }
]`)

func i18nResourcesAllEn_usJsonBytes() ([]byte, error) {
	return _i18nResourcesAllEn_usJson, nil
}

func i18nResourcesAllEn_usJson() (*asset, error) {
	bytes, err := i18nResourcesAllEn_usJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.en_US.json", size: 9044, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllEs_esJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nEscriba un número"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "Añade una frase sin sujeto que describa qué hace el comando."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "Añade un texto de instrucciones que muestre la sintaxis del comando con los parámetros y las opciones."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "Se ha producido un error al crear el archivo de registro '{{.Path}}':\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "Se ha producido un error al volcar la solicitud:\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "Se ha producido un error al volcar la respuesta:\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "Escribe con mayúscula la primera letra de la descripción."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "El comando « '{{.Name}}' » contiene un segmento « '{{.Segment}}' » que tiene menos de « {{.Count}} » caracteres. Cada palabra de un comando debe tener al menos un {{.Count}} es de caracteres."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "La orden « '{{.Name}}' » no utiliza un verbo común ni la forma plural. Considera la posibilidad de utilizar verbos como: listar, crear, actualizar, eliminar, mostrar, obtener, establecer... o sustantivos en plural como «instancias», «servicios»"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "El comando « '{{.Name}}' » no tiene descripción. Todos los comandos deben ir acompañados de una descripción clara."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "El comando « '{{.Name}}' » no tiene información de uso."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "El comando « '{{.Name}}' » tiene un nivel de « {{.Level}} », lo que supera el máximo de « {{.Level}} ». Las jerarquías de comandos muy complejas resultan difíciles de recordar y de descubrir para los usuarios."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "El uso del comando « '{{.Name}}' » incluye valores de argumentos en minúsculas: {{.Args}}. Los valores introducidos por el usuario deben estar en MAYÚSCULAS."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "El comando « '{{.Name}}' » utiliza un nombre de indicador reservado. De esto se encarga el marco CLI."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "El comando « '{{.Name}}' » utiliza la forma plural, pero la descripción no indica claramente que devuelva una lista o un conjunto de elementos."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "Convierte los valores de los argumentos a mayúsculas (por ejemplo, NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "No se ha podido leer la entrada: "
  },
  {
    "id": "Description",
    "translation": "Descripción"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "La descripción de « '{{.Name}}' » contiene {{.WordCount}} palabras. Te recomendamos que limites el texto a menos de {{.MaxWordCount}} palabras para que se vea mejor."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "La descripción de « '{{.Name}}' » debe comenzar con mayúscula."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "La descripción de « '{{.Name}}' » comienza con « '{{.Bad}}' ». Escribe una frase sin sujeto."
  },
  {
    "id": "Elapsed:",
    "translation": "Transcurrido:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Ha fallado la autenticación externa. Código de error: {{.ErrorCode}} mensaje: {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "ERROR"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "Fallido, la cabecera no se ha podido convertir a formato csv"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "Fallo, las filas no se han podido convertir a formato csv"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "Tipo de subvención no válido: "
  },
  {
    "id": "Invalid token: ",
    "translation": "Señal no válida: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) es inferior al mínimo permitido {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "Correcto"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "Especifique 'y', 'n', 'yes' o 'no'."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "Especifique un número entre 1 y {{.Count}}."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "Especifique un número flotante válido."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "Especifique un número válido."
  },
  {
    "id": "Please enter value.",
    "translation": "Especifique un valor."
  },
  {
    "id": "REQUEST:",
    "translation": "SOLICITUD:"
  },
  {
    "id": "RESPONSE:",
    "translation": "RESPUESTA:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Error del servidor remoto. Código de estado: {{.StatusCode}}, código de error: {{.ErrorCode}}, mensaje: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "Elimina « '{{.Bad}}' » y empieza directamente con la acción. Ejemplo: «Mostrar todas las instancias» en lugar de «Este comando muestra todas las instancias»."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "Eliminar los marcadores de posición del texto de instrucciones de uso"
  },
  {
    "id": "Session inactive: ",
    "translation": "Sesión inactiva: "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "Establece el parámetro `+"`"+` MinCliVersion `+"`"+` en `+"`"+` {{.Version}} `+"`"+` o superior para garantizar la compatibilidad con las versiones de la CLI compatibles."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "Acorta la descripción para que sea más concisa."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "No se ha podido guardar la configuración del plugin:"
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "Actualiza la descripción para incluir palabras como «lista», «mostrar», «visualizar», «ver», «todos» o «varios», a fin de aclarar que devuelve varios elementos."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "El uso contiene argumentos/flags"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "El uso contiene {{.UnclosedGroup}} sin cerrar entre indicies {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "El uso no tiene ningún texto de ayuda, salvo el comando"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "Utiliza verbos comunes en los nombres de los comandos, como «listar», «crear», «actualizar» o «eliminar», o emplea formas en plural para indicar operaciones de listado."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "Utiliza palabras más descriptivas, con al menos {{.Count}} caracteres por cada segmento."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} contiene los siguientes caracteres prohibidos: {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} es necesario"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} debe contener al menos {{.Param}} elemento"
  }
]`)

func i18nResourcesAllEs_esJsonBytes() ([]byte, error) {
	return _i18nResourcesAllEs_esJson, nil
}

func i18nResourcesAllEs_esJson() (*asset, error) {
	bytes, err := i18nResourcesAllEs_esJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.es_ES.json", size: 9852, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllFr_frJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nEntrez un nombre"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "Ajoutez une phrase sans sujet décrivant ce que fait la commande."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "Ajouter un texte d'aide indiquant la syntaxe de la commande avec ses paramètres et ses options."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "Erreur lors de la création du fichier journal '{{.Path}}':\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "Erreur lors de la demande de vidage :\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "Erreur lors de la réponse de vidage :\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "Mettez une majuscule à la première lettre de la description."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "La commande « '{{.Name}}' » contient un segment « '{{.Segment}}' » dont la longueur est inférieure à {{.Count}} caractères. Chaque mot d'une commande doit comporter au moins {{.Count}} caractères."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "La forme impérative « '{{.Name}}' » n'utilise ni un verbe courant ni la forme plurielle. Pensez à utiliser des verbes tels que : lister, créer, mettre à jour, supprimer, afficher, récupérer, définir... ou des noms pluriels tels que « instances », « services »"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "La commande « '{{.Name}}' » ne comporte pas de description. Toutes les commandes doivent être accompagnées d'une description claire."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "La commande « '{{.Name}}' » ne dispose d'aucune information d'utilisation."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "La commande « '{{.Name}}' » comporte {{.Level}} niveaux, ce qui dépasse le maximum de {{.Level}}. Les hiérarchies de commandes trop complexes sont difficiles à mémoriser et à comprendre pour les utilisateurs."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "L'utilisation de la commande « '{{.Name}}' » contient des valeurs d'argument en minuscules : {{.Args}}. Les valeurs saisies par l'utilisateur doivent être en MAJUSCULES."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "La commande « '{{.Name}}' » utilise un nom de drapeau réservé. Ces éléments sont gérés par le framework CLI."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "La commande « '{{.Name}}' » utilise le pluriel, mais la description n'indique pas clairement qu'elle renvoie une liste ou un ensemble d'éléments."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "Convertir les valeurs des arguments en MAJUSCULES (par exemple, NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "Lecture impossible à partir de l'entrée : "
  },
  {
    "id": "Description",
    "translation": "Description"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "La description de « '{{.Name}}' » compte {{.WordCount}} mots. Pour un meilleur affichage, pensez à limiter la longueur de vos phrases à moins d' {{.MaxWordCount}}."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "La description de '{{.Name}}' doit commencer par une majuscule."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "La description de '{{.Name}}' commence par '{{.Bad}}'. Utilisez une phrase sans sujet."
  },
  {
    "id": "Elapsed:",
    "translation": "Ecoulé :"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Échec d'authentification externe. Code d'erreur : {{.ErrorCode}}, message : {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "ECHEC"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "Échec, l'en-tête n'a pas pu être converti au format csv"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "Échec, les lignes n'ont pas pu être converties au format csv"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "Type de subvention non valide : "
  },
  {
    "id": "Invalid token: ",
    "translation": "Jeton non valide : "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) est inférieur au minimum autorisé {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "L'entrée doit être 'y', 'n', 'yes' ou 'no'."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "Entrez un nombre entre 1 et {{.Count}}."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "Entrez un nombre flottant valide."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "Entrez un nombre valide."
  },
  {
    "id": "Please enter value.",
    "translation": "Entrez une valeur."
  },
  {
    "id": "REQUEST:",
    "translation": "DEMANDE :"
  },
  {
    "id": "RESPONSE:",
    "translation": "REPONSE :"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Erreur du serveur distant. Code de statut : {{.StatusCode}}, code d'erreur : {{.ErrorCode}}, message : {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "Supprimez « '{{.Bad}}' » et passez directement à l'action. Exemple : « Lister toutes les instances » au lieu de « Cette commande liste toutes les instances »."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "Supprimer les espaces réservés du texte d'utilisation de la commande"
  },
  {
    "id": "Session inactive: ",
    "translation": "Session inactive : "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "Définissez la variable d'environnement `+"`"+` MinCliVersion `+"`"+` sur `+"`"+` {{.Version}} `+"`"+` ou une version supérieure pour garantir la compatibilité avec les versions prises en charge de l'interface CLI."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "Raccourcissez la description pour la rendre plus concise."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Impossible d'enregistrer la configuration du plug-in : "
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "Modifiez la description en y ajoutant des mots tels que « liste », « afficher », « visualiser », « tout » ou « plusieurs » afin de préciser qu'elle renvoie plusieurs éléments."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "L'utilisation contient des arguments et des drapeaux de substitution"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "L'utilisation contient des {{.UnclosedGroup}} non fermés entre les indicateurs {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "La commande ne comporte aucun texte d'aide, à l'exception de COMMAND"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "Utilisez des verbes courants dans les noms de commandes, tels que « list », « create », « update » ou « delete », ou employez le pluriel pour désigner des opérations de liste."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "Utilisez des mots plus descriptifs, avec au moins {{.Count}} caractères par segment."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} contient les caractères interdits suivants : {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} est nécessaire"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} doit contenir au moins {{.Param}} élément"
  }
]`)

func i18nResourcesAllFr_frJsonBytes() ([]byte, error) {
	return _i18nResourcesAllFr_frJson, nil
}

func i18nResourcesAllFr_frJson() (*asset, error) {
	bytes, err := i18nResourcesAllFr_frJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.fr_FR.json", size: 10019, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllIt_itJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nImmetti un numero"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "Aggiungi una frase senza soggetto che descriva cosa fa il comando."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "Aggiungere un testo descrittivo che illustri la sintassi del comando con i parametri e le opzioni."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "Si è verificato un errore durante la creazione del file di log '{{.Path}}':\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "Si è verificato un errore durante il dump della richiesta:\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "Si è verificato un errore durante il dump della risposta:\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "Scrivi la prima lettera della descrizione in maiuscolo."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "Il comando '{{.Name}}' contiene un segmento '{{.Segment}}' che ha una lunghezza inferiore a {{.Count}} caratteri. Ogni parola di un comando deve essere composta da almeno {{.Count}} caratteri."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "Il comando \" '{{.Name}}' \" non utilizza né un verbo comune né la forma plurale. Si consiglia di utilizzare verbi come: elencare, creare, aggiornare, eliminare, visualizzare, recuperare, impostare... o sostantivi plurali come «istanze», «servizi»"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "Il comando \" '{{.Name}}' \" non ha una descrizione. Tutti i comandi devono essere accompagnati da una descrizione chiara."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "Il comando \" '{{.Name}}' \" non dispone di informazioni sull'utilizzo."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "Il comando \" '{{.Name}}' \" presenta livelli di \" {{.Level}} \", superando il limite massimo di \" {{.Level}} \". Le gerarchie di comandi complesse sono difficili da ricordare e da individuare per gli utenti."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "L'uso del comando `+"`"+` '{{.Name}}' `+"`"+` prevede valori degli argomenti in minuscolo: {{.Args}}. I valori inseriti dall'utente devono essere scritti in MAIUSCOLO."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "Il comando `+"`"+` '{{.Name}}' `+"`"+` utilizza un nome di flag riservato. Questi sono gestiti dal framework CLI."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "Il comando \" '{{.Name}}' \" utilizza la forma plurale, ma la descrizione non indica chiaramente che restituisce un elenco o un gruppo di elementi."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "Converti i valori degli argomenti in MAIUSCOLO (ad es. NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "Impossibile leggere dall'input: "
  },
  {
    "id": "Description",
    "translation": "Descrizione"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "La descrizione di \" '{{.Name}}' \" contiene {{.WordCount}} parole. Si consiglia di limitare il numero di parole a meno di {{.MaxWordCount}} per una migliore visualizzazione."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "La descrizione di '{{.Name}}' dovrebbe iniziare con una maiuscola."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "La descrizione di \" '{{.Name}}' \" inizia con \" '{{.Bad}}' \". Usa una frase senza soggetto."
  },
  {
    "id": "Elapsed:",
    "translation": "Trascorso:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Autenticazione esterna non riuscita. Codice di errore: {{.ErrorCode}}, messaggio: {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "NON RIUSCITO"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "Fallito, l'intestazione non può essere convertita in formato csv"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "Fallito, non è stato possibile convertire le righe in formato csv"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "Tipo di concessione non valido: "
  },
  {
    "id": "Invalid token: ",
    "translation": "Token non valido: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) è inferiore al minimo consentito {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "Immetti 's', 'n', 'sì' o 'no'."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "Immetti un numero compreso tra 1 e {{.Count}}."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "Immetti un numero decimale valido."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "Immetti un numero valido."
  },
  {
    "id": "Please enter value.",
    "translation": "Immetti un valore."
  },
  {
    "id": "REQUEST:",
    "translation": "RICHIESTA:"
  },
  {
    "id": "RESPONSE:",
    "translation": "RISPOSTA:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Errore server remoto. Codice di stato: {{.StatusCode}}, codice di errore: {{.ErrorCode}}, messaggio: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "Elimina \" '{{.Bad}}' \" e passa direttamente all'azione. Esempio: «Elenca tutte le istanze» invece di «Questo comando elenca tutte le istanze»."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "Rimuovere i segnaposto dal testo relativo all'uso del comando"
  },
  {
    "id": "Session inactive: ",
    "translation": "Sessione inattiva: "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "Imposta il parametro `+"`"+` MinCliVersion `+"`"+` su `+"`"+` {{.Version}} `+"`"+` o una versione superiore per garantire la compatibilità con le versioni supportate della CLI."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "Accorcia la descrizione per renderla più concisa."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Impossibile salvare la configurazione del plug-in: "
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "Aggiornare la descrizione includendo termini come \"elenco\", \"mostra\", \"visualizza\", \"tutti\" o \"più\" per chiarire che restituisce più elementi."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "L'uso contiene argomenti segnaposto/flags"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "L'uso contiene {{.UnclosedGroup}} non chiuso tra gli indici {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "L'uso non prevede alcun testo descrittivo oltre a COMMAND"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "Utilizza verbi comuni nei nomi dei comandi, come \"list\", \"create\", \"update\" e \"delete\", oppure usa la forma plurale per indicare operazioni di elenco."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "Utilizza parole più descrittive, con almeno {{.Count}} caratteri per ogni segmento."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} contiene i seguenti caratteri vietati: {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} è necessario"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} deve contenere almeno {{.Param}} elemento"
  }
]`)

func i18nResourcesAllIt_itJsonBytes() ([]byte, error) {
	return _i18nResourcesAllIt_itJson, nil
}

func i18nResourcesAllIt_itJson() (*asset, error) {
	bytes, err := i18nResourcesAllIt_itJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.it_IT.json", size: 9790, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllJa_jpJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n数値を入力してください"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "そのコマンドが何をするのかを説明する、主語のない文を追加してください。"
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "パラメータやオプションを含むコマンドの構文を示す説明文を追加してください。"
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "ログ・ファイル '{{.Path}}' を作成中にエラーが発生しました:\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "要求のダンプ中にエラーが発生しました:\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "応答のダンプ中にエラーが発生しました:\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "説明文の最初の文字を大文字にしてください。"
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "コマンド `+"`"+` '{{.Name}}' `+"`"+` には、 {{.Count}} 文字未満のセグメント `+"`"+` '{{.Segment}}' `+"`"+` が含まれています。 コマンド内の各単語は、少なくとも {{.Count}} 文字以上である必要があります。"
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "「 '{{.Name}}' 」という命令には、一般的な動詞や複数形は使われません。 「list」「create」「update」「delete」「show」「get」「set」などの動詞の使用を検討してください。 または「インスタンス」や「サービス」のような複数形の名詞"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "コマンド「 '{{.Name}}' 」には説明がありません。 すべてのコマンドには、明確な説明が必要です。"
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "コマンド「 '{{.Name}}' 」には使用方法に関する情報がありません。"
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "コマンド「 '{{.Name}}' 」には、 {{.Level}} レベルのデータがあり、最大値である {{.Level}} を超えています。 階層が深いコマンド体系は、ユーザーが覚えたり見つけたりするのが難しい。"
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "コマンド「 '{{.Name}}' 」の使用例には、小文字の引数値が含まれています： {{.Args}}。 ユーザーが入力する値は、すべて大文字にしてください。"
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "コマンド「 '{{.Name}}' 」は、予約済みのフラグ名を使用しています。 これらはCLIフレームワークによって処理されます。"
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "コマンド「 '{{.Name}}' 」は複数形を使用していますが、説明文にはリストや項目のグループが返されることが明確に示されていません。"
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "引数の値を大文字に変換します（例：NAME、INSTANCE_ID、FORMAT）。"
  },
  {
    "id": "Could not read from input: ",
    "translation": "入力から読み取れませんでした。 "
  },
  {
    "id": "Description",
    "translation": "説明"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "「 '{{.Name}}' 」の説明文には、 {{.WordCount}} 語が含まれています。 表示を最適化するため、文字数を {{.MaxWordCount}} 文字未満に抑えることをご検討ください。"
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "'{{.Name}}' の説明文は、大文字で始める必要があります。"
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "'{{.Name}}' の説明は、 '{{.Bad}}' から始まります。 主語のない文を使ってください。"
  },
  {
    "id": "Elapsed:",
    "translation": "経過:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "外部認証に失敗しました。 エラーコード {{.ErrorCode}} メッセージ {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "失敗"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "ヘッダがcsv形式に変換できませんでした"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "行をcsv形式に変換できませんでした"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "無効なグラントタイプです： "
  },
  {
    "id": "Invalid token: ",
    "translation": "トークンが無効です: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) は許容最小値より低い。 {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "「y」、「n」、「yes」、または「no」を入力してください。"
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "1 から {{.Count}} までの数値を入力してください。"
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "有効な浮動小数点数を入力してください。"
  },
  {
    "id": "Please enter a valid number.",
    "translation": "有効な数値を入力してください。"
  },
  {
    "id": "Please enter value.",
    "translation": "値を入力してください。"
  },
  {
    "id": "REQUEST:",
    "translation": "要求:"
  },
  {
    "id": "RESPONSE:",
    "translation": "応答:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "リモート・サーバー・エラー。 状況コード: {{.StatusCode}}、エラー・コード: {{.ErrorCode}}、メッセージ: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "'{{.Bad}}' を削除し、アクションを直接開始してください。 例：「このコマンドはすべてのインスタンスを表示します」ではなく、「すべてのインスタンスを表示」とする。"
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "コマンドの使用方法の説明文からプレースホルダーを削除する"
  },
  {
    "id": "Session inactive: ",
    "translation": "セッションは不活発： "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "サポートされているCLIバージョンとの互換性を確保するため、 MinCliVersion を {{.Version}} 以上に設定してください。"
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "説明文を簡潔にするために、内容を短縮してください。"
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "プラグイン構成を保存できません: "
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "説明文に「一覧」、「表示」、「表示する」、「閲覧」、「すべて」、「複数」などの単語を追加し、複数の項目が返されることを明確にしてください。"
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "使用法にはプレースホルダ引数/フラグが含まれる"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "使用法には、インジケーターの間に閉じていない {{.UnclosedGroup}} が含まれる。 {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "このコマンドには、COMMAND 以外の説明文はありません"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "コマンド名には、「list」「create」「update」「delete」などの一般的な動詞を使用するか、一覧表示操作を示すために複数形を使用してください。"
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "各セグメントには、少なくとも {{.Count}} 文字の、より具体的な言葉を使用してください。"
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} には以下の禁止文字が含まれている： {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} が必要です"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} は少なくとも {{.Param}} の要素を含んでいなければならない"
  }
]`)

func i18nResourcesAllJa_jpJsonBytes() ([]byte, error) {
	return _i18nResourcesAllJa_jpJson, nil
}

func i18nResourcesAllJa_jpJson() (*asset, error) {
	bytes, err := i18nResourcesAllJa_jpJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.ja_JP.json", size: 10762, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllKo_krJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n번호 입력"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "명령어가 무엇을 하는지 설명하는, 주어가 없는 문장을 추가하세요."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "매개변수와 옵션을 포함한 명령어 구문을 설명하는 사용법을 추가하세요."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "'{{.Path}}' 로그 파일을 작성할 때 다음 오류가 발생했습니다. \n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "요청을 덤프할 때 다음 오류가 발생했습니다. \n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "응답을 덤프할 때 다음 오류가 발생했습니다. \n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "설명문의 첫 글자를 대문자로 표기하십시오."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "명령어 `+"`"+` '{{.Name}}' `+"`"+`에는 `+"`"+` {{.Count}} `+"`"+`자보다 짧은 `+"`"+` '{{.Segment}}' `+"`"+` 세그먼트가 포함되어 있습니다. 명령어의 각 단어는 최소 {{.Count}} 자 이상이어야 합니다."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "'{{.Name}}' 라는 명령어에는 일반적인 동사나 복수형이 사용되지 않습니다. 다음과 같은 동사를 사용하는 것을 고려해 보세요: 나열하다, 생성하다, 업데이트하다, 삭제하다, 표시하다, 가져오다, 설정하다... 또는 'instances', 'services'와 같은 복수 명사"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "명령어 `+"`"+` '{{.Name}}' `+"`"+`에는 설명이 없습니다. 모든 명령어에는 명확한 설명이 포함되어야 합니다."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "'{{.Name}}' 명령어에 대한 사용법이 없습니다."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "'{{.Name}}' 명령어의 레벨이 {{.Level}} 개이며, 이는 {{.Level}} 의 최대값을 초과합니다. 복잡한 명령 계층 구조는 사용자가 기억하거나 파악하기 어렵습니다."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "'{{.Name}}' 명령어의 사용법에 소문자 인수가 포함되어 있습니다: {{.Args}}. 사용자 입력 값은 대문자로 입력해야 합니다."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "'{{.Name}}' 명령어는 예약된 플래그 이름을 사용합니다. 이는 CLI 프레임워크에서 처리합니다."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "'{{.Name}}' 명령어는 복수형을 사용하지만, 설명에는 이 명령어가 항목의 목록이나 그룹을 반환한다는 점이 명확하게 명시되어 있지 않습니다."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "인수의 값을 대문자로 변환합니다(예: NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "입력에서 읽지 못함: "
  },
  {
    "id": "Description",
    "translation": "설명"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "'{{.Name}}' 에 대한 설명에는 {{.WordCount}} 개의 단어가 포함되어 있습니다. 더 나은 표시를 위해 단어 수를 {{.MaxWordCount}} 자 미만으로 제한하는 것을 고려해 보세요."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "'{{.Name}}' 의 설명은 대문자로 시작해야 합니다."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "'{{.Name}}' 에 대한 설명은 '{{.Bad}}' 에서 시작됩니다. 주어가 없는 문장을 사용하세요."
  },
  {
    "id": "Elapsed:",
    "translation": "경과 시간:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "외부 인증에 실패했습니다. 오류 코드: {{.ErrorCode}}, 메시지: {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "실패"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "실패, 헤더를 CSV 형식으로 변환할 수 없습니다"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "실패, 행을 CSV 형식으로 변환할 수 없음"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "잘못된 보조금 유형입니다: "
  },
  {
    "id": "Invalid token: ",
    "translation": "올바르지 않은 토큰: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} )가 허용된 최소값보다 낮습니다 {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "확인"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "'y', 'n', '예' 또는 '아니오'를 입력하십시오."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "1 - {{.Count}} 사이의 수를 입력하십시오."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "올바른 float 수를 입력하십시오."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "올바른 수를 입력하십시오."
  },
  {
    "id": "Please enter value.",
    "translation": "값을 입력하십시오."
  },
  {
    "id": "REQUEST:",
    "translation": "요청:"
  },
  {
    "id": "RESPONSE:",
    "translation": "응답:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "원격 서버 오류가 발생했습니다. 상태 코드: {{.StatusCode}}, 오류 코드: {{.ErrorCode}}, 메시지: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "'{{.Bad}}' 를 제거하고 바로 액션으로 넘어가세요. 예: '이 명령은 모든 인스턴스를 나열합니다' 대신 '모든 인스턴스 나열'."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "명령어 사용법 설명에서 자리 표시자를 제거합니다"
  },
  {
    "id": "Session inactive: ",
    "translation": "세션이 비활성 상태입니다: "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "지원되는 CLI 버전과의 호환성을 보장하려면 `+"`"+` MinCliVersion `+"`"+`을 `+"`"+` {{.Version}} `+"`"+` 이상으로 설정하십시오."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "설명을 더 간결하게 줄이세요."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "플러그인 구성을 저장할 수 없음:"
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "설명문에 '목록', '표시', '표시하기', '보기', '전체' 또는 '여러 개'와 같은 단어를 추가하여 여러 개의 항목이 반환됨을 명확히 하세요."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "사용법에는 자리 표시자 인수/플래그가 포함됩니다"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "사용법에는 표시 사이에 닫히지 않은 {{.UnclosedGroup}} {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "사용법에는 COMMAND 외에 다른 설명이 없습니다"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "명령어 이름에는 list, create, update, delete와 같은 일반적인 동사를 사용하거나, 목록 표시 작업을 나타내기 위해 복수형을 사용하십시오."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "각 문단마다 최소 {{.Count}} 자 이상의 설명적인 단어를 사용하세요."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} 에는 다음과 같은 금지 문자가 포함되어 있습니다: {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} 필수"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} 최소 {{.Param}} 요소를 포함해야 합니다"
  }
]`)

func i18nResourcesAllKo_krJsonBytes() ([]byte, error) {
	return _i18nResourcesAllKo_krJson, nil
}

func i18nResourcesAllKo_krJson() (*asset, error) {
	bytes, err := i18nResourcesAllKo_krJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.ko_KR.json", size: 10173, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllPt_brJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nInsira um número"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "Adicione uma frase sem sujeito que descreva o que o comando faz."
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "Adicione um texto de instruções que mostre a sintaxe do comando com parâmetros e opções."
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "Ocorreu um erro ao criar o arquivo de log '{{.Path}}':\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "Ocorreu um erro ao fazer dump da solicitação:\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "Ocorreu um erro ao fazer dump da resposta:\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "Escreva a primeira letra da descrição com maiúscula."
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "O comando `+"`"+` '{{.Name}}' `+"`"+` contém um segmento `+"`"+` '{{.Segment}}' `+"`"+` com menos de `+"`"+` {{.Count}} `+"`"+` caracteres. Cada palavra em um comando deve ter pelo menos {{.Count}} es."
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "A expressão \" '{{.Name}}' \" não utiliza um verbo comum nem a forma plural. Considere usar verbos como: listar, criar, atualizar, excluir, exibir, obter, definir... ou substantivos no plural, como \"instâncias\", \"serviços\""
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "O comando `+"`"+` '{{.Name}}' `+"`"+` não possui descrição. Todos os comandos devem ter uma descrição clara."
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "O comando `+"`"+` '{{.Name}}' `+"`"+` não possui informações de uso."
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "O comando `+"`"+` '{{.Name}}' `+"`"+` possui `+"`"+` {{.Level}} `+"`"+` níveis, excedendo o máximo de `+"`"+` {{.Level}} `+"`"+`. Hierarquias de comandos muito complexas são difíceis de memorizar e de entender para os usuários."
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "O uso do comando `+"`"+` '{{.Name}}' `+"`"+` contém valores de argumentos em minúsculas: {{.Args}}. Os valores inseridos pelo usuário devem estar em MAIÚSCULAS."
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "O comando `+"`"+` '{{.Name}}' `+"`"+` utiliza um nome de sinalizador reservado. Isso é feito pela estrutura CLI."
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "O comando '{{.Name}}' usa a forma plural, mas a descrição não indica claramente que ele retorna uma lista ou um grupo de itens."
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "Converta os valores dos argumentos para letras maiúsculas (por exemplo, NAME, INSTANCE_ID, FORMAT)."
  },
  {
    "id": "Could not read from input: ",
    "translation": "Não foi possível ler a apartir da entrada: "
  },
  {
    "id": "Description",
    "translation": "Descrição"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "A descrição de “ '{{.Name}}' ” tem {{.WordCount}} palavras. Considere limitar o texto a menos de {{.MaxWordCount}} palavras para uma melhor exibição."
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "A descrição de '{{.Name}}' deve começar com letra maiúscula."
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "A descrição de '{{.Name}}' começa com '{{.Bad}}'. Use uma frase sem sujeito."
  },
  {
    "id": "Elapsed:",
    "translation": "Decorrido:"
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Falha na autenticação externa. Código de erro: {{.ErrorCode}}, mensagem: {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "COM FALHA"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "Falha, o cabeçalho não pôde ser convertido para o formato csv"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "Falha, não foi possível converter as linhas para o formato csv"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "Tipo de concessão inválido: "
  },
  {
    "id": "Invalid token: ",
    "translation": "Token inválido: "
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) é menor do que o mínimo permitido {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "Insira 'y', 'n', 'yes' ou 'no'."
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "Insira um número entre 1 e {{.Count}}."
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "Insira um número flutuante válido."
  },
  {
    "id": "Please enter a valid number.",
    "translation": "Digite um número válido."
  },
  {
    "id": "Please enter value.",
    "translation": "Insira um valor."
  },
  {
    "id": "REQUEST:",
    "translation": "SOLICITAÇÃO:"
  },
  {
    "id": "RESPONSE:",
    "translation": "RESPOSTA:"
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "Erro do servidor remoto. Código de status: {{.StatusCode}}, código de erro: {{.ErrorCode}}, mensagem: {{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "Remova o \" '{{.Bad}}' \" e vá direto ao ponto. Exemplo: \"Listar todas as instâncias\" em vez de \"Este comando lista todas as instâncias\"."
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "Remover os marcadores de lugar do texto de instruções de uso do comando"
  },
  {
    "id": "Session inactive: ",
    "translation": "Sessão inativa: "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "Defina `+"`"+` MinCliVersion `+"`"+` como `+"`"+` {{.Version}} `+"`"+` ou superior para garantir a compatibilidade com as versões da CLI suportadas."
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "Resuma a descrição para torná-la mais concisa."
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Não é possível salvar a configuração do plug-in: "
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "Atualize a descrição para incluir palavras como \"lista\", \"mostrar\", \"exibir\", \"visualizar\", \"todos\" ou \"vários\", a fim de deixar claro que ela retorna vários itens."
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "O uso contém argumentos/flags de espaço reservado"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "O uso contém {{.UnclosedGroup}} não fechado entre os índices {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "A sintaxe não possui nenhum texto de instruções além de COMMAND"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "Use verbos comuns nos nomes dos comandos, como listar, criar, atualizar, excluir, ou utilize formas no plural para indicar operações de listagem."
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "Use palavras mais descritivas, com pelo menos {{.Count}} caracteres em cada segmento."
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} contém os seguintes caracteres proibidos: {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} é necessário"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} deve conter pelo menos {{.Param}} elemento"
  }
]`)

func i18nResourcesAllPt_brJsonBytes() ([]byte, error) {
	return _i18nResourcesAllPt_brJson, nil
}

func i18nResourcesAllPt_brJson() (*asset, error) {
	bytes, err := i18nResourcesAllPt_brJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.pt_BR.json", size: 9628, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllZh_hansJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n请输入数字"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "添加一个省略主语的句子，说明该命令的功能。"
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "添加使用说明，展示包含参数和选项的命令语法。"
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "创建日志文件“{{.Path}}”时发生错误：\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "转储请求时发生错误：\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "转储响应时发生错误：\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "描述的首字母应大写。"
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "指令 '{{.Name}}' 包含一个长度小于 {{.Count}} 个字符的片段 '{{.Segment}}'。 命令中的每个单词长度应至少为 {{.Count}} 个字符。"
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "命令“ '{{.Name}}' ”既不使用常见动词，也不采用复数形式。 建议使用以下动词：列出、创建、更新、删除、显示、获取、设置…… 或复数名词，如“实例”、“服务”"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "命令 '{{.Name}}' 没有描述。 所有命令都必须有清晰的说明。"
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "命令 '{{.Name}}' 没有用法信息。"
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "命令 `+"`"+` '{{.Name}}' `+"`"+` 包含 `+"`"+` {{.Level}} `+"`"+` 个级别，超过了 `+"`"+` {{.Level}} `+"`"+` 的最大值。 层级过深的命令结构让用户难以记忆和掌握。"
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "命令 '{{.Name}}' 的用法中包含小写参数值： {{.Args}}。 用户输入的值应为大写字母。"
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "命令 '{{.Name}}' 使用了预留的标志名称。 这些由 CLI 框架负责处理。"
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "'{{.Name}}' 命令使用了复数形式，但说明中并未明确指出它会返回一个列表或一组项目。"
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "将参数值转换为大写（例如：NAME、INSTANCE_ID、FORMAT）。"
  },
  {
    "id": "Could not read from input: ",
    "translation": "无法从输入进行读取： "
  },
  {
    "id": "Description",
    "translation": "描述"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "'{{.Name}}' 的描述包含 {{.WordCount}} 个单词。 为获得更好的显示效果，建议将字数控制在 {{.MaxWordCount}} 字以内。"
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "'{{.Name}}' 的描述应以大写字母开头。"
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "'{{.Name}}' 的说明以 '{{.Bad}}' 开头。 请用一个没有主语的句子。"
  },
  {
    "id": "Elapsed:",
    "translation": "经过时长："
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "外部认证失败。 错误代码： {{.ErrorCode}} 信息： {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "失败"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "失败，标头无法转换为 csv 格式"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "失败，数据行无法转换为 csv 格式"
  },
  {
    "id": "ID",
    "translation": "标识"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "授予类型无效： "
  },
  {
    "id": "Invalid token: ",
    "translation": "令牌无效："
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) 低于允许的最小值。 {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "确定"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "请输入“y”、“n”、“yes”或“no”。"
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "请输入 1 到 {{.Count}} 之间的数字。"
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "请输入有效的浮点数。"
  },
  {
    "id": "Please enter a valid number.",
    "translation": "请输入有效的数字。"
  },
  {
    "id": "Please enter value.",
    "translation": "请输入有效的值。"
  },
  {
    "id": "REQUEST:",
    "translation": "请求: "
  },
  {
    "id": "RESPONSE:",
    "translation": "响应: "
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "远程服务器错误。状态码：{{.StatusCode}}，错误代码：{{.ErrorCode}}，消息：{{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "删除 '{{.Bad}}' ，直接从该操作开始。 示例：使用“列出所有实例”，而不是“此命令列出所有实例”。"
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "从命令用法说明中删除占位符"
  },
  {
    "id": "Session inactive: ",
    "translation": "会议非活动： "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "将 `+"`"+` MinCliVersion `+"`"+` 设置为 `+"`"+` {{.Version}} `+"`"+` 或更高版本，以确保与受支持的 CLI 版本兼容。"
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "请将描述内容精简，使其更加简洁。"
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "无法保存插件配置："
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "更新描述，加入“列表”、“显示”、“查看”、“全部”或“多个”等词语，以明确说明该操作会返回多个项目。"
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "用法包含占位参数/标志"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "用法中包含指示符之间未封闭的 {{.UnclosedGroup}} {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "除 COMMAND 之外，该用法没有其他说明文本"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "在命令名称中使用常见动词，例如 list、create、update、delete，或使用复数形式来表示列表操作。"
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "请使用更具描述性的词语，每个段落至少包含 {{.Count}} 个字符。"
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} 包含以下禁用字符： {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} 需要"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} 必须至少包含 {{.Param}} 元素"
  }
]`)

func i18nResourcesAllZh_hansJsonBytes() ([]byte, error) {
	return _i18nResourcesAllZh_hansJson, nil
}

func i18nResourcesAllZh_hansJson() (*asset, error) {
	bytes, err := i18nResourcesAllZh_hansJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.zh_Hans.json", size: 9052, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllZh_hantJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n請輸入數字"
  },
  {
    "id": "Add a sentence without subject describing what the command does.",
    "translation": "請添加一句不帶主語的句子，說明該指令的功能。"
  },
  {
    "id": "Add usage text showing command syntax with parameters and options.",
    "translation": "新增使用說明，展示包含參數與選項的指令語法。"
  },
  {
    "id": "An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n",
    "translation": "建立日誌檔 '{{.Path}}' 時發生錯誤：\n{{.Error}}\n\n"
  },
  {
    "id": "An error occurred while dumping request:\n{{.Error}}\n",
    "translation": "傾出要求時發生錯誤：\n{{.Error}}\n"
  },
  {
    "id": "An error occurred while dumping response:\n{{.Error}}\n",
    "translation": "傾出回應時發生錯誤：\n{{.Error}}\n"
  },
  {
    "id": "Capitalize the first letter of the description.",
    "translation": "描述文字的首字母應大寫。"
  },
  {
    "id": "Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 包含一個長度小於 `+"`"+` {{.Count}} `+"`"+` 個字元的區段 `+"`"+` '{{.Segment}}' `+"`"+`。 指令中的每個單詞長度應至少為 {{.Count}} 個字元。"
  },
  {
    "id": "Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'",
    "translation": "命令「 '{{.Name}}' 」不使用常見的動詞或複數形式。 建議使用以下動詞：列出、建立、更新、刪除、顯示、取得、設定…… 或複數名詞，例如「實例」、「服務」"
  },
  {
    "id": "Command '{{.Name}}' has no description. All commands must have a clear description.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 沒有說明。 所有指令都必須附有明確的說明。"
  },
  {
    "id": "Command '{{.Name}}' has no usage information.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 沒有使用說明。"
  },
  {
    "id": "Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 包含 `+"`"+` {{.Level}} `+"`"+` 個層級，超過了 `+"`"+` {{.Level}} `+"`"+` 的最大值。 層級過深的指令結構，使用者往往難以記住和掌握。"
  },
  {
    "id": "Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 的用法包含小寫參數值： {{.Args}}。 使用者輸入的數值應使用大寫字母。"
  },
  {
    "id": "Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 使用了保留的旗標名稱。 這些由 CLI 框架負責處理。"
  },
  {
    "id": "Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.",
    "translation": "指令 `+"`"+` '{{.Name}}' `+"`"+` 使用了複數形式，但說明文字並未明確指出它會傳回一組項目。"
  },
  {
    "id": "Convert arguments values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
    "translation": "將參數值轉換為大寫字母（例如：NAME、INSTANCE_ID、FORMAT）。"
  },
  {
    "id": "Could not read from input: ",
    "translation": "無法從輸入讀取： "
  },
  {
    "id": "Description",
    "translation": "說明"
  },
  {
    "id": "Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.",
    "translation": "'{{.Name}}' 的說明文字共有 {{.WordCount}} 個字。 為獲得更好的顯示效果，建議將字數限制在 {{.MaxWordCount}} 字以內。"
  },
  {
    "id": "Description for '{{.Name}}' should start with a capital letter.",
    "translation": "'{{.Name}}' 的描述應以大寫字母開頭。"
  },
  {
    "id": "Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.",
    "translation": "'{{.Name}}' 的說明以 '{{.Bad}}' 為開頭。 請使用一個沒有主語的句子。"
  },
  {
    "id": "Elapsed:",
    "translation": "經歷時間："
  },
  {
    "id": "External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "外部鑑別失敗。 錯誤代碼： {{.ErrorCode}}, 訊息： {{.Message}}"
  },
  {
    "id": "FAILED",
    "translation": "失敗"
  },
  {
    "id": "Failed, header could not convert to csv format",
    "translation": "失敗，標頭無法轉換為 csv 格式"
  },
  {
    "id": "Failed, rows could not convert to csv format",
    "translation": "失敗，資料無法轉換成 csv 格式"
  },
  {
    "id": "ID",
    "translation": "ID"
  },
  {
    "id": "Invalid grant type: ",
    "translation": "無效的授予類型： "
  },
  {
    "id": "Invalid token: ",
    "translation": "無效的記號："
  },
  {
    "id": "MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}",
    "translation": "MinCliVersion ( {{.ProvidedMinVersion}} ) 低於允許的最小值。 {{.AllowedMinimum}}"
  },
  {
    "id": "OK",
    "translation": "確定"
  },
  {
    "id": "Please enter 'y', 'n', 'yes' or 'no'.",
    "translation": "請輸入 'y'、'n'、'yes' 或 'no'。"
  },
  {
    "id": "Please enter a number between 1 to {{.Count}}.",
    "translation": "請輸入 1 到 {{.Count}} 之間的數字。"
  },
  {
    "id": "Please enter a valid floating number.",
    "translation": "請輸入有效的浮點數。"
  },
  {
    "id": "Please enter a valid number.",
    "translation": "請輸入有效的數字。"
  },
  {
    "id": "Please enter value.",
    "translation": "請輸入值。"
  },
  {
    "id": "REQUEST:",
    "translation": "要求："
  },
  {
    "id": "RESPONSE:",
    "translation": "回應："
  },
  {
    "id": "Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
    "translation": "遠端伺服器錯誤。狀態碼：{{.StatusCode}}，錯誤碼：{{.ErrorCode}}，訊息：{{.Message}}"
  },
  {
    "id": "Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
    "translation": "移除 `+"`"+` '{{.Bad}}' `+"`"+`，並直接從該動作開始。 範例：使用「列出所有執行個體」而非「此指令會列出所有執行個體」。"
  },
  {
    "id": "Remove placeholders from command usage text",
    "translation": "從命令使用說明中移除佔位符"
  },
  {
    "id": "Session inactive: ",
    "translation": "會議非主動： "
  },
  {
    "id": "Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.",
    "translation": "請將 `+"`"+` MinCliVersion `+"`"+` 設定為 `+"`"+` {{.Version}} `+"`"+` 或更高版本，以確保與受支援的 CLI 版本相容。"
  },
  {
    "id": "Shorten the description to be more concise.",
    "translation": "請將描述簡化，使其更為簡潔。"
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "無法儲存外掛程式配置："
  },
  {
    "id": "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
    "translation": "請在描述中加入「清單」、「顯示」、「檢視」、「全部」或「多個」等字詞，以明確指出該功能會返回多個項目。"
  },
  {
    "id": "Usage contains placeholder arguments/flags",
    "translation": "用法包含占位符參數/旗標"
  },
  {
    "id": "Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}",
    "translation": "使用方式包含指標之間未封閉的 {{.UnclosedGroup}} {{.Indicies}}"
  },
  {
    "id": "Usage does not have any usage text besides COMMAND",
    "translation": "除了 COMMAND 之外，此指令沒有其他說明文字"
  },
  {
    "id": "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
    "translation": "請在指令名稱中使用常見動詞，例如 list、create、update、delete，或使用複數形式來表示清單操作。"
  },
  {
    "id": "Use more descriptive words with at least {{.Count}} characters for each segment.",
    "translation": "請使用更具描述性的詞彙，並確保每個段落至少有 {{.Count}} 個字元。"
  },
  {
    "id": "{{.Field}} contains the following forbidden characters: {{.Chars}}",
    "translation": "{{.Field}} 包含下列禁止使用的字元： {{.Chars}}"
  },
  {
    "id": "{{.Field}} is required",
    "translation": "{{.Field}} 需要"
  },
  {
    "id": "{{.Field}} must contain at least {{.Param}} element",
    "translation": "{{.Field}} 必須至少包含 {{.Param}} 元素"
  }
]`)

func i18nResourcesAllZh_hantJsonBytes() ([]byte, error) {
	return _i18nResourcesAllZh_hantJson, nil
}

func i18nResourcesAllZh_hantJson() (*asset, error) {
	bytes, err := i18nResourcesAllZh_hantJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n/resources/all.zh_Hant.json", size: 9114, mode: os.FileMode(420), modTime: time.Unix(1780525928, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"i18n/resources/all.de_DE.json":   i18nResourcesAllDe_deJson,
	"i18n/resources/all.en_US.json":   i18nResourcesAllEn_usJson,
	"i18n/resources/all.es_ES.json":   i18nResourcesAllEs_esJson,
	"i18n/resources/all.fr_FR.json":   i18nResourcesAllFr_frJson,
	"i18n/resources/all.it_IT.json":   i18nResourcesAllIt_itJson,
	"i18n/resources/all.ja_JP.json":   i18nResourcesAllJa_jpJson,
	"i18n/resources/all.ko_KR.json":   i18nResourcesAllKo_krJson,
	"i18n/resources/all.pt_BR.json":   i18nResourcesAllPt_brJson,
	"i18n/resources/all.zh_Hans.json": i18nResourcesAllZh_hansJson,
	"i18n/resources/all.zh_Hant.json": i18nResourcesAllZh_hantJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"i18n": &bintree{nil, map[string]*bintree{
		"resources": &bintree{nil, map[string]*bintree{
			"all.de_DE.json":   &bintree{i18nResourcesAllDe_deJson, map[string]*bintree{}},
			"all.en_US.json":   &bintree{i18nResourcesAllEn_usJson, map[string]*bintree{}},
			"all.es_ES.json":   &bintree{i18nResourcesAllEs_esJson, map[string]*bintree{}},
			"all.fr_FR.json":   &bintree{i18nResourcesAllFr_frJson, map[string]*bintree{}},
			"all.it_IT.json":   &bintree{i18nResourcesAllIt_itJson, map[string]*bintree{}},
			"all.ja_JP.json":   &bintree{i18nResourcesAllJa_jpJson, map[string]*bintree{}},
			"all.ko_KR.json":   &bintree{i18nResourcesAllKo_krJson, map[string]*bintree{}},
			"all.pt_BR.json":   &bintree{i18nResourcesAllPt_brJson, map[string]*bintree{}},
			"all.zh_Hans.json": &bintree{i18nResourcesAllZh_hansJson, map[string]*bintree{}},
			"all.zh_Hant.json": &bintree{i18nResourcesAllZh_hantJson, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
