# release-me

Utility to validate and generate a release note from the PR title, `Release Note` comment section in the PR description. The `validate` can be used during the PR validation to check the PR have atleast one of labels to categorize this change to the generated release notes.

## Validator

Checks the existance of labels required for generating release notes

### Valid labels

```text
"breaking", "misc", "bug", "enhancement"
```

### Generator

A release note is generated through fetching all the pull requests merged after the latest tag (release) of the repository. The release note is outputted to stdout.
