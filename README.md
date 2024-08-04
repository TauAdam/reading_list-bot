# Reading list bot 

This is a simple bot that saves links to articles and books that you want to read later.

## Features
- Save links to articles and books
- Pick random link from the list
- Start reading from the list
- Help command

### Roadmap
- [x] Handle main commands: save, start, help, pick random
- [x] Save links to articles and books
- [x] Pick random link from the list
- [ ] Finish sqlite storage

## What package includes
- Built in client for Telegram API
- Local disk storage
- Sqlite storage

## Engineering practises
- clean code
- abstract interfaces
- not depends on realization of messenger (Telegram in that case)
- using context for canceling long operations