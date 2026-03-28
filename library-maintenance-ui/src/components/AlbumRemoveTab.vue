<template>
  <div>
    <h2>Remove Album</h2>

    <input
      v-model="mbid"
      placeholder="Enter MBID"
      type="text"
    />

    <button @click="deleteAlbum" :disabled="loading">
      {{ loading ? "Sending..." : "Remove Album" }}
    </button>

    <p v-if="message">{{ message }}</p>
    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { api } from "../api/client.js";

const mbid = ref("");
const loading = ref(false);
const message = ref("");
const error = ref("");

const deleteAlbum = async () => {
  loading.value = true;
  error.value = "";
  message.value = "";

  try {
    await api.post("/album/delete", {
      mbid: mbid.value,
    });

    message.value = "Album successfully removed";
    mbid.value = "";
  } catch (err: any) {
    error.value = err.response?.data || err.message;
  } finally {
    loading.value = false;
  }
};
</script>
