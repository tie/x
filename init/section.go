package main

type Line []Token

type Section struct {
	Header Line
	Commands []Line
}
