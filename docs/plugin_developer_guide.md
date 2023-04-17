# IBM Cloud CLI Plug-in Developer's Guide

This guide introduces how to develop an IBM Cloud CLI plug-in by using utilities and libraries provided by the CLI SDK. It also covers specifications including wording, format and color of the terminal output that we highly recommend developers to follow.

You can see the API doc in [GoDoc](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk).

## 1. Plug-in Context Management

IBM Cloud CLI SDK provides a set of APIs to register and manage plug-ins. It also provides a set of utilities and libaries to simplify the plug-in development.

### 1.1. Register a New Plug-in

1.  Define a new struct for IBM Cloud CLI plug-in and associate `Run` and `GetMetadata` methods to that struct:

    ```go
    import (
        "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
    )

    type DemoPlugin struct {}

    func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string) {}

    func (demo *DemoPlugin) GetMetadata() plugin.PluginMetadata {}
    ```

2.  In main function, invoke `plugin.Start` method to register the plug-in:

    ```go
    func main() {
        plugin.Start(new(DemoPlugin))
    }
    ```

3.  Return a `plugin.PluginMetadata` struct to finalize the registration of the plug-in in `GetMetadata` method:

    ```go
    func (demo *DemoPlugin) GetMetadata() plugin.PluginMetadata {
        return plugin.PluginMetadata{
            Name: "demo-plugin",
            Version: plugin.VersionType{
                Major: 1,
                Minor: 0,
                Build: 0,
            },
            MinCliVersion: plugin.VersionType{
                Major: 0,
                Minor: 0,
                Build: 1,
            },

            PrivateEndpointSupported: true,

            IsCobraPlugin: true,

            Commands: []plugin.Command{
                {
                    Name:        "echo",
                    Alias:       "ec",
                    Description: "Echo a message on terminal.",
                    Usage:       "ibmcloud echo MESSAGE [-u]",
                    Flags: []plugin.Flag{
                        {
                            Name: "u",
                            Description: "Change the message to upper case.",
                            HasValue: false,
                        },
                    },
                },
            },
        }
    }
    ```

    **Understanding the fields in this `plugin.PluginMetadata` struct:**
    - _Name_: The name of plug-in. It will be displayed when using `ibmcloud plugin list` command or can be used to uninstall the plug-in through `ibmcloud plugin uninstall` command.
      - It is **strongly** encouraged to use a name that best describes the service the plug-in provides.
    - _Aliases_: A list of short names of the plug-in that can be used as a stand-in for installing, updating, uninstalling and using the plug-in.
      - It is strongly recommended that you have at least one alias to improve the usability of the plug-in.
    - _Version_: The version of plug-in.
    - _MinCliVersion_: The minimal version of IBM Cloud CLI required by the plug-in.
    - _PrivateEndpointSupported_: Indicates if the plug-in is designed to also be used over the private network.
    - _IsCobraPlugin_: Indicates if the plug-in is built using the Cobra framework.
      - It is **strongly** recommended that you use this framework to build your plug-in.
    - _Commands_: The array of `plugin.Commands` to register the plug-in commands.
    - _Alias_: Alias of the Alias usually is a short name of the command.
    - _Command.Flags_: The command flags (options) which will be displayed as a part of help output of the command.

4.  Add the logic of plug-in command process in `Run` method, for example:

    ```go
    func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string) {
        if args[0] == "echo" {
            // echo command logic here
        }
    }
    ```

`PluginContext` provides the most useful methods which allow you to get command line properties from CF configuration as well as IBM Cloud  specific properties.

### 1.2. Namespace

IBM Cloud  CLI introduced a new concept called "Namespace". A namespace is a category of commands which have similar functionality. Some namespaces are predefined by IBM Cloud  CLI and can be shared by plug-ins, but others are non-shared namespaces which can be defined in each plug-in. The plug-in can reference a predefined namespace in IBM Cloud  CLI or define a non-shared namespace by its own. You can also use sub-namespaces to organize commands into categories.

#### Shared namespaces
The following are shared namespaces are currently predefined in IBM Cloud  CLI and can be shared across plug-ins:

| Namespace | Description |
| --- | --- |
| account | Manage accounts, orgs, spaces and users |
| iam | Manage identities and accesses |
| catalog | Manage IBM Cloud catalog |
| app | Manage IBM Cloud applications |
| service | Manage IBM Cloud services |
| resource | Manage the IBM cloud resources |
| billing | Retrieve usage and billing information |
| cf | Run Cloud Foundry CLI with IBM Cloud context |
| plugin | Manage plug-in repositories and plug-ins |

For shared namespace, refer to the namespace in the plug-in:

```go
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

func (p *CatalogExtPlugin) GetMetadata() plugin.PluginMetadata {
    return plugin.PluginMetadata{
        ...
        Commands: []plugin.Command{
            {
                Namespace:   "catalog",
                Name:        "cmd1",
            },
            {
                Namespace:   "catalog",
                Name:        "cmd2",
            },
        },
    }
}

func (p *CatalogExtPlugin) Run(context plugin.PluginContext, args []string) {
    switch args[0] {
    case "cmd1":
        // cmd1 command logic here
    case "cmd2":
        // cmd2 command logic here
    default:
        // command not recognized
    }
}
```

