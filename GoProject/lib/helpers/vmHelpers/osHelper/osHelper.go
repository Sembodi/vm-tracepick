package osHelper

import (
  "fmt"
  "os"
  "strings"
  "bufio"
  "project1/tracepick/lib/own/owntypes"
  // "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/runners/remoterun"
)

// This function extracts AT MOST the first two terms of the distribution name, along with ONLY the first version number
func GetVMOS()  {
  var (
    programName string = "GetVMOS"
    nameCmd string = `
      awk -F= '
        /^ID=/{
          gsub(/"/, "", $2);
          print $2
        } ' /etc/os-release
    `
    //`
    // awk -F= '
    //   /^NAME=/ {
    //       gsub(/"/, "", $2);
    //       split($2, parts, " ");
    //       name = parts[1];
    //       if (length(parts) > 1) name = name " " parts[2];
    //   } END { print name;
    //   } ' /etc/os-release
    // `

    versionCmd string = `
    awk -F= '
      /^VERSION_ID=/ {
          gsub(/"/, "", $2);
          version = $2;
      } END { print version;
      } ' /etc/os-release
    `
    //   awk -F= '
    //     /^NAME=/{
    //         gsub(/"/, "", $2);
    //         split($2, nameParts, " ");
    //         name = nameParts[1];
    //         if (length(nameParts) > 1) name = name " " nameParts[2]
    //     }
    //     /^VERSION_ID=/{
    //         gsub(/"/, "", $2);
    //         split($2, verParts, ".");
    //         version = verParts[1];
    //         for (i=2; i<=length(verParts); i++) version = version " " verParts[i]
    //     }
    //     END{ print name, version }
    //   ' /etc/os-release
    // `
    outDir string = "outputs/os/"

    cmdMap owntypes.CommandMap = make(owntypes.CommandMap)
  )

  cmdMap["os"] = nameCmd
  cmdMap["version"] = versionCmd

  fsHelper.CleanPath(outDir)

  remoterun.RemoteRun(programName, cmdMap, outDir)
}

func ParseVMOS(property string) string {
  var (
    programName string = "ParseVMOS"
    path string = "outputs/os/"

    propStr string
    err error
  )
  if propStr, err = lineparse.ReadContent(programName, path + property); err != nil { return "" } //only one line, hence ReadContent suffices

  //todo add prompt to shorten name if necessary

  scanner := bufio.NewScanner(os.Stdin) //Going to need user input

  fmt.Println("Found", property, ":", propStr)

  propStr = owndialogues.EditString(scanner, propStr, false)

  return propStr
}


// KEEP COMMENTED
// var ReturnSomeString owntypes.CommandMap = owntypes.CommandMap{
//   "ubuntu": ,
//   "alpine": ,
//   "centos": ,
//   "arch": ,
//   "opensuse": ,
// } KEEP COMMENTED

var ReturnUpdateString owntypes.CommandMap = owntypes.CommandMap{
  "ubuntu": "apt update",
  "alpine": "apk update",
  "centos": "yum makecache",
  "arch": "pacman -Sy",
  "opensuse": "zypper refresh",
}

var ReturnSmallInstallString owntypes.CommandMap = owntypes.CommandMap{
    "ubuntu": "apt install -y %s",
    "alpine": "apk add %s",
    "centos": "yum install -y %s",
    "arch": "pacman -S --noconfirm %s",
    "opensuse": "zypper install -y %s",
}


var ReturnInstallString owntypes.CommandMap = owntypes.CommandMap{
    "ubuntu": `apt install -y %[1]s || { echo "deb [trusted=yes] %[2]s" | tee /etc/apt/sources.list.d/extra.list && apt update && apt install -y %[1]s; }`,
              //`apt install -y %[1]s || { codename=$(grep VERSION_CODENAME /etc/os-release | cut -d= -f2) && echo "deb [trusted=yes] %[2]s $codename main" | tee /etc/apt/sources.list.d/extra.list && apt update && apt install -y %[1]s; } || echo 'Executed with errors'`,
    "alpine": "apk add %s",
    "centos": "yum install -y %s",
    "arch": "pacman -S --noconfirm %s",
    "opensuse": "zypper install -y %s",
}

