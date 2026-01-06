<template>
  <div id="app">
    <div class="main-content">
      <div class="row header-row">
        <h1>Velo Dashboard</h1>
        <button @click="startMeasurement" :disabled="loading" class="measure-btn">
          <i v-if="loading" class="loader"></i>
          {{ loading ? 'Measuring...' : 'Start New Test' }}
        </button>
      </div>

      <card title="Performance History" sub-title="Internet speed metrics over time">
        <div class="controls-area">
          <div class="control-group">
            <span class="control-label">Type:</span>
            <div class="btn-group">
              <button 
                v-for="type in ['line', 'area', 'bar']" 
                :key="type"
                :class="{ active: selectedChartType === type }"
                @click="selectedChartType = type"
              >
                {{ type.charAt(0).toUpperCase() + type.slice(1) }}
              </button>
            </div>
          </div>

          <div class="control-group">
            <span class="control-label">Scope:</span>
            <div class="btn-group">
              <button 
                v-for="scope in scopes" 
                :key="scope.value"
                :class="{ active: selectedScope === scope.value }"
                @click="changeScope(scope.value)"
              >
                {{ scope.label }}
              </button>
            </div>
          </div>

          <div class="control-group" v-if="availableIPs.length > 0">
            <span class="control-label">IP:</span>
            <velo-dropdown 
                v-model="selectedIP" 
                :options="ipOptions" 
                @change="processHistory(fullHistory)"
                placeholder="All IPs"
            />
          </div>
        </div>

        <div class="chart-wrapper">
          <div v-if="datacollection" class="chart-inner">
            <speed-chart v-if="selectedChartType === 'line' || selectedChartType === 'area'" :chart-data="datacollection" :options="chartOptions"></speed-chart>
            <bar-chart v-else :chart-data="datacollection" :options="chartOptions"></bar-chart>
          </div>
          <div v-else class="no-data">
              <p>No data available. Please start a measurement.</p>
          </div>
        </div>
        
        <template #footer>
            <div class="stats">
                <span v-if="lastUpdated">Last updated: {{ lastUpdated }}</span>
                <span v-else>No measurements yet.</span>
            </div>
        </template>
      </card>
    </div>
  </div>
</template>

<script>
import SpeedChart from './components/SpeedChart.vue'
import BarChart from './components/BarChart.vue'
import Card from './components/Card.vue'
import VeloDropdown from './components/VeloDropdown.vue'
import moment from 'moment'

