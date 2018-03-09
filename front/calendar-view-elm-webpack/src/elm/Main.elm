import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick, on, keyCode)
import Http exposing (Response)

import Json.Encode as Encode
import Json.Decode as Decode

--import Dict exposing (..)
--import Date exposing (..)
import List

import Debug

-- TOOD split everything into modules

main : Program Flags Model Msg
main =
  Html.programWithFlags
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    }

-- DEFINITIONS

type alias Flags =
    { backend : String
    }

apiEventCreate : String -> String
apiEventCreate apiBase =
    apiBase ++ "event/create"

apiEventList : String -> String
apiEventList apiBase =
    apiBase ++ "event/list"

apiEventDelete : String -> String
apiEventDelete apiBase =
    apiBase ++ "event/delete"

-- TYPES

type alias Event =
    { summary : String
    , user : String
    , start : String
    , end : String
    , location : String
    , creationDate : String
    , uuid : String
    }

decodeEvent : Decode.Decoder Event
decodeEvent =
    Decode.map7 Event
        (Decode.field "summary" Decode.string)
        (Decode.field "user" Decode.string)
        (Decode.field "start" Decode.string)
        (Decode.field "end" Decode.string)
        (Decode.field "location" Decode.string)
        (Decode.field "creationDate" Decode.string)
        (Decode.field "uuid" Decode.string)

encodeEvent : Event -> Encode.Value
encodeEvent record =
    Encode.object
        [ ("summary",  Encode.string <| record.summary)
        , ("user",  Encode.string <| record.user)
        , ("start",  Encode.string <| record.start)
        , ("end",  Encode.string <| record.end)
        , ("location",  Encode.string <| record.location)
        , ("creationDate",  Encode.string <| record.creationDate)
        , ("uuid",  Encode.string <| record.uuid)
        ]

encodeEventDelete : String -> Encode.Value
encodeEventDelete uuid =
    Encode.object
        [ ("uuid", Encode.string <| uuid)
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

eventDeleteResponseDecoder : Decode.Decoder String
eventDeleteResponseDecoder =
    Decode.string

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

    , backend : String
    }

formEncoder : Model -> Encode.Value
formEncoder model =
    Encode.object
        [ ("summary", Encode.string model.summary)
        , ("user", Encode.string model.user)
        , ("location", Encode.string model.location)
        , ("start", Encode.string model.start)
        , ("end", Encode.string model.end)
        ]

