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
    div [ Html.Attributes.style "text-align" "center" ]
        [ Html.h1 [ Html.Attributes.style "font-size" "2em", Html.Attributes.style "color" "#333" ] [ text "Welcome to TcTurtle!" ]
        , Html.p [ Html.Attributes.style "font-size" "1.2em", Html.Attributes.style "color" "#666", Html.Attributes.style "margin-bottom" "20px" ] [ text "Start by writing the following instruction! xx" ]
        , div [ Html.Attributes.style "margin" "auto", Html.Attributes.style "width" "80%", Html.Attributes.style "display" "flex", Html.Attributes.style "justify-content" "space-between" ]
            [ div [ Html.Attributes.style "width" "65%", Html.Attributes.style "text-align" "left" ]
                [ input
                    [ placeholder "Enter your TcTurtle program here"
                    , value model.userInput
                    , onInput UpdateInput
                    , Html.Attributes.style "width" "100%"
                    ]
                    []
                , button
                    [ onClick ParseInput
                    , Html.Attributes.style "background-color" "grey"
                    , Html.Attributes.style "color" "white"
                    , Html.Attributes.style "font-size" "1em"
                    , Html.Attributes.style "padding" "10px 20px"
                    , Html.Attributes.style "border" "none"
                    , Html.Attributes.style "border-radius" "5px"
                    , Html.Attributes.style "width" "100%"
                    , Html.Attributes.style "margin-top" "10px"
                    ]
                    [ text "Parse and Draw" ]
                , div [ Html.Attributes.style "margin-top" "20px" ]
                    [ viewOutput model.parsedInstructions ]
                ]
            , div [ Html.Attributes.style "width" "30%", Html.Attributes.style "background-color" "grey", Html.Attributes.style "padding" "10px", Html.Attributes.style "border-radius" "5px", Html.Attributes.style "color" "white" ]
                [ Html.h2 [ Html.Attributes.style "font-size" "1.5em", Html.Attributes.style "color" "white" ] [ text "TcTurtle Examples" ]
                , Html.ul []
                    [ Html.li [] [ text "To draw a square:[ Repeat 4 [ Forward 100 Left 90 ] ]" ]
                    , Html.li [] [ text "To draw a star:[ Repeat 5 [ Forward 100, Right 144 ] ]" ]
                    , Html.li [] [ text "To draw a circle:[ Repeat 36 [ Forward 10, Right 10 ] ]" ]
                    , Html.li [] [ text "To draw a triangle:[ Repeat 3 [ Forward 100, Left 120 ] ]" ]
                    , Html.li [] [ text "To draw a hexagon:[ Repeat 6 [ Forward 100, Left 60 ] ]" ]
                    , Html.li [] [ text "To draw a spiral:[ Repeat 36 [ Forward 10, Right 10 ] ]" ]
                    ]
                ]
            ]
        ]





viewOutput : Result String (List Instruction) -> Html msg
viewOutput result =
    case result of
        Ok instructions ->
            svg [ width "500", height "500", viewBox "0 0 500 500" ]
                [ drawInstructions instructions ]

        Err error ->
            text (error)


-- MAIN

main : Program () Model Msg
main =
    Browser.sandbox { init = init, update = update, view = view }

