# Go Project: TracePick

## See lib/own/owntypes for own type definitions and global variables. Variable (default value):

- `UseSudo` (true): enable if 'sudo' should be prepended to remote command strings (for source VM). ArtifactGenerator checks whether root user is already used (in which case, UseSudo becomes false).
- `DoTracing` (true): disable if tracing data has already been retrieved in a previous run (saves time).
- `CachedBuild` (false): enable if we do NOT want to measure the build time (disable if we do).
- `SkipDialogue` (false): enable for batch runs (also make sure the user can do 'sudo /bin/cp' and 'sudo /bin/du' locally without password)


## Requirements:
- Make sure to have the following tools installed on your target vm:
  - `apt update`
  - `apt install bpftrace`
  - `apt install apt-rdepends`
  - `apt install wrk`
  <!-- - `apt install ffuf` -->

  <!-- Run in root directory:
  - `git clone https://github.com/danielmiessler/SecLists.git` -->

- Make sure to have the Docker engine running on the local machine.

- Run `make help` to see more details about the data pipeline.

- Make sure the following commands can run on your local machine:
  - `go`
  - `make`
  - `ssh`
  - `tar`
  - `docker`
  - `du` (evaluation purposes)
  - `python3`

- In case you want to create dependency graphs (optional), make sure you have `dot` installed on the local machine.


## Troubleshooting:
Sometimes runtime extraction won't work exactly as needed. If running against some errors, one can always go through the following debug options:
- Check the `outputs/...` folder for empty files (usually means failed command)
- Search the full project (shift + control/cmd + F) the erroneous output path (e.g. `outputs/bpfout`) to find the driving function.

### Certificate error when downloading packages

Potential fix (Debian/Ubuntu):
`wget -qO - <pgp key link> | apt-key add -`, or in Dockerfile:
`RUN apt install -y wget gnupg
RUN wget -qO - <pgp key link>| apt-key add -`

How to find the pgp key link:
- Look in the corresponding Dockerfile to find the associated website of the package.
- Look for static files on the website with the appropriate pgp keys for your package version
