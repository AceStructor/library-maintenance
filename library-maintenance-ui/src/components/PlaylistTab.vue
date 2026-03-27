<template>
  <div class="container">
    <h1>Playlists</h1>

    <div class="tabs">
      <button @click="activeTab = 'createplaylist'">Playlist hinzufügen</button>
      <button @click="activeTab = 'updateplaylist'">Playlist bearbeiten</button>
    </div>

    <PlaylistCreateTab v-if="activeTab === 'createplaylist'" />
    <PlaylistUpdateTab v-if="activeTab === 'updateplaylist'" />

    <button @click="sync" :disabled="synchronizing">
      {{ synchronizing ? "synchronisiert..." : "Synchronisieren" }}
    </button>
    
    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { api } from "../api/client.js";
import PlaylistCreateTab from "./PlaylistCreateTab.vue";
import PlaylistUpdateTab from "./PlaylistUpdateTab.vue";

const activeTab = ref<"createplaylist" | "updateplaylist">("createplaylist");
const synchronizing = ref(false);
const error = ref("");

const sync = async () => {
  synchronizing.value = true;

  try {
    await api.post("/playlist/sync");
  } catch (err: any) {
    error.value = err.message;
  } finally {
    synchronizing.value = false;
  }
};
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
