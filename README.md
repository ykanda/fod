Introduction
================

fod (meaning File Open Dialog) is a tool to select the directories and files interactively. 
will output the path of the selected directories and files to standard output.

__Usage__

		cd $(fod)
		mv $(fod)
		git add $(fod)
		git add $(fod -m --separator=' ')


Exmaple: Directory select mode
--------

For example, If you run `fod` on this project repository working copy.
If you are not given the option, it will be a mode in which you can select one of the directory.
It will see the display as follows:

		>
		[d]   ../
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/.git
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/_vendor
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/pkg-config-files

Items that focus is highlighted. 
You can change the items that you have focus in the up and down arrow keys.


Mark Item
--------

When you press Ctrl + S, it will be marked by the selected item.
To the marked item is displayed '*'.

		>
		[d]   ../
		[d] * /Users/kandayasu/.go/src/github.com/ykanda/fod/.git
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/_vendor
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/pkg-config-files

If you mark the other items, the mark of the current item is excluded.
If you want to select multiple items at the same time, use the `--multiple` option.


Exit and output selected item to STDIO
--------

When you press Ctrl + O, and then exit the selection.
The marked items are displayed in the standard output.


Change directory
--------

If you see with the left side of the list of items '[d]', that item is a directory.
In the case of the file is displayed as '[f]'.

When you press Enter in a state of focus item is a directory, you can change the directory
For example, if you press Enter in a state of focus the .git directory, it will be displayed as follows:

		>
		[d]   ../
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/branches
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/hooks  
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/info            
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/logs
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/objects
		[d]   /Users/kandayasu/.go/src/github.com/ykanda/fod/refs


Name filter
--------

__TODO__


Options and arguments
================

| **option**      | **Explaination**                         |
| --------------- | ---------------------------------------- |
| --mode, -m      | select mode (f, file, directory, dir, d) |
| --base, -b      | base dir                                 |
| --multi         | multiple selection mode                  |
| --separator, -s | path separator string (use with --multi) |


Install
================

```
go get github.com/ykanda/fod
```


Key Bindings
================

| **KEY**      | **Explaination**                 |
| ------------ | -------------------------------- |
| Enter        | open directory                   |
| Arrow Up     | move cursor down                 |
| Arrow Down   | move cursor down                 |
| Arrow Left   | move parent directory            |
| Arrow Right  | move sub directory               |
| Ctrl + S     | toggle marked / unmarked         |
| Ctrl + H     | toggle dotfile filter            |
| Ctrl + O     | OK, exit and output selcted item |
| Ctrl + C     | cancel and exit, no output       |
| Ctrl + Q     | cancel and exit, no output       |
| Esc          | cancel and exit, no output       |


Development, Contributions
================

1. Please fork repository on GitHub.
2. Execute `go get github.com/your_id_on_github/fod`, to make working copy on your computer.
3. Checkout working branch.
4. Editing code.
5. Send PR.


Dependencies
================

Too see glide.yaml file.
I thank the authors of the library.


LICENSE
================

[The MIT License](http://opensource.org/licenses/mit-license.php)

Copyright (C) 2015 Yasuhiro KANDA ([@kandayasu](https://twitter.com/kandayasu))


TODO
================

* test
* configuable key bind
* symlinks
* windows support
