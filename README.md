Introduction
================

fod (meaning File Open Dialog) is a tool to select the directories and files interactively. 
will output the path of the selected directories and files to standard output.

__Example__

    cd `fod`
    mv `fod`
    git add `fod`
    git add `fod -m` # comming soon


Options and arguments
================

| **option**   | **Explaination**                        |
| ------------ | --------------------------------------- |
| -f           | file select mode                        |
| -d (default) | directory select mode                   |
| -m           | multiple selection mode (comming soon)  |

Install
================

    go get github.com/ykanda/fod


Key Bindings
================

| **KEY**      | **Explaination**                 |
| ------------ | -------------------------------- |
| Enter        | open directory                   |
| Arrow Up     | move cursor down                 |
| Arrow Down   | move cursor down                 |
| Arrow Left   | move parent directory            |
| Arrow Right  | move sub directory               |
| Ctrl + H     | toggle dotfile filter            |
| Ctrl + O     | OK, exit and output selcted item |
| Ctrl + C     | cancel and exit, no output       |
| Ctrl + Q     | cancel and exit, no output       |
| Esc          | cancel and exit, no output       |


Development, Contributions
================

1. Please fork repository on GitHub.
2. Execute `go get github.com/your_id_on_github/fod`, to make working copy on your computer.
3. Editing code.
4. Send PR.


Dependencies
================

* [github.com/k0kubun/pp](https://github.com/k0kubun/pp)
* [github.com/nsf/termbox-go](https://github.com/nsf/termbox-go)


LICENSE
================

[The MIT License](http://opensource.org/licenses/mit-license.php)

Copyright (C) 2015 Yasuhiro KANDA ([@kandayasu](https://twitter.com/kandayasu))


TODO
================

* test
* configuable key bind
* symlinks
* multiple selection
* bookmark
* windows support
