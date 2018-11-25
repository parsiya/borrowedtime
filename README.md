# Borrowed Time
Borrowed time is my personal project management application. As a security consultant, I switch projects every week or other week. This tool helps automate some tasks such as project structure creation.

It's named after a discipline priest spell in World of Warcraft:

* https://www.wowhead.com/spell=213642/borrowed-time

Note: Borrowed Time is under heavy development. The API of sub-packages and `prompter` might be modified. `prompter` is not released yet.

## Capabilities

* Create and manage a workspace with a custom configuration.
* Create projects based on a json file representing the directory structure.
    * Project templates are json files. See below for details.
    * Each file can have a specific template using Go's template engine. It can incorporate custom items from the configuration file without having to recompile the app. See below for details.
* Use your favorite editor.
* Automatic editor detection on Windows.
* Add/remove/edit data files and project/file templates.

## Operating System Support
Borrowed Time is mainly developed and used on Windows. Currently only Windows is supported. However, no OS dependent libraries are in use so default editors/locations/paths can be added for other Operating Systems.

## Installation
Borrowed Time depends on the following external packages:

* https://github.com/basgys/goxml2json
* https://github.com/olekukonko/tablewriter
* https://github.com/c-bata/go-prompt
* https://github.com/starkriedesel/prompter. **This package is not released yet.**

## Quickstart
Execute Borrowed Time and run the `Deploy` command. It creates a directory named `borrowedtime` in your home directory and opens the configuration file in the default editor in your OS. Home is based on [go-homedir](https://github.com/mitchellh/go-homedir). Home on Windows is `C:\Users\[your-user]\`.

```
borrowedtime
│   config.json
│
├───backups
├───data
└───templates
        creds.json
        notes.json
        project-config.json
        project-structure.json
```

## Configuration File
Borrowed Time uses a configuration file to persist settings. It's a simple JSON text file. New entries can be added manually. These can be used in file/project templates. Sample configuration file on Windows:

``` json
{
    "burppath": "",
    "editor": "C:\\Program Files\\Microsoft VS Code\\Code.exe",
    "projectstructure": "project-structure",
    "workspace": "C:\\Users\\Parsia\\Desktop\\projects",
    "yourname": ""
}
```

The `config` subcommand is used to manipulate the configuration. For example, `config edit` opens the configuration file in your default editor.

![config command](.github/configcmd.png)

## Templates
Borrowed Time supports customizing project structure and generated files using templates.

### Project Templates
Project templates are JSON files in the following structure using Go's template engine.

* `Workspace` points to the root of your workspace nad `ProjectName` to the name of the directory.

``` json
{
    "path": "{{ .Workspace }}/{{ .ProjectName }}",
    "info": {
        "name": "{{ .ProjectName }}",
        "isdir": true,
        "template": ""
    },
    "children": [
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/.config.json",
            "info": {
                "name": ".config.json",
                "isdir": false,
                "template": "project-config"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@notes.md",
            "info": {
                "name": "@notes.md",
                "isdir": false,
                "template": "notes"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@pix",
            "info": {
                "name": "@pix",
                "isdir": true,
                "template": ""
            },
            "children": []
        },
    ]
}
```

### File Templates
File templates are simple text files. They can contain similar placeholders. For example, the `notes` template is:

```
# {{ .ProjectName }} Notes

## 
```

### Modifying Templates
File templates can be modified directly. Add new directories, files, and assign file templates. New templates can be added manually by dropping them into the `templates` directory. The `template add` command can also be used to add a new file as a template. It simply copies the content of the file into the directory.

You can use sub-directories to manage templates but each template name must be unique application wide. This means you cannot have two files named `notes.json` inside the `templates` directory in different sub-directories.

The template is addressed by the name of the file containing it without the extension. As a matter of convenience, all templates get the extension json although file templates are just plain text files.

### Custom Fields in File Templates
It's possible to add custom items to the configuration file and use them in file templates without rebuilding the application. The project struct has a field named `WorkspaceConfig` that is populated with a map of key/values from the configuration file.

1. Add a new key to the config file. For example, `"customkey": "value"`.
2. Inside the file template, add: `{{ index .WorkspaceConfig "customkey" }}`. Remember keys are case-sensitive. As a matter of convenience, only use lowercase keys.
3. This will be replaced by the value of `customkey` from the configuration file.

Personally, I am very proud of how this part turned out.

For more information about Go templates, please read: https://golang.org/pkg/text/template/.

## License
Opensourced under the Apache License v 2.0 license. See [LICENSE](LICENSE) for details.

## TODO:

1. Change error message in all unexported functions and remove module name. Only leave the error message?
2. ~~Add `edit` command to template and data files.~~
3. Update docs after `prompter` is released.
4. Create gifs of some commands.
5. Explain all commands.
6. Add generation of Burp project based on a base Burp configuration file.
    * Research and add custom generation of Burp config files based on project using name, credentials, etc.
7. Add `dep` support.