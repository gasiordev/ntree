# ntree

ntree is a tiny program that displays directories and files of a given path and
expands the subdirectories that are your current working directory.

Check the screenshot for the following command:

```
ntree start -r /Users/miko -w /Users/miko/Repos/gasiordev/ntree
```

![ntree screenshot](ntree.png)

It was created to be used with tmux (or similar). You run ntree in a separate
pane and it gets refreshed while you work on the terminal in other panes.
For example, you can alias `cd` command in bash profile to automatically update
ntree with the following code:

```
cd() {
  ddd="$(pwd)/$1"
  if [ -d "$ddd" ]; then
    ntree send WORKDIR "$ddd" > /dev/null 2>&1;
  fi
  builtin cd $1;
}
```

So, as you can see, working directory can be changed while the program is 
running. This is done with `ntree send`, eg. 
`ntree send WORKDIR /Users/miko/Repos/gasiordev/`.

On the screenshot next to filenames, you can spot numbers. Number on the left,
like `(3)...`, indicates that there are 3 items before which are not shown,
and number on the right like `...(5)` tells us that there are 5 items after
which are not shown. Items can be hidden because they do not fit on the
terminal. Also, the priority is to show contents of working directory. So if
necessary, any directories along the path are partly hidden. 

## Configuration
Before you start ntree, copy `sampleconfig.json` to your home directory as
`.ntree.json`. You can change Unix socket path or refresh time (called 
`loop_sleep`) in the file.


## Commands
In the above example, we mentioned that with `ntree send` you order the program
to change what's visible.

The available commands are:
* `ROOTDIR <value>` - change root directory;
* `WORKDIR <value>` - change working directory (the one that is expanded);
* `DIRS` - toggle directories visibility (you can hide directories and just show the files);
* `FILES` - toggle files visibility;
* `HIDDEN` - toggle hidden files visibility (by default they are hidden);
* `FILTER <value>` - only show files and directories that contain `<value>`;
* `HIGHLIGHT <value>` - highlight (in a colour) `<value>` in file names;
* `RESET-FILTER` - reset the filter;
* `RESET-HIGHLIGHT` - reset the highlight;
* `FREEZE` - toggle freeze (tree can be frozen so it doesn't change).


