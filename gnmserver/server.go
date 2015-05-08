package gnmserver
import (
	"net/http"
	"fmt"
	"github.com/geonetwork/gnm_collect/gnmsys"
	"bytes"
	"regexp"
	"os"
	"io/ioutil"
)

type Server struct {
	Config gnmsys.SysConfig
	Sys gnmsys.System
	Port int
}

func (s Server) index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
	<body>
		<ul>
			<li><a href="status">Status</a></li>
			<li><a href="reports">Reports</a></li>
			<li><a href="save">Save</a></li>
			<li><a href="shutdown">Shutdown</a></li>
		</ul>
	</body>
</html>`)
}

func (s Server) status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"isRunning": true}`)
}

func (s Server) shutdown(w http.ResponseWriter, r *http.Request) {
	s.Sys.SignalTerm()
	fmt.Fprintf(w, `{"quitting": true}`)
}

func (s Server) save(w http.ResponseWriter, r *http.Request) {
	s.Sys.SignalFlush()
	fmt.Fprintf(w, `{"saving": true}`)
}
func (s Server) reports(w http.ResponseWriter, r *http.Request) {
	s.Sys.SignalFlush()
	id := func(id string) string {
		re := regexp.MustCompile(` |\.|>|:|#`)
		return re.ReplaceAllString(id, "-")
	}
	categories := bytes.Buffer{}
	categorized := map[string][]string{}
	for _, report := range s.Sys.GetReports() {
		link := fmt.Sprintf(`<a id="report-link-%s" class="list-group-item report-link" onclick="report('report-link-%s', '%s', '%s')">%s</a>
		`, id(report.GetName()), id(report.GetName()), report.GetCategory(), report.GetFileName(), report.GetName())
		categorized[report.GetCategory()] = append(categorized[report.GetCategory()], link)
	}
	details := bytes.Buffer{}
	for category, links := range categorized {
		categories.WriteString(fmt.Sprintf(`<a id="%s" class="list-group-item cat-item" onclick="show('%s')">%s</a>`, id(category), id(category), category))

		details.WriteString(fmt.Sprintf(`<ul id="cat-detail-%s" class="list-group cat-group" style="display:none">`, id(category)))
		for _, link := range links {
			details.WriteString(link)
		}
		details.WriteString("</ul>")
	}

	template := `<html>
  <head>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css">
    <script src="//code.jquery.com/jquery-1.11.3.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/js/bootstrap.min.js"></script>
    <script>
      function show(cat) {
        $('.cat-item').removeClass('active');
        $('.cat-group').hide();
        $('#' + cat).addClass('active');
        $('#cat-detail-' + cat).toggle()
      }
      function report(id, cat, report) {
        $('.report-link').removeClass('active')
      	$('#' + id).addClass('active')
      	$('#report').html('<img src="report/' + cat + '/' + report + '"/>')
      }
      $( document ).ready(function(){
        var items = $('.cat-item');
        if (items.length > 0) {
          show(items.first().attr('id'))
        }
      })
    </script>
  </head>
  <body>
    <div class="row">
      <div class="col-md-4">
        <div class="panel panel-default">
          <div class="panel-heading">
            <h3 class="panel-title">Reports</h3>
          </div>
          <div class="panel-body">
            <ul class="list-group">
%s
            </ul>
          </div>
        </div>
      </div>
      <div class="col-md-8">
        <div class="panel panel-default">
          <div class="panel-heading">
            <h3 class="panel-title">Report Groups</h3>
          </div>
          <div class="panel-body">
            <ul class="list-group">
%s
            </ul>
          </div>
        </div>
      </div>
    </div>
    <div id="report"/>
  </body>
</html>`
	fmt.Fprintf(w, fmt.Sprintf(template, categories.String(), details.String()))
}

var reportPathExtractor = regexp.MustCompile(`/report/([^/]+)/(.+)`)
func (s Server) report(w http.ResponseWriter, r *http.Request) {
	matches := reportPathExtractor.FindStringSubmatch(r.URL.Path)

	if len(matches) != 3 {
		http.NotFound(w, r)
		return
	}

	category, name := matches[1], matches[2]

	var path string
	for _, report := range s.Sys.GetReports() {
		if report.GetFileName() == name && report.GetCategory() == category {
			path = s.Sys.GetReportFile(report)
		}
	}

	if path != "" {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				http.Error(w, "Error Loading report file", 500)
			} else {
				w.Header().Add("Content-Type", "Content-Type:image/png")
				w.Write(data)
				return
			}
		}
	}

	http.NotFound(w, r)
}

func (s Server) Start() {
	http.HandleFunc("/status", s.status)
	http.HandleFunc("/save", s.save)
	http.HandleFunc("/shutdown", s.shutdown)
	http.HandleFunc("/reports", s.reports)
	http.HandleFunc("/report/", s.report)
	http.HandleFunc("/index.html", s.index)
	http.HandleFunc("/", s.index)
	fmt.Printf("gnm_collect Http Server is waiting on http://localhost:%d/index.html\n", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}