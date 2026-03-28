<template>
  <div>
    <h2>Add Album</h2>

    <input
      v-model="mbid"
      placeholder="Enter MBID"
      type="text"
    />

    <button @click="addAlbum" :disabled="loading">
      {{ loading ? "Sending..." : "Add Album" }}
    </button>

    <p v-if="message">{{ message }}</p>
    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { api } from "../api/client";

const mbid = ref("");
const loading = ref(false);
const message = ref("");
const error = ref("");

const addAlbum = async () => {
  loading.value = true;
  error.value = "";
  message.value = "";

  try {
    await api.post("/album", {
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
