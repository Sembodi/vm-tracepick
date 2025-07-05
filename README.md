# See README in ./GoProject/

# Contents of /GoProject:

## /cmd
This directory contains drivers that fulfill larger-scale tasks, such as generating artifacts, building and running docker images, and exporting relevant metrics for evaluation.
## /lib
This directory contains all modules that assist the drivers in completing their tasks, such as detecting configurations, finding OS info, tracing application usage, along with many other features.
## /lib/helpers
This directory contains most components that the drivers call. From this part of the library, other modules (i.e. in `/lib/*`) are called as well.
## /bashrun
This directory contains templates of bash scripts for batch runs.
## /bashscripts
This directory contains auxiliary Bash scripts that support various automation tasks, such as querying Docker or resetting relevant entries in the known_hosts file.
## /benchmarks
This directory contains scripts and outputs used for evaluating the performance of our framework.
## /config
In this directory, users can customize extra excluded services. For future purposes, more customization could be included here.
## /connectyamls
This directory contains YAML files with connection information and identification keys to connect to the remote VM using ssh.
## /defaultdata
This directory contains default configurations that should not be altered by the user, such as default excluded services, the copying algorithm specified in the thesis, and the set of protected paths, fetched from Linux' Filesystem Hierarchy Standard (FHS) using a bashscript specified in /bashscripts.
## /visualizeGoProject
This directory is not necessary for our framework to run. However, it provides insight into how all modules are connected to one another through graphing mechanisms. One thing to consider when using this, is that including all components makes the resulting graphs very dense.
