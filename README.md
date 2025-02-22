[中文](./README-CN.md)

# Harbor Cleaner

Clean images in Harbor by policies.

## Features

- **Delete tags without side effects** As we known when we delete a tag from a repo in docker registry, the underneath manifest is deleted, so are other tags what share the same manifest. In this tool, we protect tags from such situation.
- **Delete by policies** Support delete tags by configurable policies
- **Dry run before actual cleanup** To see what would be cleaned up before performing real cleanup.
- **Cron Schedule** Schedule the cleanup regularly by cron.

## Concepts

| Image | Project | Repo | Tag |
|--|--|--|--|
| library/busybox:1.30.0 | library | busybox | 1.30.0 |  
| release/devops/tools:v1.0 | release | devops/tools | v1.0 |

## Policies

### Number Policy

Number policy will retain latest N tags for each repo and remove other old ones. Latest tags are determined by images' creation time.

```yaml
numberPolicy:
    number: 5
```

This policy takes only one argument, the number of tags to retain.

### Regex Policy

Regex policy removes images that match the given repo and tag regex patterns. A tag will be removed only when following conditions are all satisfied:

- It matches at least one repo pattern
- It matches at least one tag pattern

Regex here are `Golang` supported regex. For example `.*` matches all.

```yaml
regexPolicy:
    repos: [".*"]
    tags: [".*-alpha.*", "dev"]
```

The above policy config will remove tags from all repos that are 'dev' or contain 'alpha'. For example,

- dev
- v1.0.0-alpha
- v1.4.0-alpha.2
- 1.0-alpha.5

## Recently Not Touched Policy

This policy works depends on Harbor's access log. It collects images that are recently touched (pull, push, delete), and remove all other images that are not touched recently. It takes a time in second to configure the time period.

```yaml
notTouchedPolicy:
  time: 604800
```

## How To Use

### Get Image

```bash
$ make image VERSION=latest
```

You can also pull one from DockerHub.

```bash
$ docker pull zhangguanzhang/harbor-cleaner:v0.1.0
```

### Configure

```yaml
# Host of the Harbor
host: https://dev.cargo.io
# insecure ssl
insecure: false
# Version of the Harbor, e.g. 1.7, 1.4.0
version: 1.7
# Admin account
auth:
  user: admin
  password: Pwd123456
# Projects list to clean images for, it you want to clean images for all
# projects, leave it empty.
projects: []
# Policy to clean images
policy:
  # Policy type, e.g. "number", "regex", "recentlyNotTouched"
  type: number

  # Number policy: to retain the latest N tags for each repo
  # This configure takes effect only when 'policy.type' is set to 'number'
  numberPolicy:
    number: 5

  # Regex policy: only clean images that match the given repo patterns and tag patterns
  # This configure takes effect only when 'policy.type' is set to 'regex'
  regexPolicy:
    # Regex to match repos, a repo will be regarded as matched when it matches any regex in the list
    repos: [".*"]
    # Regex to match tags, a tag will be regarded as matched when it matches any regex in the list
    tags: [".*-alpha.*", "dev"]

  # Recently not touched policy: clean images that not touched within the given time period
  # This configure takes effect only when 'policy.type' is set to 'recentlyNotTouched'
  notTouchedPolicy:
    # Time in second that to check for images
    time: 604800

  # Tags that should be retained anyway, '?', '*' supported.
  retainTags: []
# Trigger for the cleanup, if you only want to run cleanup once, remove the 'trigger' part or leave
# the 'trigger.cron' empty
trigger:
  # Cron expression to trigger the cleanup, for example "0 0 * * *", leave it empty will disable the
  # trigger and fallback to run cleanup once. Note: you may need to quote the cron expression with double quote.
  # Time zone of the cron depends on the running environment, if run in docker container, it's UTC time.
  cron:

```

In the policy part, exact one of `numberPolicy`, `regexPolicy`, `notTouchedPolicy` should be configured according to the policy type. 

### DryRun

```bash
$ docker run -it --rm \
    -v <your-config-file>:/workspace/config.yaml \
    zhangguanzhang/harbor-cleaner:latest --dryrun=true
```

### Clean

```bash
$ docker run -it --rm \
    -v <your-config-file>:/workspace/config.yaml \
    zhangguanzhang/harbor-cleaner:latest
```

### Cron Schedule

Configure the cron trigger and run harbor cleaner container in background.

```yaml
# Trigger for the cleanup, if you only want to run cleanup once, remove the 'trigger' part or leave
# the 'trigger.cron' empty
trigger:
  # Cron expression to trigger the cleanup, for example "0 0 * * *", leave it empty will disable the
  # trigger and fallback to run cleanup once.
  cron: 0 0 * * *
```

```bash
$ docker run -d --name=harbor-cleaner --rm \
    -v <your-config-file>:/workspace/config.yaml \
    zhangguanzhang/harbor-cleaner:latest
```

## Supported Version

- v2

