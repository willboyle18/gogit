# GoGit
GoGit is a replica of the Git version control system programmed entirely in Go. It implements many of Git's fundamental features and internal architecture.

This project is designed to help developers understand exactly how Git works under the hood, from hashing and storing objects to restoring project state with checkout.

## Features
### Object Database (.gogit/objects)
The Object Database is where Git stores all of it's binary data for your repository on disk. The objects are compressed and organized by an associated SHA-1 hash.
The structure of the object database is (first two hexadecimal digits of SHA-1 hash - this is the directory)/(last 18 digits of SHA-1 hash - file that stores the binary data of the object)
  1. **Blob:** A blob is simply a binary blob of data that stores the contents of individual files. These files are compressed using zlib and stored in the database according to their SHA-1 hash.
  2. **Tree**: A tree is a binary object that stores the data of a commit. It contains the file mode, name, and the SHA-1 of the blob where the data is stored.
  3. **Commit**: A commit is a binary object that stores the commit data. It contains the SHA-1 of it's associated tree, the SHA-1 of the parent (previous) commit object, and the author data.
  4. **HEAD (.gogit/HEAD)**: The HEAD is a file that contains the SHA-1 hash of the last commit. This allows us to traverse down all of our previous commits.
### Index (.gogit/index)
The index is where our staging area exists in GoGit. it contains all of the file data of that project at that point in time. Every time we add files to be staged, this index gets updated before we insert them into the object database.
  1. **Active Cache**: This is our way of storing files into memory to prep them to be stored into the index. Each entry on the Active Cache stores a struct that contains all the file metadata as well as it's name and SHA-1 hash.

## Commands
1. **init**: The 'init' command is the first command you run when using GoGit. It's job is to initialize the object database with all the SHA-1 directories ready to store binary files.
2. **add**: The 'add' command places files from the current state of the project, turns them into blobs, and places them into the index. To use it you must do ```add (files to stage)...``` and they will overwrite their current state in the index.
3. **commit**: The 'commit' command takes the index and creates a tree object that stores the current state of the index. It then creates a commit object and updates the HEAD as the latest commit. To use do ```commit -m "commit message".
4. **log**: Displays list of all commits.
5. **latest**: Displays the latest commit.
6. **checkout**: Loads a specific commit and all of its file data. This commands lets you go back to how your project looked at a previous time. To use do ```checkout (tree SHA-1 of the commit you want to go back to)``` and your project will load all of the commit data from that point in time.
