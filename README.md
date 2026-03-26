# MEMORIZER - Flashcard learning program

## Overview

Memorizer is a flashcard learning tool for the command line, with a SRS (spaced repetition system). It doesn't implement the full SM-2 Algorithm, but an Anki-lite version.

## Features

- Add a learning deck (which will be saved as deck-name.json) - each deck gets a separate JSON file
  - decks can also be provided through environment variable "MEMORIZER_FILE" like follows: 
  - `export MEMORIZER_FILE=name.json`
- Add cards for a deck (either via flags or via interactive prompt input)
- List all the cards of a deck
- Edit the cards (e.g. edit the term or the solution or both descriptions)
- Delete cards
- Run the Review process (based on a Anki-lite SRS)

## Requirements

The [Go](https://go.dev/) toolchain is needed to build and run the program.

My favorite way to install Go: [Webi installer](https://webinstall.dev/webi/)

## Installation

To set up and run this project, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/Alex-Koch-85/memorizer
   ```

   

2. Navigate to the project directory:

   ```bash
   cd memorizer
   ```

   

3. Build the go project

   ```bash
   go build
   ```

   

## Usage:

Usage information:

```bash
./flashcard -h
```

Add a card through flags:

```bash
./flashcard -deck "Go basics" -addCard -term "slice" -solution "dynamic array"
```

Create a "default" deck and add card:

```bash
./flashcard -addCard -term "..." -solution "..." 
```

Interactive prompt input:

```bash
./flashcard -deck "Go basics" -addCard
```

List all cards in a deck:

```bash
./flashcard -deck "Go basics" -list
```

ID can be copied from -list listing; edit term and or solution through interactive prompt input:

```bash
./flashcard -deck "Go basics" -edit -id <UUID>
```

Term "array slice"` - edit term through flag (solution stays the same):

```bash
./flashcard -deck "Go basics" -edit -id <UUID>
```

Solution "dynamically sized, flexible view into the elements of an array"` - edit solution through flag (term stays the same):

```bash
./flashcard -deck "Go basics" -edit -id <UUID>
```

Delete card with provided id:

```bash
./flashcard -deck "Go basics" -delete -id <UUID>
```

Start review process for the deck:

```bash
./flashcard -deck "Go basics" -review
```



Thank you!