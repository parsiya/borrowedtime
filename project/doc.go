/*
	Package project manages individual projects on the machine.
	It is used when starting a new projects to create the directory structure.

	Default project directory location is: homedir > {{project name}}

	Inside the directory:
	{{ project name }}
		.{{project name}}-config.json <-- Project config file.
		@Notes.md <-- Main notes file.
		@Findings.md <-- Findings file.
		@Pix - directory <-- Evidence images go here.
		@ClientFiles - directory <-- All client files go here.
		@{{projectname.burpconfig}} <-- Burp config file. TODO: Create Burp module to create files.
										Later it can be used to generate a Burp file.
		@Credz.md <-- Credentials file.
		@TODO.md <-- TODO file.


	TODO:
	1. ~~Find a way to pass project info to the templates to create them? Maybe get it from path?~~ Done
	2.
*/

package project
