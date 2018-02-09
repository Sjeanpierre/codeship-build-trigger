# codeship-build-trigger
Trigger dependent codeship builds using API v2

## Setup
*Required environment variables* 

`CODESHIP_USERNAME`

`CODESHIP_PASSWORD`

## build_trigger.yml
This is the central configuration file for the build triggering process. All projects listed will be triggered for build. File must be in same directory.

simply provide the project UUID & branch name in the form if `heads/BRANCH_NAME`

see [build_trigger.yml.example](https://github.com/Sjeanpierre/codeship-build-trigger/blob/master/build_trigger.yml.example) for examples

## Usecases
* Trigger dependent builds from Codeship steps adding the following as build step in master build
```
- name: Launch additional builds
  service: app
  command: bash -c "./bin/codeship-build-trigger "
```

* Trigger multiple builds from command line by simply executing binary `./codeship-build-trigger`
