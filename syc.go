package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "html/template"
)

type Snippet struct {
   Title string
   Category string
   SubCategory string
   Problem string
   Solution []byte
   Time_complexiy string
   Space_complexity string
   Location string
}

func (s *Snippet) save() error {
    filename := s.Title + ".txt"
    return ioutil.WriteFile(filename, s.Solution, 0600)
}

func loadSnippet(title string) (*Snippet, error){
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Snippet{Title: title, Solution: body}, nil
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Snippet) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w,p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadSnippet(title)
    if err != nil {
        p = &Snippet{Title: title}
    }
    renderTemplate(w, "edit", p)
}


func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := loadSnippet(title);
    if err != nil {
        http.Redirect(w, r, "/edit/" +title, http.StatusFound)
    }
    renderTemplate(w, "view", p)
}


func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/") : ]
    solution := r.FormValue("body")
    p := &Snippet{Title: title, Solution: []byte(solution)}
    p.save()
    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, " %s", r.URL.Path[1:])
}

func main() {
    s1 := &Snippet{Title: "Test snippet", Solution: []byte("for int i = 0 do ooga")}
    s1.save()
    p2, _ := loadSnippet("Test snippet")
    http.HandleFunc("/", handler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
    fmt.Println(string(p2.Solution))
}

