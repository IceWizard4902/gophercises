# Gophercises

My attempts in solving (or mostly following the solution of `gophercises`). Inspired by [John Hammond](https://www.youtube.com/watch?v=hs2acc8AibU). 

The exercises are available on the [`gophercises` website](https://gophercises.com/)

# Some notes

- If you are using VSCode, chances are the static analyzer will scream at you if you do not initialise the `go.mod` file for each folder and the `go.work` folder for the subpackages. 
- For each folder, say `exercise-1-quiz`, to get rid of the nasty underlined red parts, do the following: 

```bash
cd exercise-1-quiz
go mod init <some package in the source code>
cd .. 
go work init
go work use ./exercise-1-quiz
```

- When you add a new folder, just omit the `go work init` part and run the same commands.