#### Non-shared namespaces
Define a non-shared namespace for your plug-in:

```go
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"


func (p *DemoPlugin) GetMetadata() plugin.PluginMetadata {
    // define plugin namespace
    demo := plugin.Namespace{
        Name: "demo",
        Description: "Demonstrate non-shared namespace.",
    }

    return plugin.PluginMetadata{
        ...
        Namespaces: []plugin.Namespace {
            demo,
        },
        Commands: []plugin.Command{
            {
                Namespace:   demo.Name,
                Name:        "start",
            },
            {
                Namespace:   demo.Name,
                Name:        "help",
            },
            {
                Namespace:   demo.Name,
                Name:        "*",
            },
        },
    }
}

func (p *DemoPlugin) Run(context plugin.PluginContext, args []string) {
    switch args[0] {
    case "start":
        // start command logic here
    case "help"
        // print help
    default:
        // default run. Wildcard "*" matched
    }
}
```

#### Sub-namespaces
If you want to organize commands into categories at different levels, you can use sub-namespace. White spaces in namespace will be treated as delimiter for sub-namespaces. For example, considering namespace "a b c", its parent namespace is "a b", and its ancestor namespace is "a".

Following is an example of using sub-namespace:

```go
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

type DemoPlugin struct{}

func (p *DemoPlugin) GetMetadata() plugin.PluginMetadata {
    return plugin.PluginMetadata{
        Namespaces: []plugin.Namespace{
            // parent namespace
            {
                Name:        "demo",
                Description: "Show parent namespace.",
            },
            // sub namespaces
            {
                Name:        "demo app",
                Description: "Show app sub namespace.",
            },
            {
                Name:        "demo service",
                Description: "Show service sub namespace.",
            },
        },

        Commands: []plugin.Command{
            {
                Namespace:   "demo app",
                Name:        "create",
                Description: "Create an application.",
            },
            {
                Namespace:   "demo app",
                Name:        "delete",
                Description: "Delete an application.",
            },
            {
                Namespace:   "demo service",
                Name:        "create",
                Description: "Create an service instance.",
            },
            {
                Namespace:   "demo service",
                Name:        "delete",
                Description: "Delete an service instance.",
            },
        },
    }
}

func (p *DemoPlugin) Run(context plugin.PluginContext, args []string) {
    namespace := context.CommandNamespace()
    switch namespace {
    case "demo app":
        switch args[0] {
        case "create":
            // create an application
        case "delete":
            // delete an application
        default:
            // unrecognized command
        }
    case "demo service":
        switch args[0] {
        case "create":
            // create service instance
        case "delete":
            // delete service instance
        default:
            // unrecognized command
        }
    }
}
```

#### Notes on namespaces

The following items should be noticed here:

- If a command is not associated with any namespace, it will be registered as a root command. For example, if a command is associated with the app namespace, it must be executed by typing "ibmcloud app xxx", but if no namespace is associated with a command, the command will be invoked by "ibmcloud xxx" directly.
- All shared namespaces can only be defined in IBM Cloud CLI and referenced by plug-ins.
- Developer can define multiple non-shared namespaces in one plug-in.
- If the plug-in only belongs to non-shared namespaces, a special command with name "\*" can be defined, which means even if a user typed in a command which is not defined in current namespace, the plug-in will still be invoked, otherwise, an error message will be displayed by IBM Cloud CLI. Taking the above code snippet as an example, if user typed in "ibmcloud demo others", then "default" logic in "Run" method will be hit.
- If multiple plug-ins define a namespace with the same name, they will follow a "first install first serve" strategy, which means the latter plug-ins can't be installed successfully due to the namespace conflict.
- Developer can define help command for non-shared namespace. If it is not defined, the help content will be generated by IBM Cloud CLI automatically based on registered commands. For shared namespaces the help command was predefined in IBM Cloud CLI.
- **Important:** No matter whether the command belongs to a namespace, args[0] will always be the command name or alias instead of the namespace name. You can get the namespace via `PluginContext.CommandNamespace()`.

### 1.3. Manage Plug-in Configuration

IBM Cloud CLI SDK provides APIs to allow you to access the plug-in's own configuration saved in JSON format. Each plug-in will have its own configuration file be created automatically on installation. Take the following code as an example:

```go
func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string){
    config := context.PluginConfig()

    // set
    err := config.Set("s", "string value")
    panic(err)

    err = config.Set("n", 123)
    panic(err)

    err = config.Set("ss", []string("foo", "bar"))
    panic(err)

    err = config.Set("m", nap[string]string{"one": "foo", "two": "bar"}})
    panic(err)

    // get
    myStr, err :=  config.GetString("s")
    panic(err)

    myNum, err := config.GetIntWithDefault("n", 100)
    panic(err)

    mySlice, err := config.GetStringSlice("ss")
    panic(err)

    myMap, err := config.GetStringMapString("m")
    panic(err)
    ...
}
```

## 2. Wording, Format and Color of Output

