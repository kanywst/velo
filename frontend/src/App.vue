<template>
  <div id="app">
    <aside class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <div class="logo-icon">V</div>
          <span class="logo-text">Velo</span>
        </div>
      </div>
      
      <nav class="sidebar-nav">
        <div class="nav-item" :class="{ active: activeTab === 'dashboard' }" @click="activeTab = 'dashboard'">
          <i class="nav-icon">📊</i>
          <span>Dashboard</span>
        </div>
        <div class="nav-item" :class="{ active: activeTab === 'history' }" @click="activeTab = 'history'">
          <i class="nav-icon">🕒</i>
          <span>History</span>
        </div>
        <div class="nav-item" :class="{ active: activeTab === 'settings' }" @click="activeTab = 'settings'">
          <i class="nav-icon">⚙️</i>
          <span>Settings</span>
        </div>
      </nav>

      <div class="sidebar-footer">
        <button @click="startMeasurement" :disabled="loading" class="measure-btn">
          <i v-if="loading" class="loader"></i>
          {{ loading ? 'Testing...' : 'Start New Test' }}
        </button>
      </div>
    </aside>

    <main class="main-content">
      <header class="content-header">
        <div class="header-title">
          <h1 v-if="activeTab === 'dashboard'">Network Dashboard</h1>
          <h1 v-else-if="activeTab === 'history'">Test History</h1>
          <h1 v-else-if="activeTab === 'settings'">Settings</h1>
          <p v-if="activeTab === 'dashboard' && lastUpdated">Last updated: {{ lastUpdated }}</p>
        </div>
        <div class="header-actions">
          <!-- Additional actions can go here -->
        </div>
      </header>

      <div class="content-body" v-if="activeTab === 'dashboard'">
        <div class="dashboard-grid">
          <card title="Performance History" sub-title="Download and Upload speeds over time">
            <template #header>
              <div class="card-header-flex">
                <div>
                  <h4 class="card-title">Performance History</h4>
                  <p class="card-category">Download and Upload speeds over time</p>
                </div>
                <div class="chart-controls">
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
              </div>
            </template>

            <div class="controls-area">
              <div class="control-group">
                <span class="control-label">Chart Type:</span>
                <div class="btn-group mini">
                  <button 
                    v-for="type in ['area', 'line', 'bar']" 
                    :key="type"
                    :class="{ active: selectedChartType === type }"
                    @click="selectedChartType = type"
                  >
                    {{ type.charAt(0).toUpperCase() + type.slice(1) }}
                  </button>
                </div>
              </div>

              <div class="control-group" v-if="availableIPs.length > 0">
                <span class="control-label">Interface IP:</span>
                <velo-dropdown 
                    v-model="selectedIP" 
                    :options="ipOptions" 
                    @change="processHistory(fullHistory)"
                    placeholder="All Interfaces"
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
          </card>
        </div>
      </div>

      <div class="content-body" v-else-if="activeTab === 'history'">
        <div class="dashboard-grid">
          <card title="Detailed History" sub-title="List of all previous network tests">
            <div class="table-responsive">
              <table class="velo-table">
                <thead>
                  <tr>
                    <th>Timestamp</th>
                    <th>Download</th>
                    <th>Upload</th>
                    <th>Latency</th>
                    <th>IP Address</th>
                    <th class="text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in sortedFullHistory" :key="item.id">
                    <td>{{ formatFullDate(item.timestamp) }}</td>
                    <td class="text-success font-weight-bold">{{ item.download_speed.toFixed(2) }} Mbps</td>
                    <td class="text-primary font-weight-bold">{{ item.upload_speed.toFixed(2) }} Mbps</td>
                    <td>{{ item.latency.toFixed(0) }} ms</td>
                    <td><code class="ip-code">{{ item.ip_address }}</code></td>
                    <td class="text-right">
                      <button class="btn-icon delete" @click="deleteItem(item.id)" title="Delete Entry">
                        <i class="icon-delete">🗑️</i>
                      </button>
                    </td>
                  </tr>
                  <tr v-if="fullHistory.length === 0">
                    <td colspan="6" class="text-center py-4">No history records found.</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </card>
        </div>
      </div>

      <div class="content-body" v-else-if="activeTab === 'settings'">
        <div class="dashboard-grid">
          <div class="settings-layout">
            <card title="General Settings">
              <div class="setting-item">
                <div class="setting-info">
                  <label>Background Measurements</label>
                  <p>Measurements are automatically taken every 1 hour while the app is running.</p>
                </div>
                <div class="setting-action">
                  <span class="badge badge-info">Active (1h)</span>
                </div>
              </div>
            </card>

            <card title="Data Management">
              <div class="setting-item">
                <div class="setting-info">
                  <label>Clear All History</label>
                  <p>Permanently delete all speed test results from the database. This action cannot be undone.</p>
                </div>
                <div class="setting-action">
                  <button class="btn-danger" @click="confirmClearHistory">Clear Database</button>
                </div>
              </div>
            </card>
          </div>
        </div>
      </div>
    </main>

    <!-- UI Components -->
    <velo-modal 
      :show="showDeleteModal" 
      title="Delete Result"
      message="Are you sure you want to delete this test result? This action cannot be undone."
      confirm-text="Delete"
      type="danger"
      @close="showDeleteModal = false"
      @confirm="executeDelete"
    />

    <velo-modal 
      :show="showClearModal" 
      title="Clear Database"
      message="WARNING: This will delete ALL test history permanently. This action is irreversible. Are you sure?"
      confirm-text="Clear All"
      type="danger"
      @close="showClearModal = false"
      @confirm="executeClearHistory"
    />

    <velo-toast ref="toast" />
  </div>
