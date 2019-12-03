// Do not edit. This file is generated by running 'make man'.
package man

const Section1 = `.\"t
.\" Automatically generated by Pandoc 2.7.3
.\"
.TH "CITOP" "1" "December 03, 2019" "" "version 0.0.0-79-g9440cfc-dirty"
.hy
.SH NAME
.PP
\f[B]citop\f[R] \[en] Continuous Integration Table Of Pipelines
.SH SYNOPSIS
.PP
\f[C]citop [-r REPOSITORY | --repository REPOSITORY] [COMMIT]\f[R]
.PP
\f[C]citop -h | --help\f[R]
.PP
\f[C]citop --version\f[R]
.SH DESCRIPTION
.PP
citop monitors the CI pipelines associated to a specific commit of a git
repository.
.PP
citop currently integrates with the following online services.
Each of the service is one or both of the following:
.IP \[bu] 2
A \[lq]source provider\[rq] that is used to list the pipelines
associated to a given commit of an online repository
.IP \[bu] 2
A \[lq]CI provider\[rq] that is used to get detailed information about
CI builds
.PP
.TS
tab(@);
lw(11.7n) lw(8.8n) lw(7.8n) lw(27.2n).
T{
Service
T}@T{
Source
T}@T{
CI
T}@T{
URL
T}
_
T{
GitHub
T}@T{
yes
T}@T{
no
T}@T{
<https://github.com/>
T}
T{
GitLab
T}@T{
yes
T}@T{
yes
T}@T{
<https://gitlab.com/>
T}
T{
AppVeyor
T}@T{
no
T}@T{
yes
T}@T{
<https://www.appveyor.com/>
T}
T{
CircleCI
T}@T{
no
T}@T{
yes
T}@T{
<https://circleci.com/>
T}
T{
Travis CI
T}@T{
no
T}@T{
yes
T}@T{
<https://travis-ci.org/> <https://travis-ci.com/>
T}
.TE
.SH POSITIONAL ARGUMENTS
.SS \f[C]COMMIT\f[R]
.PP
Specify the commit to monitor.
COMMIT is expected to be the SHA identifier of a commit, or the name of
a tag or a branch.
If this option is missing citop will monitor the commit referenced by
HEAD.
.PP
Example:
.IP
.nf
\f[C]
# Show pipelines for commit 64be3c6
citop 64be3c6
# Show pipelines for the commit referenced by the tag \[aq]0.9.0\[aq]
citop 0.9.0
# Show pipelines for the commit at the tip of a branch
citop feature/doc
\f[R]
.fi
.SH OPTIONS
.SS \f[C]-r=REPOSITORY, --repository=REPOSITORY\f[R]
.PP
Specify the git repository to work with.
REPOSITORY can be either a path to a local git repository, or the URL of
an online repository hosted at GitHub or GitLab.
Both web URLs and git URLs are accepted.
.PP
In the absence of this option, citop will work with the git repository
located in the current directory.
If there is no such repository, citop will fail.
.PP
Examples:
.IP
.nf
\f[C]
# Work with the git repository in the current directory
citop
# Work with the repository specified by a web URL
citop -r https://gitlab.com/nbedos/citop
citop -r github.com/nbedos/citop
# Git URLs are accepted
citop -r git\[at]github.com:nbedos/citop.git
# Paths to a local repository are accepted too
citop -r /home/user/repos/myrepo
\f[R]
.fi
.SS \f[C]-h, --help\f[R]
.PP
Show usage of citop
.SS \f[C]--version\f[R]
.PP
Print the version of citop being run
.SH INTERACTIVE COMMANDS
.PP
Below are the default commands for interacting with citop.
.PP
.TS
tab(@);
lw(10.7n) lw(46.7n).
T{
Key
T}@T{
Action
T}
_
T{
Up, j
T}@T{
Move cursor up by one line
T}
T{
Down, k
T}@T{
Move cursor down by one line
T}
T{
Page Up
T}@T{
Move cursor up by one screen
T}
T{
Page Down
T}@T{
Move cursor down by one screen
T}
T{
o, +
T}@T{
Open the fold at the cursor
T}
T{
O
T}@T{
Open the fold at the cursor and all sub-folds
T}
T{
c, -
T}@T{
Close the fold at the cursor
T}
T{
C
T}@T{
Close the fold at the cursor and all sub-folds
T}
T{
/
T}@T{
Open search prompt
T}
T{
Escape
T}@T{
Close search prompt
T}
T{
Enter, n
T}@T{
Move to the next match
T}
T{
N
T}@T{
Move to the previous match
T}
T{
v
T}@T{
View the log of the job at the cursor[a]
T}
T{
b
T}@T{
Open with default web browser
T}
T{
q
T}@T{
Quit
T}
T{
?
T}@T{
View manual page
T}
.TE
.IP \[bu] 2
[a] Note that if the job is still running, the log may be incomplete.
.SH CONFIGURATION FILE
.SS Location
.PP
citop follows the XDG base directory specification [2] and expects to
find the configuration file at one of the following locations depending
on the value of the two environment variables \f[C]XDG_CONFIG_HOME\f[R]
and \f[C]XDG_CONFIG_DIRS\f[R]:
.IP "1." 3
\f[C]\[dq]$XDG_CONFIG_HOME/citop/citop.toml\[dq]\f[R]
.IP "2." 3
\f[C]\[dq]$DIR/citop/citop.toml\[dq]\f[R] for every directory
\f[C]DIR\f[R] in the comma-separated list
\f[C]\[dq]$XDG_CONFIG_DIRS\[dq]\f[R]
.PP
If \f[C]XDG_CONFIG_HOME\f[R] (resp.
\f[C]XDG_CONFIG_DIRS\f[R]) is not set, citop uses the default value
\f[C]\[dq]$HOME/.config\[dq]\f[R] (resp.
\f[C]\[dq]/etc/xdg\[dq]\f[R]) instead.
.SS Format
.PP
citop uses a configuration file in TOML version
v0.5.0 (https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.5.0.md)
format.
The configuration file is made of keys grouped together in tables.
The specification of each table is given below.
.SS Table \f[C][providers]\f[R]
.PP
The ` + "`" + `providers' table is used to define credentials for accessing online
services.
citop relies on two types of providers:
.IP \[bu] 2
` + "`" + `source providers' are used for listing the CI pipelines associated to a
given commit (GitHub and GitLab are source providers)
.IP \[bu] 2
` + "`" + `CI providers' are used to get detailed information about CI pipelines
(GitLab, AppVeyor, CircleCI and Travis are CI providers)
.PP
citop requires credentials for at least one source provider and one CI
provider to run.
.SS Table \f[C][[providers.gitlab]]\f[R]
.PP
\f[C][[providers.gitlab]]\f[R] defines a GitLab account
.PP
.TS
tab(@);
lw(8.8n) lw(48.6n).
T{
Key
T}@T{
Description
T}
_
T{
name
T}@T{
Name under which this provider appears in the TUI (string, optional,
default: \[lq]gitlab\[rq])
T}
T{
url
T}@T{
URL of the GitLab instance (string, optional, default:
\[lq]gitlab.com\[rq])
T}
T{
token
T}@T{
Personal access token for the GitLab API (string, optional, default:
\[dq]\[dq])
T}
.TE
.PP
GitLab access tokens are managed at
<https://gitlab.com/profile/personal_access_tokens>
.PP
Example:
.IP
.nf
\f[C]
[[providers.gitlab]]
name = \[dq]gitlab.com\[dq]
url = \[dq]https://gitlab.com\[dq]
token = \[dq]gitlab_api_token\[dq]
\f[R]
.fi
.SS Table \f[C][[providers.github]]\f[R]
.PP
\f[C][[providers.github]]\f[R] defines a GitHub account
.PP
.TS
tab(@);
lw(7.8n) lw(50.6n).
T{
Key
T}@T{
Description
T}
_
T{
token
T}@T{
Personal access token for the GitHub API (string, optional, default:
\[dq]\[dq])
T}
.TE
.PP
GitHub access tokens are managed at <https://github.com/settings/tokens>
.PP
Example:
.IP
.nf
\f[C]
[[providers.github]]
token = \[dq]github_api_token\[dq]
\f[R]
.fi
.SS Table \f[C][[providers.travis]]\f[R]
.PP
\f[C][[providers.travis]]\f[R] defines a Travis CI account
.PP
.TS
tab(@);
lw(7.8n) lw(50.6n).
T{
Key
T}@T{
Description
T}
_
T{
name
T}@T{
Name under which this provider appears in the TUI (string, mandatory)
T}
T{
url
T}@T{
URL of the GitLab instance.
\[lq]org\[rq] and \[lq]com\[rq] can be used as shorthands for the full
URL of travis.org and travis.com (string, mandatory)
T}
T{
token
T}@T{
Personal access token for the Travis API (string, optional, default:
\[dq]\[dq])
T}
.TE
.PP
Travis access tokens are managed at the following locations:
.IP \[bu] 2
<https://travis-ci.org/account/preferences>
.IP \[bu] 2
<https://travis-ci.com/account/preferences>
.PP
Example:
.IP
.nf
\f[C]
[[providers.travis]]
name = \[dq]travis.org\[dq]
url = \[dq]org\[dq]
token = \[dq]travis_org_api_token\[dq]

[[providers.travis]]
name = \[dq]travis.com\[dq]
url = \[dq]com\[dq]
token = \[dq]travis_com_api_token\[dq]
\f[R]
.fi
.SS Table \f[C][[providers.appveyor]]\f[R]
.PP
\f[C][[providers.appveyor]]\f[R] defines an AppVeyor account
.PP
.TS
tab(@);
lw(7.8n) lw(49.6n).
T{
Key
T}@T{
Description
T}
_
T{
name
T}@T{
Name under which this provider appears in the TUI (string, optional,
default: \[lq]appveyor\[rq])
T}
T{
token
T}@T{
Personal access token for the AppVeyor API (string, optional, default:
\[dq]\[dq])
T}
.TE
.PP
AppVeyor access tokens are managed at <https://ci.appveyor.com/api-keys>
.PP
Example:
.IP
.nf
\f[C]
[[providers.appveyor]]
name = \[dq]appveyor\[dq]
token = \[dq]appveyor_api_key\[dq]
\f[R]
.fi
.SS Table \f[C][[providers.circleci]]\f[R]
.PP
\f[C][[providers.circleci]]\f[R] defines a CircleCI account
.PP
.TS
tab(@);
lw(7.8n) lw(49.6n).
T{
Key
T}@T{
Description
T}
_
T{
name
T}@T{
Name under which this provider appears in the TUI (string, optional,
default: \[lq]circleci\[rq])
T}
T{
token
T}@T{
Personal access token for the CircleCI API (string, optional, default:
\[dq]\[dq])
T}
.TE
.PP
CircleCI access tokens are managed at <https://circleci.com/account/api>
.PP
Example:
.IP
.nf
\f[C]
[[providers.circleci]]
name = \[dq]circleci\[dq]
token = \[dq]circleci_api_token\[dq]
\f[R]
.fi
.SS Examples
.PP
Here are a few examples of \f[C]citop.toml\f[R] configuration files.
.PP
Monitor pipelines on Travis CI, AppVeyor and CircleCI for a repository
hosted on GitHub:
.IP
.nf
\f[C]
[[providers.github]]
token = \[dq]github_api_token\[dq]

[[providers.travis]]
url = \[dq]org\[dq]
token = \[dq]travis_org_api_token\[dq]

[[providers.appveyor]]
token = \[dq]appveyor_api_key\[dq]

[[providers.circleci]]
token = \[dq]circleci_api_token\[dq]
\f[R]
.fi
.PP
Monitor pipelines on GitLab CI for a repository hosted on GitLab itself:
.IP
.nf
\f[C]
[[providers.gitlab]]
token = \[dq]gitlab_api_token\[dq]
\f[R]
.fi
.SH ENVIRONMENT
.SS ENVIRONMENT VARIABLES
.IP \[bu] 2
\f[C]BROWSER\f[R] is used to find the path of the default web browser
.IP \[bu] 2
\f[C]HOME\f[R], \f[C]XDG_CONFIG_HOME\f[R] and \f[C]XDG_CONFIG_DIRS\f[R]
are used to locate the configuration file
.SS LOCAL PROGRAMS
.PP
citop relies on the following local executables:
.IP \[bu] 2
\f[C]git\f[R] to translate the abbreviated SHA identifier of a commit
into a non-abbreviated SHA
.IP \[bu] 2
\f[C]less\f[R] to show job logs
.IP \[bu] 2
\f[C]man\f[R] to show the manual page
.SH EXAMPLES
.PP
Show pipelines associated to the HEAD of the current git repository
.IP
.nf
\f[C]
citop
\f[R]
.fi
.PP
Show pipelines associated to a specific commit, tag or branch
.IP
.nf
\f[C]
citop 64be3c6
citop 0.9.0
citop feature/doc
\f[R]
.fi
.PP
Show pipelines of a repository specified by a URL
.IP
.nf
\f[C]
citop -r https://gitlab.com/nbedos/citop
citop -r git\[at]github.com:nbedos/citop.git
citop -r github.com/nbedos/citop
\f[R]
.fi
.PP
Show pipelines of a local repository specified by a path
.IP
.nf
\f[C]
citop -r /home/user/repos/myrepo
\f[R]
.fi
.PP
Specify both repository and commit
.IP
.nf
\f[C]
citop -r github.com/nbedos/citop 64be3c6
\f[R]
.fi
.SH NOTES
.IP "1." 3
\f[B]citop repository\f[R]
.RS 4
.IP \[bu] 2
<https://github.com/nbedos/citop>
.RE
.IP "2." 3
\f[B]XDG base directory specification\f[R]
.RS 4
.IP \[bu] 2
<https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html>
.RE
.SH AUTHORS
Nicolas Bedos.`