export default {
  name: 'App',
  components: {
    SpeedChart,
    BarChart,
    Card,
    VeloDropdown
  },
  data() {
    return {
      loading: false,
      datacollection: null,
      fullHistory: [], // Store full data for filtering
      selectedChartType: 'area', 
      selectedScope: '24h',
      selectedIP: 'all',
      availableIPs: [],
      lastUpdated: null,
      scopes: [
          { label: '1H', value: '1h' },
          { label: '24H', value: '24h' },
          { label: '7D', value: '7d' },
          { label: '30D', value: '30d' },
          { label: 'ALL', value: 'all' }
      ],
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: true,
                position: 'bottom'
            },
            tooltip: {
                mode: 'index',
                intersect: false,
            }
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    displayFormats: {
                        hour: 'MMM D, hA',
                        day: 'MMM D'
                    },
                    tooltipFormat: 'll HH:mm'
                },
                grid: {
                    display: false
                },
                title: {
                    display: true,
                    text: 'Time'
                },
                ticks: {
                    autoSkip: true,
                    maxTicksLimit: 8
                }
            },
            y: {
                title: {
                    display: true,
                    text: 'Speed (Mbps)'
                },
                grid: {
                    borderDash: [2, 4],
                    color: "rgba(0, 0, 0, 0.1)"
                },
                beginAtZero: true
            }
        }
      }
    }
  },
  computed: {
    ipOptions() {
        const opts = [{ label: 'All IPs', value: 'all' }];
        this.availableIPs.forEach(ip => {
            opts.push({ label: ip, value: ip });
        });
        return opts;
    }
  },
  mounted() {
    if (window.runtime) {
        this.initWails();
    }
    else {
        this.initWails();
    }
  },
  methods: {
    initWails() {
        this.fetchHistory()
        
        if (window.runtime) {
            window.runtime.EventsOn("measurement_complete", (result) => {
                this.loading = false
                this.fetchHistory()
            })
        }
    },
    changeScope(scope) {
        this.selectedScope = scope;
        this.fetchHistory();
    },
    startMeasurement() {
      this.loading = true
      if (window.go && window.go.backend && window.go.backend.VeloApp) {
          window.go.backend.VeloApp.RunMeasurement()
            .then(result => {
                // Event listener will handle update
            })
            .catch(err => {
                console.error(err);
                this.loading = false;
            })
      } else {
          // Simulate for dev without backend
          setTimeout(() => { this.loading = false; }, 2000);
      }
    },
    fetchHistory() {
      if (window.go && window.go.backend && window.go.backend.VeloApp) {
          window.go.backend.VeloApp.GetHistory(this.selectedScope).then(history => {
              this.fullHistory = history || [];
              this.updateAvailableIPs(this.fullHistory);
              this.processHistory(this.fullHistory)
          })
      }
    },
    updateAvailableIPs(history) {
        const ips = new Set();
        history.forEach(h => {
            if (h.ip_address) ips.add(h.ip_address);
        });
        this.availableIPs = Array.from(ips);
        
        // If selected IP is no longer available, reset to all
        if (this.selectedIP !== 'all' && !this.availableIPs.includes(this.selectedIP)) {
            this.selectedIP = 'all';
        }
    },
    processHistory(historyData) {
        // Filter by IP if needed
        let history = historyData;
        if (this.selectedIP !== 'all') {
            history = history.filter(h => h.ip_address === this.selectedIP);
        }

        if (!Array.isArray(history) || history.length === 0) {
            this.datacollection = null;
            return;
        }

        // Sort just in case
        history.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

        const lastItem = history.slice(-1)[0];
        this.lastUpdated = moment(lastItem.timestamp).format('YYYY-MM-DD HH:mm:ss');

        const dlData = history.map(h => ({
            x: moment(h.timestamp).toDate(),
            y: h.download_speed.toFixed(2)
        }));
        
        const ulData = history.map(h => ({
            x: moment(h.timestamp).toDate(),
            y: h.upload_speed.toFixed(2)
        }));

        const isArea = this.selectedChartType === 'area';
        const isBar = this.selectedChartType === 'bar';

        this.datacollection = {
            datasets: [
                {
                    label: 'Download',
                    backgroundColor: isArea ? 'rgba(66, 185, 131, 0.2)' : (isBar ? 'rgba(66, 185, 131, 0.5)' : 'transparent'),
                    borderColor: '#42b983',
                    pointBackgroundColor: '#42b983',
                    borderWidth: 2,
                    fill: isArea,
                    data: dlData
                },
                {
                    label: 'Upload',
                    backgroundColor: isArea ? 'rgba(52, 152, 219, 0.2)' : (isBar ? 'rgba(52, 152, 219, 0.5)' : 'transparent'),
                    borderColor: '#3498db',
                    pointBackgroundColor: '#3498db',
                    borderWidth: 2,
                    fill: isArea,
                    data: ulData
                }
            ]
        }
    }
  },
  watch: {
      selectedChartType() {
          // Re-process history to update chart colors/fill based on new type
          // We need to fetch again or just re-process cached data? 
          // Since we don't cache raw data in a variable, let's just fetch (it's fast local DB)
          this.fetchHistory();
      }
  }
}
</script>

<style>
body {
    background-color: #f4f3ef;
    margin: 0;
}
#app {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  padding: 20px;
}

.main-content {
    max-width: 1000px;
    margin: 0 auto;
}

.header-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

h1 {
    font-weight: 300;
    margin: 0;
}

.measure-btn {
    background-color: #41B883;
    border: none;
    color: white;
    padding: 10px 24px;
    text-align: center;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    font-size: 16px;
    border-radius: 4px;
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    transition: all 0.3s ease;
}

.measure-btn:hover {
    background-color: #35495E;
    transform: translateY(-1px);
}

.measure-btn:disabled {
    background-color: #aebdb6;
    cursor: not-allowed;
    transform: none;
}

.controls-area {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 15px;
    margin-bottom: 20px;
    padding: 0 10px;
}

.control-group {
    display: flex;
    align-items: center;
    gap: 10px;
}

.control-label {
    font-weight: 600;
    color: #666;
    font-size: 0.9rem;
}

.btn-group {
    display: inline-flex;
    border-radius: 4px;
    overflow: hidden;
    border: 1px solid #ddd;
}

.btn-group button {
    background-color: #fff;
    border: none;
    padding: 8px 16px;
    cursor: pointer;
    border-right: 1px solid #ddd;
    font-size: 0.9rem;
    color: #555;
    transition: background-color 0.2s;
}

.btn-group button:last-child {
    border-right: none;
}

.btn-group button:hover {
    background-color: #f8f9fa;
}

.btn-group button.active {
    background-color: #35495E;
    color: white;
}

.chart-wrapper {
    height: 400px;
    position: relative;
    padding: 10px;
}

.chart-inner {
    height: 100%;
    width: 100%;
}

.no-data {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    color: #999;
}

.loader {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  width: 14px;
  height: 14px;
  animation: spin 1s linear infinite;
  margin-right: 8px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>