</template>

<script>
import SpeedChart from './components/SpeedChart.vue'
import BarChart from './components/BarChart.vue'
import Card from './components/Card.vue'
import VeloDropdown from './components/VeloDropdown.vue'
import VeloModal from './components/VeloModal.vue'
import VeloToast from './components/VeloToast.vue'
import moment from 'moment'

export default {
  name: 'App',
  components: {
    SpeedChart,
    BarChart,
    Card,
    VeloDropdown,
    VeloModal,
    VeloToast
  },
  data() {
    return {
      activeTab: 'dashboard',
      loading: false,
      datacollection: null,
      fullHistory: [], // Store full data for filtering
      selectedChartType: 'area', 
      selectedScope: '24h',
      selectedIP: 'all',
      availableIPs: [],
      lastUpdated: null,
      // Modal state
      showDeleteModal: false,
      showClearModal: false,
      itemToDelete: null,
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
                position: 'bottom',
                labels: {
                    usePointStyle: true,
                    padding: 20,
                    font: {
                        size: 12,
                        weight: '500'
                    }
                }
            },
            tooltip: {
                mode: 'index',
                intersect: false,
                backgroundColor: 'rgba(255, 255, 255, 0.9)',
                titleColor: '#1a1b1e',
                bodyColor: '#495057',
                borderColor: '#e9ecef',
                borderWidth: 1,
                padding: 12,
                boxPadding: 6,
                usePointStyle: true,
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
                ticks: {
                    autoSkip: true,
                    maxTicksLimit: 8,
                    color: '#adb5bd',
                    font: {
                        size: 11
                    }
                }
            },
            y: {
                grid: {
                    borderDash: [4, 4],
                    color: "#f1f3f5",
                    drawBorder: false
                },
                ticks: {
                    color: '#adb5bd',
                    font: {
                        size: 11
                    },
                    callback: function(value) {
                        return value + ' Mbps';
                    }
                },
                beginAtZero: true
            }
        }
      }
    }
  },
  computed: {
    ipOptions() {
        const opts = [{ label: 'All Interfaces', value: 'all' }];
        this.availableIPs.forEach(ip => {
            opts.push({ label: ip, value: ip });
        });
        return opts;
    },
    sortedFullHistory() {
      return [...this.fullHistory].sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
    }
  },
  mounted() {
    this.initWails();
  },
  methods: {
    initWails() {
        this.fetchHistory()

        if (window.runtime) {
            window.runtime.EventsOn("measurement_complete", (_result) => {
                this.loading = false
                this.fetchHistory()
                this.$refs.toast.add('Speed test completed');
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
            .then(_result => {
                // Event listener will handle update
            })
            .catch(err => {
                console.error(err);
                this.loading = false;
            })
      } else {
          // Simulate for dev without backend
          setTimeout(() => {
              this.loading = false;
              this.fetchHistory();
          }, 2000);
      }
    },
    fetchHistory() {
      if (window.go && window.go.backend && window.go.backend.VeloApp) {
          // When in History tab, always fetch 'all'
          const scope = this.activeTab === 'history' ? 'all' : this.selectedScope;
          window.go.backend.VeloApp.GetHistory(scope).then(history => {
              this.fullHistory = history || [];
              this.updateAvailableIPs(this.fullHistory);
              this.processHistory(this.fullHistory)
          })
      }
    },
    deleteItem(id) {
      this.itemToDelete = id;
      this.showDeleteModal = true;
    },
    executeDelete() {
      if (this.itemToDelete && window.go && window.go.backend && window.go.backend.VeloApp) {
        window.go.backend.VeloApp.DeleteMeasurement(this.itemToDelete).then(success => {
          this.showDeleteModal = false;
          this.itemToDelete = null;
          if (success) {
            this.fetchHistory();
            this.$refs.toast.add('Measurement deleted successfully');
          } else {
            this.$refs.toast.add('Failed to delete measurement', 'error');
          }
        });
      }
    },
    confirmClearHistory() {
      this.showClearModal = true;
    },
    executeClearHistory() {
      if (window.go && window.go.backend && window.go.backend.VeloApp) {
        window.go.backend.VeloApp.ClearHistory().then(success => {
          this.showClearModal = false;
          if (success) {
            this.fetchHistory();
            this.$refs.toast.add('All history cleared successfully');
          } else {
            this.$refs.toast.add('Failed to clear history', 'error');
          }
        });
      }
    },
    formatFullDate(ts) {
      return moment(ts).format('YYYY-MM-DD HH:mm:ss');
    },
    updateAvailableIPs(history) {
        const ips = new Set();
        history.forEach(h => {
            if (h.ip_address) ips.add(h.ip_address);
        });
        this.availableIPs = Array.from(ips);
        
        if (this.selectedIP !== 'all' && !this.availableIPs.includes(this.selectedIP)) {
            this.selectedIP = 'all';
        }
    },
    processHistory(historyData) {
        let history = historyData;
        if (this.selectedIP !== 'all') {
            history = history.filter(h => h.ip_address === this.selectedIP);
        }

        if (!Array.isArray(history) || history.length === 0) {
            this.datacollection = null;
            return;
        }

        history = [...history].sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

        const lastItem = history.slice(-1)[0];
        this.lastUpdated = moment(lastItem.timestamp).fromNow();

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
                    backgroundColor: isArea ? 'rgba(92, 124, 250, 0.1)' : (isBar ? 'rgba(92, 124, 250, 0.7)' : 'transparent'),
                    borderColor: '#5c7cfa',
                    pointBackgroundColor: '#fff',
                    pointBorderColor: '#5c7cfa',
                    pointBorderWidth: 2,
                    pointRadius: 4,
                    pointHoverRadius: 6,
                    borderWidth: 2,
                    fill: isArea,
                    tension: 0.4,
                    data: dlData
                },
                {
                    label: 'Upload',
                    backgroundColor: isArea ? 'rgba(32, 201, 151, 0.1)' : (isBar ? 'rgba(32, 201, 151, 0.7)' : 'transparent'),
                    borderColor: '#20c997',
                    pointBackgroundColor: '#fff',
                    pointBorderColor: '#20c997',
                    pointBorderWidth: 2,
                    pointRadius: 4,
                    pointHoverRadius: 6,
                    borderWidth: 2,
                    fill: isArea,
                    tension: 0.4,
                    data: ulData
                }
            ]
        }
    }
  },
  watch: {
      selectedChartType() {
          this.processHistory(this.fullHistory);
      },
      activeTab(newTab) {
        if (newTab === 'history') {
          this.fetchHistory();
        }
      }
  }
}
</script>

