# IBM Cloud CLI Interface Summary

This summarizes the CLI interface guidelines as found in the [Plug-in Developer Guide](plugin_developer_guide.md).

## Name of the plug-in

Name of the plug-in should be a name that best describes the service the plug-in provides. Use the full name of the service in all lower case with a hyphen to replace space, for example 'my-service' ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#22-name-and-decription-of-plug-in-namespace-and-command)). Shorter aliases should be used to improve usability ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#11-register-a-new-plug-in)). Use commonly accepted shortened names of services as possible.

## Names of commands ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#22-name-and-decription-of-plug-in-namespace-and-command))

- For commands that list objects, use the plural of the object name, e.g. `ibmcloud iam api-keys`. If a plural will not work for some reason, use `list`, such as `ibmcloud app list`.
- For commands that retrieve the details of an object, use the object name, e.g. `ibmcloud iam api-key`. If it does not work, use `show`, e.g. `ibmcloud app show`
- Use common verbs in the names such as add, create, bind, update, delete …
- Use `remove` to remove an object from a collection or another object it is associated with. Use `delete` to erase an object.

You have 2 options for the overall naming scheme of your commands:

Option 1: hyphenated commands ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#22-name-and-decription-of-plug-in-namespace-and-command))
- Use lower case words and hyphen
- Names should be at least 2 characters in length
- Follow a “noun-verb” sequence (verb is the final part of the command)
- Example: "object-create"

Option 2: using namespaced commands
- follows the same rules as Option 1, only the delimiter is not hyphen, but space
- Example: "object create"

## Command and namespace descriptions ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#22-name-and-decription-of-plug-in-namespace-and-command))

Use a sentence without subject to describe your plug-in or command. Limit the number of words to be less than 10 so that it can be properly displayed.

Correct description:

- `List all the virtual server instances.`
- `Manage cloud database service.`

Incorrect description:
- `Plugin to manage cloud database service.`
- `Commands to manage cloud database service.`
- `This command shows details of a sever instance.`


## Messaging standards ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#21-global-specification))

1. Do NOT use "Please" in any message.
   Correct:
   <pre>First log in by running '<span class="yellow">ibmcloud login</span>'.</pre>
   Invalid:
   <pre>Please use '<span class="yellow">ibmcloud login</span>' to login first.</pre>

2. Capitalize the first letter for all sentences and short descriptions. For example:
  <pre>Change the instance count for an app or container group.</pre>
  <pre>-i   Number of instances</pre>

3. Add "..." at the end of "in-progress" messages. For example:
  <pre>Scaling container group '<span class="cyan">abc</span>'...</pre>

4. Use "plug-in" instead of "plugin" in all places.

### Command formatting in output messages ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#22-name-and-decription-of-plug-in-namespace-and-command))

