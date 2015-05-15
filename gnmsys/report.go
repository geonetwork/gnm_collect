package gnmsys
import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"path"
	"os"
	"log"
	"time"
	"path/filepath"
	"fmt"
	"io"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

type PlotFactory func() *plot.Plot

type Report interface {
	GetCategory() string
	GetName() string
	GetFileName() string
	GetUpdateInterval() time.Duration
	Update(timeSeconds int64, metrics Json)
	Save(titleModifier string, reportDir string)
	Load(reportDir string)
}

type LineGraphReport struct {
	collectors []Collector
	name string
	report PlotFactory
	width, height vg.Length
	sampleConfig SampleConfig
}

func NewLineGraphReport(name string,
report PlotFactory,
width, height vg.Length,
sampleConfig SampleConfig,
collectors ...Collector) LineGraphReport {
	return LineGraphReport{collectors: collectors, name: name, report: report, width:width, height:height, sampleConfig: sampleConfig}
}
func (r LineGraphReport) GetCategory() string {
	return r.sampleConfig.Name
}
func (r LineGraphReport) GetName() string {
	return r.name
}
func (r LineGraphReport) GetFileName() string {
	return r.name + ".png"
}
func (r LineGraphReport) Update(timeSeconds int64, metrics Json) {
	log.Printf("Updating report: %q\n", r.GetName())

	for _, coll := range r.collectors {
		coll.AddSample(r.sampleConfig.Unit().ConvertSeconds(timeSeconds), metrics)
	}
}

func (r LineGraphReport) GetUpdateInterval() time.Duration {
	return r.sampleConfig.UpdateInterval
}
func (r LineGraphReport) Save(titleModifier string, reportDir string) {
	report := r.report()
	report.X.Label.Text = r.sampleConfig.Unit().String()
	report.Title.Text = report.Title.Text + "(" + titleModifier + ")"
	lines := make([]interface{}, len(r.collectors) * 2)
	for i, coll := range r.collectors {
		lines[i * 2] = coll.Name()
		lines[i * 2 + 1] = coll.GetXYs()
	}

	err := plotutil.AddLinePoints(report, lines...)

	if err != nil {
		panic(err)
	}
	outDir := path.Join(reportDir, r.sampleConfig.DirName)
	os.MkdirAll(outDir, os.FileMode(0755))

	outputFile := path.Join(outDir, r.GetFileName())
	log.Printf("Saving %q to %q", r.name, outputFile)
	if err := report.Save(r.width, r.height, outputFile); err != nil {
		panic(err)
	}

	stateFile, err := os.Create(stateFilePath(reportDir, r.name))
	if err != nil {
		log.Printf("ERROR: unable to save report state for %s, %s", stateFile.Name(), err.Error())
		return;
	}
	defer stateFile.Close()

	r.saveState(stateFile)
}

func (r LineGraphReport) saveState(stateFile io.Writer) {
	serialized := &bytes.Buffer{}
	serialized.WriteString("{")
	for i, coll := range r.collectors {
		serialized.WriteString(fmt.Sprintf("%q: ", coll.Name()))
		if marshalledData, err := json.Marshal(coll.GetXYs()); err != nil {
			msg := fmt.Sprintf("Failed to convert state of collector %s to json, %s\n", coll.Name(), err.Error())
			fmt.Println(msg)
			log.Println(msg, err.Error())
		} else {
			serialized.Write(marshalledData)
		}
		if i + 1 < len (r.collectors) {
			serialized.WriteString(",")
		}
	}
	serialized.WriteString("}")

	if _, err := stateFile.Write(serialized.Bytes()); err != nil {
		msg := "Failed to save report state, %s\n"
		fmt.Printf(msg, err.Error())
		log.Printf(msg, err.Error())
	}
}

func stateFilePath(reportDir, name string) string {
	return filepath.Join(reportDir, name + ".json")
}

func (r LineGraphReport) Load(reportDir string) {
	stateFile, err := os.Open(stateFilePath(reportDir, r.name))

	if err == nil || !os.IsNotExist(err) {
		defer stateFile.Close()
		r.doLoad(r.name, stateFile)
	}
}

func (r LineGraphReport) doLoad(name string, stateFile io.Reader) {
	data, err := ioutil.ReadAll(stateFile)

	if err != nil {
		msg := fmt.Sprintf("Failed to load state data from file: %v\n", err.Error())
		fmt.Printf(msg)
		log.Printf(msg)
	} else {
		var unmarshalled map[string]interface{}
		if err := json.Unmarshal(data, &unmarshalled); err != nil {
			msg := fmt.Sprintf("Loading state failed: %v\n", err.Error())
			fmt.Printf(msg)
			log.Printf(msg)
		} else {
			for _, coll := range r.collectors {
				if xys, has := unmarshalled[coll.Name()]; has {
					for _, rawXY := range xys.([]interface{}) {
						xy := rawXY.(map[string]interface{})
						coll.AddXYSample(xy["X"].(float64), xy["Y"].(float64))
					}
				} else {
					log.Printf("No saved state for collector %q\n", coll.Name())
				}
			}
		}
	}
}