<style>
:root {
    --primary-color: #5c7cfa;
    --primary-hover: #4c6ef5;
    --bg-color: #f8f9fa;
    --sidebar-bg: #ffffff;
    --border-color: #e9ecef;
    --text-main: #1a1b1e;
    --text-muted: #909296;
    --sidebar-width: 260px;
}

body {
    background-color: var(--bg-color);
    margin: 0;
    overflow: hidden;
}

#app {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--text-main);
  display: flex;
  height: 100vh;
}

/* Sidebar Styles */
.sidebar {
    width: var(--sidebar-width);
    background-color: var(--sidebar-bg);
    border-right: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    z-index: 10;
}

.sidebar-header {
    padding: 24px;
}

.logo {
    display: flex;
    align-items: center;
    gap: 12px;
}

.logo-icon {
    width: 32px;
    height: 32px;
    background-color: var(--primary-color);
    color: white;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    font-size: 1.2rem;
}

.logo-text {
    font-size: 1.4rem;
    font-weight: 700;
    letter-spacing: -0.5px;
}

.sidebar-nav {
    flex: 1;
    padding: 0 12px;
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 16px;
    border-radius: 8px;
    cursor: pointer;
    color: var(--text-muted);
    font-weight: 500;
    transition: all 0.2s ease;
    margin-bottom: 4px;
}

