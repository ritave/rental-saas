import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Http

main =
  Html.program
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    }

eventServer : String
eventServer =
    ""

-- MODEL

type alias Model =
    { summary : String
    , user : String
    , location : String
--    , start : Date
--    , end : Date
    }

model : Model
model =
  Model "" "" ""

-- INIT

init : (Model, Cmd Msg)
init =
    (model, Cmd.none)

-- UPDATE

type Msg =
    Summary String
    | User String
    | Location String
--    | StartDate Date
--    | EndDate Date
    | Submit

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        Summary summary ->
            ({model | summary = summary}, Cmd.none)
        User user ->
            ({model | user = user}, Cmd.none)
        Location location ->
            ({model | location = location}, Cmd.none)
    --    EndDate endDate ->
    --        {model | end = endDate}
    --    StartDate startDate ->
    --        {model | start = startDate}
        Submit ->
            (model, Cmd.none)


-- VIEW

view : Model -> Html Msg
view model =
  div []
    [ input [ type_ "text", placeholder "Email", onInput User ] []
    , input [ type_ "text", placeholder "Summary", onInput Summary ] []
    , input [ type_ "text", placeholder "Location", onInput Location ] []
--    , input [ type_ "date", onClick StartDate ] []
--    , input [ type_ "date", onClick EndDate ] []
    , button [ ] [ text "Send" ]
    ]


-- TODO
-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

-- HTTP

submitEventCreation =
    Http.post