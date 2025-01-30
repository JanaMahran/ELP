module Drawing exposing (drawInstructions)

import Svg exposing (Svg, line, g)
import Svg.Attributes exposing (..)
import TCParser exposing (Instruction(..))  -- Import the Instruction type


-- STATE

type alias State msg =
    { position : (Float, Float)
    , angle : Float  -- Angle in degrees
    , lines : List (Svg msg) -- Accumulated SVG lines
    }


initialState : State msg
initialState =
    { position = (250, 250) -- Start at the center of a 500x500 canvas
    , angle = 0
    , lines = []
    }


-- DRAWING FUNCTION

drawInstructions : List Instruction -> Svg msg
drawInstructions instructions =
    let
        finalState =
            List.foldl processInstruction initialState instructions
    in
    g [] finalState.lines


processInstruction : Instruction -> State msg -> State msg
processInstruction instruction state =
    case instruction of
        Forward distance ->
            let
                (x, y) =
                    state.position

                radians =
                    degreesToRadians state.angle

                newX =
                    x + toFloat distance * cos radians

                newY =
                    y - toFloat distance * sin radians

                newLine =
                    line
                        [ x1 (String.fromFloat x)
                        , y1 (String.fromFloat y)
                        , x2 (String.fromFloat newX)
                        , y2 (String.fromFloat newY)
                        , stroke "black"
                        , strokeWidth "2"
                        ]
                        []
            in
            { state
                | position = (newX, newY)
                , lines = state.lines ++ [ newLine ]
            }

        Left degrees ->
            { state | angle = state.angle - toFloat degrees }

        Right degrees ->
            { state | angle = state.angle + toFloat degrees }

        Repeat times subInstructions ->
            List.foldl
                (\_ accState -> List.foldl processInstruction accState subInstructions)
                state
                (List.repeat times ())


-- HELPER FUNCTIONS

degreesToRadians : Float -> Float
degreesToRadians degrees =
    degrees * pi / 180