.nav-icon {
    font-style: normal;
    font-size: 1.1rem;
    width: 20px;
    display: flex;
    justify-content: center;
}

.nav-item:hover {
    background-color: #f1f3f5;
    color: var(--text-main);
}

.nav-item.active {
    background-color: #edf2ff;
    color: var(--primary-color);
}

.sidebar-footer {
    padding: 24px;
    border-top: 1px solid var(--border-color);
}

/* Main Content Styles */
.main-content {
    flex: 1;
    overflow-y: auto;
    padding: 0;
    background-color: var(--bg-color);
}

.content-header {
    padding: 24px 40px;
    background-color: white;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: sticky;
    top: 0;
    z-index: 5;
}

.header-title h1 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 700;
    letter-spacing: -0.5px;
}

.header-title p {
    margin: 4px 0 0;
    color: var(--text-muted);
    font-size: 0.85rem;
}

.dashboard-grid {
    padding: 32px 40px;
    max-width: 1200px;
    margin: 0 auto;
}

/* Components */
.measure-btn {
    background-color: var(--primary-color);
    border: none;
    color: white;
    padding: 12px 20px;
    width: 100%;
    justify-content: center;
    text-align: center;
    display: inline-flex;
    align-items: center;
    font-size: 0.95rem;
    font-weight: 600;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.measure-btn:hover {
    background-color: var(--primary-hover);
    box-shadow: 0 4px 12px rgba(92, 124, 250, 0.2);
}

.measure-btn:disabled {
    background-color: #bac8ff;
    cursor: not-allowed;
    box-shadow: none;
}

.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
}

/* Table Styles */
.table-responsive {
  overflow-x: auto;
}

.velo-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.velo-table th {
  text-align: left;
  padding: 12px 16px;
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
  border-bottom: 2px solid var(--border-color);
}

.velo-table td {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-main);
  vertical-align: middle;
}

.velo-table tr:last-child td {
  border-bottom: none;
}

.ip-code {
  background-color: #f1f3f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.85rem;
  color: #495057;
}

.text-success { color: #20c997; }
.text-primary { color: #5c7cfa; }
.font-weight-bold { font-weight: 700; }
.text-right { text-align: right; }
.py-4 { padding-top: 2rem; padding-bottom: 2rem; }

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
  color: var(--text-muted);
}

.btn-icon.delete:hover {
  background-color: #fff5f5;
  color: #fa5252;
}

.icon-delete {
  font-style: normal;
  font-size: 1.2rem;
  line-height: 1;
}

/* Settings Styles */
.settings-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.setting-info label {
  display: block;
  font-weight: 600;
  margin-bottom: 4px;
  color: var(--text-main);
}

.setting-info p {
  margin: 0;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.btn-danger {
  background-color: #fa5252;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-danger:hover {
  background-color: #e03131;
}

.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
}

.badge-info {
  background-color: #e7f5ff;
  color: #228be6;
}

.controls-area {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-wrap: wrap;
    gap: 24px;
    margin-bottom: 24px;
}

.control-group {
    display: flex;
    align-items: center;
    gap: 12px;
}

.control-label {
    font-weight: 600;
    color: var(--text-muted);
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.btn-group {
    display: inline-flex;
    background-color: #f1f3f5;
    padding: 4px;
    border-radius: 8px;
}

.btn-group button {
    background-color: transparent;
    border: none;
    padding: 6px 16px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-muted);
    border-radius: 6px;
    transition: all 0.2s ease;
}

.btn-group button:hover {
    color: var(--text-main);
}

.btn-group button.active {
    background-color: white;
    color: var(--primary-color);
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.btn-group.mini button {
    padding: 4px 12px;
    font-size: 0.8rem;
}

.chart-wrapper {
    height: 450px;
    position: relative;
    padding: 0;
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
    color: var(--text-muted);
    font-style: italic;
}

.loader {
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  width: 14px;
  height: 14px;
  animation: spin 0.8s linear infinite;
  margin-right: 10px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Custom Scrollbar */
::-webkit-scrollbar {
  width: 8px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #dee2e6;
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: #ced4da;
}
</style>