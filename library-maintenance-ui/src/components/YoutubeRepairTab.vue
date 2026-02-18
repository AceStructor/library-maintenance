<template>
  <div>
    <h2>Youtube Code Repair</h2>

    <!-- Search Section -->
    <div>
      <input
        v-model="artist"
        placeholder="Artist eingeben"
        type="text"
      />
      <button @click="search" :disabled="loadingSearch">
        {{ loadingSearch ? "Suche..." : "Suchen" }}
      </button>
    </div>

    <!-- Result List -->
    <div v-if="results.length">
      <h3>Ergebnisse</h3>

      <table>
        <thead>
            <tr>
                <th>Artist</th>
                <th>Album</th>
                <th>Title</th>
                <th>Youtube Code</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="item in results"
                :key="item.track_id"
                :class="{ selected: selected && selected?.track_id === item.track_id }"
                class="hover:bg-gray-700 cursor-pointer"
                @click="selectItem(item)"
            >
                <td>{{ item.artist }}</td>
                <td>{{ item.album }}</td>
                <td>{{ item.title }}</td>
                <td>{{ item.youtube_code }}</td>
                <td>{{ item.download_status }}</td>
            </tr>
        </tbody>
      </table>
    </div>

    <!-- Detail Section -->
    <div 
        v-if="selected"
        class="details detail-section flex flex-col gap-2"
    >
        <h3>Details</h3>

        <div class="detail-row">
            <span class="label">Artist: </span>
            <span class="value">{{ selected.artist }}</span>
        </div>
        <div class="detail-row">
            <span class="label">Title: </span>
            <span class="value">{{ selected.title }}</span>
        </div>

        <div class="detail-row">
            <div class="detail-row">
                <span class="label">Album: </span>
                <span class="value">{{ selected.album }}</span>
            </div>
            <span class="value">{{ selected.album_mbid }}</span>
            <hr>
        </div>

        <label>Youtube Code</label>
        <input
        type="text"
        v-model="editedYoutubeCode"
        />

        <div class="flex items-center gap-2">
            <input type="checkbox" v-model="retryDownload" />
            <label>Download wiederholen</label>
        </div>

        <button
            @click="retry"
            :disabled="loadingRetry || !selected || (!hasChanged && !retryDownload)"
        >
            {{ loadingRetry ? "Starte..." : "Speichern & Aktion ausführen" }}
        </button>
    </div>

    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { youtubeApi } from "../api/client";

interface TrackResult {
  track_id: number;
  artist: string;
  album: string;
  album_mbid: string;
  title: string;
  track_mbid: string;
  youtube_code: string;
  download_status: string;
  file_path: string;
}

const artist = ref("");
const results = ref<TrackResult[]>([]);
const selected = ref<TrackResult | null>(null);

const retryDownload = ref(false);
const editedYoutubeCode = ref("");

const loadingSearch = ref(false);
const lastSearchArtist = ref("");
const loadingRetry = ref(false);
const error = ref("");

const search = async () => {
  if (!artist.value) return;

  loadingSearch.value = true;
  error.value = "";
  results.value = [];
  selected.value = null;

  lastSearchArtist.value = artist.value;

  try {
    const response = await youtubeApi.post("/youtube/search", {
      artist: artist.value,
    });

    results.value = response.data;
  } catch (err: any) {
    error.value = err.response?.data || err.message;
  } finally {
    loadingSearch.value = false;
  }
};

const refreshResults = async () => {
  if (!lastSearchArtist.value) return;

  try {
    const response = await youtubeApi.post("/youtube/search", {
      artist: lastSearchArtist.value,
    });

    results.value = response.data;

    // optional: Auswahl wiederherstellen
    if (selected.value) {
      const updated = results.value.find(
        r => r.track_mbid === selected.value?.track_mbid
      );

      if (updated) {
        selected.value = updated;
        editedYoutubeCode.value = updated.youtube_code ?? "";
      } else {
        selected.value = null;
      }
    }

  } catch (err: any) {
    error.value = err.response?.data || err.message;
  }
};

const selectItem = (item: TrackResult) => {
  selected.value = item;
  editedYoutubeCode.value = item.youtube_code ?? "";
};

const hasChanged = computed(() => {
  if (!selected.value) return false;
  return editedYoutubeCode.value !== (selected.value.youtube_code ?? "");
});

const retry = async () => {
  if (!selected.value) return;

  loadingRetry.value = true;
  error.value = "";

  try {
    await youtubeApi.post("/youtube/retry", {
      track_mbid: selected.value.track_mbid,
      youtube_code: editedYoutubeCode.value,
      retry_download: retryDownload.value,
    });

    alert("Aktion erfolgreich gestartet");
  } catch (err: any) {
    error.value = err.response?.data || err.message;
  } finally {
    loadingRetry.value = false;
  }

  retryDownload.value = false;
  await refreshResults();
};
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 15px;
}
td, th {
  border: 1px solid #ddd;
  padding: 8px;
}
tr:hover {
  background-color: #434547;
  cursor: pointer;
}
.selected {
  background-color: #643a3a;
}
.detail-section {
  margin-top: 20px;
}

.details input {
  margin-bottom: 0.5rem;
}
.details .checkbox-row {
  margin-top: 0.5rem;
}
</style>
