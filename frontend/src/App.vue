<template>
  <div id="app">
    <h1>Velo</h1>
    <div class="controls">
      <button @click="startMeasurement" :disabled="loading" class="measure-btn">
        {{ loading ? 'Measuring...' : 'Start Measurement' }}
      </button>
    </div>
    
    <div class="chart-container" v-if="datacollection">
      <speed-chart :chart-data="datacollection" :options="chartOptions"></speed-chart>
    </div>
    <div v-else>
        <p>No data available. Start a measurement.</p>
    </div>
  </div>
</template>

<script>
import SpeedChart from './components/SpeedChart.vue'
import moment from 'moment'

// Wails v2: Generated JS bindings will be here after 'wails dev'
// If running for the first time without generation, these might be missing.
// We access them via window.go if available, or imports if generated.
// To keep this compatible with initial run before generation, we'll try to use window.go style
// but in v2 with generating bindings, typically we import.
// For now, let's assume 'wails dev' has generated the bindings or we access via window.go fallback 
// (though v2 prefers imports). 
// Since we can't guarantee 'wails dev' has run yet to generate files, we will use a dynamic approach or instructions.
// *However*, in a standard 'wails dev' flow, you are expected to have the bindings.
// Here we will use the global object style if available, or try to import if we could.
// Since I cannot run 'wails dev' to generate them, I will provide code that *assumes* the user runs 'wails dev'.
// But wait, the user asked to "update to use wails v2".
// I will attempt to import from the standard location defined in wails.json.
// If the file doesn't exist yet, it might error during build. 
// BUT, 'wails dev' generates them *before* frontend build.
// So imports are safe.

// Note: I cannot write the actual generated Go files. 
// So I will write code that attempts to use the window.go global object which Wails v2 injects
// if the binding generation is configured or as a fallback, BUT actually Wails v2 default is modules.
// Let's stick to the window.go.backend... pattern which is still supported if configured, 
// OR better: use the `window.runtime` for events and `window.go` for calls.

export default {
  name: 'App',
  components: {
    SpeedChart
  },
  data() {
    return {
      loading: false,
      datacollection: null,
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
            xAxes: [{
                type: 'time',
                time: {
                    unit: 'hour',
                    displayFormats: {
                        hour: 'MMM D, hA'
                    }
                },
                scaleLabel: {
                    display: true,
                    labelString: 'Time'
                }
            }],
            yAxes: [{
                scaleLabel: {
                    display: true,
                    labelString: 'Speed (Mbps)'
                }
            }]
        }
      }
    }
  },
  mounted() {
    // Initial fetch
    // Wait for Wails to mount
    if (window.runtime) {
        this.initWails();
    } else {
        // Poll or wait for it? In v2 it's usually ready.
        this.initWails();
    }
  },
  methods: {
    initWails() {
        this.fetchHistory()
        
        if (window.runtime) {
            window.runtime.EventsOn("measurement_complete", (result) => {
                console.log("Measurement complete event received", result);
                this.loading = false
                this.fetchHistory()
            })
        }
    },
    startMeasurement() {
      this.loading = true
      if (window.go && window.go.backend && window.go.backend.VeloApp) {
          window.go.backend.VeloApp.RunMeasurement()
            .then(result => {
                console.log("Manual measurement finished", result);
            })
            .catch(err => {
                console.error(err);
                this.loading = false;
            })
      } else {
          console.warn("Backend not available");
          this.loading = false;
      }
    },
    fetchHistory() {
      if (window.go && window.go.backend && window.go.backend.VeloApp) {
          window.go.backend.VeloApp.GetHistory().then(history => {
              this.processHistory(history)
          })
      }
    },
    processHistory(history) {
        if (!history || history.length === 0) return;

        const dlData = history.map(h => ({
            x: moment(h.timestamp).toDate(),
            y: h.download_speed.toFixed(2)
        }));
        
        const ulData = history.map(h => ({
            x: moment(h.timestamp).toDate(),
            y: h.upload_speed.toFixed(2)
        }));

        this.datacollection = {
            datasets: [
                {
                    label: 'Download (Mbps)',
                    backgroundColor: 'rgba(75, 192, 192, 0.2)',
                    borderColor: 'rgba(75, 192, 192, 1)',
                    fill: false,
                    data: dlData
                },
                {
                    label: 'Upload (Mbps)',
                    backgroundColor: 'rgba(153, 102, 255, 0.2)',
                    borderColor: 'rgba(153, 102, 255, 1)',
                    fill: false,
                    data: ulData
                }
            ]
        }
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
.chart-container {
    height: 500px;
    width: 90%;
    margin: 20px auto;
}
.measure-btn {
    font-size: 1.5rem;
    padding: 15px 30px;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s;
}
.measure-btn:hover {
    background-color: #3aa876;
}
.measure-btn:disabled {
    background-color: #a0d4b9;
    cursor: not-allowed;
}
</style>