package config

// Default project structure.

const defaultProjectStructure = `
{
    "path": "{{ .Workspace }}/{{ .ProjectName }}",
    "info": {
        "isdir": true,
        "template": ""
    },
    "children": [
        {
            "path": ".config.json",
            "info": {
                "isdir": false,
                "template": "project-config"
            },
            "children": []
        },
        {
            "path": "@notes.md",
            "info": {
                "isdir": false,
                "template": "notes"
            },
            "children": []
        },
        {
            "path": "@findings.md",
            "info": {
                "isdir": false,
                "template": ""
            },
            "children": []
        },
        {
            "path": "@creds.md",
            "info": {
                "isdir": false,
                "template": "creds"
            },
            "children": []
        },
        {
            "path": "@pix",
            "info": {
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "@report",
            "info": {
                "isdir": true,
                "template": ""
            },
            "children": [
                {
                    "path": "report.json",
                    "info": {
                        "isdir": false,
                        "template": "report"
                    },
                    "children": []
                }
            ]
        },
        {
            "path": "@clientFiles",
            "info": {
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "@TODO.md",
            "info": {
                "isdir": false,
                "template": ""
            },
            "children": []
        },
        {
            "path": ".gitignore",
            "info": {
                "isdir": false,
                "template": ""
            },
            "children": []
        }
    ]
}`
