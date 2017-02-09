[![Build Status](https://travis-ci.org/nfisher/daywalker.svg?branch=master)](https://travis-ci.org/nfisher/daywalker)

# daywalker

Download transitive dependencies for maven co-ordinate into the current folder. Good for downloading dependencies when using build tools like "buck":https://buckbuild.com/.

## Maven Coordinates

The maven coordinates are inspired by buildr where it has the following format:

```
${GROUPID}:${ARTIFACT}:${VERSION}
```

## Sample Execution

```
daywalker com.sparkjava:spark-core:2.5.4
```

## Improvements

- parallel download of jar files once dependency graph is walked.
- track parent pom and prevent attempts to download associated jar.
- evaluation execution.
- download list of maven coordinates.
- alternative repo locations.
- download to specified folder.
- download of test dependencies in test folder or equivalent.
- generate buck entry.
