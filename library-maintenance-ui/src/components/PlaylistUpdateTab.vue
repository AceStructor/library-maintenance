<template>
  <div>
    <h3>Update Playlist</h3>

    <div v-if=!editing>
      <div v-if="allPlaylists.length">
        <table>
          <thead>
              <tr>
                  <th>Name</th>
                  <th>Month</th>
              </tr>
          </thead>
          <tbody>
              <tr
                  v-for="item in allPlaylists"
                  :key="item.playlist_id"
                  :class="{ selected: selected && selected?.playlist_id === item.playlist_id }"
                  class="hover:bg-gray-700 cursor-pointer"
                  @click="selectItem(item)"
              >
                  <td>{{ item.name }}</td>
                  <td>{{ item.month }}</td>
              </tr>
          </tbody>
        </table>
      </div>

      <div 
          v-if="selected"
          class="details detail-section flex flex-col gap-2"
      >
          <button
              @click="editPlaylist"
          >
              Edit
          </button>

          <button
              @click="deletePlaylist"
          >
              Delete
          </button>
      </div>
    </div>

    <div v-if=editing>
      <div v-if="allTracks.length">
        <table>
          <thead>
              <tr>
                  <th>Title</th>
                  <th>Artist</th>
              </tr>
          </thead>
          <tbody>
              <tr
                  v-for="item in allTracks"
                  :key="item.id"
                  :class="{ selected: selectedTrack && selectedTrack?.id === item.id }"
                  class="hover:bg-gray-700 cursor-pointer"
                  @click="selectTrack(item)"
              >
                  <td>{{ item.title }}</td>
                  <td>{{ item.artist }}</td>
              </tr>
          </tbody>
        </table>
      </div>

      <div>
        <input
          v-model="playlistLength"
          placeholder="Enter length"
          type="number"
        />

        <button
              @click="fillPlaylist"
          >
              Fill
        </button>

		<button
              @click="addTracks"
          >
              Add tracks
        </button>
      </div>

      <div 
          v-if="selectedTrack"
          class="details detail-section flex flex-col gap-2"
      >
          <button
              @click="deleteTrack"
          >
              Delete
          </button>
      </div>

	  <button
			@click="back"
		>
			Back
		</button>
    </div>

    <p v-if="error" style="color:red">{{ error }}</p>

	<AddTrackModal
      :visible="overlayVisible"
      :playlist="selected"
	  :playlistTracks="allTracks"
      @close="closeOverlay"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { api } from "../api/client.js";
import AddTrackModal from "./AddTrackModal.vue";

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

const allPlaylists = ref<Playlist[]>([]);
const allTracks = ref<PlaylistTrack[]>([]);
const selected = ref<Playlist | null>(null);
const selectedTrack = ref<PlaylistTrack | null>(null);
const editing = ref(false);
const playlistLength = ref(allTracks.value.length)
const overlayVisible = ref(false)

const error = ref("");

const selectItem = (item: Playlist) => {
    selected.value = item;
};

const selectTrack = (item: PlaylistTrack) => {
    selectedTrack.value = item;
};

function addTracks() {
	overlayVisible.value = true;
}

async function editPlaylist() {
    editing.value = true;

    await loadAllTracks();
}

async function closeOverlay() {
	overlayVisible.value = false

	await loadAllTracks();
}

function back() {
    editing.value = false;
}

async function loadAllPlaylists() {
    try {
      const res = await api.get("/playlist/all");

      allPlaylists.value = await res.data["playlists"];

    } catch (err: any) {
      error.value = err.response?.data || err.message;
    }
}

async function loadAllTracks() {
  try {
      const res = await api.get(`/playlist/${selected.value?.playlist_id}/tracks`);

      allTracks.value = await res.data["tracks"];

    } catch (err: any) {
      error.value = err.response?.data || err.message;
    }

    playlistLength.value = allTracks.value.length;
	selectedTrack.value = null;
}

async function deletePlaylist() {
    try {
      await api.delete(`/playlist/${selected.value?.playlist_id}`);

    } catch (err: any) {
      error.value = err.message;
    }

    await loadAllPlaylists();
}

async function deleteTrack() {
	try {
		await api.delete(`/playlist/${selected.value?.playlist_id}/tracks/${selectedTrack.value?.id}`);

	} catch (err: any) {
		error.value = err.message;
	}

	await loadAllTracks();
}

async function fillPlaylist() {
  if (playlistLength.value < allTracks.value.length) {
    error.value = "Target length must be higher!";
    return false;
  }
  
  try {
      await api.post(`/playlist/${selected.value?.playlist_id}/fill`, {
        "length": playlistLength.value,
      });

    } catch (err: any) {
      error.value = err.message;
    }

    await loadAllTracks();
}

onMounted(() => {
    loadAllPlaylists();
})
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

