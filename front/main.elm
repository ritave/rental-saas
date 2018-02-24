import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Http

import Json.Encode as Encode
import Json.Decode as Decode

import Dict exposing (..)
import Date exposing (..)

import Debug

-- TOOD split everything into modules

main =
  Html.program
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    }

-- DEFINITIONS

apiBase : String
apiBase =
    "https://calendar-cron.appspot.com/"
--    "http://localhost:8080/" -- TODO env variables/config -- generally webpack should be my friend

apiEventCreate : String
apiEventCreate =
    apiBase ++ "event/create"

apiEventList : String
apiEventList =
    apiBase ++ "event/list"

-- TYPES

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
    , start : String
    , end : String

    , startDate : String
    , endDate : String
    , startTime : String
    , endTime : String

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
        , ("start", Encode.string model.start)
        , ("end", Encode.string model.end)
        ]

validateForm : Model -> (Model, Cmd Msg)
validateForm model =
    if model.user == "" then ({model | error = "No user specified"}, Cmd.none) else 
    if model.startDate == "" then ({model | error = "No start date"}, Cmd.none) else
    if model.endDate == "" then ({model | error = "No end date"}, Cmd.none) else
    if model.startTime == "" then ({model | error = "No start time"}, Cmd.none) else
    if model.endTime == "" then ({model | error = "No end time"}, Cmd.none) else
    let
--        "2006-01-02T15:04:05Z07:00"
        start = model.startDate ++ "T" ++ model.startTime ++ ":00Z01:00"
        end = model.endDate ++ "T" ++ model.endTime ++ ":00Z01:00"
    in
     ({model | start = start, end = end}, eventCreatePost model)


-- INIT

startUpValue : Model
--startUpValue = Model "Summary" "radekantichrist@gmail.com" "Location" ("") ("") "2018-02-25" "2018-02-25" "09:00" "10:00" "" []
startUpValue = Model "" "" "" ("") ("") "" "" "" "" "" []

init : (Model, Cmd Msg)
init =
    (startUpValue, eventListGet)

-- UPDATE

type Msg =
    Summary String
    | User String
    | Location String
    | StartDate String
    | EndDate String
    | StartTime String
    | EndTime String
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
        EndDate endDate ->
            ({model | endDate = endDate}, Cmd.none)
        StartDate startDate ->
            ({model | startDate = startDate}, Cmd.none)
        EndTime endTime ->
            ({model | endTime = endTime}, Cmd.none)
        StartTime startTime ->
            ({model | startTime = startTime}, Cmd.none)
        SubmitForm ->
            validateForm model
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
  div [ class "main" ]
    [ inputView model
    , br [] []
    , errorView model
    , br [] []
    , eventsView model
    ]

inputView : Model -> Html Msg
inputView model =
    div [ class "form" ]
    [ input [ type_ "text", placeholder "Email", onInput User ] []
    , input [ type_ "text", placeholder "Summary", onInput Summary ] []
    , input [ type_ "text", placeholder "Location", onInput Location ] []
    , br [] []
    , input [ type_ "date", onInput StartDate ] []
    , input [ type_ "time", onInput StartTime, placeholder "09:00" ] []
    , br [] []
    , input [ type_ "date", onInput EndDate ] []
    , input [ type_ "time", onInput EndTime, placeholder "10:00" ] []
    , br [] []
    , button [ onClick SubmitForm ] [ text "Send" ]
    ]

eventsView : Model -> Html msg
eventsView model =
    ul [ class "events" ]
    (
        List.map
            (\e ->
                div [] [singleEventView e, br [] []]
            )
            model.events
    )

singleEventView : Event -> Html msg
singleEventView event =
    li [ class "event" ] [ul [class "event-element" ] 
        [ li [] [ text ("User: " ++ event.user) ]
        , li [] [ text ("Start: " ++ event.start) ]
        , li [] [ text ("End: " ++ event.end) ]
        , li [] [ text ("Location: " ++ event.location) ]
        , li [] [ text ("Summary: " ++ event.summary) ]
        , li [] [ text ("Created: " ++ event.creationDate) ]
        ]
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
        Http.BadPayload something response -> "Bad payload: " ++ something ++ response.body

