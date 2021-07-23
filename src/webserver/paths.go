package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// --> GET /spans
func get_all_spans(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM retrospect.spans;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /spans_by_trace/{id}
func get_all_spans_by_trace(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  trace_id, ok := vars["id"]

  if !ok {
    fmt.Println("trace_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.spans WHERE trace_id='%s';", trace_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /spans_by_chapter/{id}
func get_all_spans_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.spans WHERE chapter_id='%s';", chapter_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /spans_by_session/{id}
func get_all_spans_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /events
func get_all_events(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM retrospect.events;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /events_by_chapter/{id}
func get_all_events_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.events WHERE chapter_id='%s';", chapter_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /events_by_session/{id}
func get_all_events_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.events WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> POST /spans
func insert_spans(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Inserting a Span")
  // r.Body is type *http.bodyblob
  // io.ReadAll returns an array of bytes
  body, _ := io.ReadAll(r.Body)

  // format_spans takes an array of bytes
  // and returns an array of structs.CassandraSpan objects
  cspans := format_spans(body)

  for _, span := range(cspans) {
    if span == nil { continue }
    // json.Marshal returns the json encoding of the variable passed into it
    j, _ := json.Marshal(span)

    // each json-ified span is stringified and inserted into the database as a json object
    query := "INSERT INTO retrospect.spans JSON '" + string(j) + "';"
    session.Query(query).Exec()
  }

  w.WriteHeader(http.StatusOK)
}

// --> POST /events
func insert_events(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  cevent := format_event(body, r)
  if cevent == nil { return }

  j, _ := json.Marshal(cevent)
  query := "INSERT INTO retrospect.events JSON '" + string(j) + "';"
  session.Query(query).Exec()

  w.WriteHeader(http.StatusOK)
}

// --> GET /events/snapshots
func get_snapshots(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM retrospect.snapshots;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /events/snapshots_by_session/{id}
func get_all_snapshots_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.snapshots WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> POST /events/snapshots
func insert_snapshots(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Inserting a snapshot")
  body, _ := io.ReadAll(r.Body)
  snapshot := format_snapshot(body, r)
  if snapshot == nil { return }

  j, _ := json.Marshal(snapshot)
  query := "INSERT INTO retrospect.snapshots JSON '" + string(j) + "';"
  session.Query(query).Exec()

  w.WriteHeader(http.StatusOK)
}

// --> GET /trigger_routes
func get_all_trigger_routes(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON trigger_route, data FROM retrospect.spans;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /trace_ids_by_trigger
func get_all_trace_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  trigger_route := string(body)

  query := fmt.Sprintf("SELECT trace_id FROM retrospect.spans WHERE trigger_route='%s' ALLOW FILTERING;", trigger_route);

  j := enumerate_query(query)
  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /chapters_by_session/{id}
func get_all_chapter_ids_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
    // TODO: return this as error
  }

  query := fmt.Sprintf("SELECT JSON chapter_id FROM retrospect.spans WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /chapter_ids_by_trigger
func get_all_chapter_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  target := string(body)

  query := fmt.Sprintf("SELECT chapter_id FROM retrospect.spans WHERE trigger_route='%v' ALLOW FILTERING;", target);

  j := enumerate_query(query)
  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET? /span_search
func span_search_handler(w http.ResponseWriter, r *http.Request) {
  acceptedParams := []string {
    "trace_id",
    "user_id",
    "session_id",
    "chapter_id",
    "status_code",
  }

  var dynamicQuery []string

  for _, p := range acceptedParams {
    val := r.FormValue(p)
    if val != "" {
      if p != "status_code" {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v='%v'", p, val))
      } else {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v=%v", p, val))
      }
    }
  }

  dynamicQueryString := strings.Join(dynamicQuery," AND ")

  if len(dynamicQueryString) != 0 {
    dynamicQueryString = "WHERE " + dynamicQueryString + " ALLOW FILTERING"
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.spans " + dynamicQueryString + ";")

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET? /event_search
func event_search_handler(w http.ResponseWriter, r *http.Request) {
  acceptedParams := []string {
    "user_id",
    "session_id",
    "chapter_id",
  }

  var dynamicQuery []string

  for _, p := range acceptedParams {
    val := r.FormValue(p)
    if val != "" {
      if p != "status_code" {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v='%v'", p, val))
      } else {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v=%v", p, val))
      }
    }
  }

  dynamicQueryString := strings.Join(dynamicQuery," AND ")

  if len(dynamicQueryString) != 0 {
    dynamicQueryString = "WHERE " + dynamicQueryString + " ALLOW FILTERING"
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.events " + dynamicQueryString + ";")

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}
