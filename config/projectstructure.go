package config

// Default project structure.

const defaultProjectStructure = `
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
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@findings.md",
            "info": {
                "name": "@findings.md",
                "isdir": false,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@creds.md",
            "info": {
                "name": "@creds.md",
                "isdir": false,
                "template": "creds"
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
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@report",
            "info": {
                "name": "@report",
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@Report/report.json",
            "info": {
                "name": "report.json",
                "isdir": false,
                "template": "report"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@clientFiles",
            "info": {
                "name": "@clientFiles",
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@TODO.md",
            "info": {
                "name": "@TODO.md",
                "isdir": false,
                "template": ""
            },
            "children": []
        }
    ]
}`
