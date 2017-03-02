[![Build Status](https://travis-ci.org/nfisher/daywalker.svg?branch=master)](https://travis-ci.org/nfisher/daywalker)

# daywalker

Brings those transitive dependencies for maven co-ordinate into the light. No need to hide in the dark! My current use case is for downloading dependencies when using build tools like "buck":https://buckbuild.com/.

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


## Graph Relationship

project
  |--[ compile             ]--> project
  |--[ managed_compile     ]--> project
  |--[ managed_provided    ]--> project
  |--[ managed_runtime     ]--> project
  |--[ managed_system      ]--> project
  |--[ managed_test        ]--> project
  |--[ parent              ]--> project
  |--[ property            ]--> property --[ value ]--> value
  |--[ provided            ]--> project
  |--[ runtime             ]--> project
  |--[ system              ]--> project
  |--[ test                ]--> project
  |--[ unresolved_compile  ]--> project
  |--[ unresolved_provided ]--> project
  |--[ unresolved_runtime  ]--> project
  |--[ unresolved_system   ]--> project
  +--[ unresolved_test     ]--> project
