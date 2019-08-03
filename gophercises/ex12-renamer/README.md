# Exercise #12: File Renaming Tool

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/renamer)

## Exercise details

In this exercise we are going to explore ways to navigate a file system by creating an application that will rename a bunch of user files in nested directories. The exact files you rename are up to you, but I have provided a sample directory in case you need some ideas. It has the files and directories shown below.

```
sample/
  birthday_001.txt
  birthday_002.txt
  birthday_003.txt
  birthday_004.txt
  christmas 2016 (1 of 100).txt
  christmas 2016 (2 of 100).txt
  christmas 2016 (3 of 100).txt
  christmas 2016 (4 of 100).txt
  christmas 2016 (5 of 100).txt
  nested/
    n_008.txt
    n_009.txt
    n_010.txt
```

The goal of our program is to rename a specific subset of these files. For instance, we might want to take all the files that end in `_NNN.txt` and rename them to instead read `(1 of 4).txt`. Or maybe we will want to rename all of the `XXXXXX_NNN.txt` files to instead read `NNN - XXXXXX.txt`. The exact naming pattern isn't really important, but what IS important is that you can write a program that only modifies the files you specifically want modified, and is able to rename them.

It will be very tempting to try to write a generic tool that works for any naming pattern, but I recommend not doing this at first. Instead, focus on hard coding a few naming patterns into your code and getting an understanding of the requirements and intricacies of the task. Once you understand them, feel free to try to come up with a more generic solution.

## Hints and recommended libraries

You will very likely want to use the [path/filepath](https://golang.org/pkg/path/filepath) package to solve this exercise. Specifically, check out `Walk` and `WalkFunc`.

You will also need to find a way to rename files. There are a few ways to do this, such as using the `os/exec` package, but I would probably start by looking at the [os.Rename](https://golang.org/pkg/os/#Rename) function.

Lastly, if you try to make a more generic solution you will probably want to use regular expressions and then [regexp](https://golang.org/pkg/regexp/) package. 

**Full disclosure:** *I probably won't be creating a generic solution because I find that tasks like this require just enough customization between each use case that it is often easier to modify the existing source code than it is to try to create something generic that satisfies all those needs, but you are welcome to give it a go.*

## Bonus

Verify that your code works recursively, and if you want try to make a more generic program that can be used for a few different filename patterns.