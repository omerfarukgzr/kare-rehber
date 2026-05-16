<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import VChart from 'vue-echarts'
import PageHeader from '@/components/PageHeader.vue'
import { reportsApi, type CoachPerformanceRow, type CityDistributionRow, type SummaryReport, type WeekStatsRow } from '@/api/reports'

const summary = ref<SummaryReport | null>(null)
const coachPerf = ref<CoachPerformanceRow[]>([])
const cityRole = ref<'student' | 'coach'>('student')
const cityDist = ref<CityDistributionRow[]>([])
const weekStats = ref<WeekStatsRow[]>([])

async function fetchAll() {
  summary.value = await reportsApi.summary()
  coachPerf.value = await reportsApi.coachPerformance()
  cityDist.value = await reportsApi.cityDistribution(cityRole.value)
  weekStats.value = await reportsApi.weekStats()
}
onMounted(fetchAll)

async function reloadCity() {
  cityDist.value = await reportsApi.cityDistribution(cityRole.value)
}

const cityChartOption = computed(() => ({
  title: { text: cityRole.value === 'student' ? 'Öğrenci İl Dağılımı' : 'Koç İl Dağılımı', left: 'center' },
  tooltip: { trigger: 'item' },
  series: [{
    name: 'Şehir',
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: true,
    data: cityDist.value.map(d => ({ name: d.city ?? 'Belirtilmemiş', value: d.count })),
  }],
}))

const weekChartOption = computed(() => ({
  title: { text: 'Haftalık Değerlendirme Sayıları', left: 'center' },
  tooltip: { trigger: 'axis' },
  legend: { data: ['Toplam', 'Onaylanmış', 'Aktif eşleştirme'], bottom: 0 },
  grid: { left: 50, right: 30, top: 50, bottom: 50 },
  xAxis: { type: 'category', data: weekStats.value.map(w => `${w.weekNo}`) },
  yAxis: { type: 'value' },
  series: [
    { name: 'Toplam', type: 'bar', data: weekStats.value.map(w => w.totalEvaluations) },
    { name: 'Onaylanmış', type: 'bar', data: weekStats.value.map(w => w.approvedCount) },
    { name: 'Aktif eşleştirme', type: 'line', data: weekStats.value.map(w => w.activeAssignments) },
  ],
}))

const coachPerfChartOption = computed(() => {
  const top = coachPerf.value.slice(0, 15)
  return {
    title: { text: 'Koç Performansı (en aktif 15)', left: 'center' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { data: ['Değerlendirme', 'Onaylanmış'], bottom: 0 },
    grid: { left: 200, right: 30, top: 50, bottom: 50 },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: top.map(c => c.coachName).reverse(), inverse: false },
    series: [
      { name: 'Değerlendirme', type: 'bar', data: top.map(c => c.evaluationCount).reverse() },
      { name: 'Onaylanmış', type: 'bar', data: top.map(c => c.approvedCount).reverse() },
    ],
  }
})
</script>

<template>
  <div>
    <PageHeader title="Raporlar" subtitle="Genel performans ve istatistikler" />

    <el-row v-if="summary" :gutter="16" class="stats">
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary.activeStudents }}</div><div class="label">Aktif Öğrenci</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary.activeCoaches }}</div><div class="label">Aktif Koç</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary.activeAssignments }}</div><div class="label">Aktif Eşleştirme</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary.pendingRegistrations }}</div><div class="label">Bekleyen Başvuru</div></div></el-card></el-col>
      <el-col :span="6" class="row2"><el-card><div class="stat"><div class="num">{{ summary.totalEvaluations }}</div><div class="label">Toplam Değerlendirme</div></div></el-card></el-col>
      <el-col :span="6" class="row2"><el-card><div class="stat"><div class="num">{{ summary.pendingEvaluations }}</div><div class="label">Bekleyen Değerlendirme</div></div></el-card></el-col>
      <el-col :span="6" class="row2"><el-card><div class="stat"><div class="num">{{ summary.approvedEvaluations }}</div><div class="label">Onaylanmış Değerlendirme</div></div></el-card></el-col>
      <el-col :span="6" class="row2"><el-card><div class="stat"><div class="num">{{ summary.openThreads }}</div><div class="label">Açık Mesaj</div></div></el-card></el-col>
    </el-row>

    <el-row :gutter="16" class="charts">
      <el-col :span="12">
        <el-card>
          <div class="chart-toolbar">
            <el-radio-group v-model="cityRole" size="small" @change="reloadCity">
              <el-radio-button value="student">Öğrenci</el-radio-button>
              <el-radio-button value="coach">Koç</el-radio-button>
            </el-radio-group>
          </div>
          <VChart class="chart" :option="cityChartOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <VChart class="chart" :option="weekChartOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <el-card class="mt">
      <VChart class="chart-tall" :option="coachPerfChartOption" autoresize />
    </el-card>

    <el-card class="mt">
      <h3>Koç Performans Tablosu</h3>
      <el-table :data="coachPerf" stripe>
        <el-table-column label="Koç" prop="coachName" />
        <el-table-column label="Öğrenci" prop="studentCount" width="100" />
        <el-table-column label="Değerlendirme" prop="evaluationCount" width="140" />
        <el-table-column label="Onaylı" prop="approvedCount" width="100" />
        <el-table-column label="Ort. Motivasyon" width="160">
          <template #default="{ row }">{{ row.avgMotivation ? row.avgMotivation.toFixed(2) : '—' }}</template>
        </el-table-column>
        <el-table-column label="Ort. Davranış" width="140">
          <template #default="{ row }">{{ row.avgBehavior ? row.avgBehavior.toFixed(2) : '—' }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.stats { margin-bottom: 16px; }
.stats .row2 { margin-top: 12px; }
.stat { text-align: center; padding: 8px 0; }
.stat .num { font-size: 28px; font-weight: 700; color: #1f2937; }
.stat .label { color: #6b7280; font-size: 12px; margin-top: 4px; }
.charts { margin-bottom: 16px; }
.chart { height: 360px; }
.chart-tall { height: 460px; }
.chart-toolbar { display: flex; justify-content: flex-end; margin-bottom: 8px; }
.mt { margin-top: 16px; }
</style>