var ReturnInstallAddAptString owntypes.CommandMap = owntypes.CommandMap{
    "ubuntu": `apt install -y software-properties-common`,
    "alpine": "echo 'no software-properties-common installed'",
    "centos": "echo 'no software-properties-common installed'",
    "arch": "echo 'no software-properties-common installed'",
    "opensuse": "echo 'no software-properties-common installed'",
}

var ReturnPurgeAddAptString owntypes.CommandMap = owntypes.CommandMap{
    "ubuntu": `apt purge -y software-properties-common`,
    "alpine": "echo 'no software-properties-common purged'",
    "centos": "echo 'no software-properties-common purged'",
    "arch": "echo 'no software-properties-common purged'",
    "opensuse": "echo 'no software-properties-common purged'",
}

var ReturnCleanString owntypes.CommandMap = owntypes.CommandMap{
  "ubuntu": "apt clean && rm -rf /var/lib/apt/lists/*",
  "alpine": "rm -rf /var/cache/apk/*",
  "centos": "yum clean all",
  "arch": "pacman -Scc --noconfirm",
  "opensuse": "zypper clean --all",
}

var ReturnEmptyString owntypes.CommandMap = owntypes.CommandMap{
  "ubuntu": "",
  "alpine": "",
  "centos": "",
  "arch": "",
  "opensuse": "",
}


//call:
// trimItem := function{convert php8.1-fpm => php8.1}
// cmd = fmt.Sprintf(str, item)
// remoterun(cmd)
// search (FULL) service names as substrings in listed suggestions and take those. Comment non-matching suggestions, but take them along
var ReturnSearchPackageString owntypes.CommandMap = owntypes.CommandMap{
    "ubuntu": `{ dpkg -l | awk '{print $2}' | grep '%[1]s'; } || { dpkg -l | awk '{print $2}' | grep $(dpkg -S %[1]s.service | cut -d: -f1); }`, //`apt-mark showmanual | grep '%s.*'`,
    "alpine": `rpm -qa | grep -o '^%[1]s-\?[a-z]*/'`, //apk search %[1]s
    "centos": `yum search --names %[1]s | grep -o '^%[1]s-\?[a-z]*/'`,
    "arch": `pacman -Ss %[1]s | grep '^.*\/.*%[1]s' | awk '{print $1}'`,
    "opensuse": `zypper search --match-substrings %[1]s | awk 'NR > 2 {print $3}'`,
}

var ReturnPackageRepoString owntypes.CommandMap = owntypes.CommandMap{
    "ubuntu": `apt-cache policy %s | grep 'http[s]\?://[^ ]*' | sort -u | head -n1 | sed 's/^[[:space:]]*[0-9]\+[[:space:]]//' | awk '{print $1 " " $2}' | sed -E 's#(.*)/([^/]+)$#\1 \2#'`,
 // "ubuntu": `apt-cache policy %s | grep -o 'http[s]\?://[^ ]*' | sort -u | head -n1`,
    "alpine": ``,
    "centos": ``,
    "arch": ``,
    "opensuse": ``,
}


func GetOSCommand(osName string, returnMap owntypes.CommandMap) (string, error) {
  base := strings.ToLower(osName)

  switch {
    case strings.Contains(base, "ubuntu"), strings.Contains(base, "debian"):
      return returnMap["ubuntu"], nil
    case strings.Contains(base, "alpine"):
      return returnMap["alpine"], nil
    case strings.Contains(base, "centos"), strings.Contains(base, "fedora"), strings.Contains(base, "redhat"), strings.Contains(base, "rhel"):
      return returnMap["centos"], nil
    case strings.Contains(base, "arch"):
      return returnMap["arch"], nil
    case strings.Contains(base, "opensuse"):
      return returnMap["opensuse"], nil
    default:
      return "", fmt.Errorf("ERROR: osName not recognized")
  }
}

func SimplifyOSName(osName string) string {
  base := strings.ToLower(osName)

  switch {
    case strings.Contains(base, "ubuntu"), strings.Contains(base, "debian"):
      return "ubuntu"
    case strings.Contains(base, "alpine"):
      return "alpine"
    case strings.Contains(base, "centos"), strings.Contains(base, "fedora"), strings.Contains(base, "redhat"), strings.Contains(base, "rhel"):
      return "centos"
    case strings.Contains(base, "arch"):
      return "arch"
    case strings.Contains(base, "opensuse"):
      return "opensuse"
    default:
      return ""
  }
}
