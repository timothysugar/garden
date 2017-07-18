package garden

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

type plant struct {
	Name    string
	Planted time.Time
}

func init() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/garden", garden)
	http.HandleFunc("/plant/", addPlants)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello gardener")
}

func addPlants(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.Method != "PUT" {
		http.NotFound(w, r)
	} else {
		k := datastore.NewIncompleteKey(c, "Variety", gardenKey(c))
		name := r.URL.Path[len("/plant/"):]
		p := plant{
			Name:    name,
			Planted: time.Now(),
		}
		if _, err := datastore.Put(c, k, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "added %v to the garden at %v", p.Name, p.Planted)
	}
}

func garden(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	plants := make([]plant, 0, 10)

	q := datastore.NewQuery("Variety").Ancestor(gardenKey(c))
	if _, err := q.GetAll(c, &plants); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := gardenTempl.Execute(w, plants); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

var gardenTempl = template.Must(template.New("Garden").Parse(`
<html>
  <head>
    <title>Garden Journal</title>
  </head>
  <body>
	<p>My Garden</p>
    {{range .}}
		<table>
		<th>Plant name</th>
		<th>Datetime planted</th>
		<tr>
		<td>{{.Name}}</td>
		<td>{{.Planted}}</td>
		</tr>
			{{else}}
			<p>Your garden is empty</p>
    {{end}}
  </body>
</html>
`))

func gardenKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "garden", "default_garden", 0, nil)
}
