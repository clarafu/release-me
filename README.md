# release-me

CLI to generate a release note for your GitHub repository.

Example of release notes generated using this CLI: [Concourse CI releases](https://github.com/concourse/concourse/releases)

## What is release-me

It generates a release note through fetching all the pull requests merged after the latest tag of the repository. 

**The reason I created this was because I needed to have a release note generator that can create release notes for older release branches.**

For example, my latest release is 7.4.0 but I need to generate a release note and release a patch for 6.4.0 and then in the future create a 7.5.0 release. I didn't see any open-source release note generators that was able to do that for me successfully.


### How to use it?

There are two commands that you can run using this CLI: `generate` and `validate`. Both these flags require the following flags.

| Flag             | Example      | Required   | Desciptions           
| ---------------- | ------------ | ---------- | ---------------------
| `github-owner`   | `clara`      | True       | Login field of a github user or organization.
| `github-repo`    | `release-me` | True       | Name of the GitHub repository.
| `github-token`   | `60497df..`  | True       | GitHub OAuth token to authenticate with.


### Generating a release note

The `generate` command accepts the following flags

| Flag                    | Example     | Required | Desciptions           
| ----------------------- | ----------- | -------- | ---------------------
| `release-version`       | `1.0.0`     | True     | The version that the release note will be generated for.
| `github-branch`         | `master`    | False    | The branch name of the GitHub repository to pull the pull request from. Defaults to master.
| `last-commit-SHA`       | `d6cd1..`   | False    | Generates a release note using all prs merged up to this commit SHA. If it is empty, it will generate a release note until latest commit.
| `ignore-authors`        | `clara,alex`| False    | Comma separated list of github handles. Any PRs authored by these handles will be ignored.
| `ignore-release-regex`  | `1.2.*`     | False    | A regular expression indicating releases to ignore when determining the release to start generating the release note from.


For example, you can generate a release note using the following command:

```
./releaseme generate \
  --github-token=$GITHUB_TOKEN \
  --github-owner=$GITHUB_OWNER \
  --github-repo=$GITHUB_REPO \
  --github-branch=$GITHUB_BRANCH \
  --last-commit-SHA=$LAST_COMMIT_SHA \
  --release-version=$RELEASE_VERSION \
  --ignore-authors=dependabot \
```

The CLI grabs all the pull requests merged after commit that is referenced by the latest tag. Then it sorts the pull requests by number in ascending order and fetches the optional release note description from the pull request body. It uses the labels on the pull request to sort them into sections (and also priority) and uses the go templating library to construct the release note and output it to stdout.

The CLI depends on certain labels to exist on each pull request in order to group them into the correct sections. This means that the pull request reviewer must label the pull request before merging with the label(s) that they think best fit. Typically, you should only need to label it with one of the following labels but if the reviewer decides to attach more than one label, the CLI will group the pull request based off the labels' hierarchy. 

Here's an overview of the labels (and sections in the release note) that need to be considered:

`breaking`: this means that the pull request contains a breaking change or in other words backwards-incompatible.

`enhancement`: another word for feature, where we are introducing something new.

`bug`: a pull request that fixes an existing bug.

`misc`: you would add misc label if the pull request introduces no behavioural changes.

You can also add a `priority` label to the pull request if you want it to be at the top of the section. If there are multiple pull requests with the `priority` label in the same section, it will be ordered by pr number.

The release note will be generated using the *title*, *pr number*, *author* and *optional release note description*. The optional release note description will be found in the pull request description/body under the header `## Release Note`. It will be found using regex so it will also accept things like `# Release Note` or `## release notes`.

An example of the note that it will generate:

## ✈️ Features

* Add flag to ignore specific PR authors (#3) @chenbh <sub><sup><a name="3" href="#3">:link:</a></sup></sub>  
  * Can be used for pull requests created from bots, ex. dependabot

To summarize, a pull request will need to be labeled with either breaking, enhancement, bug or release/no-impact for it grouped into the corresponding section in the release note and for the test to pass in order to merge the pr.

The way that I used this to generate release notes for older versions is through having the older versions on a separate branch for the major version. For example, I would have a `master` branch, `7.x` branch and a `6.x` branch. When I release a new major version I would create a new branch for it. Then when I need to release a new version for `6.4.0`, I would run the release note generater with the additional `github-branch` flag set to `6.x` and it will grab all the prs merged to the `6.x` branch from the last release made from the branch.

### Validating the labels on a pull request

You can also validate that the pull request has valid labels through the `validate` command.

| Flag             | Example      | Required   | Desciptions           
| ---------------- | ------------ | ---------- | ---------------------
| `pr-number`      | `123`        | True       | Checks the existance of labels required to generate release note in this pr.

For example, you can validate a pull request by:

```
./releaseme validate \
  --github-token=$GITHUB_TOKEN \
  --github-owner=$GITHUB_OWNER \
  --github-repo=$GITHUB_REPO \
  --pr-number=123 \
```
