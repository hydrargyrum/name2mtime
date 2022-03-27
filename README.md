# name2mtime

name2mtime is a command-line tool to change the last-modification time of files basing on their filenames.
For example, if you have a file `IMG_20220327_123456.jpg`, it will change it last-modification time (mtime) to `2022-03-27 12:34:56`.

Usage: 

	name2mtime FILE1 [FILE2…]

This requires Go but no third-party modules.
