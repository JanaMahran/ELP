module TCParser exposing (Instruction(..), parseTurtleProgram)

import Html exposing (Html, div, text)
import Parser exposing (Parser, succeed, symbol, spaces, int, oneOf, (|=), (|.), loop, Step(..), lazy, run)


-- Types to represent instructions
type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)


-- Main parser for a Turtle program
parseTurtleProgram : Parser (List Instruction)
parseTurtleProgram =
    parseBracketed (lazy (\() -> loop [] parseRepeatingInstructions))


-- Helper to parse content wrapped in brackets
parseBracketed : Parser a -> Parser a
parseBracketed innerParser =
    succeed identity
        |. symbol "["
        |. spaces
        |= innerParser
        |. spaces
        |. symbol "]"


-- Parser for repeating a list of instructions
parseRepeatingInstructions : List Instruction -> Parser (Step (List Instruction) (List Instruction))
parseRepeatingInstructions accumulated =
    oneOf
        [ succeed (\instr -> Loop (accumulated ++ [instr]))
            |= parseInstruction
            |. parseOptionalSeparator
        , succeed (Done accumulated)
        ]


-- Handles optional separators (comma or just spaces)
parseOptionalSeparator : Parser ()
parseOptionalSeparator =
    oneOf
        [ symbol "," |. spaces
        , spaces
        ]


-- Parses an individual instruction
parseInstruction : Parser Instruction
parseInstruction =
    oneOf
        [ parseForward
        , parseLeft
        , parseRight
        , parseRepeat
        ]


-- Parsers for specific instructions
parseForward : Parser Instruction
parseForward =
    succeed Forward
        |. parseKeyword "Forward"
        |= int
        |. spaces


parseLeft : Parser Instruction
parseLeft =
    succeed Left
        |. parseKeyword "Left"
        |= int
        |. spaces


parseRight : Parser Instruction
parseRight =
    succeed Right
        |. parseKeyword "Right"
        |= int
        |. spaces


parseRepeat : Parser Instruction
parseRepeat =
    succeed Repeat
        |. parseKeyword "Repeat"
        |= int
        |. spaces
        |= parseBracketed (lazy (\() -> loop [] parseRepeatingInstructions))
        |. spaces


-- Helper for keywords
parseKeyword : String -> Parser ()
parseKeyword word =
    symbol word |. spaces



