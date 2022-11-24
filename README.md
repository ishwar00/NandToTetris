# NandToTetris

> What I cannot create, I do not understand — Richard Feynman.

## Table of contents
* [General info](#general-info)
    * [What is it?](#what-is-it)
    * [How i am doing it?](#how-i-am-doing-it)
* [Technologies](#technologies)
* [Project structure](#project-structure)
* [Status](#status)


# General info

## What is it?
This repository is me documenting, things i will be learning and building while reading [The elements computing systems](https://mitpress.mit.edu/books/elements-computing-systems).This is an attempt to make a modern comptuter out of nand gates.

quoting from book:
> *Nand to Tetris:* a hands-on journey that starts with the most elementary logic gate, called Nand, and ends up, twelve projects later, with a general  purpose computer system capable of running Tetris, as well as any other
program that comes to your mind. 
[know more](https://www.nand2tetris.org/).


## How i am doing it?
They are two ways to go about,
 1. Take online courses: [NandToTetris Part1](https://www.coursera.org/learn/build-a-computer), and [NandToTetris Part2](https://www.coursera.org/learn/nand2tetris2)
 2. Reading the **book** [The elements of computing systems: Building a Modern Computer from First Principles](https://mitpress.mit.edu/books/elements-computing-systems)

 I prefer reading, so i went with the **book** :), and sometimes [Computer Systems Design - Kamakoti | IIT Madras ](https://www.youtube.com/playlist?list=PLEAYkSg4uSQ0eDa24iKd7qJlsrvr8XcvF)

[back to table of contents?](#table-of-contents)

## Technologies
 * [The Nand to Tetris Software Suite](https://www.nand2tetris.org/software)
 * [Go 1.18.2 or later](https://go.dev/)

## Project structure
Namespace Tree of the project at depth level 3.

    +---nand2tetris
    |   +---projects
    |   |   +---01 Boolean Logic
    |   |   +---02 Boolean Arithmetic
    |   |   +---03 Memory
    |   |   +---04 Machine Language
    |   |   +---05 Computer Architecture
    |   |   +---06 Assembler
    |   |   +---07 Virtual Machine 1 (processing)
    |   |   +---08 Virtual Machine 2 (control)
    |   |   +---09 High-Level Language
    |   |   +---10 Compiler I: Syntax Analysis
    |   |   +---11
    |   |   +---12
    |   |   +---13
    |   |   +---demo
    |   +---tools
    |   |   +---bin
    |   |   +---builtInChips
    |   |   +---builtInVMCode
    |   |   +---OS
    +---utils

directory `nand2tetris/projects` has 12 projects(directories), each project is part of a chapter in the book.\
directory `nand2tetris/tools` has *The Nand to Tetris Software Suite* which will be used to simulate the computer.

[back to table of contents?](#table-of-contents)

## Status
#### HARDWARE PART
✔️  Boolean Logic\
✔️  Boolean Arithmetic\
✔️  Memory\
✔️  Machine Language\
✔️  Computer Architecture\
✔️  Assembler

Completed projects 1–6, and built a general-purpose
computer system from first principles 🙌, the computer is capable of executing only programs written in
machine language  

#### SOFTWARE PART

✔️  Virtual Machine I: Processing\
✔️  Virtual Machine II: Control\
✔️  High-Level Language\
✔️  Compiler I: Syntax Analysis([did not clear all tests ](https://github.com/ishwar00/NandToTetris/blob/main/nand2tetris/projects/10%20Compiler%20I:Syntax%20Analysis/JackAnalyzer/src/parserXML/readme.md))\
❌ Compiler II: Code generation\
❌ Operating System\
❌ More Fun to Go

[back to table of contents?](#table-of-contents)
