# Mashery V3 Connection Command-Line Utility

Mashery client (and derivative packages) need Mashery V3 access token in order to perform operations. This
utility is used to request, view, and refresh access tokens.

# Synopsis

`mash-connect [init|show|export|refresh] sub-command options`

# `init` command

The command obtains the initial access token and saves it to the file for further reuse.

# `show` command

The `show` command is used to display how much time is left in the given access token. If a token has already
expired, a clear prompt will be printed

# `export` command

The `export` sub-command is used to export the access token into a environment variable. 