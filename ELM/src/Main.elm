module Main exposing (main)

import Browser
import Html exposing (Html, div, text, button, input)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import Svg exposing (Svg, svg)
import Svg.Attributes exposing (..)
import TCParser exposing (Instruction(..), parseTurtleProgram)
import Parser exposing (run)
import Drawing exposing (drawInstructions)
import Debug exposing (toString)


-- MODEL

type alias Model =
    { userInput : String
    , parsedInstructions : Result String (List Instruction)
    }


init : Model
init =
    { userInput = ""
    , parsedInstructions = Err "Enter a program to parse and draw!" }


-- MESSAGES

type Msg
    = UpdateInput String
    | ParseInput


-- UPDATE

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput input ->
            { model | userInput = input }

        ParseInput ->
            let
                result =
                    run parseTurtleProgram model.userInput
            in
            { model | parsedInstructions = Result.mapError toString result }


-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ input
            [ placeholder "Enter your TcTurtle program here"
            , value model.userInput
            , onInput UpdateInput
            ]
            []
        , button [ onClick ParseInput ] [ text "Parse and Draw" ]
        , div [] [ viewOutput model.parsedInstructions ]
        ]


viewOutput : Result String (List Instruction) -> Html msg
viewOutput result =
    case result of
        Ok instructions ->
            svg [ width "500", height "500", viewBox "0 0 500 500" ]
                [ drawInstructions instructions ]

        Err error ->
            text ("Error: " ++ error)


-- MAIN

main : Program () Model Msg
main =
    Browser.sandbox { init = init, update = update, view = view }