To keep user experience consistent, developers of IBM Cloud CLI plug-in should apply specific wordings, formats and colors to the terminal output. IBM Cloud CLI SDK provides the utility to help plug-in developers easily format and colorize the message output. We strongly recommend developers to comply with the following specifications so that the plug-ins are consistent with each other in terms of user experience.


### 2.1. Global Specification

1. Do NOT use "Please" in any message.
   Correct:
   <pre>First log in by running '<span class="yellow">ibmcloud login</span>'.</pre>
   Invalid:
   <pre>Please use '<span class="yellow">ibmcloud login</span>' to login first.</pre>

2. Capitalize the first letter for all sentences and short descriptions. For example:
  <pre>Change the instance count for an app or container group.</pre>
  <pre>-i   Number of instances</pre>

3. Add "..." at the end of "in-progress" messages. For example:
  <pre>Scaling container group '<span class="cyan">xxx</span>'...</pre>

4. Use "plug-in" instead of "plugin" in all places.

### 2.2. Name and Decription of Plug-in, Namespace and Command

**To name the plug-in for a service**:

- Use full name of the service in all lower case, and hyphen to replace space, for example 'my-service'.
- Use abbreviation only when it’s commonly accepted, for example use 'vpn' for Virtual Private Network service.

**To name the commands**:

- Use lower case words and hyphen
- Names should be at least 2 characters in length
- Follow a “noun-verb” sequence
- For commands that list objects, use the plural of the object name, e.g. `ibmcloud iam api-keys`. If a plural will not work, or it is already described in namespace, use `list`, such as `ibmcloud app list`.
- For commands that retrieve the details of an object, use the object name, e.g. `ibmcloud iam api-key`. If it does not work, use `show`, e.g. `ibmcloud app show`
- Use common verbs such as add, create, bind, update, delete …
- Use `remove` to remove an object from a collection or another object it associated with. Use `delete` to erase an object.

The following command names are formatted correctly, and provides a possible description for the command:

- `route-map`: Map out the route from one place to another in a location application.

- `route-unmap`: Unmap a previously mapped route.

- `template-run`: Run a selected template.

- `plugin-repo-add`: Add a plug-in repository.

- `plugin-repo-remove`: Remove a plug-in repository.

- `templates`: Get a list of templates.

The following command names are invalid:

- `Map-Route`: Uses uppercase letters (M, R). Does not follow "noun-verb" word order (the route needs to be mapped, but as this is written, the order implies that the map needs to be routed).

- `map-route`: Does not follow "noun-verb" word order (the route needs to be mapped, but as this is written, the order implies that the map needs to be routed).

- `map route`: Uses a space instead of a hyphen. Does not follow "noun-verb" word order (the route needs to be mapped, but as this is written, the order implies that the map needs to be routed).

- `route map`: Uses a space instead of a hyphen.

- `route\_map`: Uses unallowed characters (`\_`) instead of a hyphen.

- `add-plugin-repo`: Does not follow "noun-verb" word order. The verb "add" should be the last part of the command name.

- `plugin-add-repo`: Does not follow the "noun-verb" word order. The verb "add" should be the last part of the command name.


**Command formatting in output messages**:

