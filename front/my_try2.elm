import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Http

import Json.Encode as Encode
import Json.Decode as Decode

import Debug

main =
  Html.program
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    }

apiEventCreate : String
apiEventCreate =
--    "https://calendar-cron.appspot.com/event/create"
    "http://localhost:8080/event/create"

-- MODEL

type alias Model =
    { summary : String
    , user : String
    , location : String
--    , start : Date
--    , end : Date
    , error : String
    }

model : Model
model = startUpValue


formEncoder : Model -> Encode.Value
formEncoder model =
    Encode.object
        [ ("summary", Encode.string model.summary)
        , ("user", Encode.string model.user)
        , ("location", Encode.string model.location)
        , ("start", Encode.string "wat")
        , ("end", Encode.string "wat")
        ]

postEventResponseDecoder : Decode.Decoder String
postEventResponseDecoder =
    Decode.string

-- INIT

startUpValue : Model
startUpValue = Model "Summary" "radekantichrist@gmail.com" "Location" ""

init : (Model, Cmd Msg)
init =
    (startUpValue, Cmd.none)

-- UPDATE

type Msg =
    Summary String
    | User String
    | Location String
--    | StartDate Date
--    | EndDate Date
    | SubmitForm
    | PostEventResponse (Result Http.Error String)
    | Error String

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
        SubmitForm ->
            (model, postEventSend model)
        PostEventResponse result ->
            let
                maybeError = extractErrorFromResult result
            in
            ({model | error = maybeError}, Cmd.none)
        Error error ->
            ({ model | error = error}, Cmd.none)


-- VIEW

view : Model -> Html Msg
view model =
  div []
    [ input [ type_ "text", placeholder "Email", onInput User ] []
    , input [ type_ "text", placeholder "Summary", onInput Summary ] []
    , input [ type_ "text", placeholder "Location", onInput Location ] []
--    , input [ type_ "date", onClick StartDate ] []
--    , input [ type_ "date", onClick EndDate ] []
    , button [ onClick SubmitForm ] [ text "Send" ]
    , br [] []
    , viewError model
    ]


-- TODO
-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

-- HTTP

postEventCreation : Model -> Http.Request String
postEventCreation model =
    let
        body =
            model
                |> formEncoder
                |> Http.jsonBody
    in
        Http.post apiEventCreate body postEventResponseDecoder

postEventSend : Model -> Cmd Msg
postEventSend model =
    Http.send PostEventResponse (postEventCreation model)

-- LOGGING

-- ERRORS
viewError : Model -> Html msg
viewError model =
  let
    (color, message) =
      if model.error == "" then
        ("green", "No errors")
      else
        ("red", model.error)
  in
    div [ style [("color", color)] ] [ text message ]

errorToString : Http.Error -> String
errorToString error =
    case error of 
        Http.BadUrl something -> "Bad url: " ++ something
        Http.Timeout -> "Timeout"
        Http.NetworkError -> "Network error"
        Http.BadStatus _ -> "Bad status"
        Http.BadPayload something _-> "Bad payload: " ++ something

extractErrorFromResult : Result Http.Error String -> String
extractErrorFromResult result =
    case result of
        Ok _ -> ""
        Err error ->
            errorToString error
