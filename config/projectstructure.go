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
            "path": "{{ .Workspace }}/{{ .ProjectName }}/.config.json",
            "info": {
                "isdir": false,
                "template": "project-config"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@notes.md",
            "info": {
                "isdir": false,
                "template": "notes"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@findings.md",
            "info": {
                "isdir": false,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@creds.md",
            "info": {
                "isdir": false,
                "template": "creds"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@pix",
            "info": {
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@report",
            "info": {
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@Report/report.json",
            "info": {
                "isdir": false,
                "template": "report"
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@clientFiles",
            "info": {
                "isdir": true,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/@TODO.md",
            "info": {
                "isdir": false,
                "template": ""
            },
            "children": []
        },
        {
            "path": "{{ .Workspace }}/{{ .ProjectName }}/.gitignore",
            "info": {
                "isdir": false,
                "template": ""
            },
            "children": []
        }
    ]
}`