Use single quotation marks (') around the command name and options in output message. The command itself should be yellow with **bold**.  Where possible, place command names at the end of the sentence, not in the middle. For example:

```
You are not logged in. Log in by running 'ibmcloud login'.
```

You can use the following APIs provided by IBM Cloud CLI SDK to print the previous example message:

```go
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

ui := terminal.NewStdUI()

ui.Say(`You are not logged in. Log in by running "%s".`,
    terminal.CommandColor("ibmcloud login"))

```
**Command and namespace description**

Use a sentence without subject to describe your plug-in or command. Limit the number of words to be less than 10 so that it can be properly displayed.

Correct description:

- `List all the virtual server instances.`
- `Manage cloud database service.`

Incorrect description:
- `Plugin to manage cloud database service.`
- `Commands to manage cloud database service.`
- `This command shows details of a sever instance.`


### 2.3. Entity Name

Add single quotation marks (') around entity names and keep the entity names in cyan with **bold**. For example:

```
Mapping route 'my-app.us-south.cf.cloud.ibm.com' to CF application 'my-app'...
```

The IBM Cloud CLI SDK also provides API to help you print the previous example message:

```go
ui.Say("Mapping route '%s' to CF application '%s'...",
    terminal.EntityNameColor("my-app.cloud.ibm.com"),
    terminal.EntityNameColor("my-app"),
)
```

### 2.4. Help of Command

Use the guidelines below to compose command help.
- Use "-" for single letter flags, and "--" for multiple letter flags, e.g. `-c ACCOUNT_ID` and `--guid`.
- All user input values should be capital letters, e.g. `ibmcloud scale RESOURCE_NAME`
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


### 2.5. Incorrect Usage

When users run a command with wrong usage (e.g. incorrect number of arguments, invalid option value, required options not specified and etc.), the message should be displayed in the following format:

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

If possible, provide details to help the users figure out what was wrong with their usage. For example:

```
Incorrect usage. The '-k' option is not valid for a container group.
```

### 2.6. Command Failure

If a command failed due to the client-side or server-side error, explain the error and provide guidance on how to resolve the issue, such as shown in the following output:

```
Creating application 'my-app'...
FAILED
An application with name 'my-app' already exists.
Use another name and try again.
```

```
Scaling container group 'xxx'...
FAILED
A server error occurred while scaling the container group.
Try again later. If the problem continues, contact IBM Cloud Support.
```

To summarize, the failure message must start with "FAILED" in red with **bold** and followed by the detailed error message in a new line as previously shown. A recovery solution must be provided, such as "Use another name and try again." or "Try again later."

IBM Cloud CLI also provides Failed method to print out failure message:

```go
func Run() error {
    ui.Say("Scaling container group '%s'...", terminal.EntityNameColor("xxx"))
    ...
    if err != nil {
        ui.Failed("A server error occurred while scaling the container group.\nTry again later. If the problem continues, contact IBM Cloud Support.")
        return err
    }
    ...
}
```

### 2.7. Command Success

When command was successful, the success message should start with "OK" in green with **bold** and followed by the optional details in new line like the following examples:

```
Creating application 'my-app'...
OK
Application 'my-app' was created.
```

The following code snippet can help you print the above message:

```go
ui.Say("Creating application '%s'...",terminal.EntityNameColor("my-app"))
...
ui.Ok()
ui.Say("Application '%s' was created.", terminal.EntityNameColor("my-app"))
```

### 2.8. Warning Message

All of command warnings should be magenta with **bold** like:

```
WARNING:
   If you enter your password as a command option, your password might be visible to others or recorded in your shell history.
```

Take the following code as an example for the warning message output:

```go
ui.Warn("WARNING:...")
```

### 2.9. Important Information

The important information displayed to the end-user should be cyan with **bold**. For example:

```
A newer version of the IBM Cloud CLI is available.
You can download it from http://xxx.xxx.xxx
```

The corresponding code snippet:

```go
ui.Say(terminal.PromptColor("A newer version of the IBM Cloud CLI ..."))
```

### 2.10. User Input Prompt

Following specifications should be followed to prompt for user input:

1. The input prompt should consist of a prompt message and a right angle bracket '<span style="color: cyan">**>**</span>' with a trailing space.
2. For password prompt, the user input must be hidden.
3. For confirmation prompt, the prompt message should end with `[Y/n]` or `[y/N]` (the capitalized letter indicates the default answer) or `[y/n]` (no default, user input is required).

Following are examples and code snippets:

#### Text prompt

```
Logging in to https://cloud.ibm.com...
Email>xxx@example.com
Password>
OK
```

Code:

```go
ui.Say("Logging in to https://cloud.ibm.com")

var email string
err := ui.Prompt("Email", nil).Resolve(&email)
if err != nil {
    panic(err)
}

var passwd string
err =  ui.Prompt("Password", &terminal.PromptOptions{HideInput: true}).Resolve(&passwd)
if err ! = nil {
    panic(err)
}
...
ui.OK()
```

#### Confirmation prompt

```
Are you sure you want to remove the file? [y/N]> y
Removing ...
OK
```

Code:

```go
confirmed := false
err := ui.Prompt("Are you sure you want to remove the file?", nil).Resolve(&confirmed)
if err != nil {
    panic(err)
}

if confirmed {
    ui.Say("Removing...")
    ...
    ui.OK()
}
```

#### Choices prompt

```
Select the plug-in to be upgraded:
1. plugin1
2. plugin2
Enter a number (1)> 2

Upgrading 'plugin2'...
```

Code:

```go
plugins = []string{"plugin1", "plugin2"}
selected := plugins[0] // default is the first plugin

err := ui.ChoicesPrompt("Select the plug-in to be upgraded:", plugins, nil).Resolve(&selected)
if err != nil {
    panic(err)
}

ui.Say("upgrading '%s'...", terminal.EntityNameColor(selected))
```

#### Prompt override
There must be a way to override the prompt from a command line switch to allow the execution of non interactive scripts.

For the configrmation [y/N] prompt use the -f force option.


### 2.11. Table Output

For consistent user experience, developers of IBM Cloud CLI plug-ins should comply with the following table output specifications:

1. Table header should be bold.
2. Each word in table header should be capitalized.
3. The entity name or resource name in table (usually the first column) should be cyan with bold.

The following is an example:

```
Name       Description
my-app     This is my application.
demo-app   This is a long long long ...
           description.
```

Using APIs provided by IBM Cloud CLI can let you easily print results in table format:
```go
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

func (demo *DemoPlugin) PrintTable() {
    ui := terminal.NewStdUI()

    table := ui.Table([]string{"Name", "Description"})
    table.Add("my-app", "This is my application.")
    table.Add("demo-app", "This is a long long long ...\ndescription.")
    table.Print()
}
```
### 2.11. Json Output

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

When the output is an empty list the plugin should produce an empty json list (not null).  For example if there were not account orgs:
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

### 2.12. Common options

Customers will be writing scripts that use multiple services.  Consistency with option names will help them be successful.

```
OPTIONS:
   --force, -f                  Force the operation without confirmation
   --instance-id                ID of the service instance.
   --output json                Format output in JSON
   --resource-group-id value    ID of the resource group. This option is mutually exclusive with --resource-group-name
   --resource-group-name value  Name of the resource group. This option is mutually exclusive with --resource-group-id
```

## 3. Tracing

IBM Cloud CLI provides utility for tracing based on "IBMCLOUD\_TRACE" environment variable. The trace will be disabled if environment variable "IBMCLOUD\_TRACE" was not set or it was set to "false" (case ignored), which means, in that case, the invocation of trace API has no effect. If "IBMCLOUD\_TRACE" was set to "true" (case ignored), the trace will be printed on the terminal. Otherwise, the value of "IBMCLOUD\_TRACE" will be treated as the path of trace file.

An example to use IBM Cloud CLI trace API:

```go
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"

func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string) {
    // first, initialize the trace logger
    trace.Logger = trace.NewLogger(context.Trace())

    // start using the trace logger
    trace.Logger.Println("Start to initialize Demo plug-in.")
    ...
    trace.Logger.Printf("%s plug-in initialized.", "Demo")
}
```

## 4. HTTP Utilities

### 4.1. HTTP tracing

TraceLoggingTransport in package `ibm-cloud-cli-sdk/bluemix/http` is a thin wrapper around HTTP transport. It dumps each HTTP request and its response by using the trace logger.

1. Initialize the trace logger.

   ```go
   trace.Logger = trace.NewLogger(pluginContext.Trace())
   ```

2. Create a HTTP client with TraceLoggingTransport to send the request:
   ```go
   import (
       "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/http"
       gohttp "net/http"
   )

   client := &gohttp.Client{
       Transport: http.NewTraceLoggingTransport(gohttp.DefaultTransport)
   }
   client.Get("http://www.example.com")
   ```

Now during each round-trip, the trace logger dumps the request and its response.

### 4.2. REST client

HTTP interaction with a remote server is a common task in both core and plug-in commands. Package `ibm-cloud-cli-sdk/common/rest` provides APIs for building a REST API request and a REST client for sending the request.

Examples for GET requests, query strings, HTTP headers, POST requests, and file uploads.

* To create a GET request:
  ```go
  import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"

  var r = rest.GetRequest("http://www.example.com")
  // equal to
  // var r = NewRequest(“http://www.example.com”).Method(“GET”)
  ```

* To add a query string in the URL:
  ```go
  r.Query("foo", "bar")
  ```

* To add or set HTTP headers in the request:
  ```go
  r.Add("Accept-Encoding","gzip")
  r.Set("Accept", "application/json")
  ```

* To create a POST request, and send form data:
  ```go
  var r = rest.PostRequest("http://www.example.com")
  r.Field("foo", "bar")
  ```

* To upload a file:
  ```go
  f, err := os.Open("1.img")
  if err != nil {
     // handle error
  }
  r.File("file1", rest.File{
      Name: f.Name(),
      Content: f,
      Type: "image/jpeg", // Optional. Default is "application/octet-stream"
  })
  ```

You can send form data and upload multiple files in a same request.

To post a JSON data in the request, you can simply pass a Go struct to the Body () method. The method automatically encodes the Go struct to a JSON string.

For example:
```go
type Foo struct
    Name string
}
var r = rest.PostRequest("http://www.example.com")
r.Body(Foo{Name: "bar"})
```

The Body() method also supports raw string and steam. The previous example can also be written as:
```go
r.Body("{\"name\": \"bar\"}")
// Or
r.Body(strings.NewReader("{\"name\": \"foo1\"}"))
```

After the request is created, you can then create a REST client to send the request. The REST client is safe for concurrent use by multiple Go routines and thus is recommend to be reused.

To create a REST client:
```go
client := rest.NewClient()
```

By default, the client use Go’s standard HTTP client. You can override it:
```go
client.HTTPClient = &http.Client{
    Timeout: 60 * time.Second,
}
```

Also, you can set default HTTP header for all outgoing requests:
```go
h := http.Header{}
h.Set("User-Agent", "IBM Cloud CLI")
Client.DefaultHeader = h
```

Now, you can invoke client’s Do() method to send the request. The method automatically unmarshals the response body to the Go struct. If server response’s status code is 2xx, successV is unmarshaled; otherwise, errorV is un- marshaled if exists. If errorV is not provided or not successfully unmarshaled, an ErrorResponse typed error is returned which has status code and response text.
```go
var successV Foo
var errorV = struct {
    Message string
}{}

resp, err := client.Do(r, &successV, &errorV)
if errorV != nil || err != nil {
    // handle error
}
```

## 5. Authentication

### 5.1 Get Access Token

To access IBM Cloud back-end API, normally a token is required. You can get the IAM access token and UAA access token from  IBM CLoud SDK as follows:

```go
func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string){
    config := context.PluginConfig()

    // get IAM access token
    iamToken := config.IAMToken()
    if iamToken == "" {
        ui.Say("IAM token is not available. Have you logged in?")
        return
    }

    // get UAA access token
    uaaToken := config.CFConfig().UAAToken()
    if iamToken == "" {
        ui.Say("UAA token is not available. Have you logged into Cloud Foundry?")
        return
    }
    ...
}
```

And you can set the `Authorization` header of the HTTP request to the token value.

```go
h := http.Header{}
h.Set("Authorization", token)
```

For more details of the API, refer to docs of [core_config](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config).

If you want to fetch the token by yourselve, refer to API docs for [iam](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam) and [uaa](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/uaa).

### 5.2 Refresh Access Token on Expiry

When an HTTP 401 or 403 is returned from back-end service, it could mean the access token is expired. You can try to refresh existing access token following the [Oauth2 spec regarding how to refresh access token](https://www.oauth.com/oauth2-servers/access-tokens/refreshing-access-tokens/).

Let's take refresh IAM token as an example.

```go
// get existing refresh token
refreshToken := config.IAMRefreshToken()

// prepare token refresh request. See https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam#RefreshTokenRequest for API details
request := iam.RefreshTokenRequest(refreshToken)

// send request to IAM endpoint to generate new tokens
c := &rest.Client{
	HTTPClient:    http.NewHTTPClient(config),
	DefaultHeader: http.DefaultHeader(config),
}
client := iam.NewClient(iam.DefaultConfig(config.IAMEndpoint()), c)

token, err := client.GetToken(request)

// get the new access token and refresh token
accessToken := token.AccessToken
newRefreshToken := token.RefreshToken

// optional, set access token and refresh token back to config
config.SetAccessToken(accessToken)
config.SetRefreshToken(newRefreshToken)

// optional, maintain session for long running workloads
request = iam.RefreshSessionRequest(token)
client.RefreshSession(token)
```

### 5.3 VPC Compute Resource Identity Authentication

#### 5.3.1 Get the IAM Access Token
The IBM Cloud CLI supports logging in as a VPC compute resource identity. The CLI will fetch a VPC instance identity token and exchange it for an IAM access token when logging in as a VPC compute resource identity. This access token is stored in configuration once a user successfully logs into the CLI.

Plug-ins can invoke `plugin.PluginContext.IsLoggedInAsCRI()` and `plugin.PluginContext.CRIType()` in the CLI SDK to detect whether the user has logged in as a VPC compute resource identity.
You can get the IAM access token resulting from the user logging in as a VPC compute resource identity from the IBM CLoud SDK as follows:

```go
func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string){
    // confirm user is logged in as a VPC compute resource identity
    isLoggedInAsCRI := context.IsLoggedInAsCRI()
    criType := context.CRIType()
    if isLoggedInAsCRI && criType == "VPC" {
        // get IAM access token
        iamToken := context.IAMToken()
        if iamToken == "" {
            ui.Say("IAM token is not available. Have you logged in?")
            return
        }
    }
    ...
}
```

This token can be used to access IBM Cloud back-end APIs. You can set the `Authorization` header of the HTTP request to the token value.

```go
h := http.Header{}
h.Set("Authorization", iamToken)
```

#### 5.3.2 Get a VPC Instance Identity Token and exchange it for an IAM Access Token
When an HTTP 401 or 403 is returned from back-end service, it could mean the access token is expired. You can manually fetch a new VPC instance identity token, and exchange it for a new IAM access token. When using the method below, the new token is stored in
configuration and will be used for subsequent IAM enabled service calls.

Example:

```go
token, err := context.RefreshIAMToken()
if err != nil {
	return err
}
```
The `RefreshIAMToken()` method will detect if the user is logged in as a VPC VSI compute resource, and perform the token exchange
using the VPC metadata service using the trusted profile information stored in configuration.

For more details of the API, refer to the docs for [vpc](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/vpc). For more details of the `RefreshIAMToken` method, refer to the docs for [RefreshIAMToken](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config#Repository).

_Note_: Currently an IAM refresh token is not supported when authenticating as a VPC compute resource identity.


## 6. Utility for Unit Testing

We highly recommended that terminal.StdUI was used in your code for output, because it can be replaced by FakeUI which is a utility provided by IBM Cloud CLI for easy unit testing.

FakeUI can intercept all of output and store the output in a buffer temporarily, so that you can compare the actual output and expected output easily. Take the following code snippet as an example:

```go
// Here is the source code to be tested:
import (
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
)

type DemoPluginstruct {}

func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string){
    ...
    if args[0] == "start" {
        demo.Start(terminal.StdUI)
    }
}

func (demo *DemoPlugin) Start(ui terminal.UI){
    ...
    ui.echo("OK")
    ...
}

// The following is the testing code:
import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

func TestStart() {
    demoPlugin := DemoPlugin{}
    fakeUI := terminal.NewFakeUI()
    demoPlugin.Start(fakeUI)

    // Now you can use any third-party unit testing framework to assert the result,
    // for example:
    assert.IsTrue(fakeUI.ContainsOutput("OK"))
    assert.IsFalse(fakeUI.ContainsOutput("FAILED"))
}
```

### 6.1 Using the Test Doubles

When writing unit tests you may need to mock parts of the cli-sdk, dev-plugin, or your own interfaces. The dev-plugin uses [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) to mock interfaces that allow you to fake their implementation.

You can use the fakes implementation to mock the return, arguments, etc. Below are a few example of showing how to stub various methods in [PluginContext](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/plugin/plugin_context.go) and [PluginConfig](https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/blob/master/plugin/plugin_config.go).


```go
import (
  "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin/pluginfakes"
)
var (
  context *pluginfakes.FakePluginContext
  config  *pluginfakes.FakePluginConfig
)

BeforeEach(func() {
  context = new(pluginfakes.FakePluginContext)
  config = new(pluginfakes.FakePluginConfig)

  // mock the IsLoggedIn method to always return true
  // NOTE: expect method signature format to be `<METHOD_NAME>Returns()`
  context.IsLoggedInReturns(true)

  // mock the second IsLoggedIn call return false
  // NOTE: expect method signature format to be `<METHOD_NAME>ReturnsOnCall()`
  context.IsLoggedInReturnsOnCall(1, false)

  // stub the arguments for the first PluginConfig.Set method
  // NOTE: expect method signature format to be `<METHOD_NAME>ArgsForCall(...args)`
  config.SetArgsForCall("region", "us-south")
})

It("should call RefreshToken more than once", func() {

  // return the number of times a method is called
  // NOTE: expect method signature format to be `<METHOD_NAME>Calls()`
  Expect(context.RefreshIAMTokenCalls()).Should(BeNumerically(">", 0))
})

```

You can find other examples in [tests](https://github.com/maxbrunsfeld/counterfeiter/blob/master/generated_fakes_test.go) found in counterfeiter.

## 7. Globalization

IBM Cloud CLI tends to be used globally. Both IBM Cloud CLI and its plug-ins should support globalization. We have enabled internationalization (i18n) for CLI's base commands with the help of the third-party tool "[go-i18n](https://github.com/nicksnyder/go-i18n)". To keep user experience consistent, we recommend plug-in developers follow the CLI's way of i18n enablement.

Please install the *go-i18n* CLI for version 1.10.0. Newer versions of the CLI are no longer compatible with translations files prior to go-i18n@2.0.0. You can install the CLI using the command: `go install github.com/nicksnyder/go-i18n/goi18n@1.10.1`

Here's the workflow:

1.  Add new strings or replace existing strings with `T()` function calls to load the translated strings.  For example:

    ```go
    T("Installing the plugin...")

    //With variable substitution
    T("Plugin '{{.Name}}' was successfully installed.",
    map[string]interface{}{"Name": pluginName})
    ```

    **Note**: `T` is type of [TranslateFunc](https://godoc.org/github.com/nicksnyder/go-i18n/i18n#TranslateFunc) which is set in i18n initialization (see Step 4).

2.  Prepare translation files.
    1.  Add all strings in en-us.all.json, and then use go-i18n CLI to generate translation files for other languages. A sample en-us.all.json is like:

        ```javascript
        [
           {
                  "id": "Installing the plug-in…",
                  "translation": "Installing the plug-in…"
           },
           {
                  "id": "Plug-in '{{.Name}}' was successfully installed.",
                  "translation": " Plug-in '{{.Name}}' was successfully installed."
           },
           ...
        ]
        ```

      2.  To generate translation files for other language, such as: `zh\_Hans` (Simplified Chinese) and `fr\_BR`:
          1.  Create empty files `zh-hans.all.json` and `fr-br.all.json`, and run:
          2.  Run:
              ```bash
              goi18n -flat=false –outdir <directory\_of\_generated\_translation\_files> en-us.all.json zh-hans.all.json fr-br.all.json
              ```

      The previous command will generate 2 output files for each language: `xx-yy.all.json` contains all strings for the language, and `xx-yy.untranslated.json` contains untranslated strings. After the strings are translated, they should be merged back into `xx-yy.all.json`. If plugin is on the ibm-cloud-cli-sdk 1.00 or above, rename the file from `xx-yy.all.json` to `all.xx-yy.json`. For more details, refer to goi18n CLI's help by 'goi18n –help'.

3.  Package translation files. IBM Cloud CLI is to be built as a stand-alone binary distribution. In order to load i18n resource files in code, we use [go-bindata](https://github.com/jteeuwen/go-bindata) to auto-generate Go source code from all i18n resource files and the compile them into the binary. You can write a script to do it automatically during build. A sample script could be like:

    ```bash
    #!/bin/bash

    set -e

    go get github.com/jteeuwen/go-bindata/...

    echo "Generating i18n resource file ..."
    $GOPATH/bin/go-bindata -pkg resources -o bluemix/resources/i18n\_resources.go bluemix/i18n/resources
    ```

    After execution, a source file `bluemix/resources/i18n\_resources.go` with package name `resources` is generated to embed all files under bluemix/i18n/resource. To access a resource file, invoke resources.Asset("bluemix/i18n/resources/<fileName>") which returns bytes of the file.

4.  Initialize i18n

    `T()` must be initialized before use. During i18n initialization in IBM Cloud CLI, user locale is used if it's set in `~/.bluemix/config.json` (plug-in can get user locale via `PluginContext.Locale()`). Otherwise, system locale is auto discovered (see [jibber\_jabber](https://github.com/cloudfoundry/jibber_jabber)) and used. If system locale is not detected or supported, default locale `en\_US` is then used. Next, we initialize the translate function with the locale. Sample code:

    ```go
    import (
      "fmt"
      "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
    )

    var T i18n.TranslateFunc = Init(core_config.NewCoreConfig(func(e error) {}), new(JibberJabberDetector))

    func Init(coreConfig core_config.Repository, detector Detector) i18n.TranslateFunc {
	    bundle = i18n.Bundle()
	    userLocale := coreConfig.Locale()
	    if userLocale != "" {
		    return initWithLocale(userLocale)
      }
		}

    func initWithLocale(locale string) i18n.TranslateFunc {
        err := loadFromAsset(locale)
        if err != nil {
            panic(err)
        }
        return i18n.MustTfunc(locale, DEFAULT_LOCALE)
    }

    // load translation asset for the given locale
    func loadFromAsset(locale string) (err error) {
        assetName := fmt.Sprintf("all.%s.json", locale)
        assetKey := filepath.Join(resourcePath, assetName)
        bytes, err := resources.Asset(assetKey)
        if err != nil {
           return
        }
        _, err = bundle.ParseMessageFileBytes(bytes, resourceKey)
        return
    }
    ```

## 8. Command Design

### 8.1. Honor Region/Resource Group Setting of CLI

When users are using CLI, they probably have already targeted region or resource group during login. It's cumbersome to ask users to re-target region or resource group in specific command again.

- By default, plugin should honor the region/resource group setting of CLI. Check `CurrentRegion`, `HasTargetedRegion`, `CurrentResourceGroup`, and `HasTargetedResourceGroup` in the [`core_config.Repository`](https://godoc.org/github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config#Repository).

    ```go
    func (demo *DemoPlugin) Run(context plugin.PluginContext, args []string){
        config := context.PluginConfig()

        // get currently targeted region
        region := ""
        if config.HasTargetedRegion() {
            region = config.CurrentRegion().Name
            ui.Say("Currently targeted region: " + region)
        }

        // get currently targeted resource group
        group := ""
        if config.HasTargetedResourceGroup() {
            group = config.CurrentResourceGroup().Name
            ui.Say("Currently targeted resource group: " + group)
        }

        ...
    }
    ```

- If no region or resource group is targeted, it means target to **all regions and resource groups**.
- [Optional]: plugin can provide options like `-r, --region` or `-g` to let users to overwrite the corresponding setting of CLI.

## 9. Private Endpoint Support

Private endpoint enables customers to connect to IBM Cloud services over IBM's private network. Plug-ins should provide the private endpoint support whenever possible. If the user chooses to use the private endpoint, all the traffic between the CLI client and IBM Cloud services must go through the private network. If private endpoint is not supported by the plug-in, the CLI should fail any requests instead of falling back to using the public network.

**Choosing private endpoint**

IBM CLoud CLI users can select to go with private network by specifying `private.cloud.ibm.com` as the API endpoint of IBM Cloud CLI, either with command `ibmcloud api private.cloud.ibm.com` or `ibmcloud login -a private.cloud.ibm.com`.

Plug-ins can invoke `plugin.PluginContext.IsPrivateEndpointEnabled()` in the CLI SDK to detect whether the user has selected private endpoint or not.

**Publishing Plug-in with private endpoint support**

There is a field `private_endpoint_supported` in the plug-in metadata file indicating whether a plug-in supports private endpoint or not. The default value is `false`.  When publishing a plug-in, it needs to be set to `true` if the plug-in has private endpoint support. Likewise, in the plug-in Golang code, the `plugin.PluginMetadata` struct needs to have the `PrivateEndpointSupported` field set the same as this field in the metadata file. Otherwise the core CLI will fail the plug-in commands if the user chooses to use private endpoint.

If the plug-in for an IBM Cloud service only has partial private endpoint support in specific regions, this field should still be set to be `true`. It is the plug-in's responsibility to get the region setting and decide whether to fail the command or not. The plug-in should not fall back to the public endpoint if the region does not have private endpoint support.

**Private endpoints of platform services**

The CLI SDK provides an API to retrieve both the public endpoint and private endpoint of core platform services.

`plugin.PluginContext.ConsoleEndpoint()` returns the public endpoint of IBM Cloud console API if the user selects to go with public endpoint. It returns private endpoint of IBM Cloud console API if the user chooses private endpoint when logging into IBM Cloud.

`plugin.PluginContext.ConsoleEndpoints()` returns both the public and private endpoints of IBM Cloud console API.


`plugin.PluginContext.IAMEndpoint()` returns the public endpoint of IAM token service API if user selects to go with public endpoint. It returns private endpoint of of IAM token service API if the user chooses private endpoint when logging into IBM Cloud.

`plugin.PluginContext.IAMEndpoints()` returns both the public and private endpoints of the IAM token service API.


`plugin.PluginContext.GetEndpoint(endpoints.Service)` returns both the public and private endpoints for the following platform services

- global-search
- global-tagging
- account-management
- user-management
- billing
- enterprise
- global-catalog


## 10. Deprecation Policy

A CLI plugin should maintain backward compatibility within a major version, but MAY introduce incompatible changes in command format, parameters, or behavior in a major release of the plugin.

When a new major version of a plugin introduces incompatible changes, support for the prior major version of the plugin may only be withdrawn one year from an official deprecation notice for that version of the plugin.
