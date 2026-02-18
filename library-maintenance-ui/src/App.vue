<template>
  <div class="container">
    <h1>Library Maintenance</h1>

    <div class="tabs">
      <button @click="activeTab = 'album'">Album hinzufügen</button>
      <button @click="activeTab = 'youtube'">Youtube Repair</button>
    </div>

    <AlbumAddTab v-if="activeTab === 'album'" />
    <YoutubeRepairTab v-if="activeTab === 'youtube'" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import AlbumAddTab from "./components/AlbumAddTab.vue";
import YoutubeRepairTab from "./components/YoutubeRepairTab.vue";

const activeTab = ref<"album" | "youtube">("album");
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
