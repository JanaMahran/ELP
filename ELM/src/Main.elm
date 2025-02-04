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
    div [ Html.Attributes.style "display" "flex"
        , Html.Attributes.style "flex-direction" "column"
        , Html.Attributes.style "align-items" "center"
        , Html.Attributes.style "justify-content" "center"
        , Html.Attributes.style "min-height" "100vh"
        , Html.Attributes.style "background-color" "#f0f4f8"
        , Html.Attributes.style "font-family" "Arial, sans-serif"
        , Html.Attributes.style "color" "#333"
        , Html.Attributes.style "padding" "20px"
        ]
        [ Html.h1 [ Html.Attributes.style "font-size" "2.5em"
                  , Html.Attributes.style "color" "#1e3a8a"
                  , Html.Attributes.style "margin-bottom" "5px"
                  ]
            [ text "Welcome to TcTurtle!" ]
        , Html.p [ Html.Attributes.style "font-size" "1.1em"
                 , Html.Attributes.style "color" "#4b5563"
                 , Html.Attributes.style "margin-bottom" "20px"
                 , Html.Attributes.style "font-weight" "normal"
                 ]
            [ text "Enter a program to parse and draw below:" ]
        , div [ Html.Attributes.style "display" "flex"
              , Html.Attributes.style "gap" "20px"
              , Html.Attributes.style "width" "100%"
              , Html.Attributes.style "max-width" "900px"
              , Html.Attributes.style "height" "75vh"
              ]
            [ div [ Html.Attributes.style "flex" "2"
                  , Html.Attributes.style "padding" "20px"
                  , Html.Attributes.style "background-color" "white"
                  , Html.Attributes.style "box-shadow" "0 4px 8px rgba(0, 0, 0, 0.1)"
                  , Html.Attributes.style "border-radius" "12px"
                  , Html.Attributes.style "overflow-y" "auto"
                  , Html.Attributes.style "max-height" "100%"
                  ]
                [ input
                    [ placeholder "Enter your TcTurtle program here"
                    , value model.userInput
                    , onInput UpdateInput
                    , Html.Attributes.style "width" "95%"
                    , Html.Attributes.style "padding" "12px"
                    , Html.Attributes.style "font-size" "1em"
                    , Html.Attributes.style "border" "1px solid #d1d5db"
                    , Html.Attributes.style "border-radius" "8px"
                    , Html.Attributes.style "margin-bottom" "12px"
                    , Html.Attributes.style "background-color" "#f9fafb"
                    ]
                    []
                , button
                    [ onClick ParseInput
                    , Html.Attributes.style "background-color" "#1e40af"
                    , Html.Attributes.style "color" "white"
                    , Html.Attributes.style "font-size" "1em"
                    , Html.Attributes.style "padding" "12px"
                    , Html.Attributes.style "border" "none"
                    , Html.Attributes.style "border-radius" "8px"
                    , Html.Attributes.style "width" "100%"
                    , Html.Attributes.style "cursor" "pointer"
                    , Html.Attributes.style "box-shadow" "0 4px 6px rgba(0, 0, 0, 0.1)"
                    ]
                    [ text "Parse and Draw" ]
                , div [ Html.Attributes.style "margin-top" "20px"
                      , Html.Attributes.style "overflow-y" "auto"
                      , Html.Attributes.style "max-height" "60vh"
                      , Html.Attributes.style "padding" "10px"
                      , Html.Attributes.style "border" "1px solid #d1d5db"
                      , Html.Attributes.style "border-radius" "8px"
                      , Html.Attributes.style "background-color" "#f9fafb"
                      ]
                    [ viewOutput model.parsedInstructions ]
                ]
            , div [ Html.Attributes.style "flex" "1"
                  , Html.Attributes.style "padding" "20px"
                  , Html.Attributes.style "background-color" "#1e3a8a"
                  , Html.Attributes.style "color" "white"
                  , Html.Attributes.style "border-radius" "12px"
                  , Html.Attributes.style "box-shadow" "0 4px 8px rgba(0, 0, 0, 0.1)"
                  , Html.Attributes.style "overflow-y" "auto"
                  , Html.Attributes.style "max-height" "100%"
                  ]
                [ Html.h2 [ Html.Attributes.style "font-size" "1.8em", Html.Attributes.style "margin-bottom" "16px" ] [ text "TcTurtle Examples" ]
                , Html.ul [ Html.Attributes.style "list-style-type" "none", Html.Attributes.style "padding" "0" ]
                    [ Html.li [ Html.Attributes.style "margin-bottom" "16px" ] [ text "Square: [ Repeat 4 [ Forward 100, Left 90 ] ]" ]
                    , Html.li [ Html.Attributes.style "margin-bottom" "16px" ] [ text "Star: [ Repeat 5 [ Forward 100, Right 144 ] ]" ]
                    , Html.li [ Html.Attributes.style "margin-bottom" "16px" ] [ text "Circle: [ Repeat 36 [ Forward 10, Right 10 ] ]" ]
                    , Html.li [ Html.Attributes.style "margin-bottom" "16px" ] [ text "Triangle: [ Repeat 3 [ Forward 100, Left 120 ] ]" ]
                    , Html.li [ Html.Attributes.style "margin-bottom" "16px" ] [ text "Hexagon: [ Repeat 6 [ Forward 100, Left 60 ] ]" ]
                    , Html.li [ Html.Attributes.style "margin-bottom" "16px" ] [ text "Spiral: [ Repeat 36 [ Forward 10, Right 10 ] ]" ]
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

