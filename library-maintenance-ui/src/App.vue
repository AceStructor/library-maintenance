<template>
  <div class="container">
    <h1>Library Maintenance</h1>

    <div class="tabs">
      <button @click="activeTab = 'addalbum'">Album hinzufügen</button>
      <button @click="activeTab = 'delalbum'">Album entfernen</button>
      <button @click="activeTab = 'youtube'">Youtube Repair</button>
      <button @click="activeTab = 'genre'">Genre Repair</button>
      <button @click="activeTab = 'playlist'">Playlisten</button>
    </div>

    <AlbumAddTab v-if="activeTab === 'addalbum'" />
    <AlbumRemoveTab v-if="activeTab === 'delalbum'" />
    <YoutubeRepairTab v-if="activeTab === 'youtube'" />
    <GenreRepairTab v-if="activeTab === 'genre'" />
    <PlaylistTab v-if="activeTab === 'playlist'" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import AlbumAddTab from "./components/AlbumAddTab.vue";
import AlbumRemoveTab from "./components/AlbumRemoveTab.vue";
import YoutubeRepairTab from "./components/YoutubeRepairTab.vue";
import GenreRepairTab from "./components/GenreRepairTab.vue";
import PlaylistTab from "./components/PlaylistTab.vue";

const activeTab = ref<"addalbum" | "delalbum" | "youtube" | "genre" | "playlist">("addalbum");
</script>

<style>
.container {
  max-width: 900px;
  margin: auto;
}
.tabs {
  margin-bottom: 20px;
}
.tabs button {
  margin-right: 10px;
}
</style>
