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
    "id": "Could not read from input: ",
    "translation": "Lesen der Eingabedaten nicht möglich: "
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
    "id": "Session inactive: ",
    "translation": "Sitzung inaktiv: "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Speichern der Plug-in-Konfiguration nicht möglich: "
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} darf nicht gleich ' {{.Param}} ' sein"
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

	info := bindataFileInfo{name: "i18n/resources/all.de_DE.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllEn_usJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nEnter a number"
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
    "id": "Could not read from input: ",
    "translation": "Could not read from input: "
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
    "id": "Session inactive: ",
    "translation": "Session inactive: "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Unable to save plugin config: "
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} must not equal '{{.Param}}'"
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

	info := bindataFileInfo{name: "i18n/resources/all.en_US.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllEs_esJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nEscriba un número"
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
    "id": "Could not read from input: ",
    "translation": "No se ha podido leer la entrada: "
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
    "id": "Session inactive: ",
    "translation": "Sesión inactiva: "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "No se ha podido guardar la configuración del plugin:"
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} no debe ser igual a ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.es_ES.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllFr_frJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nEntrez un nombre"
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
    "id": "Could not read from input: ",
    "translation": "Lecture impossible à partir de l'entrée : "
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
    "id": "Session inactive: ",
    "translation": "Session inactive : "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Impossible d'enregistrer la configuration du plug-in : "
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} ne doit pas être égal à ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.fr_FR.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllIt_itJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nImmetti un numero"
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
    "id": "Could not read from input: ",
    "translation": "Impossibile leggere dall'input: "
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
    "id": "Session inactive: ",
    "translation": "Sessione inattiva: "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Impossibile salvare la configurazione del plug-in: "
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} non deve essere uguale a ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.it_IT.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllJa_jpJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n数値を入力してください"
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
    "id": "Could not read from input: ",
    "translation": "入力から読み取れませんでした。 "
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
    "id": "Session inactive: ",
    "translation": "セッションは不活発： "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "プラグイン構成を保存できません: "
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} は ' {{.Param}} ' と同じであってはならない"
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

	info := bindataFileInfo{name: "i18n/resources/all.ja_JP.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllKo_krJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n번호 입력"
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
    "id": "Could not read from input: ",
    "translation": "입력에서 읽지 못함: "
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
    "id": "Session inactive: ",
    "translation": "세션이 비활성 상태입니다: "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "플러그인 구성을 저장할 수 없음:"
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.ko_KR.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllPt_brJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\nInsira um número"
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
    "id": "Could not read from input: ",
    "translation": "Não foi possível ler a apartir da entrada: "
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
    "id": "Session inactive: ",
    "translation": "Sessão inativa: "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "Não é possível salvar a configuração do plug-in: "
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} não deve ser igual a ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.pt_BR.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllZh_hansJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n请输入数字"
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
    "id": "Could not read from input: ",
    "translation": "无法从输入进行读取： "
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
    "id": "Session inactive: ",
    "translation": "会议非活动： "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "无法保存插件配置："
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} 不得等于 ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.zh_Hans.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _i18nResourcesAllZh_hantJson = []byte(`[
  {
    "id": "\nEnter a number",
    "translation": "\n請輸入數字"
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
    "id": "Could not read from input: ",
    "translation": "無法從輸入讀取： "
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
    "id": "Session inactive: ",
    "translation": "會議非主動： "
  },
  {
    "id": "Unable to save plugin config: ",
    "translation": "無法儲存外掛程式配置："
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
  },
  {
    "id": "{{.Field}} must not equal '{{.Param}}'",
    "translation": "{{.Field}} 不得等於 ' {{.Param}} '"
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

	info := bindataFileInfo{name: "i18n/resources/all.zh_Hant.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
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
