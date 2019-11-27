# Contributing

We are pleased to welcome all contributors who willing to join with us in our journey.


## Pull requests

To submit a new feature or bug fix,

* Fork the repository.
* Create a new branch for your changes. If the changes are related to GitHub Issue you can use 'issue-<id>' as branch name
* Develop the feature or bug fix.
* Add new test cases
* Add/Modify the documentation if necessary.
* Submit the PR.


## Issues


We use [GitHub issues](https://github.com/wso2/cellery-controller/issues/new) to track bugs, feature requests and questions.

When reporting a bug please include at least Mesh Controller version number and the git commit.

You can use following command to get version information about the Mesh Controller

```bash
kubectl logs -n cellery-system -l app=controller | grep 'Mesh Controller'
```
