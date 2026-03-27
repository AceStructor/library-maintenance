<template>
  <div v-if="visible" class="overlay">
    <div class="modal">

      <h2 class="mb-4">Tracks hinzufügen</h2>

      <p class="mb-2">Playlist: {{ playlist?.name }}</p>

        <select
        v-model.number="selectedArtistId"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value="">Artist auswählen</option>
        <option
            v-for="ar in allArtists"
            :key="ar.artist_id"
            :value="ar.artist_id"
        >
            {{ ar.name }}
        </option>
        </select>

        <select
        v-if="selectedArtist"
        v-model.number="selectedAlbumId"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value="">Album auswählen</option>
        <option
            v-for="al in allAlbums"
            :key="al.album_id"
            :value="al.album_id"
        >
            {{ al.title }}
        </option>
        </select>

        <select
        v-if="selectedAlbum"
        v-model.number="selectedTrackId"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value="">Track auswählen</option>
        <option
            v-for="tr in allTracks"
            :key="tr.track_id"
            :value="tr.track_id"
        >
            {{ tr.title }}
        </option>
        </select>

        <button
            @click="addTrack" :disabled="!selectedTrack"
        >
            Hinzufügen
        </button>

        <div v-if="selectedTracks.length">
        <table>
        <thead>
            <tr>
                <th>Title</th>
                <th>Album</th>
                <th>Artist</th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="item in selectedTracks"
                :key="item.track_id"
            >
                <td>{{ item.title }}</td>
                <td>{{ item.album_name }}</td>
                <td>{{ item.artist_name }}</td>
            </tr>
        </tbody>
        </table>
        </div>

      <div class="actions">
        <button @click="submit" :disabled="!selectedTracks.length">{{ submitting ? "Speichert..." : "Speichern" }}</button>
        <button @click="close">Abbrechen</button>
      </div>

      <p v-if="error" style="color:red">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from "vue"
import { api } from "../api/client.js";


const error = ref("");

interface Playlist {
    playlist_id: number
    name: string
    navidromeId: string
    month: string
}

interface PlaylistTrack {
    id: number,
    title: string,
    artist: string,
    mbid: string,
    navidrome_id: string
}

interface Artist {
    artist_id: number
    name: string
}

interface Album {
    album_id: number
    title: string
}

interface Track {
    track_id: number
    title: string
}

interface TrackList {
    track_id: number
    artist_name: string
    album_name: string
    title: string
}

const props = defineProps<{
  visible: boolean
  playlist: Playlist | null
  playlistTracks: PlaylistTrack[] | []
}>()

const emit = defineEmits<{
  (e: "close"): void
}>()

const allArtists = ref<Artist[]>([]);
const selectedArtistId = ref<number | null>(null);
const selectedArtist = computed(() => {
    if (selectedArtistId.value === null) return null;
    return allArtists.value.find(a => a.artist_id === selectedArtistId.value) || null;
});
const allAlbums = ref<Album[]>([]);
const selectedAlbumId = ref<number | null>(null);
const selectedAlbum = computed(() => {
    if (selectedAlbumId.value === null) return null;
    return allAlbums.value.find(a => a.album_id === selectedAlbumId.value) || null;
});
const allTracks = ref<Track[]>([]);
const selectedTrackId = ref<number | null>(null);
const selectedTrack = computed(() => {
    if (selectedTrackId.value === null) return null;
    return allTracks.value.find(t => t.track_id === selectedTrackId.value) || null;
});
const submitting = ref(false);
const selectedTracks = ref<TrackList[]>([]);

watch(() => props.playlist, () => {
  selectedArtistId.value = null;
  selectedAlbumId.value = null;
  selectedTrackId.value = null;
  selectedTracks.value = [];
})

watch(() => selectedArtist.value, async () => {
    if (selectedArtist.value) {
        await loadArtistAlbums()
    }
    else {
        allAlbums.value = [];
    }
    allTracks.value = [];
    selectedAlbumId.value = null;
    selectedTrackId.value = null;
})

watch(() => selectedAlbum.value, async () => {
    if (selectedAlbum.value) {
        await loadAlbumTracks()
    }
    else {
        allTracks.value = [];
    }
    selectedTrackId.value = null;
})

function close() {
    selectedArtistId.value = null;
    selectedAlbumId.value = null;
    selectedTrackId.value = null;
    selectedTracks.value = [];
    allTracks.value = [];
    allAlbums.value = [];

  emit("close")
}

function addTrack() {
    if (!selectedTrack.value || !selectedArtist.value || !selectedAlbum.value) {
        return;
    }
    selectedTracks.value.push({
        "track_id": selectedTrack.value.track_id,
        "artist_name": selectedArtist.value.name,
        "album_name": selectedAlbum.value.title,
        "title": selectedTrack.value.title
    })

    selectedTrackId.value = null;
}

async function submit() {
    if (!props.playlist) {
        return;
    }
    const trackIds = new Set(selectedTracks.value.map(t => t.track_id));

    selectedTracks.value = [];
    submitting.value = true;

    try {
      await api.post(`/playlist/${props.playlist.playlist_id}/tracks/add`, {
        "tracks": Array.from(trackIds),
      });

    } catch (err: any) {
      error.value = err.message;
    }

    submitting.value = false;
    close()
}

async function loadAllArtists() {
    try {
      const res = await api.get(`/artist/all`);

      allArtists.value = await res.data["artists"];

    } catch (err: any) {
      error.value = err.response?.data || err.message;
    }
}

async function loadArtistAlbums() {
    if (!selectedArtist.value) {
        return;
    }
    try {
      const res = await api.get(`/album/${selectedArtist.value.artist_id}/all`);

      allAlbums.value = await res.data["albums"];

    } catch (err: any) {
      error.value = err.response?.data || err.message;
    }
}

async function loadAlbumTracks() {
    if (!selectedAlbum.value) {
        return;
    }
    try {
      const res = await api.get(`/track/${selectedAlbum.value.album_id}/all`);

      allTracks.value = await res.data["tracks"];

    } catch (err: any) {
      error.value = err.response?.data || err.message;
    }

    allTracks.value = subtractTracks(allTracks.value, props.playlistTracks)
}

function subtractTracks(
  tracks: Track[],
  playlistTracks: PlaylistTrack[]
): Track[] {
  const playlistIds = new Set(playlistTracks.map(t => t.id));

  return tracks.filter(track => !playlistIds.has(track.track_id));
}

onMounted(() => {
    loadAllArtists();
})
</script>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #1f1f1f;
  padding: 20px;
  border-radius: 8px;
  width: 400px;
}

.input {
  width: 100%;
  margin-bottom: 1rem;
  padding: 8px;
  background: #333;
  border: 1px solid #555;
  color: white;
}

.actions {
  display: flex;
  gap: 10px;
}
</style>