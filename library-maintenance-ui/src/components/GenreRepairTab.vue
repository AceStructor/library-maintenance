<template>
  <div>
    <h2>Genre Repair</h2>

    <!-- Search Section -->
    <div>
        <div class="flex gap-4 mb-4">
            <label class="flex items-center gap-2">
                <input type="radio" value="search" v-model="mode" />
                Search artist
            </label>

            <label class="flex items-center gap-2">
                <input type="radio" value="no_genre" v-model="mode" />
                Artists without genre
            </label>
        </div>

        <div v-if="mode === 'search'" class="flex gap-2 mb-4">
            <input
                v-model="searchTerm"
                placeholder="Artist Name..."
                class="bg-gray-800 border border-gray-600 px-2 py-1 rounded w-64"
            />
            <button
                @click="searchArtists"
                class="bg-blue-600 px-3 py-1 rounded hover:bg-blue-700"
            >
                Search
            </button>
            </div>

            <div v-else class="mb-4">
            <button
                @click="loadArtistsWithoutGenre"
                class="bg-blue-600 px-3 py-1 rounded hover:bg-blue-700"
            >
                Load
            </button>
        </div>
    </div>

    <!-- Result List -->
    <div v-if="artists.length">
      <h3>Results</h3>

        <table class="w-full border-collapse mb-6">
            <thead>
                <tr class="bg-gray-800 text-left">
                <th class="p-2 border-b border-gray-700">Artist</th>
                <th class="p-2 border-b border-gray-700">Genres</th>
                </tr>
            </thead>
            <tbody>
                <tr
                v-for="artist in artists"
                :key="artist.id"
                @click="selectArtist(artist)"
                :class="{ selected: selectedArtist && selectedArtist?.id === artist.id }"
                class="hover:bg-gray-700 cursor-pointer"
                >
                <td class="p-2 border-b border-gray-800">
                    {{ artist.name }}
                </td>
                <td class="p-2 border-b border-gray-800">
                    {{ artist.genres.join(', ') }}
                </td>
                </tr>
            </tbody>
        </table>
    </div>

    <div v-if="selectedArtist" class="bg-gray-800 p-4 rounded shadow-lg">

        <h3 class="text-lg mb-4">Edit genre</h3>

        <!-- Artist Info -->
        <div class="grid grid-cols-[120px_1fr] gap-y-2 gap-x-4 mb-4">
            <div class="text-gray-400">Artist</div>
            <div>{{ selectedArtist.name }}</div>

            <div class="text-gray-400">Current genres</div>
            <div class="flex flex-wrap gap-2">
            <span class="bg-gray-700 px-2 py-1 rounded text-sm">
                {{ selectedArtist.genres.join(', ') }}
            </span>
            </div>
        </div>

        <!-- Edit Mode -->
        <div class="flex gap-6 mb-4">
            <label>
            <input type="radio" value="add" v-model="editMode" />
            Add genre
            </label>

            <label>
            <input type="radio" value="remove" v-model="editMode" />
            Remove genre
            </label>
        </div>

        <!-- Remove -->
        <div v-if="editMode === 'remove'" class="flex gap-2 items-center">
            <select
            v-model="selectedGenreToRemove"
            class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
            >
            <option disabled value="">Select genre</option>
            <option
                v-for="g in selectedArtist.genres"
                :key="g"
                :value="g"
            >
                {{ g }}
            </option>
            </select>

            <button
            @click="removeGenre"
            class="bg-red-600 px-3 py-1 rounded hover:bg-red-700"
            >
            Remove
            </button>
        </div>

        <!-- Add -->
        <div v-if="editMode === 'add'" class="flex flex-col gap-3">

            <!-- Bestehendes Genre -->
            <div class="flex gap-2 items-center">
            <select
                v-model="selectedGenreToAdd"
                class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
            >
                <option disabled value="">Select existing genre</option>
                <option
                v-for="g in allGenres"
                :key="g.id"
                :value="g.name"
                >
                {{ g.name }}
                </option>
            </select>

            <button
                @click="assignExistingGenre"
                class="bg-green-600 px-3 py-1 rounded hover:bg-green-700"
            >
                Add
            </button>
            </div>

            <!-- New genre -->
            <div class="flex gap-2 items-center">
            <input
                v-model="newGenreName"
                placeholder="New genre..."
                class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
            />

            <button
                @click="createAndAssignGenre"
                class="bg-green-600 px-3 py-1 rounded hover:bg-green-700"
            >
                Create new
            </button>
            </div>

        </div>

    </div>

    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { localApi } from "../api/client";

interface Artist {
    id: number
    name: string
    genres: string[]
}

interface Genre {
    id: number
    name: string
}

const mode = ref<'search' | 'no_genre'>('search');
const editMode = ref<'add' | 'remove'>('add');

const searchTerm = ref('');
const artists = ref<Artist[]>([]);
const selectedArtist = ref<Artist | null>(null);

const allGenres = ref<Genre[]>([]);

const selectedGenreToRemove = ref('');
const selectedGenreToAdd = ref<number | ''>('');
const newGenreName = ref('');
const error = ref("");

// =======================
// API Calls
// =======================

async function searchArtists() {
    const res = await localApi.post("/artistgenres", {
      artist: searchTerm.value,
    });
    artists.value = await res.data;
}

async function loadArtistsWithoutGenre() {
    const res = await localApi.post("/artistgenres/all/nogenre");
    artists.value = await res.data;
}

async function loadAllGenres() {
    const res = await localApi.post("/genres");
    allGenres.value = await res.data;
}

function selectArtist(artist: Artist) {
    selectedArtist.value = artist;
    selectedGenreToRemove.value = '';
    selectedGenreToAdd.value = '';
    newGenreName.value = '';
}

async function reloadSelectedArtist() {
    if (!selectedArtist.value) return;
    const res = await localApi.post("/artistgenres", {
      artist: selectedArtist.value.name,
    });
    artists.value = await res.data;
}

async function removeGenre() {
    if (!selectedArtist.value || !selectedGenreToRemove.value) return;

    await localApi.post("/artistgenres/deletebyname", {
        artist_id: selectedArtist.value.id,
        genre: selectedGenreToRemove.value,
    });

    await reloadSelectedArtist();
}

async function assignExistingGenre() {
    if (!selectedArtist.value || !selectedGenreToAdd.value) return;

    await localApi.post("/artistgenres/addbyname", {
        artist_id: selectedArtist.value.id,
        genre: selectedGenreToAdd.value,
    });

    await reloadSelectedArtist();
}

async function createAndAssignGenre() {
    if (!selectedArtist.value || !newGenreName.value) return;

    await localApi.post("/artistgenres/addbyname", {
        artist_id: selectedArtist.value.id,
        genre: newGenreName.value,
    });

    await reloadSelectedArtist();
}

onMounted(() => {
    loadAllGenres();
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
