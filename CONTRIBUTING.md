# Contributing to Supertext for the Portland Opera to Go

How to set up, code, test, review, and release so contributions meet our Definition of Done.

## Code of Conduct

This project follows the [Oregon State Code of Student Conduct](https://scs.oregonstate.edu/code-student-conduct).

Please report any violations to the current Project Manager. If the concern is regarding the PM, contact the instruction team.

## Branching & Workflow

This project follows Trunk-based development workflows. Smaller, frequent changes are preferred to larger changes.
Pull request all changes to the `main` branch. Create separate branches following the format `[name]/[feature-or-bug-fix]`. For
example, `jamie/fix-cicd`.

Rebase or squash should be used for merging.

## Issues & Planning

Issues are tracked via Github Issues on the project. You may assign a project member if you know they are responsible for the code
needing an issue. Issues must be labeled with either `bug` or `enhancement` to assist in triage, with `bug` receiving a priority value.

## Commit Messages

Please follow the [Scoped Commits](https://scopedcommits.com/) standard for commit messages. Ensure that
contextual information such as the modules changes is provided. For example, `websockets: add support for multiple users`.

A summary of the changes must be included within the commit body. This should explain roughly how a problem was fixed, not any context regarding
why.

Issues must be referenced with `Fixes: #123, #456` within the commit body.

## Code Style, Linting & Formatting

Go provides a built-in formatter with `gofmt`.

Run `go fmt` within the repository before merging any pull requests.

## Testing

Go provides a built in testing framework. Tests relating to backend services will be run through `go test`

## Definition of Done

All CI jobs must pass before a pull request is merged, unless there is extenuating circumstances proving otherwise (such as a broken CI job).
Each pull request must be reviewed by one other member who was not responsible for writing the code or docs.

The primary CI jobs that must pass: the test suite, formatting, commit format. Note that these gates do not exist yet, as there is not yet any code to check.
As soon as any code is written, the Quality Owner is responsible for providing proper CI/CD coverage.

## Contribution Norm

Every project member must land at least one piece of accepted work per sprint with exception for extrenuating circumstances.
Code and documentation will be provided within the main repository, the Github issue tracker tracks bug and feature progress. 

## Pull Requests & Reviews

A pull request must be reviewable. To ensure this, make sure the PR body is equivalent to the commit body or feature ticket title.
At least one reviewer who was not responsible for the code must approve the PR. While this is not required, reviewers with familiarity
of the section of code being modified are encouraged to approve.

Approval can only be granted if CI jobs passed.

Reviewers must ensure that the PR follows all guidelines within this document, as well as ensuring that the application functions
as expected. This may require building and testing a feature manually where automated testing is not possible.

## CI/CD

All pipelines exist within the `.github/workflows` directory and are maintained by the Quality Owner. All pipelines must pass before a PR is merged.

## Security & Secrets

Do not hard code secrets anywhere within the repository. Do not hard code confidential data within the repository. This may include client or patron information,
student information, or local environment information. Use configuration files and environment variables.

Dependency updates should be done for any major security fixes or minor releases. This project uses Dependabot to ensure that dependency updates
are performed.

Security vulnerabilites must be submitted through Github.

## Documentation Expectations

Documentation should be updated whenever a breaking change is made justifying the change. Documentation should be updated at the same time
as any breaking change, do not submit multiple PRs.

The README should remain static unless it becomes out of date.

`docs/` is a living repository of its own and may be changed regardless to code status. This includes documentation such as RFCs and ADRs. 

CHANGELOG is created by a script by the Quality Owner for a release.

## Release Process

This project uses semantic versioning. Tags will be produced on the last commit for a release (generating a changelog).

Immediately after a release, a version bump commit must be made for the next release.

## Support & Contact

Maintainers are active on Discord and should respond within 24 hours.
