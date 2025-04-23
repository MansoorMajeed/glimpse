package server

import (
	"embed"
	"html/template"
	"net/http"
	"sort"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
)

//go:embed templates/*.html
var templateFS embed.FS
var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

type DashboardAgent struct {
	Hostname        string
	OS              string
	LastSeenAgo     string
	FormattedUptime string
	Metrics         *pb.AgentMetrics
}

func StartHTTPServer(store *ServerStore) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "layout.html", nil)
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {

		rawAgents := store.GetAllAgents()
		agentList := make([]DashboardAgent, 0, len(rawAgents))

		sort.Slice(rawAgents, func(i, j int) bool {
			return rawAgents[i].Hostname < rawAgents[j].Hostname
		})

		for _, a := range rawAgents {
			latest := a.Latest()
			if latest == nil {
				continue
			}
			agentList = append(agentList, DashboardAgent{
				Hostname:        a.Hostname,
				OS:              a.OS,
				LastSeenAgo:     formatRelative(a.LastSeen),
				Metrics:         latest,
				FormattedUptime: formatUptime(latest.Uptime),
			})
		}

		err := templates.ExecuteTemplate(w, "agents.html", agentList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	go func() {
		if err := http.ListenAndServe("0.0.0.0:5000", nil); err != nil {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
}
