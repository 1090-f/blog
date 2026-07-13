<template>
  <article class="weather-card card" aria-label="天气预报">
    <header class="weather-card-header">
      <div class="weather-title-wrap">
        <span class="weather-title-accent" aria-hidden="true"></span>
        <h4 class="weather-title">天气预报</h4>
      </div>
      <svg class="weather-sun" viewBox="0 0 48 48" aria-hidden="true">
        <g stroke="currentColor" stroke-width="3" stroke-linecap="round">
          <line x1="24" y1="3" x2="24" y2="9" />
          <line x1="24" y1="39" x2="24" y2="45" />
          <line x1="3" y1="24" x2="9" y2="24" />
          <line x1="39" y1="24" x2="45" y2="24" />
          <line x1="9.2" y1="9.2" x2="13.5" y2="13.5" />
          <line x1="34.5" y1="34.5" x2="38.8" y2="38.8" />
          <line x1="34.5" y1="13.5" x2="38.8" y2="9.2" />
          <line x1="9.2" y1="38.8" x2="13.5" y2="34.5" />
        </g>
        <circle cx="24" cy="24" r="11" fill="currentColor" />
      </svg>
    </header>

    <div class="weather-location">
      <span class="weather-pin" aria-hidden="true"><span></span></span>
      <span>{{ weather.location }}</span>
    </div>

    <div class="weather-current">
      <span class="weather-temperature">{{ weather.temperature }}<sup>°</sup></span>
      <span class="weather-condition">{{ weather.condition }}</span>
    </div>

    <div class="weather-range" aria-label="今日温度范围">
      <span class="weather-high">高温 {{ weather.high }}°C</span>
      <span class="weather-divider">/</span>
      <span class="weather-low">低温 {{ weather.low }}°C</span>
    </div>
  </article>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'

const WEATHER_REFRESH_INTERVAL = 30 * 60 * 1000
const FALLBACK_LOCATION = '新乡市 红旗区'
const FALLBACK_COORDINATES = { latitude: 35.303, longitude: 113.926 }

const weather = ref({
  location: FALLBACK_LOCATION,
  temperature: 36,
  condition: '多云',
  high: 37,
  low: 26
})
const coordinates = ref(FALLBACK_COORDINATES)

const weatherConditions = {
  0: '晴朗',
  1: '晴间多云',
  2: '多云',
  3: '阴天',
  45: '雾',
  48: '雾',
  51: '小雨',
  53: '小雨',
  55: '小雨',
  61: '小雨',
  63: '中雨',
  65: '大雨',
  71: '小雪',
  73: '中雪',
  75: '大雪',
  80: '阵雨',
  81: '阵雨',
  82: '强阵雨',
  95: '雷雨',
  96: '雷雨',
  99: '雷雨'
}

async function loadWeather() {
  try {
    const params = new URLSearchParams({
      latitude: String(coordinates.value.latitude),
      longitude: String(coordinates.value.longitude),
      current: 'temperature_2m,weather_code',
      daily: 'temperature_2m_max,temperature_2m_min',
      timezone: 'Asia/Shanghai',
      forecast_days: '1'
    })
    const response = await fetch(`https://api.open-meteo.com/v1/forecast?${params}`)
    if (!response.ok) throw new Error('weather request failed')

    const data = await response.json()
    weather.value = {
      ...weather.value,
      temperature: Math.round(data.current.temperature_2m),
      condition: weatherConditions[data.current.weather_code] || '多云',
      high: Math.round(data.daily.temperature_2m_max[0]),
      low: Math.round(data.daily.temperature_2m_min[0])
    }
  } catch {
    // Keep the design fallback when the public weather service is unavailable.
  }
}

function requestCurrentLocation() {
  if (!navigator.geolocation) return

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      coordinates.value = {
        latitude: coords.latitude,
        longitude: coords.longitude
      }
      weather.value = { ...weather.value, location: '当前位置' }
      loadWeather()
    },
    () => {
      // Keep the configured city when location permission is denied or unavailable.
    },
    {
      enableHighAccuracy: false,
      maximumAge: WEATHER_REFRESH_INTERVAL,
      timeout: 10000
    }
  )
}

let refreshTimer

onMounted(() => {
  loadWeather()
  requestCurrentLocation()
  refreshTimer = window.setInterval(loadWeather, WEATHER_REFRESH_INTERVAL)
})

onUnmounted(() => {
  window.clearInterval(refreshTimer)
})
</script>

<style scoped>
.weather-card {
  --weather-accent: var(--accent);
  padding: 22px;
  overflow: hidden;
  border-radius: var(--radius);
  color: var(--text-primary);
  background: var(--bg-card);
  border-color: var(--border);
}

.weather-card:hover {
  background: var(--bg-card-hover);
}

.weather-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.weather-title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.weather-title-accent {
  width: 4px;
  height: 21px;
  flex: 0 0 auto;
  border-radius: 999px;
  background: var(--weather-accent);
}

.weather-title {
  font-size: 19px;
  font-weight: 600;
}

.weather-sun {
  width: 36px;
  height: 36px;
  margin-top: 0;
  flex: 0 0 auto;
  color: var(--weather-accent);
}

.weather-location {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 16px;
  color: var(--text-secondary);
  font-size: 15px;
  font-weight: 600;
}

.weather-pin {
  position: relative;
  width: 12px;
  height: 16px;
  flex: 0 0 auto;
  border-radius: 50% 50% 50% 0;
  background: var(--weather-accent);
  transform: rotate(-45deg);
}

.weather-pin span {
  position: absolute;
  width: 4px;
  height: 4px;
  top: 4px;
  left: 4px;
  border-radius: 50%;
  background: var(--bg-card);
}

.weather-current {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-top: 22px;
}

.weather-temperature {
  color: var(--text-primary);
  font-size: 48px;
  font-weight: 600;
  line-height: .95;
  letter-spacing: -2px;
}

.weather-temperature sup {
  margin-left: 3px;
  color: var(--text-muted);
  font-size: 18px;
  font-weight: 400;
  letter-spacing: 0;
  vertical-align: top;
}

.weather-condition {
  color: var(--text-secondary);
  font-size: 20px;
  font-weight: 400;
}

.weather-range {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 22px;
  font-size: 14px;
  white-space: nowrap;
}

.weather-high {
  color: var(--weather-accent);
}

.weather-low {
  color: var(--text-secondary);
}

.weather-divider {
  color: var(--text-muted);
}

@media (max-width: 1200px) {
  .weather-card {
    padding-inline: 16px;
  }

  .weather-range {
    gap: 6px;
    font-size: 14px;
  }
}
</style>