Use single quotation marks (') around the command name and options in output messages. The command itself should be yellow with **bold**.  Where possible, place command names at the end of the sentence, not in the middle. For example:

```
You are not logged in. Log in by running 'ibmcloud login'.
```

### Entity formatting in output messages: ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#23-entity-name))

Add single quotation marks (') around entity names and keep the entity names in cyan with **bold**. For example:

```
Mapping route 'my-app.us-south.cf.cloud.ibm.com' to CF application 'my-app'...
```


### Command Help ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#24-help-of-command))

- Use "-" for single letter flags, and "--" for multiple letter flags, e.g. `-c ACCOUNT_ID` and `--guid`.
- All user input values should be shown as capital letters, e.g. `ibmcloud scale RESOURCE_NAME`
- For optional parameters and flags, surround them with "[...]", e.g. `ibmcloud account orgs [--guid]`.
- For exclusive parameters and flags, group them together by "(...)" and separate by "|".
  - Example: `ibmcloud test create (NAME | --uuid ID)`
- "[...]" and "(...)" can be nested.
- If a value accepts multiple type of inputs, it's recommended that for file type the file name should start with "@".
- If a command has multiple paradigms and it's hard to describe them together, specify each of them in separate lines,
  e.g.
  ```bash
  USAGE:
    ibmcloud test command foo.....
    ibmcloud test command bar.....
  ```

The following gives an example of the output of the `help` command:

```
NAME:
   scale - Change the instance count for an app or container group.
USAGE:
   ibmcloud scale RESOURCE_NAME [-i INSTANCES] [-k DISK] [-m MEMORY] [-f]
   RESOURCE_NAME is the name of the app or container group to be scaled.
OPTIONS:
   -i value  Number of instances.
   -k value  Disk limit (e.g. 256M, 1024M, 1G). Valid only for scaling an app, not a container group.
   -m value  Memory limit (e.g. 256M, 1024M, 1G). Valid only for scaling an app, not a container group.
   -f        Force restart of CF application without prompt. Valid only for scaling an app, not a container group.
```

### Common options ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#212-common-options))

Users use multiple services. Consistency with option names simplifies this experience.

```
OPTIONS:
   --force, -f                  Force the operation without confirmation
   --instance-id                ID of the service instance.
   --output json                Format output in JSON
   --resource-group-id value    ID of the resource group. This option is mutually exclusive with --resource-group-name
   --resource-group-name value  Name of the resource group. This option is mutually exclusive with --resource-group-id
```

### Command Output
#### Incorrect Usage ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#25-incorrect-usage))

When the user runs a command with the wrong usage (e.g. incorrect number of arguments, invalid option value, required options not specified, etc.), the message should be displayed and include help for the user for the command as below:

```
Incorrect usage.
NAME:
   scale - Change the instance count for an app or container group.
USAGE:
   ibmcloud scale RESOURCE_NAME [-i INSTANCES] [-k DISK] [-m MEMORY] [-f]
   RESOURCE_NAME is the name of the app or container group to be scaled.
OPTIONS:
   -i value  Number of instances.
   -k value  Disk limit (e.g. 256M, 1024M, 1G). Valid only for scaling an app, not a container group.
   -m value  Memory limit (e.g. 256M, 1024M, 1G). Valid only for scaling an app, not a container group.
   -f        Force restart of CF application without prompt. Valid only for scaling an app, not a container group.
```

Provide any details possible to guide and inform the users. For example:

```
Incorrect usage. The '-k' option is not valid for a container group.
```

#### Command Failure ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#26-command-failure))

The failure message must start with "FAILED" in red with **bold**, followed by the detailed error message in a new line. A recovery solution should be provided, such as "Use another name and try again." or "Try again later."

In the message, explain the error and provide guidance on how to resolve the issue such as shown in the following output:

```
Creating application 'my-app'...
FAILED
An application with name 'my-app' already exists.
Use another name and try again.
```

#### Command Success ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#27-command-success))

When a command is successful, the success message should start with "OK" in green with **bold** and followed by the optional details in new line like the following example:

```
Creating application 'my-app'...
OK
Application 'my-app' was created.
```

#### Warning Message ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#28-warning-message))

All of command warnings should be magenta with **bold** like:

```
WARNING:
   If you enter your password as a command option, your password might be visible to others or recorded in your shell history.
```

#### Important Information for the User ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#29-important-information))

The important information displayed to the end-user should be cyan with **bold**. For example:

```
A newer version of the IBM Cloud CLI is available.
You can download it from http://abc.abc.abc
```

### User Input Prompts ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#210-user-input-prompt))

Prompting for user input:

1. The input prompt should consist of a prompt message and a right angle bracket '<span style="color: cyan">**>**</span>' with a trailing space.
2. For a password prompt, the user input must be hidden.
3. For a confirmation prompt, the prompt message should end with `[Y/n]` or `[y/N]` (the capitalized letter indicates the default answer) or `[y/n]` (no default, user input is required).

Following are examples:

#### Text prompt

```
Logging in to https://cloud.ibm.com...
Email>abc@example.com
Password>
OK
```

#### Confirmation prompt

```
Are you sure you want to remove the file? [y/N]> y
Removing ...
OK
```

#### Choices prompt

```
Select the plug-in to be upgraded:
1. plugin1
2. plugin2
Enter a number (1)> 2

Upgrading 'plugin2'...
```

#### Prompt override ("force" option)
Users need a way to override prompts to allow the execution of non interactive scripts.

For the configrmation [y/N] prompt use the -f force option.


### Table Output ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#211-table-output))

For a consistent user experience, developers of IBM Cloud CLI plug-ins should comply with the following table output specifications:

1. Table header should be bold.
2. Each word in table header should be capitalized and translated.
3. The entity name or resource name in table (usually the first column) should be cyan with bold.

The following is an example:

```
Name       Description
my-app     This is my application.
demo-app   This is a long long long ...
           description.
```

### JSON Output ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#211-json-output))

Use flag `--output json` to show the json representation of resource(s) if the command is to list resources, or retrieve details of a resource. If this flag is used, don't show any informational messages or prompts but just the JSON string so that it can be easily parsed with other tools like jq. For example:

```
$ibmcloud account orgs --output json
[
    {
        "OrgGuid": "ef6e9345-2155-41bb-bd3d-ff7c10e5f071",
        "OrgName": "example-org",
        "Region": "us-south",
        "AccountOwner": "user@example.com",
        "AccountGuid": "8d63fb1cc5e99e86dd7229dddf9e5b7b",
        "Status": "active"
    }
]
```

When the output is an empty list the plugin should produce an empty json list (not null).  For example if there were no account orgs:
```
$ibmcloud account orgs --output json
[]
```

If the output is expected to be an object but no object is returned the plugin should return the normal error: For example, if a service-id could not be found then an error is returned:
```
$ibmcloud iam service-id 15a15a0f-725e-453a-b3ac-755280ad7300 --output json
FAILED
Service ID '15a15a0f-725e-453a-b3ac-755280ad7300' was not found.
```

### Globalization ([reference](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/docs/plugin_developer_guide.md#7-globalization))

IBM Cloud CLI tends to be used globally. Both IBM Cloud CLI and its plug-ins should support globalization. 

