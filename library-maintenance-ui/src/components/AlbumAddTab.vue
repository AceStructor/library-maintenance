<template>
  <div>
    <h2>Album hinzufügen</h2>

    <input
      v-model="mbid"
      placeholder="MBID eingeben"
      type="text"
    />

    <button @click="addAlbum" :disabled="loading">
      {{ loading ? "Sende..." : "Album hinzufügen" }}
    </button>

    <p v-if="message">{{ message }}</p>
    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { albumApi } from "../api/client";

const mbid = ref("");
const loading = ref(false);
const message = ref("");
const error = ref("");

const addAlbum = async () => {
  loading.value = true;
  error.value = "";
  message.value = "";

  try {
    await albumApi.post("/album", {
      mbid: mbid.value,
    });

    message.value = "Album erfolgreich hinzugefügt";
    mbid.value = "";
  } catch (err: any) {
    error.value = err.response?.data || err.message;
  } finally {
    loading.value = false;
  }
};
</script>
