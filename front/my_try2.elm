import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Http

import Json.Encode as Encode
import Json.Decode as Decode

import Dict exposing (..)

import Debug

main =
  Html.program
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    }

apiBase : String
apiBase =
    "http://localhost:8080/"

apiEventCreate : String
apiEventCreate =
--    "https://calendar-cron.appspot.com/event/create"
    apiBase ++ "event/create"

apiEventList : String
apiEventList =
    apiBase ++ "event/list"

-- TYPES

--type alias Event =
--    Dict String String

type alias Event =
    { summary : String
    , user : String
    , start : String
    , end : String
    , location : String
    , creationDate : String
    }

decodeEvent : Decode.Decoder Event
decodeEvent =
    Decode.map6 Event
        (Decode.field "summary" Decode.string)
        (Decode.field "user" Decode.string)
        (Decode.field "start" Decode.string)
        (Decode.field "end" Decode.string)
        (Decode.field "location" Decode.string)
        (Decode.field "creationDate" Decode.string)

encodeEvent : Event -> Encode.Value
encodeEvent record =
    Encode.object
        [ ("summary",  Encode.string <| record.summary)
        , ("user",  Encode.string <| record.user)
        , ("start",  Encode.string <| record.start)
        , ("end",  Encode.string <| record.end)
        , ("location",  Encode.string <| record.location)
        , ("creationDate",  Encode.string <| record.creationDate)
        ]

eventCreateResponseDecoder : Decode.Decoder String
eventCreateResponseDecoder =
    Decode.string

eventListResponseDecoder : Decode.Decoder (List Event)
eventListResponseDecoder =
    Decode.list (decodeEvent)

eventListDecoder : String -> List Event
eventListDecoder rawString =
    let
        response = Decode.decodeString eventListResponseDecoder rawString
    in
    case response of
        Ok result -> result
        Err _ -> []

-- MODEL

type alias Model =
    { summary : String
    , user : String
    , location : String
--    , start : Date
--    , end : Date
    , error : String
    , events : List Event
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

-- INIT

startUpValue : Model
startUpValue = Model "Summary" "radekantichrist@gmail.com" "Location" "" []

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
    | EventCreateResponse (Result Http.Error String)
    | EventListResponse (Result Http.Error (List Event))
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
            (model, eventCreatePost model)
        EventCreateResponse response ->
            case response of
                Ok trueResponse ->
                    let
                        _ = log "EventCreate" trueResponse -- WILL IT BLEND?
                    in
                    ({model | error = ""}, eventListGet)
                Err error ->
                    let
                        errorMsg = errorToString error
                    in
                    ({model | error = errorMsg}, Cmd.none)
        EventListResponse response ->
            case response of
                Ok events ->
                    let
                        _ = log "Events" events
                    in
                    ({model | events = events}, Cmd.none)
                Err error ->
                    let
                        errorMsg = errorToString error
                    in
                    ({model | error = errorMsg}, Cmd.none)
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
    , errorView model
    , br [] []

    ]

eventsView : Model -> Html msg
eventsView model =
    div [] []

singleEventView : Event -> Html msg
singleEventView event =
    pre []
    [
    ]

-- TODO?
-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

-- HTTP

eventCreateRequestBuilder : Model -> Http.Request String
eventCreateRequestBuilder model =
    let
        body =
            model
                |> formEncoder
                |> Http.jsonBody
    in
        Http.post apiEventCreate body eventCreateResponseDecoder

eventCreatePost : Model -> Cmd Msg
eventCreatePost model =
    Http.send EventCreateResponse (eventCreateRequestBuilder model)

eventListGet : Cmd Msg
eventListGet =
    Http.send EventListResponse (Http.get apiEventList eventListResponseDecoder)

-- LOGGING

log = Debug.log

-- ERRORS
errorView : Model -> Html msg
errorView model =
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

