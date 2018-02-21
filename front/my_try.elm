import Date exposing (Date)
import Html exposing (Html, button, div, text, input)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)


main =
  Html.beginnerProgram { model = model, view = view, update = update }

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


-- UPDATE

type Msg = 
    Summary String
    | User String
    | Location String
--    | StartDate Date
--    | EndDate Date

update : Msg -> Model -> Model
update msg model =
  case msg of
    Summary summary ->
        {model | summary = summary}
    User user ->
        {model | user = user}
    Location location ->
        {model | location = location}
--    EndDate endDate ->
--        {model | end = endDate}
--    StartDate startDate ->
--        {model | start = startDate}
        


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