--TODO true variable timezone (I'm not going to change it in a few months by hand)
timeZone : String
--timeZone = "Z"
timeZone = "+01:00"

validateForm : Model -> (Model, Cmd Msg)
validateForm model =
    if model.user == "" then ({model | error = "No user specified"}, Cmd.none) else 
    if model.location == "" then ({model | error = "No location specified"}, Cmd.none) else
    if model.startDate == "" then ({model | error = "No start date"}, Cmd.none) else
    if model.endDate == "" then ({model | error = "No end date"}, Cmd.none) else
    if model.startTime == "" then ({model | error = "No start time"}, Cmd.none) else
    if model.endTime == "" then ({model | error = "No end time"}, Cmd.none) else
    let
--        "2006-01-02T15:04:05Z07:00"
        start = model.startDate ++ "T" ++ model.startTime ++ ":00" ++ timeZone
        end = model.endDate ++ "T" ++ model.endTime ++ ":00" ++ timeZone
        updatedModel = {model | start = start, end = end}
    in
     ({model | start = start, end = end}, eventCreatePost updatedModel)


-- INIT

startUpValue : String -> Model
--startUpValue = Model "Summary" "radekantichrist@gmail.com" "Location" ("") ("") "2018-02-25" "2018-02-25" "09:00" "10:00" "" []
startUpValue = Model "" "" "" ("") ("") "" "" "" "" "" []

init : Flags -> (Model, Cmd Msg)
init flags =
    (startUpValue flags.backend, eventListGet (startUpValue flags.backend))

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
    | EventDelete String
    | EventDeleteResponse (Result Http.Error String)
    | KeyDown Int

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
                    ({model | error = ""}, eventListGet model)
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
        EventDelete uuid ->
            let
                _ = log "Button clicked!" uuid
            in
                (model, eventDelete model uuid)
        EventDeleteResponse response ->
            case response of
                Ok trueResponse ->
                    let
                        _ = log "EventDelete" trueResponse
                    in
                    ({model | error = ""}, eventListGet model)
                Err error ->
                    let
                        errorMsg = errorToString error
                    in
                    ({model | error = errorMsg}, Cmd.none)
        KeyDown key ->
            if key == 13 then
                validateForm model
            else
                (model, Cmd.none)


-- VIEW

view : Model -> Html Msg
view model =
  div [ class "container" ]
    [ inputView model
    , errorView model
    , eventsView model
    ]

inputView : Model -> Html Msg
inputView model =
    colSm12
        [ formGroupInputWithLabel "email" "Email" "Email" "e@mail.com" [onInput User, onKeyDown KeyDown]
        , formGroupInputWithLabel "text" "Summary" "Summary" "Description of the event" [onInput Summary, onKeyDown KeyDown]
        , formGroupInputWithLabel "text" "Location" "Location" "Where is it going to take place?" [onInput Location, onKeyDown KeyDown]
        , colSm6ColSm6
            (formGroupInputWithLabel "date" "Start date" "sd" "" [onInput StartDate, onKeyDown KeyDown])
            (formGroupInputWithLabel "time" "Start time" "st" "09:00" [onInput StartTime, onKeyDown KeyDown])
        , colSm6ColSm6
            (formGroupInputWithLabel "date" "End date" "ed" "" [onInput EndDate, onKeyDown KeyDown])
            (formGroupInputWithLabel "time" "End time" "et" "10:00" [onInput EndTime, onKeyDown KeyDown])
        , button [ class "btn btn-default", onClick SubmitForm ] [ text "Send" ]
        ]

formGroupInputWithLabel : String -> String -> String -> String -> List(Attribute msg) -> Html msg
formGroupInputWithLabel tp lbl nm plch onInp =
    div [ class "form-group" ]
    [ label [ for nm ] [ text lbl ]
    , input ( List.append [ type_ tp, class "form-control", id nm, placeholder plch, name nm ] onInp ) []
    ]

colSm6ColSm6 : Html msg -> Html msg -> Html msg
colSm6ColSm6 first second =
   div [ class "row" ]
   [ div [ class "col-sm-6" ] [ first ]
   , div [ class "col-sm-6" ] [ second ]
   ]

eventsView : Model -> Html Msg
eventsView model =
    colSm12
    [
        table [ class "table table-bordered" ]
        [ eventsViewHead
        , tbody []
            ( List.map
                (\e ->
                    singleEventView e
                )
                model.events
            )
        ]
    ]

eventsViewHead : Html msg
eventsViewHead =
    thead []
    [ tr []
        [ td [] [ text "User" ]
        , td [] [ text "Start" ]
        , td [] [ text "End" ]
        , td [] [ text "Location" ]
        , td [] [ text "Summary" ]
        , td [] [ text "Created" ]
        , td [] []
        ]
    ]

singleEventView : Event -> Html Msg
singleEventView event =
    tr []
    [ td [] [ text (event.user) ]
    , td [] [ text (event.start) ]
    , td [] [ text (event.end) ]
    , td [] [ text (event.location) ]
    , td [] [ text (event.summary) ]
    , td [] [ text (event.creationDate) ]
    , td [] [ buttonDeleteEvent event.uuid ]
    ]

colSm12 : List(Html msg) -> Html msg
colSm12 whatever =
    div [ class "row", style [("margin-top", "30px")] ] [ div [ class "col-sm-12" ] whatever ]

buttonDeleteEvent : String -> Html Msg
buttonDeleteEvent uuid =
    button [ class "btn btn-danger btn-xs", onClick (EventDelete uuid)] [ span [ class "glyphicon glyphicon-trash" ] [] ]

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
        Http.post (apiEventCreate model.backend) body eventCreateResponseDecoder

eventCreatePost : Model -> Cmd Msg
eventCreatePost model =
    Http.send EventCreateResponse (eventCreateRequestBuilder model)

eventListGet : Model -> Cmd Msg
eventListGet model =
    Http.send EventListResponse (Http.get (apiEventList model.backend) eventListResponseDecoder)

eventDeleteRequestBuilder : Model -> String -> Http.Request String
eventDeleteRequestBuilder model uuid =
    let
        body =
            uuid
                |> encodeEventDelete
                |> Http.jsonBody
    in
        Http.post (apiEventDelete model.backend) body eventDeleteResponseDecoder

eventDelete : Model -> String -> Cmd Msg
eventDelete model uuid =
    Http.send EventDeleteResponse (eventDeleteRequestBuilder model uuid)

-- LOGGING

log : String -> a -> a
log = Debug.log

-- ERRORS
errorView : Model -> Html msg
errorView model =
  let
    (color, message) =
      if model.error == "" then
        ("success", "No errors")
      else
        ("danger", model.error)
  in
    colSm12
    [ div [ class ("alert alert-" ++ color) ] [ text message ] ]

errorToString : Http.Error -> String
errorToString error =
    case error of 
        Http.BadUrl something -> "Bad url: " ++ something
        Http.Timeout -> "Timeout"
        Http.NetworkError -> "Network error"
        Http.BadStatus response -> "Bad status: " ++ response.body
        Http.BadPayload something response -> "Bad payload: " ++ something ++ response.body

-- LOLS
onKeyDown : (Int -> msg) -> Attribute msg
onKeyDown tagger =
  on "keydown" (Decode.map tagger keyCode)